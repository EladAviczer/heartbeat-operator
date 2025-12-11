package main

import (
	"context"
	"log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

// RunReconcile orchestrates the update process
func RunReconcile(clientset *kubernetes.Clientset, cfg Config) {
	// 1. Check external state
	isHealthy := CheckHeavyDependency(cfg.DependencyURL)

	targetStatus := corev1.ConditionFalse
	if isHealthy {
		targetStatus = corev1.ConditionTrue
	}

	ctx := context.Background()

	// 2. List Pods
	pods, err := clientset.CoreV1().Pods(cfg.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: cfg.TargetLabel,
	})
	if err != nil {
		log.Printf("Failed to list pods: %v", err)
		return
	}

	// 3. Update Pods
	for _, pod := range pods.Items {
		// Optimization: Don't update if already correct
		if isConditionAlreadySet(&pod, targetStatus, cfg.GateName) {
			continue
		}

		// Use retry logic for safe updates
		retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			p, err := clientset.CoreV1().Pods(cfg.Namespace).Get(ctx, pod.Name, metav1.GetOptions{})
			if err != nil {
				return err
			}

			// Helper from pod_utils.go
			updatePodCondition(p, targetStatus, cfg.GateName)

			_, err = clientset.CoreV1().Pods(cfg.Namespace).UpdateStatus(ctx, p, metav1.UpdateOptions{})
			return err
		})

		if retryErr != nil {
			log.Printf("Failed to update pod %s: %v", pod.Name, retryErr)
		} else {
			log.Printf("Updated pod %s gate to %s", pod.Name, targetStatus)
		}
	}
}

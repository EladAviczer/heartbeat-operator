package main

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func updatePodCondition(pod *corev1.Pod, status corev1.ConditionStatus, gateName string) {
	for i, c := range pod.Status.Conditions {
		if c.Type == corev1.PodConditionType(gateName) {
			if c.Status != status {
				pod.Status.Conditions[i].Status = status
				pod.Status.Conditions[i].LastTransitionTime = metav1.Now()
				pod.Status.Conditions[i].Message = "Updated by Health Sentinel"
			}
			return
		}
	}
	pod.Status.Conditions = append(pod.Status.Conditions, corev1.PodCondition{
		Type:               corev1.PodConditionType(gateName),
		Status:             status,
		LastTransitionTime: metav1.Now(),
		Reason:             "SentinelCheck",
		Message:            "Updated by Health Sentinel",
	})
}

func isConditionAlreadySet(pod *corev1.Pod, status corev1.ConditionStatus, gateName string) bool {
	for _, c := range pod.Status.Conditions {
		if c.Type == corev1.PodConditionType(gateName) {
			return c.Status == status
		}
	}
	return false
}

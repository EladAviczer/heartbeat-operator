package main

import (
	"log"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	// 1. Load Config
	cfg := LoadConfig()
	log.Printf("Starting Health Sentinel with config: %+v", cfg)

	// 2. Connect to Cluster
	k8sConfig, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Error getting cluster config: %v", err)
	}
	clientset, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		log.Fatalf("Error creating clientset: %v", err)
	}

	log.Println("Health Sentinel started. Watching for:", cfg.TargetLabel)

	// 3. Continuous Loop
	ticker := time.NewTicker(cfg.CheckInterval)
	for range ticker.C {
		// logic is now in reconciler.go
		RunReconcile(clientset, cfg)
	}
}

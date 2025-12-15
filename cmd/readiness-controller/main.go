package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"readiness-controller/internal/config"
	"readiness-controller/internal/controller"
	"readiness-controller/internal/prober"
	"readiness-controller/internal/ui"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	// 1. Load Rules from JSON file
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "/etc/config/gates.json"
	}

	rules, err := config.LoadRules(configPath)
	if err != nil {
		log.Fatalf("Failed to load config from %s: %v", configPath, err)
	}
	log.Printf("Loaded %d gate rules", len(rules))

	// 2. Setup K8s Client
	k8sConfig, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Failed to get k8s config: %v", err)
	}
	clientset, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		log.Fatalf("Failed to create k8s client: %v", err)
	}

	// 3. Start UI
	ui.Start("8080")

	// 4. Start Controllers (One per Rule)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup

	for _, rule := range rules {
		wg.Add(1)

		// Capture variable for closure
		r := rule

		go func() {
			defer wg.Done()

			// Factory
			var p prober.Prober
			switch r.CheckType {
			case "http":
				p = prober.NewHttpProber(r.CheckTarget)
			case "tcp":
				p = prober.NewTcpProber(r.CheckTarget)
			case "exec":
				p = prober.NewExecProber(r.CheckTarget)
			default:
				log.Printf("[%s] Unknown CheckType '%s', skipping rule", r.Name, r.CheckType)
				return
			}

			ctrl := controller.New(clientset, r, p)
			ctrl.Start(ctx)
		}()
	}

	// 5. Wait for Shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("Shutting down...")
	cancel()
	wg.Wait()
}

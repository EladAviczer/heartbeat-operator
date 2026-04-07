package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"heartbeat-operator/internal/controller"
	"heartbeat-operator/internal/ui"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	k8sConfig, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Failed to get k8s config: %v", err)
	}

	// Maximize Kubernetes Client limits to handle high-throughput probe updates
	// The default is QPS=5 and Burst=10. We raise this significantly for scale.
	k8sConfig.QPS = 100.0
	k8sConfig.Burst = 250

	clientset, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		log.Fatalf("Failed to create k8s client: %v", err)
	}

	log.Println("Initializing Event Broadcaster...")
	eventBroadcaster := record.NewBroadcaster()

	eventBroadcaster.StartStructuredLogging(0)

	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{
		Interface: clientset.CoreV1().Events(""),
	})
	defer eventBroadcaster.Shutdown()

	recorder := eventBroadcaster.NewRecorder(
		scheme.Scheme,
		corev1.EventSource{Component: "heartbeat-operator"},
	)

	ui.Start("8080")

	metricsAddr := os.Getenv("METRICS_ADDR")
	if metricsAddr == "" {
		metricsAddr = ":9090"
	}

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		server := &http.Server{
			Addr:    metricsAddr,
			Handler: mux,
		}
		log.Printf("Starting metrics server on %s", metricsAddr)
		if err := server.ListenAndServe(); err != nil {
			log.Printf("Metrics server failed: %v", err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	crdClient, err := controller.NewCrdClient(k8sConfig, "") // Empty string watches all namespaces!
	if err != nil {
		log.Fatalf("Failed to create CRD client: %v", err)
	}

	manager := controller.NewManager(clientset, crdClient, recorder)
	go manager.Start(ctx)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("Shutting down...")
	cancel()
	time.Sleep(1 * time.Second)
}

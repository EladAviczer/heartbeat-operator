package main

import (
	"log"
	"os"
	"time"
)

// Config holds our dynamic settings
type Config struct {
	GateName      string
	TargetLabel   string
	DependencyURL string
	Namespace     string
	CheckInterval time.Duration
}

// LoadConfig reads from Env Vars with defaults
func LoadConfig() Config {
	intervalStr := getEnv("CHECK_INTERVAL", "5s")
	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		log.Printf("Invalid interval format '%s', defaulting to 5s", intervalStr)
		interval = 5 * time.Second
	}

	return Config{
		// UPDATED: Your specific naming convention
		GateName: getEnv("GATE_NAME", "controller.globus/swagger-ready"),

		// UPDATED: The label for your 16 replicas
		TargetLabel: getEnv("TARGET_LABEL", "app=dummy-stratego"),

		DependencyURL: getEnv("DEPENDENCY_URL", "http://swagger-api-server:3000/readiness"),
		Namespace:     getEnv("TARGET_NAMESPACE", "default"),
		CheckInterval: interval,
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

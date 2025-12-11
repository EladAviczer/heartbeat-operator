package main

import (
	"log"
	"net/http"
	"time"
)

// We define the target URL as a constant for now,
// but you could also move this to the Config struct if you want it dynamic.
// const dependencyURL = "http://swagger-api-server:3000/readiness"

// CheckHeavyDependency performs an HTTP GET to the Swagger service
func CheckHeavyDependency(url string) bool {
	// Create a client with a strict timeout.
	// If Swagger is slow, we treat it as "down" to be safe.
	client := http.Client{
		Timeout: 2 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Dependency Check FAILED: Could not reach %s: %v", url, err)
		return false
	}
	defer resp.Body.Close()

	// We assume a 200 OK means healthy.
	// Adjust this if your /readiness endpoint returns something else (e.g., JSON).
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return true
	}

	log.Printf("Dependency Check FAILED: %s returned status %d", url, resp.StatusCode)
	return false
}

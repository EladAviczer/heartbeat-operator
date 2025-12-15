package config

import (
	"encoding/json"
	"os"
	"time"
)

// GateRule represents one protected dependency
type GateRule struct {
	Name        string `json:"name"`        // Unique ID for the rule
	GateName    string `json:"gateName"`    // The string injected into Pods
	TargetLabel string `json:"targetLabel"` // Which pods to gate
	Namespace   string `json:"namespace"`   // Namespace of those pods
	CheckType   string `json:"checkType"`   // "http", "tcp", "exec"
	CheckTarget string `json:"checkTarget"` // URL or Address
	Interval    string `json:"interval"`    // "5s", "10s"
}

// LoadRules reads the JSON config file
func LoadRules(path string) ([]GateRule, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var rules []GateRule
	if err := json.Unmarshal(file, &rules); err != nil {
		return nil, err
	}
	return rules, nil
}

// Helper to parse duration strings safely
func ParseInterval(durationStr string) time.Duration {
	d, err := time.ParseDuration(durationStr)
	if err != nil {
		return 5 * time.Second // Default fallback
	}
	return d
}

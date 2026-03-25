package config

import (
	"encoding/json"
	"os"
	"time"
)

type ProbeRule struct {
	Name        string `json:"name"`
	Namespace   string `json:"namespace"`
	CheckType   string `json:"checkType"`
	CheckTarget string `json:"checkTarget"`
	Interval    string `json:"interval"`
	Timeout     string `json:"timeout"`
}

func LoadRules(path string) ([]ProbeRule, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var rules []ProbeRule
	if err := json.Unmarshal(file, &rules); err != nil {
		return nil, err
	}
	return rules, nil
}

func ParseInterval(durationStr string) time.Duration {
	d, err := time.ParseDuration(durationStr)
	if err != nil {
		return 5 * time.Second
	}
	return d
}

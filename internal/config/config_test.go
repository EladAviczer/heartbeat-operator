package config

import (
	"testing"
	"time"
)

func TestParseInterval(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Duration
	}{
		{"10s", 10 * time.Second},
		{"1m", 1 * time.Minute},
		{"invalid", 5 * time.Second}, // Default fallback
		{"", 5 * time.Second},        // Default fallback
	}

	for _, tt := range tests {
		got := ParseInterval(tt.input)
		if got != tt.expected {
			t.Errorf("ParseInterval(%q) = %v; want %v", tt.input, got, tt.expected)
		}
	}
}

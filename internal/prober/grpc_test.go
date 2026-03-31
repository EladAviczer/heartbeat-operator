package prober

import (
	"testing"
	"time"
)

func TestGrpcProber_InvalidTarget(t *testing.T) {
	p := NewGrpcProber("localhost:11111", 1*time.Second)
	if p.Check() {
		t.Errorf("Expected connection to invalid port to fail")
	}
}

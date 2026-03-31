package prober

import (
	"testing"
	"time"
)

func TestTlsProber(t *testing.T) {
	p := NewTlsProber("google.com:443", 5*time.Second)
	if !p.Check() {
		t.Errorf("Expected google.com:443 to succeed TLS check")
	}

	pInvalid := NewTlsProber("google.com:80", 2*time.Second)
	if pInvalid.Check() {
		t.Errorf("Expected google.com:80 to fail TLS check")
	}
}

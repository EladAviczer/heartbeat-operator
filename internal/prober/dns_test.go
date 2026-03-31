package prober

import (
	"testing"
	"time"
)

func TestDnsProber(t *testing.T) {
	p := NewDnsProber("google.com", 5*time.Second)
	if !p.Check() {
		t.Errorf("Expected google.com to resolve successfully")
	}

	pInvalid := NewDnsProber("this-domain-should-not-exist-12345.local", 2*time.Second)
	if pInvalid.Check() {
		t.Errorf("Expected invalid domain to fail resolution")
	}
}

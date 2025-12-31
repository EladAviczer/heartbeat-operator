package prober

import (
	"net"
	"testing"
)

func TestTcpProber_Check(t *testing.T) {
	// Test Success: Start a real listener
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test listener: %v", err)
	}
	defer ln.Close()

	prober := NewTcpProber(ln.Addr().String())
	if !prober.Check() {
		t.Errorf("expected true for reachable tcp port, got false")
	}

	// Test Failure: Connect to a closed port (hopefully free)
	// Using port 0 usually gives a free port, but we want a closed one.
	// Let's close the listener and try to connect to the same address.
	addr := ln.Addr().String()
	ln.Close()

	proberClosed := NewTcpProber(addr)
	if proberClosed.Check() {
		t.Errorf("expected false for closed tcp port, got true")
	}
}

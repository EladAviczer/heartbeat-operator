package prober

import (
	"net"
	"testing"
	"time"
)

func TestTcpProber_Check(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test listener: %v", err)
	}
	defer ln.Close()

	prober := NewTcpProber(ln.Addr().String(), 2*time.Second)
	if !prober.Check() {
		t.Errorf("expected true for reachable tcp port, got false")
	}

	addr := ln.Addr().String()
	ln.Close()

	proberClosed := NewTcpProber(addr, 2*time.Second)
	if proberClosed.Check() {
		t.Errorf("expected false for closed tcp port, got true")
	}
}

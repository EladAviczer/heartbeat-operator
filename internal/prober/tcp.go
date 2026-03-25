package prober

import (
	"log"
	"net"
	"time"
)

type TcpProber struct {
	Address string
	Timeout time.Duration
}

func NewTcpProber(addr string, timeout time.Duration) *TcpProber {
	return &TcpProber{Address: addr, Timeout: timeout}
}

func (p *TcpProber) Check() bool {
	conn, err := net.DialTimeout("tcp", p.Address, p.Timeout)
	if err != nil {
		log.Printf("[TCP] Connection to %s failed: %v", p.Address, err)
		return false
	}
	if conn != nil {
		conn.Close()
	}
	return true
}

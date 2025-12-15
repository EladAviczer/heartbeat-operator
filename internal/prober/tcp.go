package prober

import (
	"log"
	"net"
	"time"
)

type TcpProber struct {
	Address string
}

func NewTcpProber(addr string) *TcpProber {
	return &TcpProber{Address: addr}
}

func (p *TcpProber) Check() bool {
	conn, err := net.DialTimeout("tcp", p.Address, 2*time.Second)
	if err != nil {
		log.Printf("[TCP] Connection to %s failed: %v", p.Address, err)
		return false
	}
	if conn != nil {
		conn.Close()
	}
	return true
}

package prober

import (
	"crypto/tls"
	"log"
	"net"
	"time"
)

type TlsProber struct {
	HostPort string
	Timeout  time.Duration
}

func NewTlsProber(hostPort string, timeout time.Duration) *TlsProber {
	return &TlsProber{HostPort: hostPort, Timeout: timeout}
}

func (p *TlsProber) Check() bool {
	dialer := &net.Dialer{Timeout: p.Timeout}
	
	conn, err := tls.DialWithDialer(dialer, "tcp", p.HostPort, &tls.Config{
		InsecureSkipVerify: false,
	})
	if err != nil {
		log.Printf("[TLS] Check failed for %s: %v", p.HostPort, err)
		return false
	}
	defer conn.Close()

	state := conn.ConnectionState()
	for _, cert := range state.PeerCertificates {
		if time.Until(cert.NotAfter) < 7*24*time.Hour {
			log.Printf("[TLS] Certificate for %s expires soon: %v", p.HostPort, cert.NotAfter)
			return false
		}
	}

	return true
}

package prober

import (
	"context"
	"log"
	"net"
	"time"
)

type DnsProber struct {
	Host    string
	Timeout time.Duration
}

func NewDnsProber(host string, timeout time.Duration) *DnsProber {
	return &DnsProber{Host: host, Timeout: timeout}
}

func (p *DnsProber) Check() bool {
	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout)
	defer cancel()

	addrs, err := net.DefaultResolver.LookupHost(ctx, p.Host)
	if err != nil || len(addrs) == 0 {
		log.Printf("[DNS] Resolution failed for %s: %v", p.Host, err)
		return false
	}
	return true
}

package prober

import (
	"log"
	"net/http"
	"time"
)

type HttpProber struct {
	URL     string
	Timeout time.Duration
}

func NewHttpProber(url string, timeout time.Duration) *HttpProber {
	return &HttpProber{URL: url, Timeout: timeout}
}

func (p *HttpProber) Check() bool {
	client := http.Client{Timeout: p.Timeout}
	resp, err := client.Get(p.URL)
	if err != nil {
		log.Printf("[HTTP] Check failed for %s: %v", p.URL, err)
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

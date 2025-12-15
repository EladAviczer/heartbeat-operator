package prober

import (
	"log"
	"net/http"
	"time"
)

type HttpProber struct {
	URL string
}

func NewHttpProber(url string) *HttpProber {
	return &HttpProber{URL: url}
}

func (p *HttpProber) Check() bool {
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(p.URL)
	if err != nil {
		log.Printf("[HTTP] Check failed for %s: %v", p.URL, err)
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

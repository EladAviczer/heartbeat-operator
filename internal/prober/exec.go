package prober

import (
	"log"
	"os/exec"
)

type ExecProber struct {
	CommandStr string
}

func NewExecProber(cmdStr string) *ExecProber {
	return &ExecProber{CommandStr: cmdStr}
}

func (p *ExecProber) Check() bool {
	if p.CommandStr == "" {
		return false
	}
	cmd := exec.Command("sh", "-c", p.CommandStr)
	err := cmd.Run()
	if err != nil {
		log.Printf("[Exec] Command failed: %v", err)
		return false
	}
	return true
}

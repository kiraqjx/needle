package pkg

import (
	"fmt"
	"time"

	"github.com/go-ping/ping"
)

type Ping struct {
	target string
}

func NewPing(target string) *Ping {
	return &Ping{
		target: target,
	}
}

func (p *Ping) Probe() bool {
	pinger, err := ping.NewPinger(p.target)
	if err != nil {
		fmt.Println(err)
		return false
	}

	pinger.Count = 3
	pinger.Timeout = 500 * time.Millisecond
	pinger.SetPrivileged(true)
	err = pinger.Run()
	if err != nil {
		fmt.Println(err)
		return false
	}

	return pinger.Statistics().PacketsRecv != 0
}

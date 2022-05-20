package pkg

import (
	"net"
	"time"
)

type Telnet struct {
	target string
}

func NewTelnet(target string) *Telnet {
	return &Telnet{
		target: target,
	}
}

func (t *Telnet) Probe() bool {
	for i := 0; i < 3; i++ {
		conn, err := net.DialTimeout("tcp", t.target, time.Second)
		if err == nil {
			if conn != nil {
				conn.Close()
				return true
			}
		}
	}
	return false
}

package pkg

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"sync"
	"time"
)

type Function interface {
	SuccessInit()
	ErrorInit()
	SuccessToError()
	ErrorToSuccess()
}

type CmdFunction struct {
	config *JobConfig
}

func (sf *CmdFunction) SuccessInit() {
	sf.exec(sf.config.SuccessInit)
}

func (sf *CmdFunction) ErrorInit() {
	sf.exec(sf.config.SuccessInit)
}

func (sf *CmdFunction) SuccessToError() {
	sf.exec(sf.config.SuccessInit)
}

func (sf *CmdFunction) ErrorToSuccess() {
	sf.exec(sf.config.SuccessInit)
}

func (sf *CmdFunction) exec(cmdString string) {
	cmd := exec.Command(cmdString)
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}
}

type Job struct {
	config   *JobConfig
	function Function
	probe    Probe
	nowState bool
}

type ProbeType string

const (
	PING   ProbeType = "ping"
	TELNET ProbeType = "telnet"
)

type JobType string

const (
	CMD JobType = "cmd"
)

type Probe interface {
	Probe() bool
}

type JobConfig struct {
	Name           string    `yaml:"name"`
	Target         string    `yaml:"target"`
	Interval       int       `yaml:"interval"`
	ProbeType      ProbeType `yaml:"probe-type"`
	JobType        JobType   `yaml:"job-type"`
	SuccessInit    string    `yaml:"success-init"`
	ErrorInit      string    `yaml:"error-init"`
	SuccessToError string    `yaml:"success-to-error"`
	ErrorToSuccess string    `yaml:"error-to-success"`
}

func NewJob(jobConfig *JobConfig) (*Job, error) {
	var function Function
	if jobConfig.JobType == CMD {
		function = &CmdFunction{
			config: jobConfig,
		}
	} else {
		return nil, errors.New("the job type is not supported")
	}

	var probe Probe
	if jobConfig.ProbeType == PING {
		probe = NewPing(jobConfig.Target)
	} else {
		return nil, errors.New("the probe type is not supported")
	}

	return &Job{
		config:   jobConfig,
		function: function,
		probe:    probe,
	}, nil
}

func (j *Job) Run(ctx context.Context, wg *sync.WaitGroup) {
	fmt.Printf("start job: %s\n", j.config.Name)

	defer wg.Done()

	ticker := time.NewTicker(time.Duration(j.config.Interval) * time.Second)
	j.nowState = j.probe.Probe()

	// state init event
	if j.nowState {
		fmt.Println("init state: success")
		doWithConfig(j.config.SuccessInit, j.function.SuccessInit)
	} else {
		fmt.Println("init state: success")
		doWithConfig(j.config.ErrorInit, j.function.ErrorInit)
	}

LOOP:
	for {
		select {
		case <-ticker.C:
			nextState := j.probe.Probe()
			// state change
			if nextState != j.nowState {
				fmt.Printf("the state is change, true means success, false means error: %t -> %t\n", j.nowState, nextState)
				if j.nowState {
					doWithConfig(j.config.SuccessToError, j.function.SuccessToError)
				} else {
					doWithConfig(j.config.ErrorToSuccess, j.function.ErrorToSuccess)
				}
			}
			j.nowState = nextState
			continue
		case <-ctx.Done():
			break LOOP
		}
	}
}

func doWithConfig(config string, doing func()) {
	if len(config) != 0 {
		doing()
	}
}

package main

import (
	"context"
	"flag"
	"fmt"
	"needle/pkg"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"gopkg.in/yaml.v3"
)

var configPath string

var jobsConfig Jobs

type Jobs struct {
	Jobs []pkg.JobConfig `yaml:"jobs"`
}

func init() {
	flag.StringVar(&configPath, "config-file", "./config/jobs.yaml", "the config file path")
	flag.Parse()

	file, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&jobsConfig)
	if err != nil {
		panic(err)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	fmt.Println(jobsConfig.Jobs)
	for _, jobConfig := range jobsConfig.Jobs {
		job, err := pkg.NewJob(&jobConfig)
		if err != nil {
			fmt.Println(err)
			continue
		}
		wg.Add(1)
		go job.Run(ctx, &wg)
	}
	signalExit()
	cancel()
	wg.Wait()
}

func signalExit() {
	osc := make(chan os.Signal, 1)
	signal.Notify(osc, syscall.SIGTERM, syscall.SIGINT)
	<-osc
}

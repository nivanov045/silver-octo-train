package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nivanov045/silver-octo-train/cmd/agent/metricsagent"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

func main() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT)
	agent := metricsagent.New(pollInterval, reportInterval)
	agent.Start()
	<-sigc
}

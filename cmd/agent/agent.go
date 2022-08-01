package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v6"

	"github.com/nivanov045/silver-octo-train/cmd/agent/metricsagent"
)

func main() {
	log.Println("agent::main: started")

	var cfg metricsagent.Config

	flag.StringVar(&cfg.Address, "a", "127.0.0.1:8080", "address")
	log.Println(cfg.Address)
	flag.DurationVar(&cfg.PollInterval, "p", 5*time.Second, "poll interval")
	flag.DurationVar(&cfg.ReportInterval, "r", 10*time.Second, "report interval")
	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatalln("agent::main: error in env parsing:", err)
	}

	log.Println("agent::main: cfg:", cfg)
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT)
	agent := metricsagent.New(cfg)
	agent.Start()
	<-sigc
}

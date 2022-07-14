package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/nivanov045/silver-octo-train/cmd/agent/metricsperformer"
	"github.com/nivanov045/silver-octo-train/cmd/agent/requester"
	"github.com/nivanov045/silver-octo-train/internal/metrics"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

func updateMetrics(ch chan metrics.Metrics) {
	ticker := time.NewTicker(pollInterval)
	for {
		<-ticker.C
		m := <-ch
		metricsperformer.New().UpdateMetrics(m)
		ch <- m
	}
}

func sendMetrics(ch chan metrics.Metrics) {
	ticker := time.NewTicker(reportInterval)
	for {
		<-ticker.C
		m := <-ch
		for key, val := range m.GaugeMetrics {
			err := requester.New().Send("update/gauge/" + key + "/" + strconv.FormatFloat(float64(val), 'f', -1, 64))
			if err != nil {
				log.Fatal(err)
			}
		}
		err := requester.New().Send("update/counter/PollCount/" + strconv.FormatInt(int64(m.CounterMetrics["PollCount"]), 10))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Reported ", m.CounterMetrics["PollCount"])
		ch <- m
	}
}

func main() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT)
	metricsVal := metrics.Metrics{
		GaugeMetrics:   map[string]metrics.Gauge{},
		CounterMetrics: map[string]metrics.Counter{},
	}
	c := make(chan metrics.Metrics, 1)
	go updateMetrics(c)
	go sendMetrics(c)
	c <- metricsVal
	<-sigc
}

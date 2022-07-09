package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/nivanov045/silver-octo-train/cmd/agent/metrics_performer"
	"github.com/nivanov045/silver-octo-train/cmd/agent/requester"
	met "github.com/nivanov045/silver-octo-train/internal/metrics"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

func updateMetrics(ch chan met.Metrics) {
	ticker := time.NewTicker(pollInterval)
	for {
		<-ticker.C
		m := <-ch
		metrics_performer.New().UpdateMetrics(m)
		ch <- m
	}
}

func sendMetrics(ch chan met.Metrics) {
	ticker := time.NewTicker(reportInterval)
	for {
		<-ticker.C
		m := <-ch
		for key, val := range m.Gms {
			err := requester.New().Send("http://127.0.0.1:8080/" + "update/gauge/" + key + "/" + strconv.FormatFloat(float64(val), 'f', -1, 64))
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
		}
		err := requester.New().Send("http://127.0.0.1:8080/" + "update/counter/PollCount/" + strconv.FormatInt(int64(m.Cms["PollCount"]), 10))
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		//fmt.Println("Reported ", m.Cms["PollCount"])
		ch <- m
	}
}

func main() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT)
	var metricsVal met.Metrics
	metricsVal.Gms = make(map[string]met.Gauge)
	metricsVal.Cms = make(map[string]met.Counter)
	c := make(chan met.Metrics, 1)
	go updateMetrics(c)
	go sendMetrics(c)
	c <- metricsVal
	<-sigc
}

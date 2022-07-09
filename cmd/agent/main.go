package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

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
		met.UpdateMetrics(m)
		ch <- m
	}
}

func sendMetrics(ch chan met.Metrics) {
	ticker := time.NewTicker(reportInterval)
	for {
		<-ticker.C
		m := <-ch
		f := func(a string) {
			client := &http.Client{}
			request, err := http.NewRequest(http.MethodPost, a, nil)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			request.Header.Add("Content-Type", "text/plain")
			_, err = client.Do(request)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		for key, val := range m.Gms {
			f("http://127.0.0.1:8080/" + "update/gauge/" + key + "/" + strconv.FormatFloat(float64(val), 'f', -1, 64))
		}
		f("http://127.0.0.1:8080/" + "update/counter/PollCount/" + strconv.FormatInt(int64(m.Cms["PollCount"]), 10))
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

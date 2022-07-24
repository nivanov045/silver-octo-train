package metricsagent

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nivanov045/silver-octo-train/cmd/agent/metricsperformer"
	"github.com/nivanov045/silver-octo-train/cmd/agent/requester"
	"github.com/nivanov045/silver-octo-train/internal/metrics"
)

type metricsagent struct {
	Metrics        metrics.Metrics
	pollInterval   time.Duration
	reportInterval time.Duration
}

func (a *metricsagent) updateMetrics() {
	ticker := time.NewTicker(a.pollInterval)
	for {
		<-ticker.C
		metricsperformer.New().UpdateMetrics(a.Metrics)
	}
}

func (a *metricsagent) sendMetrics() {
	ticker := time.NewTicker(a.reportInterval)
	for {
		<-ticker.C
		for key, val := range a.Metrics.GaugeMetrics {
			asfloat := float64(val)
			metricForSend := metrics.MetricsInterface{
				ID:    key,
				MType: "gauge",
				Delta: nil,
				Value: &asfloat,
			}
			marshalled, err := json.Marshal(metricForSend)
			if err != nil {
				fmt.Println("Fatal 1")
				log.Fatal(err)
			}
			err = requester.New().Send(marshalled)
			if err != nil {
				fmt.Println("sendMetrics::46")
				log.Fatal(err)
			}
		}
		pc := a.Metrics.CounterMetrics["PollCount"]
		asint := int64(pc)
		metricForSend := metrics.MetricsInterface{
			ID:    "PollCount",
			MType: "counter",
			Delta: &asint,
			Value: nil,
		}
		marshalled, err := json.Marshal(metricForSend)
		if err != nil {
			fmt.Println("sendMetrics::59")
			log.Fatal(err)
		}
		err = requester.New().Send(marshalled)
		if err != nil {
			fmt.Println("sendMetrics::64")
			log.Fatal(err)
		}
		//fmt.Println("Reported ", a.Metrics.CounterMetrics["PollCount"])
	}
}

func (a *metricsagent) Start() {
	go a.updateMetrics()
	go a.sendMetrics()
}

func New(pollInterval time.Duration, reportInterval time.Duration) *metricsagent {
	return &metricsagent{
		Metrics: metrics.Metrics{
			GaugeMetrics:   map[string]metrics.Gauge{},
			CounterMetrics: map[string]metrics.Counter{},
		},
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
	}
}

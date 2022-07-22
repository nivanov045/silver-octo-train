package metricsagent

import (
	"fmt"
	"log"
	"strconv"
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
			err := requester.New().Send("update/gauge/" + key + "/" + strconv.FormatFloat(float64(val), 'f', -1, 64))
			if err != nil {
				log.Fatal(err)
			}
		}
		err := requester.New().Send("update/counter/PollCount/" + strconv.FormatInt(int64(a.Metrics.CounterMetrics["PollCount"]), 10))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Reported ", a.Metrics.CounterMetrics["PollCount"])
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

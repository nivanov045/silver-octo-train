package metricsperformer

import (
	"fmt"
	"math/rand"
	"runtime"

	met "github.com/nivanov045/silver-octo-train/internal/metrics"
)

type metricsPerformer struct{}

func New() *metricsPerformer {
	return &metricsPerformer{}
}

func (*metricsPerformer) UpdateMetrics(m met.Metrics) {
	var memStat runtime.MemStats
	runtime.ReadMemStats(&memStat)
	// Runtime metrics.
	m.Gms["Alloc"] = met.Gauge(memStat.Alloc)
	m.Gms["BuckHashSys"] = met.Gauge(memStat.BuckHashSys)
	m.Gms["Frees"] = met.Gauge(memStat.Frees)
	m.Gms["GCCPUFraction"] = met.Gauge(memStat.GCCPUFraction)
	m.Gms["GCSys"] = met.Gauge(memStat.GCSys)
	m.Gms["HeapAlloc"] = met.Gauge(memStat.HeapAlloc)
	m.Gms["HeapIdle"] = met.Gauge(memStat.HeapIdle)
	m.Gms["HeapInuse"] = met.Gauge(memStat.HeapInuse)
	m.Gms["HeapObjects"] = met.Gauge(memStat.HeapObjects)
	m.Gms["HeapReleased"] = met.Gauge(memStat.HeapReleased)
	m.Gms["HeapSys"] = met.Gauge(memStat.HeapSys)
	m.Gms["LastGC"] = met.Gauge(memStat.LastGC)
	m.Gms["Lookups"] = met.Gauge(memStat.Lookups)
	m.Gms["MCacheInuse"] = met.Gauge(memStat.MCacheInuse)
	m.Gms["MSpanInuse"] = met.Gauge(memStat.MSpanInuse)
	m.Gms["MSpanSys"] = met.Gauge(memStat.MSpanSys)
	m.Gms["Mallocs"] = met.Gauge(memStat.Mallocs)
	m.Gms["NextGC"] = met.Gauge(memStat.NextGC)
	m.Gms["NumForcedGC"] = met.Gauge(memStat.NumForcedGC)
	m.Gms["NumGC"] = met.Gauge(memStat.NumGC)
	m.Gms["OtherSys"] = met.Gauge(memStat.OtherSys)
	m.Gms["PauseTotalNs"] = met.Gauge(memStat.PauseTotalNs)
	m.Gms["StackInuse"] = met.Gauge(memStat.StackInuse)
	m.Gms["StackSys"] = met.Gauge(memStat.StackSys)
	m.Gms["Sys"] = met.Gauge(memStat.Sys)
	m.Gms["TotalAlloc"] = met.Gauge(memStat.TotalAlloc)

	// Other metrics.
	m.Gms["RandomValue"] = met.Gauge(rand.Float64())
	m.Cms["PollCount"]++
	fmt.Println("Updated ", m.Cms["PollCount"])
}

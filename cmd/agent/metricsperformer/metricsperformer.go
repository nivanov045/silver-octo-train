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
	for _, val := range met.KnownMetrics {
		switch val {
		case "Alloc":
			m.Gms["Alloc"] = met.Gauge(memStat.Alloc)
		case "BuckHashSys":
			m.Gms["BuckHashSys"] = met.Gauge(memStat.BuckHashSys)
		case "Frees":
			m.Gms["Frees"] = met.Gauge(memStat.Frees)
		case "GCCPUFraction":
			m.Gms["GCCPUFraction"] = met.Gauge(memStat.GCCPUFraction)
		case "GCSys":
			m.Gms["GCSys"] = met.Gauge(memStat.GCSys)
		case "HeapAlloc":
			m.Gms["HeapAlloc"] = met.Gauge(memStat.HeapAlloc)
		case "HeapIdle":
			m.Gms["HeapIdle"] = met.Gauge(memStat.HeapIdle)
		case "HeapInuse":
			m.Gms["HeapInuse"] = met.Gauge(memStat.HeapInuse)
		case "HeapObjects":
			m.Gms["HeapObjects"] = met.Gauge(memStat.HeapObjects)
		case "HeapReleased":
			m.Gms["HeapReleased"] = met.Gauge(memStat.HeapReleased)
		case "HeapSys":
			m.Gms["HeapSys"] = met.Gauge(memStat.HeapSys)
		case "LastGC":
			m.Gms["LastGC"] = met.Gauge(memStat.LastGC)
		case "Lookups":
			m.Gms["Lookups"] = met.Gauge(memStat.Lookups)
		case "MCacheInuse":
			m.Gms["MCacheInuse"] = met.Gauge(memStat.MCacheInuse)
		case "MSpanInuse":
			m.Gms["MSpanInuse"] = met.Gauge(memStat.MSpanInuse)
		case "MSpanSys":
			m.Gms["MSpanSys"] = met.Gauge(memStat.MSpanSys)
		case "Mallocs":
			m.Gms["Mallocs"] = met.Gauge(memStat.Mallocs)
		case "NextGC":
			m.Gms["NextGC"] = met.Gauge(memStat.NextGC)
		case "NumForcedGC":
			m.Gms["NumForcedGC"] = met.Gauge(memStat.NumForcedGC)
		case "NumGC":
			m.Gms["NumGC"] = met.Gauge(memStat.NumGC)
		case "OtherSys":
			m.Gms["OtherSys"] = met.Gauge(memStat.OtherSys)
		case "PauseTotalNs":
			m.Gms["PauseTotalNs"] = met.Gauge(memStat.PauseTotalNs)
		case "StackInuse":
			m.Gms["StackInuse"] = met.Gauge(memStat.StackInuse)
		case "StackSys":
			m.Gms["StackSys"] = met.Gauge(memStat.StackSys)
		case "Sys":
			m.Gms["Sys"] = met.Gauge(memStat.Sys)
		case "TotalAlloc":
			m.Gms["TotalAlloc"] = met.Gauge(memStat.TotalAlloc)
		case "RandomValue":
			m.Gms["RandomValue"] = met.Gauge(rand.Float64())
		case "PollCount":
			m.Cms["PollCount"]++
			fmt.Println("Updated ", m.Cms["PollCount"])
		}
	}
}

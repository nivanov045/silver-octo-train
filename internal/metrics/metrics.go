package metrics

import (
	"fmt"
	"math/rand"
	"runtime"
)

type Gauge float64
type Counter int64

type Metrics struct {
	Gms map[string]Gauge
	Cms map[string]Counter
}

func UpdateMetrics(m Metrics) {
	var memStat runtime.MemStats
	runtime.ReadMemStats(&memStat)
	// Runtime metrics.
	m.Gms["Alloc"] = Gauge(memStat.Alloc)
	m.Gms["BuckHashSys"] = Gauge(memStat.BuckHashSys)
	m.Gms["Frees"] = Gauge(memStat.Frees)
	m.Gms["GCCPUFraction"] = Gauge(memStat.GCCPUFraction)
	m.Gms["GCSys"] = Gauge(memStat.GCSys)
	m.Gms["HeapAlloc"] = Gauge(memStat.HeapAlloc)
	m.Gms["HeapIdle"] = Gauge(memStat.HeapIdle)
	m.Gms["HeapInuse"] = Gauge(memStat.HeapInuse)
	m.Gms["HeapObjects"] = Gauge(memStat.HeapObjects)
	m.Gms["HeapReleased"] = Gauge(memStat.HeapReleased)
	m.Gms["HeapSys"] = Gauge(memStat.HeapSys)
	m.Gms["LastGC"] = Gauge(memStat.LastGC)
	m.Gms["Lookups"] = Gauge(memStat.Lookups)
	m.Gms["MCacheInuse"] = Gauge(memStat.MCacheInuse)
	m.Gms["MSpanInuse"] = Gauge(memStat.MSpanInuse)
	m.Gms["MSpanSys"] = Gauge(memStat.MSpanSys)
	m.Gms["Mallocs"] = Gauge(memStat.Mallocs)
	m.Gms["NextGC"] = Gauge(memStat.NextGC)
	m.Gms["NumForcedGC"] = Gauge(memStat.NumForcedGC)
	m.Gms["NumGC"] = Gauge(memStat.NumGC)
	m.Gms["OtherSys"] = Gauge(memStat.OtherSys)
	m.Gms["PauseTotalNs"] = Gauge(memStat.PauseTotalNs)
	m.Gms["StackInuse"] = Gauge(memStat.StackInuse)
	m.Gms["StackSys"] = Gauge(memStat.StackSys)
	m.Gms["Sys"] = Gauge(memStat.Sys)
	m.Gms["TotalAlloc"] = Gauge(memStat.TotalAlloc)

	// Other metrics.
	m.Gms["RandomValue"] = Gauge(rand.Float64())
	m.Cms["PollCount"]++
	fmt.Println("Updated ", m.Cms["PollCount"])
}

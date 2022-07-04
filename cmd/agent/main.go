package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

type gauge float64
type counter int64

type metrics struct {
	gms map[string]gauge
	cm  struct {
		name  string
		value counter
	}
}

func updateMetrics(ch chan metrics) {
	ticker := time.NewTicker(pollInterval)
	for {
		<-ticker.C
		m := <-ch
		var memStat runtime.MemStats
		runtime.ReadMemStats(&memStat)
		// Runtime metrics.
		m.gms["Alloc"] = gauge(memStat.Alloc)
		m.gms["BuckHashSys"] = gauge(memStat.BuckHashSys)
		m.gms["Frees"] = gauge(memStat.Frees)
		m.gms["GCCPUFraction"] = gauge(memStat.GCCPUFraction)
		m.gms["GCSys"] = gauge(memStat.GCSys)
		m.gms["HeapAlloc"] = gauge(memStat.HeapAlloc)
		m.gms["HeapIdle"] = gauge(memStat.HeapIdle)
		m.gms["HeapInuse"] = gauge(memStat.HeapInuse)
		m.gms["HeapObjects"] = gauge(memStat.HeapObjects)
		m.gms["HeapReleased"] = gauge(memStat.HeapReleased)
		m.gms["HeapSys"] = gauge(memStat.HeapSys)
		m.gms["LastGC"] = gauge(memStat.LastGC)
		m.gms["Lookups"] = gauge(memStat.Lookups)
		m.gms["MCacheInuse"] = gauge(memStat.MCacheInuse)
		m.gms["MSpanInuse"] = gauge(memStat.MSpanInuse)
		m.gms["MSpanSys"] = gauge(memStat.MSpanSys)
		m.gms["Mallocs"] = gauge(memStat.Mallocs)
		m.gms["NextGC"] = gauge(memStat.NextGC)
		m.gms["NumForcedGC"] = gauge(memStat.NumForcedGC)
		m.gms["NumGC"] = gauge(memStat.NumGC)
		m.gms["OtherSys"] = gauge(memStat.OtherSys)
		m.gms["PauseTotalNs"] = gauge(memStat.PauseTotalNs)
		m.gms["StackInuse"] = gauge(memStat.StackInuse)
		m.gms["StackSys"] = gauge(memStat.StackSys)
		m.gms["Sys"] = gauge(memStat.Sys)
		m.gms["TotalAlloc"] = gauge(memStat.TotalAlloc)
		//fmt.Println("Updated ", m.cm.value)

		// Other metrics.
		m.cm.value++
		m.gms["RandomValue"] = gauge(rand.Float64())
		ch <- m
	}
}

func sendMetrics(ch chan metrics) {
	ticker := time.NewTicker(reportInterval)
	for {
		<-ticker.C
		m := <-ch
		for key, val := range m.gms {
			client := &http.Client{}
			address := "http://127.0.0.1:8080/" + "update/gauge/" + key + "/" + strconv.FormatFloat(float64(val), 'f', -1, 64)
			request, err := http.NewRequest(http.MethodPost, address, nil)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			request.Header.Add("Content-Type", "text/plain")
			//fmt.Println("Reported ", m.cm.value, " to ", address)
			_, err = client.Do(request)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		client := &http.Client{}
		address := "http://127.0.0.1:8080/" + "update/gauge/PollCount/" + strconv.FormatInt(int64(m.cm.value), 10)
		request, err := http.NewRequest(http.MethodPost, address, nil)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		request.Header.Add("Content-Type", "text/plain")
		//fmt.Println("Reported ", m.cm.value, " to ", address)
		_, err = client.Do(request)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		ch <- m
	}
}

func main() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT)
	gms := map[string]gauge{"Alloc": 0.0, "BuckHashSys": 0.0, "RandomValue": 0.0}
	metricsVal := metrics{gms, struct {
		name  string
		value counter
	}{"PollCount", 0}}
	c := make(chan metrics, 1)
	go updateMetrics(c)
	go sendMetrics(c)
	c <- metricsVal
	<-sigc
}

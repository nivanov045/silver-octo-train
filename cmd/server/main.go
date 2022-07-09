package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	met "github.com/nivanov045/silver-octo-train/internal/metrics"
)

var M met.Metrics

func requestMetricsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Only \"text/plain\" Content-Type is allowed!", http.StatusBadRequest)
		return
	}
	w.Header().Set("content-type", "application/json")
	s := r.URL.Path
	s = strings.Trim(s, "/update/")
	ss := strings.Split(s, "/")
	if len(ss) != 3 {
		http.Error(w, "Wrong metrics type!", http.StatusBadRequest)
		return
	}
	if ss[0] == "gauge" {
		val, err := strconv.ParseFloat(ss[2], 64)
		if err != nil {
			http.Error(w, "Can't parse gauge value!", http.StatusBadRequest)
			return
		}
		M.Gms[ss[1]] = met.Gauge(val)
	} else if ss[0] == "counter" {
		val, err := strconv.ParseInt(ss[2], 10, 64)
		if err != nil {
			http.Error(w, "Can't parse counter value!", http.StatusBadRequest)
			return
		}
		M.Cms[ss[1]] += met.Counter(val)
		fmt.Println("counter = ", M.Cms["PollCount"])
	} else {
		http.Error(w, "Wrong metrics type!", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func main() {
	M.Cms = make(map[string]met.Counter)
	M.Gms = make(map[string]met.Gauge)
	http.HandleFunc("/update/", requestMetricsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

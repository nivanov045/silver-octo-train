package api

import (
	"fmt"
	"net/http"
	"strings"
)

type Service interface {
	ParseAndSet(data string) error
}

type api struct {
	service Service
}

func New(service Service) *api {
	return &api{service: service}
}

func (a *api) requestMetricsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("content-type", "application/json")
	s := r.URL.Path
	s = strings.Trim(s, "/update")
	if err := a.service.ParseAndSet(s); err == nil {
		w.WriteHeader(http.StatusOK)
	} else if err.Error() == "wrong metrics type" {
		w.WriteHeader(http.StatusNotImplemented)
	} else {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
	}
}

func (a *api) Run() error {
	http.HandleFunc("/update/", a.requestMetricsHandler)
	return http.ListenAndServe(":8080", nil)
}

type API interface {
	Run() error
}

var _ API = &api{}

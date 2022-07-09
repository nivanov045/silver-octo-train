package api

import (
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
	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Only \"text/plain\" Content-Type is allowed!", http.StatusBadRequest)
		return
	}
	w.Header().Set("content-type", "application/json")
	s := r.URL.Path
	s = strings.Trim(s, "/update/")
	if a.service.ParseAndSet(s) == nil {
		w.WriteHeader(http.StatusOK)
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

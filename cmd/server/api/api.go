package api

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"log"
	"net/http"
)

type Service interface {
	ParseAndSave([]byte) error
	ParseAndGet([]byte) ([]byte, error)
	GetKnownMetrics() []string
}

type api struct {
	service Service
}

func New(service Service) *api {
	return &api{service: service}
}

func (a *api) updateMetricsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("api::updateMetricsHandler: started ", r)
	w.Header().Set("content-type", "application/json")
	w.Write([]byte("{}"))
	defer r.Body.Close()
	respBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("api::updateMetricsHandler: can't read response body with", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err := a.service.ParseAndSave(respBody); err == nil {
		log.Println("api::updateMetricsHandler: parsed and saved")
		w.WriteHeader(http.StatusOK)
	} else {
		log.Println("api::updateMetricsHandler: error in parsing:", err)
		if err.Error() == "wrong metrics type" {
			w.WriteHeader(http.StatusNotImplemented)
		} else if err.Error() == "can't parse value" {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
	log.Println("api::updateMetricsHandler: response:", w)
}

func (a *api) getMetricsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("api::getMetricsHandler: started", r)
	w.Header().Set("content-type", "application/json")
	defer r.Body.Close()
	respBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("api::getMetricsHandler: can't read response body with", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if val, err := a.service.ParseAndGet(respBody); err == nil {
		log.Println("api::getMetricsHandler: parsed and get")
		w.WriteHeader(http.StatusOK)
		w.Write(val)
	} else {
		log.Println("api::getMetricsHandler: error in parsing:", err)
		if err.Error() == "wrong metrics type" {
			w.WriteHeader(http.StatusNotImplemented)
		} else if err.Error() == "no such metric" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func (a *api) rootHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("api::rootHandler: started")
	w.Header().Set("content-type", "text/html")
	for _, val := range a.service.GetKnownMetrics() {
		w.Write([]byte(val + "\n"))
	}
}

func (a *api) Run() error {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/update/", a.updateMetricsHandler)
	r.Get("/", a.rootHandler)
	r.Post("/value/", a.getMetricsHandler)
	return http.ListenAndServe(":8080", r)
}

type API interface {
	Run() error
}

var _ API = &api{}

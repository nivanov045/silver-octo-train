package api

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"net/http"
)

type Service interface {
	ParseAndSave(data string) error
	ParseAndGet(data string) (string, error)
	GetKnownMetrics() []string
}

type api struct {
	service Service
}

func New(service Service) *api {
	return &api{service: service}
}

func (a *api) updateMetricsHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("updateMetricsHandler")
	w.Header().Set("content-type", "application/json")
	respBody, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	if err := a.service.ParseAndSave(string(respBody)); err == nil {
		w.WriteHeader(http.StatusOK)
	} else if err.Error() == "wrong metrics type" {
		w.WriteHeader(http.StatusNotImplemented)
	} else if err.Error() == "can't parse value" {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusNotFound)
		//fmt.Println(err.Error())
	}
}

func (a *api) getMetricsHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("getMetricsHandler")
	w.Header().Set("content-type", "application/json")
	respBody, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	if val, err := a.service.ParseAndGet(string(respBody)); err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(val))
	} else if err.Error() == "wrong metrics type" {
		w.WriteHeader(http.StatusNotImplemented)
	} else if err.Error() == "no such metric" {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (a *api) rootHandler(w http.ResponseWriter, r *http.Request) {
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

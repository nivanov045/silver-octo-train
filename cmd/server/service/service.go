package service

import (
	"errors"
	"strconv"
	"strings"

	met "github.com/nivanov045/silver-octo-train/internal/metrics"
)

type Storage interface {
	SetGaugeMetrics(name string, val met.Gauge)
	SetCounterMetrics(name string, val met.Counter)
	GetCounterMetrics(name string) met.Counter
}

type service struct {
	storage Storage
}

// ParseAndSet implements api.Service
func (ser *service) ParseAndSet(s string) error {
	ss := strings.Split(s, "/")
	if len(ss) != 3 {
		return errors.New("wrong query")
	}
	if ss[0] == "gauge" {
		val, err := strconv.ParseFloat(ss[2], 64)
		if err != nil {
			return errors.New("can't parse gauge value")
		}
		ser.storage.SetGaugeMetrics(ss[1], met.Gauge(val))
	} else if ss[0] == "counter" {
		val, err := strconv.ParseInt(ss[2], 10, 64)
		if err != nil {
			return errors.New("can't parse counter value")
		}
		ser.storage.SetCounterMetrics(ss[1], met.Counter(int64(ser.storage.GetCounterMetrics(ss[1]))+val))
	} else {
		return errors.New("wrong metrics type")
	}
	return nil
}

func New(storage Storage) *service {
	return &service{storage: storage}
}

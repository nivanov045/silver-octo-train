package service

import (
	"errors"
	"strconv"
	"strings"

	met "github.com/nivanov045/silver-octo-train/internal/metrics"
)

type Storage interface {
	SetGaugeMetrics(name string, val met.Gauge)
	GetGaugeMetrics(name string) (met.Gauge, bool)
	SetCounterMetrics(name string, val met.Counter)
	GetCounterMetrics(name string) (met.Counter, bool)
	GetKnownMetrics() []string
}

type service struct {
	storage Storage
}

func (ser *service) ParseAndSet(s string) error {
	ss := strings.Split(s, "/")
	if len(ss) != 3 {
		return errors.New("wrong query")
	}
	if ss[0] == "gauge" {
		val, err := strconv.ParseFloat(ss[2], 64)
		if err != nil {
			return errors.New("can't parse value")
		}
		ser.storage.SetGaugeMetrics(ss[1], met.Gauge(val))
	} else if ss[0] == "counter" {
		val, err := strconv.ParseInt(ss[2], 10, 64)
		if err != nil {
			return errors.New("can't parse value")
		}
		ex_val, ok := ser.storage.GetCounterMetrics(ss[1])
		if !ok {
			ser.storage.SetCounterMetrics(ss[1], met.Counter(val))
		} else {
			ser.storage.SetCounterMetrics(ss[1], met.Counter(int64(ex_val)+val))
		}
	} else {
		return errors.New("wrong metrics type")
	}
	return nil
}

func (ser *service) ParseAndGet(s string) (string, error) {
	ss := strings.Split(s, "/")
	if len(ss) != 2 {
		return "", errors.New("wrong query")
	}
	if ss[0] == "gauge" {
		val, ok := ser.storage.GetGaugeMetrics(ss[1])
		if !ok {
			return "", errors.New("no such metric")
		}
		return strconv.FormatFloat(float64(val), 'f', 6, 64), nil
	} else if ss[0] == "counter" {
		val, ok := ser.storage.GetCounterMetrics(ss[1])
		if !ok {
			return "", errors.New("no such metric")
		}
		asint := int64(val)
		return strconv.Itoa(int(asint)), nil
	}
	return "", errors.New("wrong metrics type")
}

func (ser *service) GetKnownMetrics() []string {
	return ser.storage.GetKnownMetrics()
}

func New(storage Storage) *service {
	return &service{storage: storage}
}

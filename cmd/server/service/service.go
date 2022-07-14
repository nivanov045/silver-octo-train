package service

import (
	"errors"
	"strconv"
	"strings"

	"github.com/nivanov045/silver-octo-train/internal/metrics"
)

type Storage interface {
	SetGaugeMetrics(name string, val metrics.Gauge)
	GetGaugeMetrics(name string) (metrics.Gauge, bool)
	SetCounterMetrics(name string, val metrics.Counter)
	GetCounterMetrics(name string) (metrics.Counter, bool)
	GetKnownMetrics() []string
}

type service struct {
	storage Storage
}

const (
	gauge   string = "gauge"
	counter string = "counter"
)

func (ser *service) ParseAndSave(s string) error {
	ss := strings.Split(s, "/")
	if len(ss) != 3 {
		return errors.New("wrong query")
	}
	metricType := ss[0]
	metricName := ss[1]
	metricValue := ss[2]
	if metricType == gauge {
		val, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return errors.New("can't parse value")
		}
		ser.storage.SetGaugeMetrics(metricName, metrics.Gauge(val))
	} else if metricType == counter {
		val, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			return errors.New("can't parse value")
		}
		exVal, ok := ser.storage.GetCounterMetrics(metricName)
		if !ok {
			ser.storage.SetCounterMetrics(metricName, metrics.Counter(val))
		} else {
			ser.storage.SetCounterMetrics(metricName, metrics.Counter(int64(exVal)+val))
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
	metricType := ss[0]
	metricName := ss[1]
	if metricType == gauge {
		val, ok := ser.storage.GetGaugeMetrics(metricName)
		if !ok {
			return "", errors.New("no such metric")
		}
		return strconv.FormatFloat(float64(val), 'f', -1, 64), nil
	} else if metricType == counter {
		val, ok := ser.storage.GetCounterMetrics(metricName)
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

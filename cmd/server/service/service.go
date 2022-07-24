package service

import (
	"encoding/json"
	"errors"
	"fmt"

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
	var m metrics.MetricsInterface
	err := json.Unmarshal([]byte(s), &m)
	if err != nil {
		fmt.Println(err.Error())
		return errors.New("wrong query")
	}
	metricType := m.MType
	metricName := m.ID
	if metricType == gauge {
		value := m.Value
		if value == nil {
			return errors.New("wrong query")
		}
		ser.storage.SetGaugeMetrics(metricName, metrics.Gauge(*value))
	} else if metricType == counter {
		exVal, ok := ser.storage.GetCounterMetrics(metricName)
		if m.Delta == nil {
			return errors.New("wrong query")
		}
		value := *m.Delta
		if !ok {
			ser.storage.SetCounterMetrics(metricName, metrics.Counter(value))
		} else {
			ser.storage.SetCounterMetrics(metricName, metrics.Counter(int64(exVal)+value))
		}
	} else {
		return errors.New("wrong metrics type")
	}
	return nil
}

func (ser *service) ParseAndGet(s string) (string, error) {
	var m metrics.MetricsInterface
	err := json.Unmarshal([]byte(s), &m)
	if err != nil {
		return "", errors.New("wrong query")
	}
	metricType := m.MType
	fmt.Println("metricType: ", metricType)
	metricName := m.ID
	fmt.Println("metricaName: ", metricName)
	if metricType == gauge {
		val, ok := ser.storage.GetGaugeMetrics(metricName)
		if !ok {
			return "", errors.New("no such metric")
		}
		asFloat := float64(val)
		m.Value = &asFloat
		marshal, err := json.Marshal(m)
		if err != nil {
			return "", err
		}
		return string(marshal), nil
	} else if metricType == counter {
		val, ok := ser.storage.GetCounterMetrics(metricName)
		if !ok {
			return "", errors.New("no such metric")
		}
		asint := int64(val)
		m.Delta = &asint
		marshal, err := json.Marshal(m)
		if err != nil {
			return "", err
		}
		return string(marshal), nil
	}
	return "", errors.New("wrong metrics type")
}

func (ser *service) GetKnownMetrics() []string {
	return ser.storage.GetKnownMetrics()
}

func New(storage Storage) *service {
	return &service{storage: storage}
}

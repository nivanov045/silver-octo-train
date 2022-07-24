package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nivanov045/silver-octo-train/internal/metrics"
	"sync"
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
	mu      sync.Mutex
}

const (
	gauge   string = "gauge"
	counter string = "counter"
)

func (ser *service) ParseAndSave(s []byte) error {
	ser.mu.Lock() // Блокирует мьютекс
	defer ser.mu.Unlock()
	fmt.Println("ParseAndSave")
	var m metrics.MetricsInterface
	err := json.Unmarshal(s, &m)
	if err != nil {
		fmt.Println("wrong query")
		return errors.New("wrong query")
	}
	metricType := m.MType
	metricName := m.ID
	fmt.Println("ParseAndSave: ", metricName)
	if metricType == gauge {
		fmt.Println("gauge")
		value := m.Value
		if value == nil {
			fmt.Println("wrong query")
			return errors.New("wrong query")
		}
		ser.storage.SetGaugeMetrics(metricName, metrics.Gauge(*value))
	} else if metricType == counter {
		fmt.Println("counter")
		exVal, ok := ser.storage.GetCounterMetrics(metricName)
		if m.Delta == nil {
			fmt.Println("wrong query")
			return errors.New("wrong query")
		}
		value := *m.Delta
		if !ok {
			fmt.Println("!ok")
			ser.storage.SetCounterMetrics(metricName, metrics.Counter(value))
		} else {
			fmt.Println("else")
			ser.storage.SetCounterMetrics(metricName, metrics.Counter(int64(exVal)+value))
		}
	} else {
		fmt.Println("wrong metrics type")
		return errors.New("wrong metrics type")
	}
	return nil
}

func (ser *service) ParseAndGet(s []byte) (string, error) {
	fmt.Println("ParseAndGet")
	var m metrics.MetricsInterface
	err := json.Unmarshal(s, &m)
	if err != nil {
		fmt.Println("wrong query")
		return "", errors.New("wrong query")
	}
	metricType := m.MType
	metricName := m.ID
	fmt.Println("ParseAndGet: ", metricName)
	if metricType == gauge {
		fmt.Println("gauge")
		val, ok := ser.storage.GetGaugeMetrics(metricName)
		if !ok {
			fmt.Println("no such metric")
			return "", errors.New("no such metric")
		}
		asFloat := float64(val)
		m.Value = &asFloat
		marshal, err := json.Marshal(m)
		if err != nil {
			fmt.Println("err != nil")
			return "", err
		}
		return string(marshal), nil
	} else if metricType == counter {
		fmt.Println("counter")
		val, ok := ser.storage.GetCounterMetrics(metricName)
		if !ok {
			fmt.Println("no such metric")
			return "", errors.New("no such metric")
		}
		asint := int64(val)
		m.Delta = &asint
		marshal, err := json.Marshal(m)
		if err != nil {
			fmt.Println("err != nil")
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

package storage

import (
	"github.com/nivanov045/silver-octo-train/internal/metrics"
)

type storage struct {
	Metrics metrics.Metrics
}

func (s *storage) SetCounterMetrics(name string, val metrics.Counter) {
	s.Metrics.CounterMetrics[name] = val
}

func (s *storage) GetCounterMetrics(name string) (metrics.Counter, bool) {
	if val, ok := s.Metrics.CounterMetrics[name]; ok {
		return val, true
	}
	return 0, false
}

func (s *storage) SetGaugeMetrics(name string, val metrics.Gauge) {
	s.Metrics.GaugeMetrics[name] = val
}

func (s *storage) GetGaugeMetrics(name string) (metrics.Gauge, bool) {
	if val, ok := s.Metrics.GaugeMetrics[name]; ok {
		return val, true
	}
	return 0, false
}

func (s *storage) GetKnownMetrics() []string {
	var res []string
	for key := range s.Metrics.CounterMetrics {
		res = append(res, key)
	}
	for key := range s.Metrics.GaugeMetrics {
		res = append(res, key)
	}
	return res
}

func New() *storage {
	return &storage{
		Metrics: metrics.Metrics{
			GaugeMetrics:   map[string]metrics.Gauge{},
			CounterMetrics: map[string]metrics.Counter{},
		},
	}
}

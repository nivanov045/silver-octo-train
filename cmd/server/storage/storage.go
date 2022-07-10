package storage

import (
	met "github.com/nivanov045/silver-octo-train/internal/metrics"
)

type storage struct {
	M met.Metrics
}

func (s *storage) SetCounterMetrics(name string, val met.Counter) {
	s.M.Cms[name] = val
}

func (s *storage) GetCounterMetrics(name string) (met.Counter, bool) {
	if val, ok := s.M.Cms[name]; ok {
		return val, true
	}
	return 0, false
}

func (s *storage) SetGaugeMetrics(name string, val met.Gauge) {
	s.M.Gms[name] = val
}

func (s *storage) GetGaugeMetrics(name string) (met.Gauge, bool) {
	if val, ok := s.M.Gms[name]; ok {
		return val, true
	}
	return 0, false
}

func (s *storage) GetKnownMetrics() []string {
	var res []string
	for key := range s.M.Cms {
		res = append(res, key)
	}
	for key := range s.M.Gms {
		res = append(res, key)
	}
	return res
}

func New() *storage {
	return &storage{
		M: met.Metrics{
			Gms: map[string]met.Gauge{},
			Cms: map[string]met.Counter{},
		},
	}
}

package storage

import (
	"fmt"

	met "github.com/nivanov045/silver-octo-train/internal/metrics"
)

type storage struct {
	M met.Metrics
}

func (s *storage) SetCounterMetrics(name string, val met.Counter) {
	s.M.Cms[name] = val
	fmt.Println("counter = ", s.M.Cms["PollCount"])
}

func (s *storage) GetCounterMetrics(name string) met.Counter {
	return s.M.Cms[name]
}

func (s *storage) SetGaugeMetrics(name string, val met.Gauge) {
	s.M.Gms[name] = val
}

func New() *storage {
	// M.Cms = make(map[string]met.Counter)
	// M.Gms = make(map[string]met.Gauge)
	return &storage{
		M: met.Metrics{
			Gms: map[string]met.Gauge{},
			Cms: map[string]met.Counter{},
		},
	}
}

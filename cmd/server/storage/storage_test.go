package storage

import (
	"testing"

	met "github.com/nivanov045/silver-octo-train/internal/metrics"
)

func Test_storage_SetGetCounterMetrics(t *testing.T) {
	type args struct {
		name string
		val  met.Counter
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "set and get",
			args: args{
				name: "someMetric",
				val:  123,
			},
		},
	}
	s := storage{
		M: met.Metrics{
			Gms: map[string]met.Gauge{},
			Cms: map[string]met.Counter{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.SetCounterMetrics(tt.args.name, tt.args.val)
			val, ok := s.GetCounterMetrics(tt.args.name)
			if !ok || tt.args.val != val {
				t.Errorf("storage.SetCounterMetrics() error with name %v, val %v", tt.args.name, tt.args.val)
			}
		})
	}
}

func Test_storage_SetGaugeMetrics(t *testing.T) {
	type args struct {
		name string
		val  met.Gauge
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "set and check",
			args: args{
				name: "someMetric",
				val:  123,
			},
		},
	}
	s := storage{
		M: met.Metrics{
			Gms: map[string]met.Gauge{},
			Cms: map[string]met.Counter{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.SetGaugeMetrics(tt.args.name, tt.args.val)
		})
	}
}

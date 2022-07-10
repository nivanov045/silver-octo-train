package storage

import (
	"testing"

	met "github.com/nivanov045/silver-octo-train/internal/metrics"
	"github.com/stretchr/testify/assert"
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

func Test_storage_SetGetGaugeMetrics(t *testing.T) {
	type args struct {
		name string
		val  met.Gauge
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "set and get",
			args: args{
				name: "someMetric",
				val:  123.456,
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
			val, ok := s.GetGaugeMetrics(tt.args.name)
			if !ok || tt.args.val != val {
				t.Errorf("storage.SetCounterMetrics() error with name %v, val %v", tt.args.name, tt.args.val)
			}
		})
	}
}

func Test_storage_GetKnownMetrics(t *testing.T) {
	tests := []struct {
		name string
		s    *storage
		want []string
	}{
		{
			name: "no metrics",
			s: &storage{
				M: met.Metrics{
					Gms: map[string]met.Gauge{},
					Cms: map[string]met.Counter{},
				},
			},
			want: []string(nil),
		},
		{
			name: "one metric",
			s: &storage{
				M: met.Metrics{
					Gms: map[string]met.Gauge{"TestMetric": 0.0},
					Cms: map[string]met.Counter{},
				},
			},
			want: []string{"TestMetric"},
		},
		{
			name: "several metrics",
			s: &storage{
				M: met.Metrics{
					Gms: map[string]met.Gauge{"TestMetricG": 0.5},
					Cms: map[string]met.Counter{"TestMetricC": 10},
				},
			},
			want: []string{"TestMetricC", "TestMetricG"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.s.GetKnownMetrics(), tt.want)
		})
	}
}

package metricsperformer

import (
	"testing"

	met "github.com/nivanov045/silver-octo-train/internal/metrics"
	"github.com/stretchr/testify/assert"
)

func Test_metricsPerformer_UpdateMetrics(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "update",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := met.Metrics{
				Gms: map[string]met.Gauge{},
				Cms: map[string]met.Counter{},
			}
			mp := metricsPerformer{}
			assert.Equal(t, len(m.Gms), 0)
			assert.Equal(t, len(m.Cms), 0)
			mp.UpdateMetrics(m)
			assert.Equal(t, len(m.Gms), 27)
			assert.Equal(t, len(m.Cms), 1)
			assert.Equal(t, m.Cms["PollCount"], met.Counter(1))
			mp.UpdateMetrics(m)
			assert.Equal(t, m.Cms["PollCount"], met.Counter(2))
			mp.UpdateMetrics(m)
			assert.Equal(t, m.Cms["PollCount"], met.Counter(3))
		})
	}
}

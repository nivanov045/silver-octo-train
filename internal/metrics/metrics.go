package metrics

type Gauge float64
type Counter int64

type Metrics struct {
	Gms map[string]Gauge
	Cms map[string]Counter
}

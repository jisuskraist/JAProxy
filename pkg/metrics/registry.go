package metrics

import "github.com/prometheus/client_golang/prometheus"

// Registry represents a metric registry.
type Registry struct {
	Histogram *prometheus.HistogramVec
}

func NewRegistry() *Registry {
	return &Registry{
		Histogram: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "proxy_duration_seconds",
				Help: "Time that took to send and return the request",
			},
			[]string{"code"}),
	}
}

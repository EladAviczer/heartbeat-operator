package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ProbeSuccess = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "probe_success",
		Help: "Current status of the probe (1 for success, 0 for failure)",
	}, []string{"name", "target", "type"})

	ProbeDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "probe_duration_seconds",
		Help:    "Duration of the probe execution in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"name", "target", "type"})

	ProbeLastTimestamp = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "probe_last_timestamp_seconds",
		Help: "Timestamp of the last probe execution",
	}, []string{"name", "target", "type"})
)

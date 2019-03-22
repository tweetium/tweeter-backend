package testutil

import (
	"github.com/prometheus/client_golang/prometheus"
	prometheusTestutil "github.com/prometheus/client_golang/prometheus/testutil"
)

// PrometheusCollectorDiff returns the difference between the collector value
// after the func() provided is ran.
func PrometheusCollectorDiff(c prometheus.Collector, f func()) float64 {
	before := prometheusTestutil.ToFloat64(c)
	f()
	after := prometheusTestutil.ToFloat64(c)
	return after - before
}

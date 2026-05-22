// Package metrics defines a Prometheus-style metrics abstraction for
// exposing gauges, counters, and observers (histograms/summaries).
package metrics

// MetricName is a unique name identifying a metric.
type MetricName string

// Labels is a set of label key-value pairs attached to a metric instance.
type Labels map[string]string

// Gauge is a metric that represents a single numerical value that can
// go up and down.
type Gauge interface {
	Inc()
	Dec()
	Add(float64)
	Set(float64)
}

// Counter is a metric that represents a monotonically increasing value.
type Counter interface {
	Inc()
	Add(float64)
}

// Observer is a metric that observes values for distribution statistics
// (e.g. histograms or summaries).
type Observer interface {
	Observe(float64)
}

// Metrics is the top-level metrics factory. It creates or retrieves metric
// instances by name and labels, following the Prometheus pattern where
// label sets produce unique metric instances.
type Metrics interface {
	// Counter returns a Counter with the given name and labels.
	Counter(name MetricName, labels Labels) Counter
	// Gauge returns a Gauge with the given name and labels.
	Gauge(name MetricName, labels Labels) Gauge
	// Observer returns an Observer with the given name and labels.
	Observer(name MetricName, labels Labels) Observer
}

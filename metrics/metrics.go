package metrics

type MetricName string

type Labels map[string]string

type Gauge interface {
	Inc()
	Dec()
	Add(float64)
	Set(float64)
}

type Counter interface {
	Inc()
	Add(float64)
}

type Observer interface {
	Observe(float64)
}

type Metrics interface {
	Counter(name MetricName, labels Labels) Counter
	Gauge(name MetricName, labels Labels) Gauge
	Observer(name MetricName, labels Labels) Observer
}

package metrics

var (
	nopGauge    = &noopGauge{}
	nopCounter  = &noopCounter{}
	nopObserver = &noopObserver{}

	noop Metrics = &noopMetrics{}
)

type noopMetrics struct{}

func Noop() Metrics {
	return noop
}

func (m *noopMetrics) Counter(name MetricName, labels Labels) Counter {
	return nopCounter
}

func (m *noopMetrics) Gauge(name MetricName, labels Labels) Gauge {
	return nopGauge
}

func (m *noopMetrics) Observer(name MetricName, labels Labels) Observer {
	return nopObserver
}

type noopGauge struct{}

func (*noopGauge) Inc()          {}
func (*noopGauge) Dec()          {}
func (*noopGauge) Add(v float64) {}
func (*noopGauge) Set(v float64) {}

type noopCounter struct{}

func (*noopCounter) Inc()          {}
func (*noopCounter) Add(v float64) {}

type noopObserver struct{}

func (*noopObserver) Observe(v float64) {}

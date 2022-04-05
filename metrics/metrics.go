package metrics

type MetricName string

const (
	MetricServicesGauge                     MetricName = "gost_services"
	MetricServiceRequestsCounter            MetricName = "gost_service_requests_total"
	MetricServiceRequestsInFlightGauge      MetricName = "gost_service_requests_in_flight"
	MetricServiceRequestsDurationObserver   MetricName = "gost_service_request_duration_seconds"
	MetricServiceTransferInputBytesCounter  MetricName = "gost_service_transfer_input_bytes_total"
	MetricServiceTransferOutputBytesCounter MetricName = "gost_service_transfer_output_bytes_total"
	MetricNodeConnectDurationObserver       MetricName = "gost_chain_node_connect_duration_seconds"
	MetricServiceHandlerErrorsCounter       MetricName = "gost_service_handler_errors_total"
	MetricChainErrorsCounter                MetricName = "gost_chain_errors_total"
)

type Labels map[string]string

var (
	global Metrics = Noop()
)

func SetGlobal(m Metrics) {
	if m != nil {
		global = m
	} else {
		global = Noop()
	}
}

func Global() Metrics {
	return global
}

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

func GetCounter(name MetricName, labels Labels) Counter {
	return global.Counter(name, labels)
}

func GetGauge(name MetricName, labels Labels) Gauge {
	return global.Gauge(name, labels)
}

func GetObserver(name MetricName, labels Labels) Observer {
	return global.Observer(name, labels)
}

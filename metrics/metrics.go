package metrics

import (
	"github.com/go-gost/core/common/util"
)

type MetricName string

const (
	// Number of services. Labels: host.
	MetricServicesGauge MetricName = "gost_services"
	// Total service requests. Labels: host, service.
	MetricServiceRequestsCounter MetricName = "gost_service_requests_total"
	// Number of in-flight requests. Labels: host, service.
	MetricServiceRequestsInFlightGauge MetricName = "gost_service_requests_in_flight"
	// Request duration historgram. Labels: host, service.
	MetricServiceRequestsDurationObserver MetricName = "gost_service_request_duration_seconds"
	// Total service input data transfer size in bytes. Labels: host, service.
	MetricServiceTransferInputBytesCounter MetricName = "gost_service_transfer_input_bytes_total"
	// Total service output data transfer size in bytes. Labels: host, service.
	MetricServiceTransferOutputBytesCounter MetricName = "gost_service_transfer_output_bytes_total"
	// Chain node connect duration histogram. Labels: host, chain, node.
	MetricNodeConnectDurationObserver MetricName = "gost_chain_node_connect_duration_seconds"
	// Total service handler errors. Labels: host, service.
	MetricServiceHandlerErrorsCounter MetricName = "gost_service_handler_errors_total"
	// Total chain connect errors. Labels: host, chain, node.
	MetricChainErrorsCounter MetricName = "gost_chain_errors_total"
)

type Labels map[string]string

var (
	global Metrics = Noop()
	globals        = make(map[string]Metrics)
)

func SetGlobal(m Metrics) {
	var metrics Metrics

	if m != nil {
		metrics = m
	} else {
		metrics = Noop()
	}
	global = metrics
	globals[util.GetGoroutineID()] = metrics
}

func Global() Metrics {
	goglobal, exists := globals[util.GetGoroutineID()]

	if exists {
		return goglobal
	}
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
	return Global().Counter(name, labels)
}

func GetGauge(name MetricName, labels Labels) Gauge {
	return Global().Gauge(name, labels)
}

func GetObserver(name MetricName, labels Labels) Observer {
	return Global().Observer(name, labels)
}

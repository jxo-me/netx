package metrics

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

var (
	global IMetrics = Noop()
)

func Init(m IMetrics) {
	if m != nil {
		global = m
	} else {
		global = Noop()
	}
}

func IsEnabled() bool {
	return global != Noop()
}

func GetCounter(name MetricName, labels Labels) ICounter {
	return global.Counter(name, labels)
}

func GetGauge(name MetricName, labels Labels) IGauge {
	return global.Gauge(name, labels)
}

func GetObserver(name MetricName, labels Labels) IObserver {
	return global.Observer(name, labels)
}

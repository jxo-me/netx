package metrics

var (
	nopGauge    = &noopGauge{}
	nopCounter  = &noopCounter{}
	nopObserver = &noopObserver{}

	noop IMetrics = &noopMetrics{}
)

type noopMetrics struct{}

func Noop() IMetrics {
	return noop
}

func (m *noopMetrics) Counter(name MetricName, labels Labels) ICounter {
	return nopCounter
}

func (m *noopMetrics) Gauge(name MetricName, labels Labels) IGauge {
	return nopGauge
}

func (m *noopMetrics) Observer(name MetricName, labels Labels) IObserver {
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

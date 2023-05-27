package metrics

type MetricName string

type Labels map[string]string

type Gauge interface {
	Inc()
	Dec()
	Add(float64)
	Set(float64)
}

type ICounter interface {
	Inc()
	Add(float64)
}

type IObserver interface {
	Observe(float64)
}

type IMetrics interface {
	Counter(name MetricName, labels Labels) ICounter
	Gauge(name MetricName, labels Labels) Gauge
	Observer(name MetricName, labels Labels) IObserver
}

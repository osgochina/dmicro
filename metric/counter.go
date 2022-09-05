package metric

import "github.com/prometheus/client_golang/prometheus"

type CounterInterface interface {
	Add(delta int)
	With(labelValues ...string) CounterInterface
}

type Counter struct {
	counter     prometheus.Counter
	labelValues prometheus.Labels
}

func NewCounter(namespace string, name string, help string, labelNames prometheus.Labels) *Counter {
	opt := prometheus.CounterOpts{
		Namespace:   namespace,
		Name:        name,
		Help:        help,
		ConstLabels: labelNames,
	}
	c := &Counter{
		counter: prometheus.NewCounter(opt),
	}
	return c
}

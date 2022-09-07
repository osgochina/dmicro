package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

type CounterVec interface {
	Inc(labels ...string)
	Add(v float64, labels ...string)
	Close() bool
}

type CounterVecOpts VectorOpts

type counterVec struct {
	counter *prometheus.CounterVec
}

func NewCounterVec(c *CounterVecOpts) CounterVec {
	if c == nil {
		return nil
	}
	opt := prometheus.CounterOpts{
		Namespace: c.Namespace,
		Subsystem: c.Subsystem,
		Name:      c.Name,
		Help:      c.Help,
	}
	cv := prometheus.NewCounterVec(opt, c.Labels)
	prometheus.MustRegister(cv)

	return &counterVec{
		counter: cv,
	}
}

func (that *counterVec) Inc(labels ...string) {
	that.counter.WithLabelValues(labels...).Inc()
}

func (that *counterVec) Add(v float64, labels ...string) {
	that.counter.WithLabelValues(labels...).Add(v)
}

func (that *counterVec) Close() bool {
	return prometheus.Unregister(that.counter)
}

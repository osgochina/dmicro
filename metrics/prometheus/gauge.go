package prometheus

import "github.com/prometheus/client_golang/prometheus"

type GaugeVecOpts VectorOpts

type GaugeVec interface {
	Set(v float64, labels ...string)
	Inc(labels ...string)
	Add(v float64, labels ...string)
	Close() bool
}

type gaugeVec struct {
	gauge *prometheus.GaugeVec
}

func NewGaugeVec(opt *GaugeVecOpts) GaugeVec {
	if opt == nil {
		return nil
	}
	vec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: opt.Namespace,
			Subsystem: opt.Subsystem,
			Name:      opt.Name,
			Help:      opt.Help,
		}, opt.Labels)
	prometheus.MustRegister(vec)
	gv := &gaugeVec{
		gauge: vec,
	}

	return gv
}

func (that *gaugeVec) Inc(labels ...string) {
	that.gauge.WithLabelValues(labels...).Inc()
}

func (that *gaugeVec) Add(v float64, labels ...string) {
	that.gauge.WithLabelValues(labels...).Add(v)
}

func (that *gaugeVec) Set(v float64, labels ...string) {
	that.gauge.WithLabelValues(labels...).Set(v)
}

func (that *gaugeVec) Close() bool {
	return prometheus.Unregister(that.gauge)
}

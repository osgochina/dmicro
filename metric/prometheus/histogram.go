package prometheus

import "github.com/prometheus/client_golang/prometheus"

type HistogramVecOpts struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
	Labels    []string
	Buckets   []float64
}

type HistogramVec interface {
	Observe(v int64, labels ...string)
}

type histogramVec struct {
	histogram *prometheus.HistogramVec
}

func NewHistogramVec(opt *HistogramVecOpts) HistogramVec {
	if opt == nil {
		return nil
	}
	vec := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: opt.Namespace,
		Subsystem: opt.Subsystem,
		Name:      opt.Name,
		Help:      opt.Help,
		Buckets:   opt.Buckets,
	}, opt.Labels)
	prometheus.MustRegister(vec)

	hv := &histogramVec{
		histogram: vec,
	}

	return hv
}

func (that *histogramVec) Observe(v int64, labels ...string) {
	that.histogram.WithLabelValues(labels...).Observe(float64(v))
}

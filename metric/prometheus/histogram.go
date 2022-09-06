package prometheus

type HistogramInterface interface {
	Put(sample float64)
	With(labelValues ...string) HistogramInterface
}

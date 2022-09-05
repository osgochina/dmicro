package metric

type HistogramInterface interface {
	Put(sample float64)
	With(labelValues ...string) HistogramInterface
}

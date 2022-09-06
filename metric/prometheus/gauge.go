package prometheus

type GaugeInterface interface {
	Add(delta float64)
	Set(value float64)
	With(labelValues ...string) GaugeInterface
}

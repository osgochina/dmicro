package prometheus

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/container/gtype"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"sync"
)

type VectorOpts struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
	Labels    []string
}

var (
	once    sync.Once
	enabled = gtype.NewBool(false)
)

type PromMetrics struct {
	options metrics.Options
}

var _ metrics.Metrics = new(PromMetrics)

func NewPromMetrics(opts ...metrics.Option) *PromMetrics {
	p := &PromMetrics{
		options: metrics.Options{},
	}
	p.options.Plugins = []drpc.Plugin{NewPrometheusPlugin(p)}
	p.configure(opts...)
	return p
}

func (that *PromMetrics) Init(option ...metrics.Option) {
	that.configure(option...)
}

func (that *PromMetrics) Options() metrics.Options {
	return that.options
}

func (that *PromMetrics) configure(opts ...metrics.Option) {
	for _, o := range opts {
		o(&that.options)
	}
	if len(that.options.Host) == 0 {
		that.options.Host = "0.0.0.0"
	}
	if that.options.Port <= 0 {
		that.options.Port = 9101
	}
	if len(that.options.Path) == 0 {
		that.options.Path = "/metrics"
	}
}

func (that *PromMetrics) Enabled() bool {
	return enabled.Val()
}

func (that *PromMetrics) Start() {
	once.Do(func() {
		enabled.Cas(false, true)
		go func() {
			http.Handle(that.options.Path, promhttp.Handler())
			addr := fmt.Sprintf("%s:%d", that.options.Host, that.options.Port)
			logger.Infof(context.TODO(), "Starting prometheus agent at http://%s%s", addr, that.options.Path)
			if err := http.ListenAndServe(addr, nil); err != nil {
				logger.Error(context.TODO(), err)
			}
		}()
	})
}

// Shutdown prometheus 组件不需要停止
func (that *PromMetrics) Shutdown() {

}

func (that *PromMetrics) String() string {
	return "prometheus_metrics"
}

package prometheus

import (
	"fmt"
	"github.com/gogf/gf/container/gtype"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/metric"
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

type PromMetric struct {
	options metric.Options
}

var _ metric.Metrics = new(PromMetric)

func NewPromMetric(opts ...metric.Option) *PromMetric {
	p := &PromMetric{
		options: metric.Options{},
	}
	p.options.Plugins = []drpc.Plugin{NewPrometheusPlugin(p)}
	p.configure(opts...)
	return p
}

func (that *PromMetric) Init(option ...metric.Option) {
	that.configure(option...)
}

func (that *PromMetric) Options() metric.Options {
	return that.options
}

func (that *PromMetric) configure(opts ...metric.Option) {
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

func (that *PromMetric) Enabled() bool {
	return enabled.Val()
}

func (that *PromMetric) Start() {
	once.Do(func() {
		enabled.Cas(false, true)
		go func() {
			http.Handle(that.options.Path, promhttp.Handler())
			addr := fmt.Sprintf("%s:%d", that.options.Host, that.options.Port)
			logger.Infof("Starting prometheus agent at http://%s%s", addr, that.options.Path)
			if err := http.ListenAndServe(addr, nil); err != nil {
				logger.Error(err)
			}
		}()
	})
}

func (that *PromMetric) String() string {
	return "prometheus_metric"
}

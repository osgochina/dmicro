package prometheus

import (
	"github.com/gogf/gf/util/gconv"
	"github.com/osgochina/dmicro/drpc"
)

var serverNamespace = "rpc_server"

var metricServerReplyCodeTotal = NewCounterVec(&CounterVecOpts{
	Namespace: serverNamespace,
	Subsystem: "reply",
	Name:      "code_total",
	Help:      "rpc server call reply code count.",
	Labels:    []string{"path", "code"},
})

type prometheusPlugin struct {
	metric *PromMetric
}

var (
	_ drpc.BeforeWriteReplyPlugin  = new(prometheusPlugin)
	_ drpc.AfterReadCallBodyPlugin = new(prometheusPlugin)
)

func NewPrometheusPlugin(metric *PromMetric) *prometheusPlugin {
	return &prometheusPlugin{
		metric: metric,
	}
}

func (that *prometheusPlugin) Name() string {
	return "metric_prometheus"
}

func (that *prometheusPlugin) BeforeWriteReply(ctx drpc.WriteCtx) *drpc.Status {
	if !that.metric.Enabled() {
		return nil
	}
	readCtx := ctx.(drpc.ReadCtx)
	path := readCtx.ServiceMethod()
	code := gconv.String(ctx.Status().Code())
	metricServerReplyCodeTotal.Inc(path, code)
	return nil
}

func (that *prometheusPlugin) AfterReadCallBody(ctx drpc.ReadCtx) *drpc.Status {
	return nil
}

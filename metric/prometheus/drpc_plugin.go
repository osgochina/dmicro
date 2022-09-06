package prometheus

import (
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"github.com/osgochina/dmicro/drpc"
	"time"
)

var serverNamespace = "rpc_server"

var metricServerReplyCodeTotal = NewCounterVec(&CounterVecOpts{
	Namespace: serverNamespace,
	Subsystem: "call",
	Name:      "code_total",
	Help:      "rpc server call reply code count.",
	Labels:    []string{"name", "path", "code"},
})

var metricServerReplyDur = NewHistogramVec(&HistogramVecOpts{
	Namespace: serverNamespace,
	Subsystem: "call",
	Name:      "duration_ms",
	Help:      "rpc server call reply duration(ms).",
	Labels:    []string{"name", "path"},
	Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
})

type prometheusPlugin struct {
	metric *PromMetric
}

var _ drpc.BeforeWriteReplyPlugin = new(prometheusPlugin)

func NewPrometheusPlugin(metric *PromMetric) *prometheusPlugin {
	return &prometheusPlugin{
		metric: metric,
	}
}

func (that *prometheusPlugin) Name() string {
	return "metric_prometheus"
}

// BeforeWriteReply 回复消息之前调用，在写入客户端之前
func (that *prometheusPlugin) BeforeWriteReply(ctx drpc.WriteCtx) *drpc.Status {
	if !that.metric.Enabled() {
		return nil
	}
	readCtx := ctx.(drpc.ReadCtx)
	path := readCtx.ServiceMethod()
	code := gconv.String(ctx.Status().Code())
	metricServerReplyCodeTotal.Inc(that.metric.Options().ServiceName, path, code)
	metricServerReplyDur.Observe((gtime.Now().UnixNano()-readCtx.StartTime())/int64(time.Millisecond), that.metric.Options().ServiceName, path)
	return nil
}

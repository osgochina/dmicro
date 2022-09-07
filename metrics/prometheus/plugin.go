package prometheus

import (
	"github.com/gogf/gf/util/gconv"
	"github.com/osgochina/dmicro/drpc"
	"time"
)

var serverNamespace = "rpc_server"

// 服务端回复消息的code统计
var metricsReplyCodeTotal = NewCounterVec(&CounterVecOpts{
	Namespace: serverNamespace,
	Subsystem: "reply",
	Name:      "code_total",
	Help:      "rpc server call reply code count.",
	Labels:    []string{"name", "path", "code"},
})

// 服务端处理请求的耗时统计
var metricsReplyDur = NewHistogramVec(&HistogramVecOpts{
	Namespace: serverNamespace,
	Subsystem: "reply",
	Name:      "duration_ms",
	Help:      "rpc server call reply duration(ms).",
	Labels:    []string{"name", "path"},
	Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
})

// 客户端收到响应code的统计
var metricsCallCodeTotal = NewCounterVec(&CounterVecOpts{
	Namespace: serverNamespace,
	Subsystem: "call",
	Name:      "code_total",
	Help:      "rpc server call  code count.",
	Labels:    []string{"name", "path", "code"},
})

// 客户端请求耗时统计
var metricsCallDur = NewHistogramVec(&HistogramVecOpts{
	Namespace: serverNamespace,
	Subsystem: "call",
	Name:      "duration_ms",
	Help:      "rpc server call  duration(ms).",
	Labels:    []string{"name", "path"},
	Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
})

type prometheusPlugin struct {
	metrics *PromMetrics
}

var _ drpc.BeforeWriteReplyPlugin = new(prometheusPlugin)
var _ drpc.AfterCloseEndpointPlugin = new(prometheusPlugin)
var _ drpc.AfterReadReplyBodyPlugin = new(prometheusPlugin)

// NewPrometheusPlugin 创建插件
func NewPrometheusPlugin(metrics *PromMetrics) *prometheusPlugin {
	return &prometheusPlugin{
		metrics: metrics,
	}
}

// Name 插件名称
func (that *prometheusPlugin) Name() string {
	return "metrics_prometheus"
}

// BeforeWriteReply 回复消息之前调用，在写入客户端之前,作为服务端使用生效
func (that *prometheusPlugin) BeforeWriteReply(ctx drpc.WriteCtx) *drpc.Status {
	if !that.metrics.Enabled() {
		return nil
	}
	readCtx := ctx.(drpc.ReadCtx)
	path := readCtx.ServiceMethod()
	code := gconv.String(ctx.Status().Code())
	metricsReplyCodeTotal.Inc(that.metrics.Options().ServiceName, path, code)
	metricsReplyDur.Observe(int64(readCtx.CostTime()/time.Millisecond), that.metrics.Options().ServiceName, path)
	return nil
}

// AfterReadReplyBody 收到回复的消息以后，作为客户端使用生效
func (that *prometheusPlugin) AfterReadReplyBody(ctx drpc.ReadCtx) *drpc.Status {
	if !that.metrics.Enabled() {
		return nil
	}
	path := ctx.ServiceMethod()
	code := gconv.String(ctx.Status().Code())
	metricsCallCodeTotal.Inc(that.metrics.Options().ServiceName, path, code)
	metricsCallDur.Observe(int64(ctx.CostTime()/time.Millisecond), that.metrics.Options().ServiceName, path)
	return nil
}

// AfterCloseEndpoint endpoint关闭，取消metrics的注册
func (that *prometheusPlugin) AfterCloseEndpoint(endpoint drpc.Endpoint, err error) error {
	metricsReplyCodeTotal.Close()
	metricsReplyDur.Close()
	return nil
}

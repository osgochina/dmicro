package main

import (
	"github.com/gogf/gf/util/grand"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/metrics"
	"github.com/osgochina/dmicro/metrics/prometheus"
	"github.com/osgochina/dmicro/server"
	"time"
)

func main() {
	svr := server.NewRpcServer("test_one",
		server.OptEnableHeartbeat(true),
		server.OptListenAddress("127.0.0.1:8199"),
		server.OptMetrics(prometheus.NewPromMetrics(
			metrics.OptHost("0.0.0.0"),
			metrics.OptPort(9101),
			metrics.OptPath("/metrics"),
			metrics.OptServiceName("test_one"),
		)),
	)
	svr.RouteCall(new(Math))
	svr.ListenAndServe()
}

// Math rpc请求的最终处理器，必须集成drpc.CallCtx
type Math struct {
	drpc.CallCtx
}

func (m *Math) Add(arg *[]int) (int, *drpc.Status) {
	// add
	var r int
	for _, a := range *arg {
		r += a
	}
	time.Sleep(time.Duration(grand.Intn(100)) * time.Millisecond)
	// response
	return r, nil
}

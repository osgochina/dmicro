package main

import (
	"github.com/gogf/gf/util/grand"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/metric"
	"github.com/osgochina/dmicro/metric/prometheus"
	"github.com/osgochina/dmicro/server"
	"time"
)

func main() {
	s := server.NewRpcServer("test_one",
		server.OptEnableHeartbeat(true),
		server.OptListenAddress("127.0.0.1:8199"),
		server.OptMetric(prometheus.NewPromMetric(
			metric.OptHost("0.0.0.0"),
		)),
	)
	s.RouteCall(new(Math))
	s.ListenAndServe()
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
	return r, drpc.NewStatus(int32(grand.Intn(10)), "ok")
}

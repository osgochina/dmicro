package main

import (
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/logger"
)

func main() {
	//开启信号监听
	go drpc.GraceSignal()
	// 创建一个rpc服务
	svr := drpc.NewEndpoint(drpc.EndpointConfig{
		CountTime:   true,
		LocalIP:     "127.0.0.1",
		ListenPort:  9091,
		PrintDetail: true,
	})
	//注册处理方法
	svr.RouteCall(new(Math))
	//启动监听
	err := svr.ListenAndServe()
	logger.Warning(err)
}

// Math rpc请求的最终处理器，必须集成drpc.CallCtx
type Math struct {
	drpc.CallCtx
}

func (m *Math) Add(arg *[]int) (int, *drpc.Status) {
	// test meta
	logger.Infof("author: %s", m.PeekMeta("author"))
	// add
	var r int
	for _, a := range *arg {
		r += a
	}
	// response
	return r, nil
}

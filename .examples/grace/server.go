package main

import (
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/logger"
	"time"
)

func main() {
	//开启信号监听
	go drpc.GraceSignal()

	svr := drpc.NewEndpoint(drpc.EndpointConfig{
		CountTime:   true,
		LocalIP:     "127.0.0.1",
		ListenPort:  9091,
		PrintDetail: true,
	})

	svr.RouteCall(new(Grace))
	logger.Warning(svr.ListenAndServe())
	time.Sleep(30 * time.Second)

}

type Grace struct {
	drpc.CallCtx
}

func (m *Grace) Sleep(arg *int) (string, *drpc.Status) {
	logger.Infof("sleep %d", *arg)
	if *arg > 0 {
		sleep := *arg
		time.Sleep(time.Duration(sleep) * time.Second)
	}
	// response
	return "sleep", nil
}

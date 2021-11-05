package main

import (
	"github.com/gogf/gf/util/gconv"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/utils/gracefulv2"
	"github.com/osgochina/dmicro/utils/inherit"
	"time"
)

func main() {
	gracefulv2.GetGraceful().SetModel(gracefulv2.GracefulMasterWorker)
	addr := inherit.NewFakeAddr("tcp", "127.0.0.1", gconv.String(9091))
	err := gracefulv2.GetGraceful().InheritedListener(addr, nil)
	if err != nil {
		logger.Error(err)
		return
	}
	gracefulv2.GetGraceful().MasterWorkerModelStart()
	go gracefulv2.GetGraceful().GraceSignal()
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

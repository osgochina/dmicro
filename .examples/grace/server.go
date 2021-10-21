package main

import (
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/util/gconv"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/plugin/ignorecase"
	"github.com/osgochina/dmicro/logger"
	"os"
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
	}, ignorecase.NewIgnoreCase())

	svr.RouteCall(new(Grace))
	pidFile := gfile.SelfDir() + "/server.pid"
	logger.Info(os.Getpid())
	err := gfile.PutContents(pidFile, gconv.String(os.Getpid()))
	if err != nil {
		logger.Error(err)
	}
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

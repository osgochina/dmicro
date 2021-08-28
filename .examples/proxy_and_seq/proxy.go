package main

import (
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/plugin/proxy"
	"github.com/osgochina/dmicro/logger"
	"time"
)

//go:generate go build $GOFILE

func main() {
	srv := drpc.NewEndpoint(
		drpc.EndpointConfig{
			ListenPort: 8080,
		},
		newProxyPlugin(),
	)
	srv.ListenAndServe()
}

func newProxyPlugin() drpc.Plugin {
	cli := drpc.NewEndpoint(drpc.EndpointConfig{RedialTimes: 3})
	var sess drpc.Session
	var stat *drpc.Status
DIAL:
	sess, stat = cli.Dial(":9090")
	if !stat.OK() {
		logger.Warningf("%v", stat)
		time.Sleep(time.Second * 3)
		goto DIAL
	}
	return proxy.NewPlugin(func(*proxy.Label) proxy.Forwarder {
		return sess
	})
}

package main

import (
	"fmt"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/tfilter"
	"github.com/osgochina/dmicro/logger"
	"time"
)

type GzipCall struct {
	drpc.CallCtx
}

func (that *GzipCall) Zip(str *string) (*string, *drpc.Status) {
	fmt.Println(*str)
	return str, nil
}

func main() {
	tfilter.RegGzip(3)
	go ClientGzip()
	cfg := drpc.EndpointConfig{ListenIP: "127.0.0.1", ListenPort: 8888}
	endpoint := drpc.NewEndpoint(cfg)
	endpoint.RouteCall(new(GzipCall))
	err := endpoint.ListenAndServe()
	if err != nil {
		logger.Fatal(err)
	}
}

func ClientGzip() {
	time.Sleep(3 * time.Second)
	cfg := drpc.EndpointConfig{}
	endpoint := drpc.NewEndpoint(cfg)
	sess, stat := endpoint.Dial("127.0.0.1:8888")
	if !stat.OK() {
		logger.Fatal(stat)
	}
	var str = "liuzhiming"
	var result string
	for true {
		stat = sess.Call("/gzip_call/zip", str, &result, drpc.WithTFilterPipe(tfilter.GzipId)).Status()
		if !stat.OK() {
			logger.Warning(stat)
		}
		time.Sleep(3 * time.Second)
	}
}

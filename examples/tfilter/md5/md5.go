package main

import (
	"context"
	"fmt"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/tfilter"
	"github.com/osgochina/dmicro/logger"
	"time"
)

type Md5Call struct {
	drpc.CallCtx
}

func (that *Md5Call) Md5(str *string) (*string, *drpc.Status) {
	fmt.Println(*str)
	return str, nil
}

func main() {
	tfilter.RegMD5()
	go ClientMd5()
	cfg := drpc.EndpointConfig{ListenIP: "127.0.0.1", ListenPort: 8888}
	endpoint := drpc.NewEndpoint(cfg)
	endpoint.RouteCall(new(Md5Call))
	err := endpoint.ListenAndServe()
	if err != nil {
		logger.Fatal(context.TODO(), err)
	}
}

func ClientMd5() {
	time.Sleep(3 * time.Second)
	cfg := drpc.EndpointConfig{}
	endpoint := drpc.NewEndpoint(cfg)
	sess, stat := endpoint.Dial("127.0.0.1:8888")
	if !stat.OK() {
		logger.Fatal(context.TODO(), stat)
	}
	var str = "liuzhiming"
	var result string
	for true {
		stat = sess.Call("/md5_call/md5", str, &result, drpc.WithTFilterPipe(tfilter.Md5Id)).Status()
		if !stat.OK() {
			logger.Warning(context.TODO(), stat)
		}
		time.Sleep(3 * time.Second)
	}
}

package main

import (
	"context"
	"fmt"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/tfilter"
	"github.com/osgochina/dmicro/logger"
	"time"
)

type MainCall struct {
	drpc.CallCtx
}

func (that *MainCall) Echo(str *string) (*string, *drpc.Status) {
	fmt.Println(*str)
	return str, nil
}

func main() {
	tfilter.RegMD5()
	tfilter.RegGzip(3)
	tfilter.RegAES([]byte("1234567890123456"))
	go Client()
	cfg := drpc.EndpointConfig{ListenIP: "127.0.0.1", ListenPort: 8888}
	endpoint := drpc.NewEndpoint(cfg)
	endpoint.RouteCall(new(MainCall))
	err := endpoint.ListenAndServe()
	if err != nil {
		logger.Fatal(context.TODO(), err)
	}
}

func Client() {
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
		stat = sess.Call(
			"/main_call/echo",
			str,
			&result,
			drpc.WithTFilterPipe(tfilter.AesId, tfilter.Md5Id, tfilter.GzipId),
		).Status()
		if !stat.OK() {
			logger.Warning(context.TODO(), stat)
		}
		time.Sleep(3 * time.Second)
	}
}

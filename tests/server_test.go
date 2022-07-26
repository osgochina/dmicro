package tests

import (
	"github.com/gogf/gf/test/gtest"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/proto/pbproto"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/tests/benchmark"
	"runtime"
	"testing"
)

type MyCall struct {
	drpc.CallCtx
}

func (that *MyCall) Echo(args *benchmark.BenchmarkMessage) (*benchmark.BenchmarkMessage, *drpc.Status) {
	s := "OK"
	var i int32 = 100
	args.Field1 = s
	args.Field2 = i
	runtime.Gosched()
	return args, nil
}

func TestServer(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		logger.SetDebug(false)
		server := drpc.NewEndpoint(drpc.EndpointConfig{
			DefaultBodyCodec: "protobuf",
			ListenPort:       8199,
		})
		server.RouteCall(new(MyCall))
		server.ListenAndServe(pbproto.NewPbProtoFunc())
	})
}

package tests

import (
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/codec"
	"github.com/osgochina/dmicro/drpc/proto/rawproto"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/tests/benchmark"
	"testing"
	"time"
)

func serverRaw() {
	logger.SetDebug(false)
	svr := drpc.NewEndpoint(drpc.EndpointConfig{
		DefaultBodyCodec: codec.JsonName,
		ListenPort:       8199,
	})
	svr.RouteCall(new(MyCall))
	svr.ListenAndServe(rawproto.NewRawProtoFunc())
}

var clientRaw = drpc.NewEndpoint(drpc.EndpointConfig{DefaultBodyCodec: codec.JsonName})

func BenchmarkClientRaw(b *testing.B) {

	once.Do(func() {
		go serverRaw()
	})
	time.Sleep(2 * time.Second)
	b.ResetTimer()
	b.StartTimer()
	logger.SetDebug(false)
	b.ResetTimer()
	serviceMethod := "/my_call/echo"
	args := prepareArgs()
	b.RunParallel(func(pb *testing.PB) {
		sess, err := clientRaw.Dial("127.0.0.1:8199", rawproto.NewRawProtoFunc())
		if !err.OK() {
			b.Fatal(err)
		}
		for pb.Next() {
			var reply benchmark.BenchmarkMessage
			if !sess.Call(serviceMethod, args, &reply).StatusOK() {
				b.Fatal("调用出错")
			}
		}
	})
	b.ReportAllocs()
}

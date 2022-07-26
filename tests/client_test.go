package tests

import (
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/proto/pbproto"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/tests/benchmark"
	"testing"
)

var client = drpc.NewEndpoint(drpc.EndpointConfig{DefaultBodyCodec: "protobuf"})

func BenchmarkClient(b *testing.B) {

	logger.SetDebug(false)
	b.ResetTimer()
	serviceMethod := "/my_call/echo"
	args := prepareArgs()
	b.RunParallel(func(pb *testing.PB) {
		sess, err := client.Dial("127.0.0.1:8199", pbproto.NewPbProtoFunc())
		if !err.OK() {
			b.Fatal(err)
		}
		for pb.Next() {
			var reply benchmark.BenchmarkMessage
			if sess.Call(serviceMethod, args, &reply).StatusOK() {

			}
		}
	})
	b.ReportAllocs()
}

func call(client drpc.Endpoint) {

}

func prepareArgs() *benchmark.BenchmarkMessage {

	var i int32 = 100000
	var s = "如果爱，请深爱，老男人情话大赏!"

	var args benchmark.BenchmarkMessage

	args.Field1 = s
	args.Field2 = i
	args.Field3 = i
	args.Field4 = s
	args.Field5 = []string{s, s, s}
	args.Field6 = i
	args.Field7 = s
	args.Field8 = s
	args.Field9 = true
	args.Field10 = true
	args.Field11 = true
	args.Field12 = true
	args.Field13 = i
	args.Field14 = s
	args.Field15 = i
	args.Field16 = 9999
	args.Field17 = false
	args.Field18 = 11111
	args.Field20 = false
	args.Field21 = false
	args.Field22 = 111111111
	args.Field23 = 22222222
	args.Field24 = 22222222
	args.Field25 = false
	args.Field26 = false
	args.Field27 = false
	args.Field28 = 2223333
	args.Field29 = 45778888
	args.Field30 = s
	args.Field31 = "aaaaaa"
	args.Field32 = 12222
	args.Field33 = i
	args.Field33 = i
	args.Field34 = "fdsfsdfsafadsfsfa"
	args.Field35 = i
	args.Field36 = 23456
	args.Field37 = i
	args.Field38 = i
	args.Field39 = i
	args.Field40 = i

	return &args
}

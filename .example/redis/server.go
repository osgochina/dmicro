package main

import (
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/codec/redis_codec"
	"github.com/osgochina/dmicro/drpc/plugin/redis"
	"github.com/osgochina/dmicro/drpc/proto/redisproto"
)

func main() {
	//开启信号监听
	go drpc.GraceSignal()

	svr := drpc.NewEndpoint(drpc.EndpointConfig{
		CountTime:   true,
		LocalIP:     "127.0.0.1",
		ListenPort:  9091,
		PrintDetail: true,
	}, redis.NewRedisPlugin())

	svr.RouteCall(new(Redis))

	svr.ListenAndServe(redisproto.RedisProtoFunc)
}

type Redis struct {
	drpc.CallCtx
}

func (m *Redis) Get(arg *redis_codec.CmdLine) (redis_codec.Msg, *drpc.Status) {

	// response
	return redis_codec.MakeBulkMsg([]byte("liuizhiming")), nil
}

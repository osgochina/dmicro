package main

import (
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/codec/redis_codec"
	"github.com/osgochina/dmicro/drpc/plugin/redis"
	"github.com/osgochina/dmicro/drpc/proto/redisproto"
	"github.com/osgochina/dmicro/utils/graceful"
)

var db *gmap.AnyAnyMap

func main() {
	//开启信号监听
	go graceful.GraceSignal()
	db = gmap.New(true)
	svr := drpc.NewEndpoint(drpc.EndpointConfig{
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
	if len(*arg) > 0 {
		key := gconv.String((*arg)[0])
		val := db.GetVar(key)
		if val.IsNil() {
			return redis_codec.MakeNullBulkMsg(), nil
		}
		return redis_codec.MakeBulkMsg(val.Bytes()), nil
	}
	// response
	return redis_codec.MakeErrorMsg("参数错误"), nil
}

func (m *Redis) Set(arg *redis_codec.CmdLine) (redis_codec.Msg, *drpc.Status) {
	if len(*arg) == 2 {
		key := gconv.String((*arg)[0])
		val := gconv.String((*arg)[1])
		db.Set(key, val)
		return redis_codec.MakeSuccessMsg("ok"), nil
	}
	// response
	return redis_codec.MakeErrorMsg("参数错误"), nil
}

package socket

import (
	"github.com/osgochina/dmicro/drpc/proto"
	"github.com/osgochina/dmicro/drpc/proto/rawproto"
)

type Proto = proto.Proto
type ProtoFunc = proto.ProtoFunc

//默认传输编码协议
var defaultProtoFunc = rawproto.RawProtoFunc

// DefaultProtoFunc 获取默认的传输编码协议
func DefaultProtoFunc() ProtoFunc {
	return defaultProtoFunc
}

// SetDefaultProtoFunc 设置默认的传输编码协议
func SetDefaultProtoFunc(protoFunc ProtoFunc) {
	defaultProtoFunc = protoFunc
}

//获取指定的传输编码协议
func getProto(protoFuncList []proto.ProtoFunc, rw proto.IOWithReadBuffer) proto.Proto {
	if len(protoFuncList) > 0 && protoFuncList[0] != nil {
		return protoFuncList[0](rw)
	}
	return defaultProtoFunc(rw)
}

// Package pbSubProto 实现PROTOBUF套接字通信协议的。
package pbSubProto

import (
	"github.com/gogf/gf/util/gconv"
	"github.com/osgochina/dmicro/drpc/codec"
	"github.com/osgochina/dmicro/drpc/mixer/websocket/pbSubProto/pb"
	"github.com/osgochina/dmicro/drpc/proto"
	"io/ioutil"
	"sync"
)

// NewPbSubProtoFunc 创建pb子协议
func NewPbSubProtoFunc() proto.ProtoFunc {
	return func(rw proto.IOWithReadBuffer) proto.Proto {
		return &pbSubProto{
			id:   'p',
			name: "protobuf",
			rw:   rw,
		}
	}
}

type pbSubProto struct {
	id   byte
	name string
	rw   proto.IOWithReadBuffer
	rMu  sync.Mutex
}

// Version 返回协议的id和名称
func (psp *pbSubProto) Version() (byte, string) {
	return psp.id, psp.name
}

// Pack 打包并写入数据
func (psp *pbSubProto) Pack(m proto.Message) error {
	// 先编码消息体
	bodyBytes, err := m.MarshalBody()
	if err != nil {
		return err
	}
	// 使用管道过滤器过滤消息体数据
	bodyBytes, err = m.PipeTFilter().OnPack(bodyBytes)
	if err != nil {
		return err
	}
	// 使用pb协议编码数据
	b, err := codec.ProtoMarshal(&pb.Payload{
		Seq:           m.Seq(),
		Mtype:         int32(m.MType()),
		ServiceMethod: m.ServiceMethod(),
		Meta:          gconv.Bytes(m.Meta().String()),
		BodyCodec:     int32(m.BodyCodec()),
		Body:          bodyBytes,
		PipeTFilter:   m.PipeTFilter().IDs(),
	})
	if err != nil {
		return err
	}

	_ = m.SetSize(uint32(len(b)))

	_, err = psp.rw.Write(b)
	return err
}

// Unpack 读取数据并解包
func (psp *pbSubProto) Unpack(m proto.Message) error {
	psp.rMu.Lock()
	defer psp.rMu.Unlock()
	// 读取消息数据
	b, err := ioutil.ReadAll(psp.rw)
	if err != nil {
		return err
	}

	_ = m.SetSize(uint32(len(b)))

	// 使用pb协议解码数据
	s := &pb.Payload{}
	err = codec.ProtoUnmarshal(b, s)
	if err != nil {
		return err
	}

	// 使用消息使用的管道过滤器
	for _, r := range s.PipeTFilter {
		_ = m.PipeTFilter().Append(r)
	}

	m.SetBodyCodec(byte(s.BodyCodec))
	// 读取消息体并使用管道过滤器解包
	bodyBytes, err := m.PipeTFilter().OnUnpack(s.Body)
	if err != nil {
		return err
	}

	// 设置其他的属性
	m.SetSeq(s.Seq)
	m.SetMType(byte(s.Mtype))
	m.SetServiceMethod(s.ServiceMethod)
	_ = m.Meta().UnmarshalJSON(s.Meta)

	// 解包消息体
	err = m.UnmarshalBody(bodyBytes)
	return err
}

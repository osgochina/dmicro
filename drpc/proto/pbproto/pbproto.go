package pbproto

import (
	"encoding/binary"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/osgochina/dmicro/drpc/codec"
	"github.com/osgochina/dmicro/drpc/proto"
	"github.com/osgochina/dmicro/drpc/proto/pbproto/pb"
	"github.com/osgochina/dmicro/utils/dbuffer"
	"io"
	"sync"
)

/**
	Protobuf协议的格式 使用网络字节序，大端
	{4 bytes 表示整个消息的长度}
	{1 byte  表示传输管道过滤器id的长度}
	{传输管道过滤器id序列化后的内容}
	# 以下的内容都是经过传输管道过滤器处理过的数据
	Body: Protobuf bytes
**/

// NewPbProtoFunc 创建protobuf协议方法
func NewPbProtoFunc() proto.ProtoFunc {
	return func(rw proto.IOWithReadBuffer) proto.Proto {
		return &protoPB{
			id:   'p',
			name: "protobuf",
			rw:   rw,
		}
	}
}

type protoPB struct {
	rw   proto.IOWithReadBuffer
	rMu  sync.Mutex
	name string
	id   byte
}

func (that *protoPB) Version() (byte, string) {
	return that.id, that.name
}

// Pack 打包
func (that *protoPB) Pack(m proto.Message) error {
	//消息体打包成byte数组
	bodyBytes, err := m.MarshalBody()
	if err != nil {
		return err
	}
	b, err := codec.ProtoMarshal(&pb.Payload{
		Seq:           m.Seq(),
		Mtype:         int32(m.MType()),
		ServiceMethod: m.ServiceMethod(),
		Status:        gconv.Bytes(m.Status(true).String()),
		Meta:          gconv.Bytes(m.Meta().String()),
		BodyCodec:     int32(m.BodyCodec()),
		Body:          bodyBytes,
	})
	if err != nil {
		return err
	}

	// 通过管道处理打包，并把原始数据替换成处理后的数据
	b, err = m.PipeTFilter().OnPack(b)
	if err != nil {
		return err
	}
	pipeTFilterLen := m.PipeTFilter().Len()
	// 设置消息长度
	err = m.SetSize(uint32(1 + pipeTFilterLen + len(b)))
	if err != nil {
		return err
	}

	// 打包
	var all = make([]byte, m.Size()+4)
	binary.BigEndian.PutUint32(all, m.Size())
	all[4] = byte(pipeTFilterLen)
	copy(all[4+1:], m.PipeTFilter().IDs())
	copy(all[4+1+pipeTFilterLen:], b)
	_, err = that.rw.Write(all)
	return err
}

// Unpack 解包
func (that *protoPB) Unpack(m proto.Message) error {
	that.rMu.Lock()
	defer that.rMu.Unlock()

	var size uint32
	err := binary.Read(that.rw, binary.BigEndian, &size)
	if err != nil {
		return err
	}
	if err = m.SetSize(size); err != nil {
		return err
	}
	if m.Size() == 0 {
		return nil
	}
	bb := dbuffer.GetByteBuffer()
	defer dbuffer.ReleaseByteBuffer(bb)

	bb.ChangeLen(int(m.Size()))
	_, err = io.ReadFull(that.rw, bb.B)
	if err != nil {
		return err
	}

	// 使用管道过滤器处理内容
	var pipeTFilterLen = bb.B[0]
	bb.B = bb.B[1:]
	if pipeTFilterLen > 0 {
		err = m.PipeTFilter().Append(bb.B[:pipeTFilterLen]...)
		if err != nil {
			return err
		}
		bb.B = bb.B[pipeTFilterLen:]
		// do transfer pipe
		bb.B, err = m.PipeTFilter().OnUnpack(bb.B)
		if err != nil {
			return err
		}
	}

	s := &pb.Payload{}
	err = codec.ProtoUnmarshal(bb.B, s)
	if err != nil {
		return err
	}

	// 设置参数
	m.SetSeq(s.Seq)
	m.SetMType(byte(s.Mtype))
	m.SetServiceMethod(s.ServiceMethod)
	err = m.Status(true).UnmarshalJSON(s.Status)
	if err != nil {
		return err
	}
	err = m.Meta().UnmarshalJSON(s.Meta)
	if err != nil {
		return err
	}

	// 解析消息体
	m.SetBodyCodec(byte(s.BodyCodec))
	err = m.UnmarshalBody(s.Body)

	return err
}

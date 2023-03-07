package redisproto

import (
	"bufio"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/osgochina/dmicro/drpc/codec/redis_codec"
	"github.com/osgochina/dmicro/drpc/message"
	"github.com/osgochina/dmicro/drpc/proto"
	"github.com/osgochina/dmicro/utils/dbuffer"
	"io"
	"sync"
)

func NewRedisProtoFunc() proto.ProtoFunc {
	return RedisProtoFunc
}

var RedisProtoFunc = func(rw proto.IOWithReadBuffer) proto.Proto {
	return &redisProto{
		id:   'r',
		name: "redis",
		r:    rw,
		w:    rw,
	}
}

var _ proto.Proto = new(redisProto)

type redisProto struct {
	r    io.Reader
	w    io.Writer
	rMu  sync.Mutex
	name string
	id   byte
}

func (that *redisProto) Version() (byte, string) {
	return that.id, that.name
}

func (that *redisProto) Pack(m proto.Message) error {
	bb := dbuffer.GetByteBuffer()
	defer dbuffer.ReleaseByteBuffer(bb)

	if m.MType() == message.TypeReply {
		if !m.StatusOK() {
			_, _ = bb.Write(redis_codec.MakeErrorMsg(m.Status().Msg()).Bytes())
			goto write
		}
		bodyBytes, err := m.MarshalBody()
		if err != nil {
			return err
		}
		_, _ = bb.Write(bodyBytes)
		goto write
	}

write:
	//写入数据到链接
	_, err := that.w.Write(bb.B)
	if err != nil {
		return err
	}

	return nil
}

func (that *redisProto) Unpack(m proto.Message) error {

	bb := dbuffer.GetByteBuffer()
	defer dbuffer.ReleaseByteBuffer(bb)
	request, err := parseMsgByIO(bufio.NewReader(that.r))
	if err != nil {
		return err
	}
	//m.SetSeq(0)
	//m.SetServiceMethod("")
	//_ = m.Status(true).SetMsg("")
	//m.SetBodyCodec('s')
	//m.UnmarshalBody([]byte(""))
	switch v := request.(type) {
	case *redis_codec.SuccessMsg:
		m.SetMType(message.TypeReply)
	case *redis_codec.ErrorMsg:
		m.SetMType(message.TypeReply)
	case *redis_codec.NumberMsg:
		m.SetMType(message.TypeReply)
	case *redis_codec.BulkMsg:
	case *redis_codec.MultiBulkMsg:
		m.SetMType(message.TypeCall)
		m.SetBodyCodec('r')
		m.SetServiceMethod(string(v.Args[0]))
		s, _ := gjson.Encode(v.Args[1:])
		_ = m.UnmarshalBody(s)
	case *redis_codec.EmptyMultiBulkMsg:
	case *redis_codec.NullBulkMsg:
	default:

	}

	return nil
}

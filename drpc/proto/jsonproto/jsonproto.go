package jsonproto

import (
	"bytes"
	"encoding/binary"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/osgochina/dmicro/drpc/proto"
	"github.com/osgochina/dmicro/utils/dbuffer"
	"io"
	"strconv"
	"sync"
)

/**
	json协议的格式 使用网络字节序，大端
	{4 bytes 表示整个消息的长度}
	{1 byte  表示传输管道过滤器id的长度}
	{传输管道过滤器id序列化后的内容}
	# 以下的内容都是经过传输管道过滤器处理过的数据
	Body: json string
**/

func NewJSONProtoFunc() proto.ProtoFunc {
	return func(rw proto.IOWithReadBuffer) proto.Proto {
		return &jsonProto{
			id:   'j',
			name: "json",
			rw:   rw,
		}
	}
}

type jsonProto struct {
	rw   proto.IOWithReadBuffer
	rMu  sync.Mutex
	name string
	id   byte
}

func (that *jsonProto) Version() (byte, string) {
	return that.id, that.name
}

var (
	msg1 = []byte(`{"seq":`)
	msg2 = []byte(`,"mtype":`)
	msg3 = []byte(`,"serviceMethod":`)
	msg4 = []byte(`,"status":`)
	msg5 = []byte(`,"meta":`)
	msg6 = []byte(`,"bodyCodec":`)
	msg7 = []byte(`,"body":"`)
	msg8 = []byte(`"}`)
)

func (that *jsonProto) Pack(m proto.Message) error {
	bb := dbuffer.GetByteBuffer()
	defer dbuffer.ReleaseByteBuffer(bb)

	bodyBytes, err := m.MarshalBody()
	if err != nil {
		return err
	}
	_, _ = bb.Write(msg1)
	_, _ = bb.WriteString(strconv.FormatInt(int64(m.Seq()), 10))
	_, _ = bb.Write(msg2)
	_, _ = bb.WriteString(strconv.FormatInt(int64(m.MType()), 10))
	_, _ = bb.Write(msg3)
	_, _ = bb.WriteString(strconv.Quote(m.ServiceMethod()))
	_, _ = bb.Write(msg4)
	_, _ = bb.WriteString(strconv.Quote(m.Status(true).String()))
	_, _ = bb.Write(msg5)
	_, _ = bb.WriteString(strconv.Quote(m.Meta().String()))
	_, _ = bb.Write(msg6)
	_, _ = bb.WriteString(strconv.FormatInt(int64(m.BodyCodec()), 10))
	_, _ = bb.Write(msg7)
	_, _ = bb.Write(bytes.Replace(bodyBytes, []byte{'"'}, []byte{'\\', '"'}, -1))
	_, _ = bb.Write(msg8)

	//对消息内容进行处理
	b, err := m.PipeTFilter().OnPack(bb.B)
	if err != nil {
		return err
	}
	pipeLen := m.PipeTFilter().Len()
	//设置消息长度 1 是存管道id的长度，pipeLen大小的内容是存管道id的具体内容，其他就是实际body的长度了
	_ = m.SetSize(uint32(1 + pipeLen + len(b)))
	//新建
	var all = make([]byte, m.Size()+4)        // 4是为了在头部4个字节存储整个消息的长度
	binary.BigEndian.PutUint32(all, m.Size()) //写入消息的长度
	all[4] = byte(pipeLen)
	copy(all[4+1:], m.PipeTFilter().IDs())
	copy(all[4+1+pipeLen:], b)
	_, err = that.rw.Write(all)
	return err
}

func (that *jsonProto) Unpack(m proto.Message) error {
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

	// transfer pipe
	var pipeLen = bb.B[0]
	bb.B = bb.B[1:]
	if pipeLen > 0 {
		err = m.PipeTFilter().Append(bb.B[:pipeLen]...)
		if err != nil {
			return err
		}
		bb.B = bb.B[pipeLen:]
		// do transfer pipe
		bb.B, err = m.PipeTFilter().OnUnpack(bb.B)
		if err != nil {
			return err
		}
	}

	j := gjson.New(string(bb.B))
	// read other
	m.SetSeq(j.Get("seq").Int32())
	m.SetMType(byte(j.Get("mtype").Int8()))
	m.SetServiceMethod(j.Get("serviceMethod").String())
	_ = m.Status(true).UnmarshalJSON(j.Get("status").Bytes())
	_ = m.Meta().UnmarshalJSON(j.Get("meta").Bytes())

	// read body
	m.SetBodyCodec(byte(j.Get("bodyCodec").Int8()))
	body := j.Get("body").Bytes()
	err = m.UnmarshalBody(body)
	return err
}

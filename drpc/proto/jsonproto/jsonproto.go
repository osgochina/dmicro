package jsonproto

import (
	"encoding/binary"
	"encoding/json"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/util/gconv"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/proto"
	"github.com/osgochina/dmicro/utils/dbuffer"
	"io"
	"sync"
)

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

type packMsg struct {
	Seq           int32        `json:"seq"`
	MType         byte         `json:"mtype"`
	ServiceMethod string       `json:"serviceMethod"`
	Status        *drpc.Status `json:"status"`
	Meta          *gmap.Map    `json:"meta"`
	BodyCodec     byte         `json:"bodyCodec"`
	Body          interface{}  `json:"body"`
}

func (that *jsonProto) Pack(m proto.Message) error {
	bb := dbuffer.GetByteBuffer()
	defer dbuffer.ReleaseByteBuffer(bb)
	//组装消息
	pMsg := &packMsg{}
	pMsg.Seq = m.Seq()
	pMsg.MType = m.MType()
	pMsg.ServiceMethod = m.ServiceMethod()
	pMsg.Status = m.Status(true)
	pMsg.Meta = m.Meta()
	pMsg.BodyCodec = m.BodyCodec()
	pMsg.Body = m.Body()

	msgByte, err := json.Marshal(pMsg)
	if err != nil {
		return err
	}
	_, _ = bb.Write(msgByte)

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
	m.SetSeq(j.GetInt32("seq"))
	m.SetMType(byte(j.GetInt8("mtype")))
	m.SetServiceMethod(j.GetString("serviceMethod"))
	_ = m.Status(true).UnmarshalJSON(j.GetBytes("status"))
	_ = m.Meta().UnmarshalJSON(j.GetBytes("meta"))

	// read body
	m.SetBodyCodec(byte(j.GetInt8("bodyCodec")))
	err = m.UnmarshalBody(gconv.Bytes(j.GetString("body")))
	return err
}

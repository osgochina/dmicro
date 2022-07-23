package jsonrpcproto

import (
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/util/gconv"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/codec"
	"github.com/osgochina/dmicro/drpc/message"
	"github.com/osgochina/dmicro/drpc/proto"
	"github.com/osgochina/dmicro/utils/dbuffer"
	"io"
	"sync"
)

// 注意，json rpc协议实现了标准的jsonrpc协议，
// 不支持tfilter和其他drpc支持的一些标准，使用的时候需要注意

func NewJSONRPCProtoFunc() proto.ProtoFunc {
	return func(rw proto.IOWithReadBuffer) proto.Proto {
		return &jsonRPCProto{
			id:   3,
			name: "jsonrpc",
			r:    rw,
			w:    rw,
		}
	}
}

var _ proto.Proto = new(jsonRPCProto)

type jsonRPCProto struct {
	r    io.Reader
	w    io.Writer
	rMu  sync.Mutex
	name string
	id   byte
}

func (that *jsonRPCProto) Version() (byte, string) {
	return that.id, that.name
}

type Request struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Id      int64       `json:"id"`
	Context interface{} `json:"context"`
}

type Response struct {
	Jsonrpc string      `json:"jsonrpc"`
	Id      int64       `json:"id"`
	Result  interface{} `json:"result"`
	Context interface{} `json:"context"`
}

type ErrorData struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	Class   string `json:"class"`
}

type ErrorStruct struct {
	Code    int32       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ErrorResponse struct {
	Jsonrpc string      `json:"jsonrpc"`
	Id      int64       `json:"id"`
	Context interface{} `json:"context"`
	Error   ErrorStruct `json:"error"`
}

// Pack 打包
func (that *jsonRPCProto) Pack(m proto.Message) error {
	bb := dbuffer.GetByteBuffer()
	defer dbuffer.ReleaseByteBuffer(bb)
	if m.MType() == message.TypeCall {
		body := &Request{
			Jsonrpc: "2.0",
			Method:  m.ServiceMethod(),
			Params:  m.Body(),
			Id:      gconv.Int64(m.Seq()),
			Context: m.Meta(),
		}
		m.SetBody(body)
		bodyBytes, err := m.MarshalBody()
		if err != nil {
			return err
		}
		_, err = bb.Write(bodyBytes)
		_ = bb.WriteByte('\r')
		_ = bb.WriteByte('\n')
		if err != nil {
			return err
		}
	}
	if m.MType() == message.TypeReply {
		m.SetBodyCodec(codec.JsonId)
		if m.StatusOK() {
			body := &Response{
				Jsonrpc: "2.0",
				Id:      gconv.Int64(m.Seq()),
				Result:  m.Body(),
				Context: m.Meta(),
			}
			m.SetBody(body)
		} else {
			body := &ErrorResponse{
				Jsonrpc: "2.0",
				Id:      gconv.Int64(m.Seq()),
				Context: m.Meta(),
				Error: ErrorStruct{
					Code:    m.Status().Code(),
					Message: m.Status().Msg(),
					Data:    m.Body(),
				},
			}
			m.SetBody(body)
		}
		bodyBytes, err := m.MarshalBody()
		if err != nil {
			return err
		}
		_, err = bb.Write(bodyBytes)
		_ = bb.WriteByte('\r')
		_ = bb.WriteByte('\n')
		if err != nil {
			return err
		}
	}
	// 写入消息
	_, err := that.w.Write(bb.B)
	if err != nil {
		return err
	}
	return err
}

// Unpack 解包
func (that *jsonRPCProto) Unpack(m proto.Message) error {
	bb := dbuffer.GetByteBuffer()
	defer dbuffer.ReleaseByteBuffer(bb)
	that.rMu.Lock()
	defer that.rMu.Unlock()
	data, err := that.RecvLine()
	if err != nil {
		return err
	}
	if err = m.SetSize(gconv.Uint32(len(data))); err != nil {
		return err
	}
	json := gjson.New(data)
	if json.IsNil() {
		return gerror.New("解包失败")
	}
	m.SetSeq(json.GetInt32("id"))
	for k, v := range json.GetMap("context") {
		m.Meta().Set(k, v)
	}
	method := json.GetString("method")
	if len(method) > 0 {
		m.SetMType(message.TypeCall)
		m.SetServiceMethod(method)
	} else {
		m.SetMType(message.TypeReply)
	}
	m.SetBodyCodec(codec.JsonId)
	if m.Status(false) == nil {
		m.Status(true)
	}

	if m.MType() == message.TypeCall {
		bt, _ := json.GetJson("params").ToJson()
		return m.UnmarshalBody(bt)
	} else if m.MType() == message.TypeReply {
		if json.GetJson("result").IsNil() {
			errorJson := json.GetJson("error")
			m.SetStatus(drpc.NewStatus(errorJson.GetInt32("code"), errorJson.GetString("message")))
			return m.UnmarshalBody(errorJson.GetBytes("data"))
		} else {
			bt, _ := json.GetJson("result").ToJson()
			return m.UnmarshalBody(bt)
		}
	}
	return nil
}

func (that *jsonRPCProto) RecvLine() ([]byte, error) {
	var err error
	var buffer []byte
	var n int
	data := make([]byte, 0)
	buffer = make([]byte, 1)
	for {
		n, err = io.ReadFull(that.r, buffer)
		if n > 0 {
			if buffer[0] == '\n' {
				data = append(data, buffer[:len(buffer)-1]...)
				break
			} else {
				data = append(data, buffer...)
			}
		}
		if err != nil {
			break
		}
	}
	return data, err
}

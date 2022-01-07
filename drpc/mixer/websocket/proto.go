package websocket

import (
	"bytes"
	"github.com/osgochina/dmicro/drpc/mixer/websocket/jsonSubProto"
	"github.com/osgochina/dmicro/drpc/proto"
	"github.com/osgochina/dmicro/drpc/socket"
	"github.com/osgochina/dmicro/utils/dbuffer"
	"golang.org/x/net/websocket"
)

// 默认支持的协议是json
var defaultProto = jsonSubProto.NewJSONSubProtoFunc()

// NewWsProtoFunc 创建websocket组件支持的编解码处理器
// subProto: 真实的编解码协议，对endpoint来说，它打包和解包是使用了wsproto，实际wsproto是调用了子协议的打包和解包功能
func NewWsProtoFunc(subProto ...proto.ProtoFunc) proto.ProtoFunc {
	return func(rw proto.IOWithReadBuffer) socket.Proto {
		//如果socket对象不是websocket链接，则直接使用子协议打包和解包
		connIFace := rw.(socket.UnsafeSocket).RawLocked()
		conn, ok := connIFace.(*websocket.Conn)
		if !ok {
			if len(subProto) > 0 {
				return subProto[0](rw)
			}
			return defaultProto(rw)
		}
		// 如果是websocket协议，则需要使用虚拟桥接的方式打包和解包
		subConn := newVirtualConn()
		p := &wsProto{
			id:      'w',
			name:    "websocket",
			conn:    conn,
			subConn: subConn,
		}
		if len(subProto) > 0 {
			p.subProto = subProto[0](subConn)
		} else {
			p.subProto = defaultProto(subConn)
		}
		return p
	}
}

// websocket的rpc打包和解包协议
type wsProto struct {
	id       byte
	name     string
	conn     *websocket.Conn
	subProto socket.Proto
	subConn  *virtualConn
}

func (that *wsProto) Version() (byte, string) {
	return that.id, that.name
}

// Pack 打包
func (that *wsProto) Pack(m proto.Message) error {
	//先使用子协议打包消息，获取到了打包后的数据在只用websocket的消息发送方法传输数据
	that.subConn.w.Reset()
	err := that.subProto.Pack(m)
	if err != nil {
		return err
	}
	return websocket.Message.Send(that.conn, that.subConn.w.Bytes())
}

// Unpack 解包
func (that *wsProto) Unpack(m proto.Message) error {
	// 先从websocket中读取原始二进制消息，再使用子协议解包消息
	err := websocket.Message.Receive(that.conn, that.subConn.rBytes)
	if err != nil {
		return err
	}
	that.subConn.r = bytes.NewBuffer(*that.subConn.rBytes)
	return that.subProto.Unpack(m)
}

// 虚拟链接，读写都是在操作内存
type virtualConn struct {
	rBytes *[]byte
	r      *bytes.Buffer
	w      *dbuffer.ByteBuffer
}

func newVirtualConn() *virtualConn {
	buf := new([]byte)
	return &virtualConn{
		rBytes: buf,
		r:      bytes.NewBuffer(*buf),
		w:      dbuffer.GetByteBuffer(),
	}
}

func (that *virtualConn) Read(p []byte) (int, error) {
	return that.r.Read(p)
}

func (that *virtualConn) Write(p []byte) (int, error) {
	return that.w.Write(p)
}

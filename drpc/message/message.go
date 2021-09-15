// Package message 消息对象
package message

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/util/gconv"
	"github.com/osgochina/dmicro/drpc/codec"
	"github.com/osgochina/dmicro/drpc/status"
	"github.com/osgochina/dmicro/drpc/tfilter"
	"strconv"
)

// Header 消息头
type Header interface {
	// Seq 序列号
	Seq() int32
	// SetSeq 设置序列号
	SetSeq(int32)
	// MType 消息类型 有五种：CALL,REPLY,PUSH,AUTH_CALL,AUTH_REPLY
	MType() byte
	// SetMType 设置消息类型 有五种：CALL,REPLY,PUSH,AUTH_CALL,AUTH_REPLY
	SetMType(byte)
	// ServiceMethod 请求的服务方法名称 长度必须小于255字节 max <= 255
	ServiceMethod() string
	// SetServiceMethod 设置请求的服务方法名
	SetServiceMethod(string)
	// StatusOK 判断当前消息是否是 OK
	StatusOK() bool
	// Status 返回当前消息的状态，包含code，msg，cause或者stack，
	//如果消息是nil或者autoInit传入了true，则返回一个 code为OK的新对象
	Status(autoInit ...bool) *status.Status
	// SetStatus 设置消息的状态
	SetStatus(*status.Status)
	// Meta 获取消息的元数据，数据在传输的时候是使用了序列化串，最大长度为 max len ≤ 65535
	Meta() *gmap.Map
}

// Body 消息体
type Body interface {
	// BodyCodec 消息体编码格式
	BodyCodec() byte
	// SetBodyCodec 设置消息体编码格式
	SetBodyCodec(bodyCodec byte)
	// Body 返回消息体内容
	Body() interface{}
	// SetBody 设置消息体内容
	SetBody(body interface{})
	// SetNewBody 设置一个函数，该函数根据消息头生成一个新的消息体
	SetNewBody(NewBodyFunc)
	// MarshalBody 编码消息体
	MarshalBody() ([]byte, error)
	// UnmarshalBody 解码消息体
	UnmarshalBody(bodyBytes []byte) error
}

// NewBodyFunc 根据消息头，生成消息体，这个函数只会在读取connection上的
type NewBodyFunc func(Header) interface{}

// Message 消息
type Message interface {
	Reset(settings ...MsgSetting) Message
	Header
	Body
	// PipeTFilter 针对传入的数据做
	PipeTFilter() *tfilter.PipeTFilter

	// Size 消息长度
	Size() uint32
	// SetSize 设置消息长度，如果长度超长了，则返回错误
	SetSize(size uint32) error
	// Context 返回消息的上下文对象
	Context() context.Context
	// String 把消息转换成可打印的字符串
	String() string
	// AsHeader 把消息转换成header接口
	AsHeader() Header
	// AsBody 把消息转换成body接口
	AsBody() Body
	//防止消息在包外实现
	messageIdentity() *message
}

type message struct {
	serviceMethod string
	status        *status.Status
	meta          *gmap.Map
	body          interface{}
	newBodyFunc   NewBodyFunc
	pipeTFilter   *tfilter.PipeTFilter
	ctx           context.Context
	size          uint32
	seq           int32
	mType         byte
	bodyCodec     byte
}

// AsHeader 把消息转换成header对象
func (that *message) AsHeader() Header { return that }

// AsBody 把消息转换成body对象
func (that *message) AsBody() Body { return that }

//防止在包外实现消息接口
func (*message) messageIdentity() *message { return nil }

// Reset 重置消息并返回它自己
func (that *message) Reset(settings ...MsgSetting) Message {
	that.body = nil
	that.status = nil
	that.meta.Clear()
	that.pipeTFilter.Reset()
	that.newBodyFunc = nil
	that.seq = 0
	that.mType = 0
	that.serviceMethod = ""
	that.size = 0
	that.ctx = nil
	that.bodyCodec = codec.NilCodecID
	that.doSetting(settings...)
	return that
}

//针对消息执行一些操作
func (that *message) doSetting(settings ...MsgSetting) {
	for _, fn := range settings {
		if fn != nil {
			fn(that)
		}
	}
}

// Context 获取消息的上下文
func (that *message) Context() context.Context {
	if that.ctx == nil {
		return context.Background()
	}
	return that.ctx
}

// Seq 返回消息序列号
func (that *message) Seq() int32 {
	return that.seq
}

// SetSeq 设置消息序列号
func (that *message) SetSeq(seq int32) {
	that.seq = seq
}

// MType 返回消息类型，如下：CALL, REPLY, PUSH.
func (that *message) MType() byte {
	return that.mType
}

// SetMType 设置消息类型，可选值如下：CALL, REPLY, PUSH.
func (that *message) SetMType(mType byte) {
	that.mType = mType
}

// ServiceMethod 返回消息中的服务器接口名称
func (that *message) ServiceMethod() string {
	return that.serviceMethod
}

// SetServiceMethod 设置服务器接口名
func (that *message) SetServiceMethod(serviceMethod string) {
	that.serviceMethod = serviceMethod
}

// StatusOK 判断消息状态是否是OK
func (that *message) StatusOK() bool {
	return that.status.OK()
}

// Status 返回消息状态，如果当前消息状态是空，或者设置了自动初始化状态 autoInit = true,则创建一个值为OK的状态
func (that *message) Status(autoInit ...bool) *status.Status {
	if that.status == nil && len(autoInit) > 0 && autoInit[0] {
		that.status = new(status.Status)
	}
	return that.status
}

// SetStatus 设置消息状态
func (that *message) SetStatus(stat *status.Status) {
	that.status = stat
}

// Meta 返回元数据
func (that *message) Meta() *gmap.Map {
	return that.meta
}

// BodyCodec 获取消息的包体编码格式
func (that *message) BodyCodec() byte {
	return that.bodyCodec
}

// SetBodyCodec 设置消息的包体编码格式
func (that *message) SetBodyCodec(bodyCodec byte) {
	that.bodyCodec = bodyCodec
}

// Body 返回消息的包体内容
func (that *message) Body() interface{} {
	return that.body
}

// SetBody 设置消息的包体内容
func (that *message) SetBody(body interface{}) {
	that.body = body
}

// SetNewBody 设置自定义包体创建函数
func (that *message) SetNewBody(newBodyFunc NewBodyFunc) {
	that.newBodyFunc = newBodyFunc
}

// MarshalBody 编码包体
func (that *message) MarshalBody() ([]byte, error) {
	switch body := that.body.(type) {
	default:
		c, err := codec.Get(that.bodyCodec)
		if err != nil {
			return []byte{}, err
		}
		return c.Marshal(body)
	case nil:
		return []byte{}, nil
	case *[]byte:
		if body == nil {
			return []byte{}, nil
		}
		return *body, nil
	case []byte:
		return body, nil
	}
}

// UnmarshalBody 解码包体
func (that *message) UnmarshalBody(bodyBytes []byte) error {
	//如果包体创建函数存在，则根据业务情况自己实现一个包体解码逻辑
	if that.body == nil && that.newBodyFunc != nil {
		that.body = that.newBodyFunc(that)
	}
	length := len(bodyBytes)
	if length == 0 {
		return nil
	}
	switch body := that.body.(type) {
	default:
		c, err := codec.Get(that.bodyCodec)
		if err != nil {
			return err
		}
		return c.Unmarshal(bodyBytes, that.body)
	case nil:
		return nil
	case *[]byte:
		if cap(*body) < length {
			*body = make([]byte, length)
		} else {
			*body = (*body)[:length]
		}
		copy(*body, bodyBytes)
		return nil
	}
}

// PipeTFilter 返回消息的处理管道
func (that *message) PipeTFilter() *tfilter.PipeTFilter {
	return that.pipeTFilter
}

// Size 获取消息的长度
func (that *message) Size() uint32 {
	return that.size
}

// SetSize 设置消息的长度
func (that *message) SetSize(size uint32) error {
	err := checkMessageSize(size)
	if err != nil {
		return err
	}
	that.size = size
	return nil
}

const messageFormat = `
{
  "seq": %d,
  "mType": %d,
  "serviceMethod": %q,
  "status": %q,
  "meta": %q,
  "bodyCodec": %d,
  "body": %s,
  "pipeTFilter": %s,
  "size": %d
}`

func (that *message) String() string {
	var pipeTFilterIDs = make([]int, that.pipeTFilter.Len())
	for i, id := range that.pipeTFilter.IDs() {
		pipeTFilterIDs[i] = int(id)
	}
	idsBytes, _ := json.Marshal(pipeTFilterIDs)
	b, _ := json.Marshal(that.body)
	return fmt.Sprintf(messageFormat,
		that.seq,
		that.mType,
		that.serviceMethod,
		that.status.String(),
		that.meta.String(),
		that.bodyCodec,
		b,
		idsBytes,
		that.size,
	)
}

// GetAcceptBodyCodec 获取发送方希望接受的正文编解码器。
////注意:如果指定的编解码器无效，接收器将忽略配对数据
func GetAcceptBodyCodec(meta *gmap.Map) (byte, bool) {
	s := gconv.Bytes(meta.Get(MetaAcceptBodyCodec))
	if len(s) == 0 || len(s) > 3 {
		return 0, false
	}
	b, err := strconv.ParseUint(gconv.String(s), 10, 8)
	if err != nil {
		return 0, false
	}
	c := byte(b)
	return c, c != codec.NilCodecID
}

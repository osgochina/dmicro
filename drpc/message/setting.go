package message

import (
	"context"
	"errors"
	"github.com/gogf/gf/util/gconv"
	"github.com/osgochina/dmicro/drpc/status"
	"math"
)

//消息类型
const (
	TypeUndefined byte = 0 //未知类型
	TypeCall      byte = 1 // call
	TypeReply     byte = 2 // reply to call
	TypePush      byte = 3
	TypeAuthCall  byte = 4
	TypeAuthReply byte = 5
)

func TypeText(typ byte) string {
	switch typ {
	case TypeCall:
		return "CALL"
	case TypeReply:
		return "REPLY"
	case TypePush:
		return "PUSH"
	case TypeAuthCall:
		return "AUTH_CALL"
	case TypeAuthReply:
		return "AUTH_REPLY"
	default:
		return "Undefined"
	}
}

const (
	// MetaRealIP real IP metadata key
	MetaRealIP = "X-Real-IP"
	// MetaAcceptBodyCodec the key of body codec that the sender wishes to accept
	MetaAcceptBodyCodec = "X-Accept-Body-Codec"
)

var (
	messageSizeLimit uint32 = math.MaxUint32
	// ErrExceedMessageSizeLimit error
	ErrExceedMessageSizeLimit = errors.New("size of package exceeds limit")
)

type MsgSetting func(Message)

// WithNothing 什么也不做
func WithNothing() MsgSetting {
	return func(Message) { return }
}

// WithContext 设置消息的上下文对象
func WithContext(ctx context.Context) MsgSetting {
	return func(m Message) {
		m.(*message).ctx = ctx
	}
}

// WithServiceMethod 设置消息的服务器接口名
func WithServiceMethod(serviceMethod string) MsgSetting {
	return func(m Message) {
		m.SetServiceMethod(serviceMethod)
	}
}

// WithStatus 设置消息的状态
func WithStatus(stat *status.Status) MsgSetting {
	return func(m Message) {
		m.SetStatus(stat)
	}
}

// WithSetMeta 添加消息的元数据
func WithSetMeta(key, value string) MsgSetting {
	return func(m Message) {
		m.Meta().Set(key, value)
	}
}

// WithSetMetas 使用数组添加元数据
func WithSetMetas(metas map[string]interface{}) MsgSetting {
	return func(m Message) {
		for key, value := range metas {
			m.Meta().Set(key, gconv.String(value))
		}
	}
}

// WithDelMeta 删除消息元数据
func WithDelMeta(key string) MsgSetting {
	return func(m Message) {
		m.Meta().Remove(key)
	}
}

// WithBodyCodec 设置消息的消息体编码格式
func WithBodyCodec(bodyCodec byte) MsgSetting {
	return func(m Message) {
		m.SetBodyCodec(bodyCodec)
	}
}

// WithBody 设置消息体的内容
func WithBody(body interface{}) MsgSetting {
	return func(m Message) {
		m.SetBody(body)
	}
}

// WithNewBody 设置创建消息体的函数
func WithNewBody(newBodyFunc NewBodyFunc) MsgSetting {
	return func(m Message) {
		m.SetNewBody(newBodyFunc)
	}
}

// WithMType 设置消息类型，改方法会改变整个消息的处理逻辑，请明确知道你要做什么的时候使用
func WithMType(mType byte) MsgSetting {
	return func(m Message) {
		m.SetMType(mType)
	}
}

// WithXFerPipe 设置消息的管道类型
func WithXFerPipe(filterID ...byte) MsgSetting {
	return func(m Message) {
		if err := m.PipeTFilter().Append(filterID...); err != nil {
			panic(err)
		}
	}
}

// MsgSizeLimit 获取消息的最大长度
func MsgSizeLimit() uint32 {
	return messageSizeLimit
}

// SetMsgSizeLimit 设置消息的最大长度
func SetMsgSizeLimit(maxMessageSize uint32) {
	if maxMessageSize <= 0 {
		messageSizeLimit = math.MaxUint32
	} else {
		messageSizeLimit = maxMessageSize
	}
}

//检查消息的最大长度
func checkMessageSize(messageSize uint32) error {
	if messageSize > messageSizeLimit {
		return ErrExceedMessageSizeLimit
	}
	return nil
}

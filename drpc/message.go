package drpc

import (
	"github.com/osgochina/dmicro/drpc/codec"
	"github.com/osgochina/dmicro/drpc/message"
	"strconv"
)

var (
	TypeText = message.TypeText

	MetaRealIP          = message.MetaRealIP
	MetaAcceptBodyCodec = message.MetaAcceptBodyCodec

	TypeUndefined = message.TypeUndefined
	TypeCall      = message.TypeCall
	TypeReply     = message.TypeReply
	TypePush      = message.TypePush
	TypeAuthCall  = message.TypeAuthCall
	TypeAuthReply = message.TypeAuthReply
)

var (
	GetMessage = message.GetMessage
	PutMessage = message.PutMessage
)

var (
	// WithNothing nothing to do.
	//  func WithNothing() MessageSetting
	WithNothing = message.WithNothing
	// WithStatus sets the message status.
	// TYPE:
	//  func WithStatus(stat *Status) MessageSetting
	WithStatus = message.WithStatus
	// WithContext sets the message handling context.
	//  func WithContext(ctx context.Context) MessageSetting
	WithContext = message.WithContext
	// WithServiceMethod sets the message service method.
	// SUGGEST: max len ≤ 255!
	//  func WithServiceMethod(serviceMethod string) MessageSetting
	WithServiceMethod = message.WithServiceMethod
	// WithSetMeta sets 'key=value' metadata argument.
	// SUGGEST: urlencoded string max len ≤ 65535!
	//  func WithSetMeta(key, value string) MessageSetting
	WithSetMeta = message.WithSetMeta

	WithSetMetas = message.WithSetMetas
	// WithDelMeta deletes metadata argument.
	//   func WithDelMeta(key string) MessageSetting
	WithDelMeta = message.WithDelMeta
	// WithBodyCodec sets the body codec.
	//  func WithBodyCodec(bodyCodec byte) MessageSetting
	WithBodyCodec = message.WithBodyCodec
	// WithBody sets the body object.
	//  func WithBody(body interface{}) MessageSetting
	WithBody = message.WithBody
	// WithNewBody resets the function of geting body.
	//  NOTE: newBodyFunc is only for reading form connection.
	//  func WithNewBody(newBodyFunc socket.NewBodyFunc) MessageSetting
	WithNewBody = message.WithNewBody

	// WithTFilterPipe 设置传输过滤器.
	// 提示: 如果filterID未注册，则会产生Panic错误。
	// 建议: 最大长度不能超过255!
	//  func WithTFilterPipe(filterID ...byte) MessageSetting
	WithTFilterPipe = message.WithTFilterPipe

	GetAcceptBodyCodec = message.GetAcceptBodyCodec
)

func WithRealIP(ip string) message.MsgSetting {
	return message.WithSetMeta(MetaRealIP, ip)
}

// WithAcceptBodyCodec sets the body codec that the sender wishes to accept.
// NOTE: If the specified codec is invalid, the receiver will ignore the mate data.
func WithAcceptBodyCodec(bodyCodec byte) message.MsgSetting {
	if bodyCodec == codec.NilCodecID {
		return WithNothing()
	}
	return message.WithSetMeta(MetaAcceptBodyCodec, strconv.FormatUint(uint64(bodyCodec), 10))
}

func withMType(mType byte) message.MsgSetting {
	return func(m message.Message) {
		m.SetMType(mType)
	}
}

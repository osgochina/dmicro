package securebody

import (
	"crypto/aes"
	"fmt"
	"github.com/gogf/gf/crypto/gmd5"
	"github.com/gogf/gf/util/gconv"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/internal"
	"github.com/osgochina/dmicro/drpc/message"
)

// NewSecureBodyPlugin 创建插件
// cipherKey: 加密key
// statCode:  自定义错误码
func NewSecureBodyPlugin(cipherKey string, statCode ...int32) drpc.Plugin {
	b := gconv.Bytes(cipherKey)
	if _, err := aes.NewCipher(b); err != nil {
		internal.Fatalf("secure: %v", err)
	}
	version, err := gmd5.EncryptBytes(b)
	if err != nil {
		internal.Fatalf("md5 error: %v", err)
	}
	var code = drpc.CodeConflict
	if len(statCode) > 0 {
		code = statCode[0]
	}
	return &secureBodyPlugin{
		version:   version,
		cipherKey: b,
		statCode:  code,
	}
}

// WithSecureMeta 强制要求传输加密
func WithSecureMeta() message.MsgSetting {
	return func(message message.Message) {
		message.Meta().Set(SecureMetaKey, "true")
	}
}

// WithReplySecureMeta 要求服务端回复消息的时候加密
func WithReplySecureMeta(secure bool) message.MsgSetting {
	s := fmt.Sprintf("%v", secure)
	return func(message message.Message) {
		message.Meta().Set(ReplySecureMetaKey, s)
	}
}

// EnforceSecure 强制加密
func EnforceSecure(output message.Message) {
	output.Meta().Set(SecureMetaKey, "true")
}

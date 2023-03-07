package securebody

import (
	"fmt"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/crypto/gaes"
	"github.com/gogf/gf/v2/encoding/gbase64"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/osgochina/dmicro/drpc"
)

const (
	// SecureMetaKey 所有消息体都需要加密
	SecureMetaKey = "X-Secure-Body"
	// ReplySecureMetaKey 发送的消息不需要加密，回复的消息需要加密
	ReplySecureMetaKey = "X-Reply-Secure-Body"
	// RawBody 原始的消息结构体对象，使用该插件以后，传输层的消息结构体是Encrypt，但是用户使用的是业务层定义的对象，需要保存该对象，方便解码后给到用户层
	RawBody = "secure_body_raw_body"
	// 如果设置了 ReplySecureMetaKey，则做好标记
	replyEncrypt = "secure_body_reply_encrypt"
)

type secureBodyPlugin struct {
	version   string
	cipherKey []byte
	statCode  int32
}

// 写入消息之前
var _ drpc.BeforeWriteCallPlugin = new(secureBodyPlugin)
var _ drpc.BeforeWritePushPlugin = new(secureBodyPlugin)
var _ drpc.BeforeWriteReplyPlugin = new(secureBodyPlugin)

// 读取消息内容之前
var _ drpc.BeforeReadCallBodyPlugin = new(secureBodyPlugin)
var _ drpc.BeforeReadPushBodyPlugin = new(secureBodyPlugin)
var _ drpc.BeforeReadReplyBodyPlugin = new(secureBodyPlugin)

// 读取消息内容之后
var _ drpc.AfterReadCallBodyPlugin = new(secureBodyPlugin)
var _ drpc.AfterReadPushBodyPlugin = new(secureBodyPlugin)
var _ drpc.AfterReadReplyBodyPlugin = new(secureBodyPlugin)

// Name 插件名
func (that *secureBodyPlugin) Name() string {
	return "secure(encrypt/decrypt)"
}

// BeforeWriteCall 写入Call消息之前执行
func (that *secureBodyPlugin) BeforeWriteCall(ctx drpc.WriteCtx) *drpc.Status {
	// 如果写入之前消息已经已经出错，则不进行任何处理
	if ctx.Status() != nil {
		return nil
	}
	// 判断消息是否需要加密
	if !isSecureBody(ctx.Output().Meta()) {
		//如果请求的消息未加密，但是请求端告诉服务端，回复消息需要加密，则强制回复消息加密
		_, found := ctx.Swap().Search(replyEncrypt)
		if !found {
			return nil
		}
		// 强制加密
		EnforceSecure(ctx.Output())
	}
	// 先对消息内容编码
	bodyBytes, err := ctx.Output().MarshalBody()
	if err != nil {
		return drpc.NewStatus(that.statCode, "编码原始消息内容失败", err.Error())
	}
	// 对消息体加密
	ciphertext, err := gaes.Encrypt(bodyBytes, that.cipherKey)
	if err != nil {
		return drpc.NewStatus(that.statCode, "加密消息失败", err.Error())
	}
	// 设置加密后的消息内容结构体
	ctx.Output().SetBody(&Encrypt{
		CipherVersion: that.version,
		Ciphertext:    gbase64.EncodeToString(ciphertext),
	})
	return nil
}

// BeforeWritePush 写入Push消息之前执行
func (that *secureBodyPlugin) BeforeWritePush(ctx drpc.WriteCtx) *drpc.Status {
	return that.BeforeWriteCall(ctx)
}

// BeforeWriteReply 写入Reply消息之前执行
func (that *secureBodyPlugin) BeforeWriteReply(ctx drpc.WriteCtx) *drpc.Status {
	return that.BeforeWriteCall(ctx)
}

// BeforeReadCallBody 读取CALL消息内容之前执行
func (that *secureBodyPlugin) BeforeReadCallBody(ctx drpc.ReadCtx) *drpc.Status {
	b := ctx.PeekMeta(ReplySecureMetaKey)
	reply := gconv.String(b)
	useDecrypt := isSecureBody(ctx.Input().Meta())
	// 如果请求的消息不是加密消息
	if !useDecrypt {
		if reply == "true" {
			ctx.Swap().Set(replyEncrypt, 1)
		}
		return nil
	}
	// 如果消息本身是加密过的，且请求方还强制消息为加密，则设置回复强制加密
	if reply != "false" {
		ctx.Swap().Set(replyEncrypt, 1)
	}

	ctx.Swap().Set(RawBody, ctx.Input().Body())
	ctx.Input().SetBody(new(Encrypt))
	return nil
}

// BeforeReadPushBody 读取Push消息内容之前执行
func (that *secureBodyPlugin) BeforeReadPushBody(ctx drpc.ReadCtx) *drpc.Status {
	return that.BeforeReadCallBody(ctx)
}

// BeforeReadReplyBody 读取Reply消息内容之前执行
func (that *secureBodyPlugin) BeforeReadReplyBody(ctx drpc.ReadCtx) *drpc.Status {
	return that.BeforeReadCallBody(ctx)
}

// AfterReadCallBody 读取CALL消息内容之后执行
func (that *secureBodyPlugin) AfterReadCallBody(ctx drpc.ReadCtx) *drpc.Status {
	rawBody, found := ctx.Swap().Search(RawBody)
	if !found {
		return nil
	}
	var obj = ctx.Input().Body().(*Encrypt)
	var version = obj.GetCipherVersion()
	var bodyBytes []byte
	var err error

	if len(version) > 0 {
		if version != that.version {
			return drpc.NewStatus(that.statCode, "解密消息内容失败", fmt.Sprintf("加密key的版本不一致, get:%q, want:%q", obj.GetCipherVersion(), that.version))
		}
		bodyBytes, err = gaes.Decrypt(gbase64.MustDecodeString(obj.GetCiphertext()), that.cipherKey)
		if err != nil {
			return drpc.NewStatus(that.statCode, "解密消息内容失败", err.Error())
		}
	}
	ctx.Swap().Remove(RawBody)
	ctx.Input().SetBody(rawBody)
	err = ctx.Input().UnmarshalBody(bodyBytes)
	if err != nil {
		return drpc.NewStatus(that.statCode, "无法解密原始消息", err.Error())
	}
	return nil
}

// AfterReadPushBody 读取PUSH消息内容之后执行
func (that *secureBodyPlugin) AfterReadPushBody(ctx drpc.ReadCtx) *drpc.Status {
	return that.AfterReadCallBody(ctx)
}

// AfterReadReplyBody 读取REPLY消息内容之后执行
func (that *secureBodyPlugin) AfterReadReplyBody(ctx drpc.ReadCtx) *drpc.Status {
	return that.AfterReadCallBody(ctx)
}

// 判断当前消息是否是加密消息
func isSecureBody(meta *gmap.Map) bool {
	v := meta.GetVar(SecureMetaKey)
	if !v.IsNil() && v.Bool() {
		return true
	}
	// 注意，如果消息元数据中的 SecureMetaKey 不为true，则不需要解密该消息
	if !v.IsNil() {
		meta.Remove(SecureMetaKey)
	}
	return false
}

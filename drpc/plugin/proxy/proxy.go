package proxy

import (
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/message"
)

type (
	Forwarder interface {
		CallForwarder
		PushForwarder
	}

	CallForwarder interface {
		Call(uri string, arg interface{}, result interface{}, setting ...message.MsgSetting) drpc.CallCmd
	}
	PushForwarder interface {
		Push(uri string, arg interface{}, setting ...message.MsgSetting) *drpc.Status
	}
	Label struct {
		SessionID     string
		RealIP        string
		ServiceMethod string
	}
	proxy struct {
		callForwarder func(*Label) CallForwarder
		pushForwarder func(*Label) PushForwarder
	}
)

// NewProxyCallPlugin 创建call方法的代理插件
func NewProxyCallPlugin(fn func(*Label) CallForwarder) drpc.Plugin {
	return &proxy{callForwarder: fn}
}

// NewProxyPushPlugin 创建push方法的代理插件
func NewProxyPushPlugin(fn func(*Label) PushForwarder) drpc.Plugin {
	return &proxy{pushForwarder: fn}
}

// NewProxyPlugin 创建代理插件，支持call方法和push方法
func NewProxyPlugin(fn func(*Label) Forwarder) drpc.Plugin {
	return &proxy{
		callForwarder: func(label *Label) CallForwarder {
			return fn(label)
		},
		pushForwarder: func(label *Label) PushForwarder {
			return fn(label)
		},
	}
}

var (
	_ drpc.AfterNewEndpointPlugin = new(proxy)
)

// Name 插件名称
func (that *proxy) Name() string {
	return "proxy"
}

// AfterNewEndpoint 创建Endpoint后触发
func (that *proxy) AfterNewEndpoint(peer drpc.EarlyEndpoint) error {
	// 请求的call方法不存在时候
	if that.callForwarder != nil {
		peer.SetUnknownCall(that.call)
	}
	// 请求的push方法不存在时候
	if that.pushForwarder != nil {
		peer.SetUnknownPush(that.push)
	}
	return nil
}

// 请求的call方法不存在时候，则执行该方法
func (that *proxy) call(ctx drpc.UnknownCallCtx) (interface{}, *drpc.Status) {
	var (
		label    Label
		settings = make([]message.MsgSetting, 0, 16)
	)
	label.SessionID = ctx.Session().ID()
	ctx.VisitMeta(func(key, value interface{}) bool {
		settings = append(settings, drpc.WithSetMeta(gconv.String(key), gconv.String(value)))
		return true
	})
	var (
		result      []byte
		realIPBytes = ctx.PeekMeta(drpc.MetaRealIP)
	)
	if len(gconv.String(realIPBytes)) == 0 {
		label.RealIP = ctx.IP()
		settings = append(settings, drpc.WithSetMeta(drpc.MetaRealIP, label.RealIP))
	} else {
		label.RealIP = gconv.String(realIPBytes)
	}
	label.ServiceMethod = ctx.ServiceMethod()
	callCmd := that.callForwarder(&label).Call(label.ServiceMethod, ctx.InputBodyBytes(), &result, settings...)
	callCmd.InputMeta().Iterator(func(key, value interface{}) bool {
		ctx.SetMeta(gconv.String(key), gconv.String(value))
		return true
	})
	stat := callCmd.Status()
	if !stat.OK() && stat.Code() < 200 && stat.Code() > 99 {
		stat.SetCode(drpc.CodeBadGateway)
		stat.SetMsg(drpc.CodeText(drpc.CodeBadGateway))
	}
	return result, stat
}

// 请求的push方法不存在时候，则执行该方法
func (that *proxy) push(ctx drpc.UnknownPushCtx) *drpc.Status {
	var (
		label    Label
		settings = make([]message.MsgSetting, 0, 16)
	)
	label.SessionID = ctx.Session().ID()
	ctx.VisitMeta(func(key, value interface{}) bool {
		settings = append(settings, drpc.WithSetMeta(gconv.String(key), gconv.String(value)))
		return true
	})
	if realIPBytes := ctx.PeekMeta(drpc.MetaRealIP); len(gconv.String(realIPBytes)) == 0 {
		label.RealIP = ctx.IP()
		settings = append(settings, drpc.WithSetMeta(drpc.MetaRealIP, label.RealIP))
	} else {
		label.RealIP = gconv.String(realIPBytes)
	}
	label.ServiceMethod = ctx.ServiceMethod()
	stat := that.pushForwarder(&label).Push(label.ServiceMethod, ctx.InputBodyBytes(), settings...)
	if !stat.OK() && stat.Code() < 200 && stat.Code() > 99 {
		stat.SetCode(drpc.CodeBadGateway)
		stat.SetMsg(drpc.CodeText(drpc.CodeBadGateway))
	}
	return stat
}

package event

import (
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/logger"
)

const (
	OnAcceptEvent  = "endpoint_accept"
	OnReceiveEvent = "endpoint_receive"
	OnConnectEvent = "endpoint_connect"
	OnCloseEvent   = "endpoint_close"
)

type eventPlugin struct{}

var (
	_ drpc.AfterAcceptPlugin        = new(eventPlugin)
	_ drpc.AfterDialPlugin          = new(eventPlugin)
	_ drpc.AfterDisconnectPlugin    = new(eventPlugin)
	_ drpc.AfterReadCallBodyPlugin  = new(eventPlugin)
	_ drpc.AfterReadPushBodyPlugin  = new(eventPlugin)
	_ drpc.AfterReadReplyBodyPlugin = new(eventPlugin)
)

func NewEventPlugin() *eventPlugin {
	return &eventPlugin{}
}

func (that *eventPlugin) Name() string {
	return "endpoint-event"
}

// AfterAccept 作为服务端角色，接受客户端链接后触发该事件
func (that *eventPlugin) AfterAccept(sess drpc.EarlySession) *drpc.Status {
	eventBus := sess.Endpoint().EventBus()
	if eventBus.HasListeners(OnAcceptEvent) {
		err := eventBus.Publish(newOnAccept(sess))
		if err != nil {
			logger.Warning(err)
		}
	}
	return nil
}

// AfterDial 链接成功后，触发onConnect事件
func (that *eventPlugin) AfterDial(sess drpc.EarlySession, isRedial bool) *drpc.Status {
	eventBus := sess.Endpoint().EventBus()
	if eventBus.HasListeners(OnConnectEvent) {
		err := eventBus.Publish(newOnConnect(sess, isRedial))
		if err != nil {
			logger.Warning(err)
		}
	}
	return nil
}

// AfterDisconnect 链接断开触发该事件
func (that *eventPlugin) AfterDisconnect(sess drpc.BaseSession) *drpc.Status {
	eventBus := sess.Endpoint().EventBus()
	if eventBus.HasListeners(OnCloseEvent) {
		err := eventBus.Publish(newOnClose(sess))
		if err != nil {
			logger.Warning(err)
		}
	}
	return nil
}

// AfterReadCallBody 读取CALL消息的body之后触发该事件
func (that *eventPlugin) AfterReadCallBody(ctx drpc.ReadCtx) *drpc.Status {
	eventBus := ctx.Endpoint().EventBus()
	if eventBus.HasListeners(OnReceiveEvent) {
		err := eventBus.Publish(newOnReceive(ctx))
		if err != nil {
			logger.Warning(err)
		}
	}
	return nil
}

// AfterReadPushBody 读取PUSH消息body之后触发该事件
func (that *eventPlugin) AfterReadPushBody(ctx drpc.ReadCtx) *drpc.Status {
	eventBus := ctx.Endpoint().EventBus()
	if eventBus.HasListeners(OnReceiveEvent) {
		err := eventBus.Publish(newOnReceive(ctx))
		if err != nil {
			logger.Warning(err)
		}
	}
	return nil
}

// AfterReadReplyBody 读取REPLY消息body之后触发该事件
func (that *eventPlugin) AfterReadReplyBody(ctx drpc.ReadCtx) *drpc.Status {
	eventBus := ctx.Endpoint().EventBus()
	if eventBus.HasListeners(OnReceiveEvent) {
		err := eventBus.Publish(newOnReceive(ctx))
		if err != nil {
			logger.Warning(err)
		}
	}
	return nil
}

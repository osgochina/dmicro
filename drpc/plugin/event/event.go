package event

import (
	"context"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/internal"
	"github.com/osgochina/dmicro/eventbus"
)

const (
	OnAcceptEvent  = "endpoint_accept"
	OnReceiveEvent = "endpoint_receive"
	OnConnectEvent = "endpoint_connect"
	OnCloseEvent   = "endpoint_close"
)

type eventPlugin struct {
	bus *eventbus.EventBus
}

var (
	_ drpc.AfterAcceptPlugin        = new(eventPlugin)
	_ drpc.AfterDialPlugin          = new(eventPlugin)
	_ drpc.AfterDialFailPlugin      = new(eventPlugin)
	_ drpc.AfterDisconnectPlugin    = new(eventPlugin)
	_ drpc.AfterReadCallBodyPlugin  = new(eventPlugin)
	_ drpc.AfterReadPushBodyPlugin  = new(eventPlugin)
	_ drpc.AfterReadReplyBodyPlugin = new(eventPlugin)
)

func NewEventPlugin(bus *eventbus.EventBus) *eventPlugin {
	return &eventPlugin{
		bus: bus,
	}
}

func (that *eventPlugin) Name() string {
	return "endpoint-event"
}

// AfterAccept 作为服务端角色，接受客户端链接后触发该事件
func (that *eventPlugin) AfterAccept(sess drpc.EarlySession) *drpc.Status {
	if that.bus.HasListeners(OnAcceptEvent) {
		err := that.bus.Publish(newOnAccept(sess))
		if err != nil {
			internal.Warning(context.TODO(), err)
		}
	}
	return nil
}

// AfterDial 链接成功后，触发onConnect事件
func (that *eventPlugin) AfterDial(sess drpc.EarlySession, isRedial bool) *drpc.Status {
	if that.bus.HasListeners(OnConnectEvent) {
		err := that.bus.Publish(newOnConnect(true, sess, isRedial, nil))
		if err != nil {
			internal.Warning(context.TODO(), err)
		}
	}
	return nil
}

// AfterDialFail 链接失败后，触发onConnect事件
func (that *eventPlugin) AfterDialFail(sess drpc.EarlySession, err error, isRedial bool) *drpc.Status {
	if that.bus.HasListeners(OnConnectEvent) {
		err = that.bus.Publish(newOnConnect(false, sess, isRedial, err))
		if err != nil {
			internal.Warning(context.TODO(), err)
		}
	}
	return nil
}

// AfterDisconnect 链接断开触发该事件
func (that *eventPlugin) AfterDisconnect(sess drpc.BaseSession) *drpc.Status {
	if that.bus.HasListeners(OnCloseEvent) {
		err := that.bus.Publish(newOnClose(sess))
		if err != nil {
			internal.Warning(context.TODO(), err)
		}
	}
	return nil
}

// AfterReadCallBody 读取CALL消息的body之后触发该事件
func (that *eventPlugin) AfterReadCallBody(ctx drpc.ReadCtx) *drpc.Status {
	if that.bus.HasListeners(OnReceiveEvent) {
		err := that.bus.Publish(newOnReceive(ctx))
		if err != nil {
			internal.Warning(context.TODO(), err)
		}
	}
	return nil
}

// AfterReadPushBody 读取PUSH消息body之后触发该事件
func (that *eventPlugin) AfterReadPushBody(ctx drpc.ReadCtx) *drpc.Status {
	return that.AfterReadCallBody(ctx)
}

// AfterReadReplyBody 读取REPLY消息body之后触发该事件
func (that *eventPlugin) AfterReadReplyBody(ctx drpc.ReadCtx) *drpc.Status {
	return that.AfterReadCallBody(ctx)
}

package drpc

import (
	"github.com/osgochina/dmicro/drpc/status"
	"github.com/osgochina/dmicro/logger"
	"net"
)

// Plugin 插件的基础对象
type Plugin interface {
	Name() string
}

// BeforeNewEndpointPlugin 创建Endpoint之前触发该事件
type BeforeNewEndpointPlugin interface {
	Plugin
	BeforeNewEndpoint(*EndpointConfig, *PluginContainer) error
}

// beforeNewEndpoint 在创建endpoint之前执行已定义的插件。
func (that *PluginContainer) beforeNewEndpoint(endpointConfig *EndpointConfig) {
	var err error
	for _, plugin := range that.plugins {
		if _plugin, ok := plugin.(BeforeNewEndpointPlugin); ok {
			if err = _plugin.BeforeNewEndpoint(endpointConfig, that); err != nil {
				logger.Fatalf("[BeforeNewEndpoint:%s] %s", plugin.Name(), err.Error())
				return
			}
		}
	}
}

// AfterNewEndpointPlugin 创建Endpoint之后触发该事件
type AfterNewEndpointPlugin interface {
	Plugin
	AfterNewEndpoint(EarlyEndpoint) error
}

// afterNewEndpoint 创建Endpoint之后执行已定义的插件
func (that *PluginContainer) afterNewEndpoint(e EarlyEndpoint) {
	var err error
	for _, plugin := range that.plugins {
		if _plugin, ok := plugin.(AfterNewEndpointPlugin); ok {
			if err = _plugin.AfterNewEndpoint(e); err != nil {
				logger.Fatalf("[AfterNewEndpoint:%s] %s", plugin.Name(), err.Error())
				return
			}
		}
	}
}

// BeforeCloseEndpointPlugin 关闭Endpoint之前触发该事件
type BeforeCloseEndpointPlugin interface {
	Plugin
	BeforeCloseEndpoint(Endpoint) error
}

// beforeNewEndpoint 在关闭endpoint之前执行已定义的插件。
func (that *PluginContainer) beforeCloseEndpoint(endpoint Endpoint) (err error) {
	for _, plugin := range that.plugins {
		if _plugin, ok := plugin.(BeforeCloseEndpointPlugin); ok {
			if err = _plugin.BeforeCloseEndpoint(endpoint); err != nil {
				logger.Fatalf("[BeforeCloseEndpoint:%s] %s", plugin.Name(), err.Error())
				return err
			}
		}
	}
	return nil
}

// AfterCloseEndpointPlugin 关闭Endpoint之后触发该事件
type AfterCloseEndpointPlugin interface {
	Plugin
	AfterCloseEndpoint(Endpoint, error) error
}

// afterCloseEndpoint 关闭Endpoint之后执行已定义的插件
func (that *PluginContainer) afterCloseEndpoint(endpoint Endpoint, e error) {
	var err error
	for _, plugin := range that.plugins {
		if _plugin, ok := plugin.(AfterCloseEndpointPlugin); ok {
			if err = _plugin.AfterCloseEndpoint(endpoint, e); err != nil {
				logger.Fatalf("[AfterNewEndpoint:%s] %s", plugin.Name(), err.Error())
				return
			}
		}
	}
}

// AfterRegRouterPlugin 路由注册成功触发该事件
type AfterRegRouterPlugin interface {
	Plugin
	AfterRegRouter(*Handler) error
}

// afterRegRouter 路由注册成功触发该事件
func (that *pluginSingleContainer) afterRegRouter(h *Handler) {
	var err error
	for _, plugin := range that.plugins {
		if _plugin, ok := plugin.(AfterRegRouterPlugin); ok {
			if err = _plugin.AfterRegRouter(h); err != nil {
				logger.Fatalf("[AfterRegRouter:%s] register handler:%s %s, error:%s", plugin.Name(), h.RouterTypeName(), h.Name(), err.Error())
				return
			}
		}
	}
}

// AfterListenPlugin 服务端监听以后触发该事件
type AfterListenPlugin interface {
	Plugin
	AfterListen(net.Addr) error
}

// 该事件在listen之后，accept之前触发
func (that *pluginSingleContainer) afterListen(addr net.Addr) {
	var err error
	for _, plugin := range that.plugins {
		if _plugin, ok := plugin.(AfterListenPlugin); ok {
			if err = _plugin.AfterListen(addr); err != nil {
				logger.Fatalf("[AfterListenPlugin:%s] network:%s, addr:%s, error:%s", plugin.Name(), addr.Network(), addr.String(), err.Error())
				return
			}
		}
	}
	return
}

// BeforeDialPlugin 作为客户端链接到服务端之前调用该事件
type BeforeDialPlugin interface {
	Plugin
	BeforeDial(sess EarlySession, isRedial bool) *Status
}

// 作为客户端角色，链接到远程服务端之前，触发该事件，并返回状态
func (that *pluginSingleContainer) beforeDial(sess EarlySession, isRedial bool) (stat *Status) {
	var pluginName string
	defer func() {
		if p := recover(); p != nil {
			logger.Errorf("[BeforeDialPlugin:%s]  panic:%v\n%s", pluginName, p, status.PanicStackTrace())
			stat = statDialFailed.Copy(p)
		}
	}()
	for _, plugin := range that.plugins {
		if _plugin, ok := plugin.(BeforeDialPlugin); ok {
			pluginName = plugin.Name()
			if stat = _plugin.BeforeDial(sess, isRedial); !stat.OK() {
				logger.Debugf("[BeforeDialPlugin:%s] is_redial:%v, error:%s",
					pluginName, isRedial, stat.String(),
				)
				return stat
			}
		}
	}
	return nil
}

// AfterDialPlugin 作为客户端链接到服务端成功以后触发该事件
type AfterDialPlugin interface {
	Plugin
	AfterDial(sess EarlySession, isRedial bool) *Status
}

// 作为客户端角色，链接到远程服务端成功以后，触发该事件，并返回状态
func (that *pluginSingleContainer) afterDial(sess EarlySession, isRedial bool) (stat *Status) {
	var pluginName string
	defer func() {
		if p := recover(); p != nil {
			logger.Errorf("[AfterDialPlugin:%s] network:%s, addr:%s, panic:%v\n%s", pluginName, sess.RemoteAddr().Network(), sess.RemoteAddr().String(), p, status.PanicStackTrace())
			stat = statDialFailed.Copy(p)
		}
	}()
	for _, plugin := range that.plugins {
		if _plugin, ok := plugin.(AfterDialPlugin); ok {
			pluginName = plugin.Name()
			if stat = _plugin.AfterDial(sess, isRedial); !stat.OK() {
				logger.Debugf("[AfterDialPlugin:%s] network:%s, addr:%s, is_redial:%v, error:%s",
					pluginName, sess.RemoteAddr().Network(), sess.RemoteAddr().String(), isRedial, stat.String(),
				)
				return stat
			}
		}
	}
	return nil
}

// AfterDialFailPlugin 作为客户端链接到服务端失败以后触发该事件
type AfterDialFailPlugin interface {
	Plugin
	AfterDialFail(sess EarlySession, err error, isRedial bool) *Status
}

// 作为客户端角色，链接到远程服务端失败以后，触发该事件，并返回状态
func (that *pluginSingleContainer) afterDialFail(sess EarlySession, err error, isRedial bool) (stat *Status) {
	var pluginName string
	defer func() {
		if p := recover(); p != nil {
			logger.Errorf("[AfterDialFailPlugin:%s] network:%s, addr:%s, panic:%v\n%s", pluginName, sess.RemoteAddr().Network(), sess.RemoteAddr().String(), p, status.PanicStackTrace())
			stat = statDialFailed.Copy(p)
		}
	}()
	for _, plugin := range that.plugins {
		if _plugin, ok := plugin.(AfterDialFailPlugin); ok {
			pluginName = plugin.Name()
			if stat = _plugin.AfterDialFail(sess, err, isRedial); !stat.OK() {
				logger.Debugf("[AfterDialFailPlugin:%s] network:%s, addr:%s, is_redial:%v, error:%s",
					pluginName, sess.RemoteAddr().Network(), sess.RemoteAddr().String(), isRedial, stat.String(),
				)
				return stat
			}
		}
	}
	return nil
}

// AfterAcceptPlugin 作为服务端，接收到客户端的链接后触发该事件
type AfterAcceptPlugin interface {
	Plugin
	AfterAccept(EarlySession) *Status
}

// 接收到accept后，执行该事件，并返回status接收到accept后，执行该事件，并返回status
func (that *pluginSingleContainer) afterAccept(sess EarlySession) (stat *Status) {
	var pluginName string
	defer func() {
		if p := recover(); p != nil {
			logger.Errorf("[AfterAcceptPlugin:%s] network:%s, addr:%s, panic:%v\n%s", pluginName, sess.RemoteAddr().Network(), sess.RemoteAddr().String(), p, status.PanicStackTrace())
			stat = statInternalServerError.Copy(p)
		}
	}()
	for _, plugin := range that.plugins {
		if _plugin, ok := plugin.(AfterAcceptPlugin); ok {
			pluginName = plugin.Name()
			if stat = _plugin.AfterAccept(sess); !stat.OK() {
				logger.Debugf("[PostAcceptPlugin:%s] network:%s, addr:%s, error:%s", pluginName, sess.RemoteAddr().Network(), sess.RemoteAddr().String(), stat.String())
				return stat
			}
		}
	}
	return nil
}

// BeforeWriteCallPlugin 写入CALL消息之前触发该事件
type BeforeWriteCallPlugin interface {
	Plugin
	BeforeWriteCall(WriteCtx) *Status
}

// 写入 CALL 消息之前执行该事件
func (that *pluginSingleContainer) beforeWriteCall(ctx WriteCtx) *Status {
	var stat *Status
	for _, plugin := range that.plugins {
		if _plugin, ok := plugin.(BeforeWriteCallPlugin); ok {
			if stat = _plugin.BeforeWriteCall(ctx); !stat.OK() {
				logger.Debugf("[BeforeWriteCallPlugin:%s] %s", plugin.Name(), stat.String())
				return stat
			}
		}
	}
	return nil
}

// AfterWriteCallPlugin 写入CALL消息成功之后触发该事件
type AfterWriteCallPlugin interface {
	Plugin
	AfterWriteCall(WriteCtx) *Status
}

// 写入CALL消息之后执行该事件
func (that *pluginSingleContainer) afterWriteCall(ctx WriteCtx) *Status {
	var stat *Status
	for _, plugin := range that.plugins {
		if _plugin, ok := plugin.(AfterWriteCallPlugin); ok {
			if stat = _plugin.AfterWriteCall(ctx); !stat.OK() {
				logger.Errorf("[AfterWriteCallPlugin:%s] %s", plugin.Name(), stat.String())
				return stat
			}
		}
	}
	return nil
}

// BeforeWriteReplyPlugin 写入Reply消息之前触发该事件
type BeforeWriteReplyPlugin interface {
	Plugin
	BeforeWriteReply(WriteCtx) *Status
}

// 写入REPLY消息之前执行该事件
func (that *pluginSingleContainer) beforeWriteReply(ctx WriteCtx) {
	var stat *Status
	for _, plugin := range that.plugins {
		if _plugin, ok := plugin.(BeforeWriteReplyPlugin); ok {
			if stat = _plugin.BeforeWriteReply(ctx); !stat.OK() {
				logger.Errorf("[BeforeWriteReplyPlugin:%s] %s", plugin.Name(), stat.String())
				return
			}
		}
	}
}

// AfterWriteReplyPlugin 写入Reply消息成功之后触发该事件
type AfterWriteReplyPlugin interface {
	Plugin
	AfterWriteReply(WriteCtx) *Status
}

// 写入Reply消息之后执行该事件
func (that *pluginSingleContainer) afterWriteReply(ctx WriteCtx) {
	var stat *Status
	for _, plugin := range that.plugins {
		if _plugin, ok := plugin.(AfterWriteReplyPlugin); ok {
			if stat = _plugin.AfterWriteReply(ctx); !stat.OK() {
				logger.Errorf("[AfterWriteReplyPlugin:%s] %s", plugin.Name(), stat.String())
				return
			}
		}
	}
}

// BeforeWritePushPlugin 写入PUSH消息之前触发该事件
type BeforeWritePushPlugin interface {
	Plugin
	BeforeWritePush(WriteCtx) *Status
}

// 写入PUSH消息之前执行该事件
func (that *pluginSingleContainer) beforeWritePush(ctx WriteCtx) *Status {
	var stat *Status
	for _, plugin := range that.plugins {
		if _plugin, ok := plugin.(BeforeWritePushPlugin); ok {
			if stat = _plugin.BeforeWritePush(ctx); !stat.OK() {
				logger.Debugf("[BeforeWritePushPlugin:%s] %s", plugin.Name(), stat.String())
				return stat
			}
		}
	}
	return nil
}

// AfterWritePushPlugin 写入PUSH消息成功之后触发该事件
type AfterWritePushPlugin interface {
	Plugin
	AfterWritePush(WriteCtx) *Status
}

// 写入 PUSH消息之后执行该事件
func (that *pluginSingleContainer) afterWritePush(ctx WriteCtx) *Status {
	var stat *Status
	for _, plugin := range that.plugins {
		if _plugin, ok := plugin.(AfterWritePushPlugin); ok {
			if stat = _plugin.AfterWritePush(ctx); !stat.OK() {
				logger.Errorf("[AfterWritePushPlugin:%s] %s", plugin.Name(), stat.String())
				return stat
			}
		}
	}
	return nil
}

// BeforeReadHeaderPlugin 执行读取Header之前触发该事件
type BeforeReadHeaderPlugin interface {
	Plugin
	BeforeReadHeader(EarlyCtx) error
}

// 读取消息头之前执行该事件
func (that *pluginSingleContainer) beforeReadHeader(ctx EarlyCtx) error {
	var err error
	for _, plugin := range that.plugins {
		if _plugin, ok := plugin.(BeforeReadHeaderPlugin); ok {
			if err = _plugin.BeforeReadHeader(ctx); err != nil {
				logger.Debugf("[BeforeReadHeaderPlugin:%s] disconnected when reading: %s", plugin.Name(), err.Error())
				return err
			}
		}
	}
	return nil
}

// AfterReadCallHeaderPlugin 读取CALL消息的Header之后触发该事件
type AfterReadCallHeaderPlugin interface {
	Plugin
	AfterReadCallHeader(ReadCtx) *Status
}

// 读取CALL消息头之后执行该事件
func (that *pluginSingleContainer) afterReadCallHeader(ctx ReadCtx) *Status {
	var stat *Status
	for _, plugin := range that.plugins {
		if _plugin, ok := plugin.(AfterReadCallHeaderPlugin); ok {
			if stat = _plugin.AfterReadCallHeader(ctx); !stat.OK() {
				logger.Errorf("[AfterReadCallHeaderPlugin:%s] %s", plugin.Name(), stat.String())
				return stat
			}
		}
	}
	return nil
}

// BeforeReadCallBodyPlugin 读取CALL消息的body之前触发该事件
type BeforeReadCallBodyPlugin interface {
	Plugin
	BeforeReadCallBody(ReadCtx) *Status
}

// 读取CALL消息体之前执行该事件
func (that *pluginSingleContainer) beforeReadCallBody(ctx ReadCtx) *Status {
	var stat *Status
	for _, plugin := range that.plugins {
		if _plugin, ok := plugin.(BeforeReadCallBodyPlugin); ok {
			if stat = _plugin.BeforeReadCallBody(ctx); !stat.OK() {
				logger.Errorf("[BeforeReadCallBodyPlugin:%s] %s", plugin.Name(), stat.String())
				return stat
			}
		}
	}
	return nil
}

// AfterReadCallBodyPlugin 读取CALL消息的body之后触发该事件
type AfterReadCallBodyPlugin interface {
	Plugin
	AfterReadCallBody(ReadCtx) *Status
}

// 读取CALL体之后执行该事件
func (that *pluginSingleContainer) afterReadCallBody(ctx ReadCtx) *Status {
	var stat *Status
	for _, plugin := range that.plugins {
		if _plugin, ok := plugin.(AfterReadCallBodyPlugin); ok {
			if stat = _plugin.AfterReadCallBody(ctx); !stat.OK() {
				logger.Errorf("[AfterReadCallBodyPlugin:%s] %s", plugin.Name(), stat.String())
				return stat
			}
		}
	}
	return nil
}

// AfterReadPushHeaderPlugin 读取PUSH消息Header之后触发该事件
type AfterReadPushHeaderPlugin interface {
	Plugin
	AfterReadPushHeader(ReadCtx) *Status
}

// 读取PUSH消息头之后执行该事件
func (that *pluginSingleContainer) afterReadPushHeader(ctx ReadCtx) *Status {
	var stat *Status
	for _, plugin := range that.plugins {
		if _plugin, ok := plugin.(AfterReadPushHeaderPlugin); ok {
			if stat = _plugin.AfterReadPushHeader(ctx); !stat.OK() {
				logger.Errorf("[AfterReadPushHeaderPlugin:%s] %s", plugin.Name(), stat.String())
				return stat
			}
		}
	}
	return nil
}

// BeforeReadPushBodyPlugin 读取PUSH消息body之前触发该事件
type BeforeReadPushBodyPlugin interface {
	Plugin
	BeforeReadPushBody(ReadCtx) *Status
}

//读取PUSH消息体之前执行该事件
func (that *pluginSingleContainer) beforeReadPushBody(ctx ReadCtx) *Status {
	var stat *Status
	for _, plugin := range that.plugins {
		if _plugin, ok := plugin.(BeforeReadPushBodyPlugin); ok {
			if stat = _plugin.BeforeReadPushBody(ctx); !stat.OK() {
				logger.Errorf("[BeforeReadPushBodyPlugin:%s] %s", plugin.Name(), stat.String())
				return stat
			}
		}
	}
	return nil
}

// AfterReadPushBodyPlugin 读取PUSH消息body之后触发该事件
type AfterReadPushBodyPlugin interface {
	Plugin
	AfterReadPushBody(ReadCtx) *Status
}

// 读取PUSH消息体之后执行该事件
func (that *pluginSingleContainer) afterReadPushBody(ctx ReadCtx) *Status {
	var stat *Status
	for _, plugin := range that.plugins {
		if _plugin, ok := plugin.(AfterReadPushBodyPlugin); ok {
			if stat = _plugin.AfterReadPushBody(ctx); !stat.OK() {
				logger.Errorf("[AfterReadPushBodyPlugin:%s] %s", plugin.Name(), stat.String())
				return stat
			}
		}
	}
	return nil
}

// AfterReadReplyHeaderPlugin 读取REPLY消息Header之前触发该事件
type AfterReadReplyHeaderPlugin interface {
	Plugin
	AfterReadReplyHeader(ReadCtx) *Status
}

// 读取Reply消息头之后执行该事件
func (that *pluginSingleContainer) afterReadReplyHeader(ctx ReadCtx) *Status {
	var stat *Status
	for _, plugin := range that.plugins {
		if _plugin, ok := plugin.(AfterReadReplyHeaderPlugin); ok {
			if stat = _plugin.AfterReadReplyHeader(ctx); !stat.OK() {
				logger.Errorf("[AfterReadReplyHeaderPlugin:%s] %s", plugin.Name(), stat.String())
				return stat
			}
		}
	}
	return nil
}

// BeforeReadReplyBodyPlugin 读取REPLY消息body之前触发该事件
type BeforeReadReplyBodyPlugin interface {
	Plugin
	BeforeReadReplyBody(ReadCtx) *Status
}

// 读取Reply消息体之前执行该事件
func (that *pluginSingleContainer) beforeReadReplyBody(ctx ReadCtx) *Status {
	var stat *Status
	for _, plugin := range that.plugins {
		if _plugin, ok := plugin.(BeforeReadReplyBodyPlugin); ok {
			if stat = _plugin.BeforeReadReplyBody(ctx); !stat.OK() {
				logger.Errorf("[BeforeReadReplyBodyPlugin:%s] %s", plugin.Name(), stat.String())
				return stat
			}
		}
	}
	return nil
}

// AfterReadReplyBodyPlugin 读取REPLY消息body之后触发该事件
type AfterReadReplyBodyPlugin interface {
	Plugin
	AfterReadReplyBody(ReadCtx) *Status
}

// 读取Reply消息体之后执行该事件
func (that *pluginSingleContainer) afterReadReplyBody(ctx ReadCtx) *Status {
	var stat *Status
	for _, plugin := range that.plugins {
		if _plugin, ok := plugin.(AfterReadReplyBodyPlugin); ok {
			if stat = _plugin.AfterReadReplyBody(ctx); !stat.OK() {
				logger.Errorf("[AfterReadReplyBodyPlugin:%s] %s", plugin.Name(), stat.String())
				return stat
			}
		}
	}
	return nil
}

// AfterDisconnectPlugin 断开会话以后触发该事件
type AfterDisconnectPlugin interface {
	Plugin
	AfterDisconnect(BaseSession) *Status
}

// 会话关闭以后执行该事件
func (that *pluginSingleContainer) afterDisconnect(sess BaseSession) *Status {
	var stat *Status
	for _, plugin := range that.plugins {
		if _plugin, ok := plugin.(AfterDisconnectPlugin); ok {
			if stat = _plugin.AfterDisconnect(sess); !stat.OK() {
				logger.Errorf("[AfterDisconnectPlugin:%s] %s", plugin.Name(), stat.String())
				return stat
			}
		}
	}
	return nil
}

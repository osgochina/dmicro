package auth

import (
	"fmt"
	"github.com/gogf/gf/container/gmap"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/message"
	"net"
	"sync/atomic"
)

type (

	// Bearer 生成认证票据，并发送到对应的端
	Bearer func(sess Session, fn SendOnce) *drpc.Status
	// SendOnce 发送到对应端
	SendOnce func(info, retReCv interface{}) *drpc.Status

	// Checker 检票方法，检查对端传入的票据是否合法
	Checker func(sess Session, fn ReCvOnce) (ret interface{}, stat *drpc.Status)
	// ReCvOnce 检票
	ReCvOnce func(infoReCv interface{}) *drpc.Status

	// Session 在权限认证组件中，会话只暴露以下方法
	Session interface {
		Endpoint() drpc.Endpoint //会话所在endpoint
		SetID(newID string)      // 设置会话id
		LocalAddr() net.Addr     // 当前会话监听地址端口
		RemoteAddr() net.Addr    // 远端地址端口
		Swap() *gmap.Map         // 会话的缓冲区
	}
)

// NewBearerPlugin 创建票据生成插件
func NewBearerPlugin(fn Bearer, setting ...message.MsgSetting) drpc.Plugin {
	return &authBearerPlugin{
		bearerFunc: fn,
		msgSetting: setting,
	}
}

// NewCheckerPlugin 创建票据检查插件
func NewCheckerPlugin(fn Checker, setting ...message.MsgSetting) drpc.Plugin {
	return &authCheckerPlugin{
		checkerFunc: fn,
		msgSetting:  setting,
	}
}

//生成票据，组件
type authBearerPlugin struct {
	bearerFunc Bearer
	msgSetting []message.MsgSetting
}

// 票据检查组件
type authCheckerPlugin struct {
	checkerFunc Checker
	msgSetting  []message.MsgSetting
}

var (
	_ drpc.AfterDialPlugin   = new(authBearerPlugin)
	_ drpc.AfterAcceptPlugin = new(authCheckerPlugin)
)

func (that *authBearerPlugin) Name() string {
	return "auth-bearer"
}

func (that *authCheckerPlugin) Name() string {
	return "auth-checker"
}

// MultiSendErr authBearerPlugin组件并发发送限制
var MultiSendErr = drpc.NewStatus(
	drpc.CodeWriteFailed,
	"auth-bearer plugin usage is incorrect",
	"multiple call SendOnce function",
)

// MultiReCvErr authCheckerPlugin 并发认证限制
var MultiReCvErr = drpc.NewStatus(
	drpc.CodeInternalServerError,
	"auth-checker plugin usage is incorrect",
	"multiple call ReCvOnce function",
)

// AfterDial 客户端链接到服务端成功后触发该调用
func (that *authBearerPlugin) AfterDial(sess drpc.EarlySession, _ bool) *drpc.Status {
	if that.bearerFunc == nil {
		return nil
	}
	var called int32
	return that.bearerFunc(sess, func(info, retReCv interface{}) *drpc.Status {
		if !atomic.CompareAndSwapInt32(&called, 0, 1) {
			return MultiSendErr
		}
		//往服务端发送认证包
		stat := sess.EarlySend(drpc.TypeAuthCall, "", info, nil, that.msgSetting...)
		if !stat.OK() {
			return stat
		}
		//等待接收服务端的回包，并把回包的值放到retReCv这个变量的地址上
		retMsg := sess.EarlyReceive(func(header message.Header) interface{} {
			if header.MType() != drpc.TypeAuthReply {
				return nil
			}
			return retReCv
		})
		//回包的状态是否成功
		if !retMsg.StatusOK() {
			return retMsg.Status()
		}
		//如果回包的消息类型不是 TypeAuthReply
		if retMsg.MType() != drpc.TypeAuthReply {
			return drpc.NewStatus(
				drpc.CodeUnauthorized,
				drpc.CodeText(drpc.CodeUnauthorized),
				fmt.Sprintf("auth message(1st) expect: AUTH_REPLY, but received: %s",
					drpc.TypeText(retMsg.MType())),
			)
		}
		return nil
	})
}

func (that *authCheckerPlugin) AfterAccept(sess drpc.EarlySession) *drpc.Status {
	if that.checkerFunc == nil {
		return nil
	}
	var called int32
	ret, stat := that.checkerFunc(sess, func(infoReCv interface{}) *drpc.Status {
		//并发限制
		if !atomic.CompareAndSwapInt32(&called, 0, 1) {
			return MultiReCvErr
		}
		//获取客户端端发送过来的票据
		infoMsg := sess.EarlyReceive(func(header message.Header) interface{} {
			//如果消息类型不是权限认证，则跳过，如果是，则把票据内容写入infoReCv对应的地址
			if header.MType() != drpc.TypeAuthCall {
				return nil
			}
			return infoReCv
		})
		if !infoMsg.StatusOK() {
			return infoMsg.Status()
		}
		//如果消息类型不是权限认证
		if infoMsg.MType() != drpc.TypeAuthCall {
			return drpc.NewStatus(
				drpc.CodeUnauthorized,
				drpc.CodeText(drpc.CodeUnauthorized),
				fmt.Sprintf("auth message(1st) expect: AUTH_CALL, but received: %s",
					drpc.TypeText(infoMsg.MType())),
			)
		}
		return nil
	})
	//如果错误是并发限制，则还是会给客户端发送认证成功的消息，但是返回值是空
	if stat == MultiReCvErr {
		sess.EarlySend(drpc.TypeAuthReply, "", nil, stat, that.msgSetting...)
		return stat
	}
	//发送认证结果给客户端
	stat2 := sess.EarlySend(drpc.TypeAuthReply, "", ret, stat, that.msgSetting...)
	if !stat2.OK() {
		return stat2
	}
	return stat
}

package heartbeat

import (
	"github.com/gogf/gf/util/gconv"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/internal"
	"strconv"
	"time"
)

// Pong 心跳响应
type Pong interface {
	// Name returns name.
	Name() string
	// AfterNewEndpoint runs ping woker.
	AfterNewEndpoint(endpoint drpc.EarlyEndpoint) error
	// AfterWriteCall updates heartbeat information.
	AfterWriteCall(ctx drpc.WriteCtx) *drpc.Status
	// AfterWritePush updates heartbeat information.
	AfterWritePush(ctx drpc.WriteCtx) *drpc.Status
	// AfterReadCallHeader updates heartbeat information.
	AfterReadCallHeader(ctx drpc.ReadCtx) *drpc.Status
	// AfterReadPushHeader updates heartbeat information.
	AfterReadPushHeader(ctx drpc.ReadCtx) *drpc.Status
}

var (
	_ drpc.AfterNewEndpointPlugin    = Pong(nil)
	_ drpc.AfterWriteCallPlugin      = Pong(nil)
	_ drpc.AfterWritePushPlugin      = Pong(nil)
	_ drpc.AfterReadCallHeaderPlugin = Pong(nil)
	_ drpc.AfterReadPushHeaderPlugin = Pong(nil)
)

func NewPong() Pong {
	return new(heartPong)
}

type heartPong struct{}

func (that *heartPong) Name() string {
	return "heart-pong"
}

// AfterNewEndpoint 端点启动的时候，开启一个协程，遍历会话，检查超时时间
func (that *heartPong) AfterNewEndpoint(endpoint drpc.EarlyEndpoint) error {
	//为心跳注册处理方法
	endpoint.RouteCallFunc((*pongCall).heartbeat)
	endpoint.RoutePushFunc((*pongPush).heartbeat)

	rangeSession := endpoint.RangeSession

	const initial = time.Second*minRateSecond - 1
	interval := initial
	go func() {
		for {
			time.Sleep(interval)
			rangeSession(func(sess drpc.Session) bool {
				info, ok := getHeartbeatInfo(sess.Swap())
				if !ok {
					return true
				}
				cp := info.elemCopy()
				//检查会话是否健康，并且最后心跳时间是否超时
				if sess.Health() && cp.last.Add(cp.rate*2).Before(time.Now()) {
					_ = sess.Close()
				}
				// 时间间隔使用ping端发起传入的时间
				if cp.rate < interval || interval == initial {
					interval = cp.rate
				}
				return true
			})
		}
	}()
	return nil
}

// AfterWriteCall 写入CALL消息后执行该事件
func (that *heartPong) AfterWriteCall(ctx drpc.WriteCtx) *drpc.Status {
	return that.AfterWritePush(ctx)
}

// AfterWritePush 写入PUSH消息后执行该事件
func (that *heartPong) AfterWritePush(ctx drpc.WriteCtx) *drpc.Status {
	sess := ctx.Session()
	if !sess.Health() {
		return nil
	}
	updateHeartbeatInfo(sess.Swap(), 0)
	return nil
}

// AfterReadCallHeader 读取CALL消息的头后执行该事件
func (that *heartPong) AfterReadCallHeader(ctx drpc.ReadCtx) *drpc.Status {
	that.update(ctx)
	return nil
}

// AfterReadPushHeader 读取PUSH消息的头后执行该事件
func (that *heartPong) AfterReadPushHeader(ctx drpc.ReadCtx) *drpc.Status {
	return that.AfterReadCallHeader(ctx)
}

// 更新心跳元数据
func (that *heartPong) update(ctx drpc.ReadCtx) {
	if ctx.ServiceMethod() == heartbeatServiceMethod {
		return
	}
	sess := ctx.Session()
	if !sess.Health() {
		return
	}
	updateHeartbeatInfo(sess.Swap(), 0)
}

// 注册call方法的pong
type pongCall struct {
	drpc.CallCtx
}

func (that *pongCall) heartbeat(_ *struct{}) (*struct{}, *drpc.Status) {
	return nil, handheldHeartbeat(that.Session(), that.PeekMeta)
}

//注册push方法的pong
type pongPush struct {
	drpc.PushCtx
}

func (that *pongPush) heartbeat(_ *struct{}) *drpc.Status {
	return handheldHeartbeat(that.Session(), that.PeekMeta)
}

//处理心跳
func handheldHeartbeat(sess drpc.CtxSession, peekMeta func(string) interface{}) *drpc.Status {
	rateStr := gconv.String(peekMeta(heartbeatMetaKey))
	rateSecond := parseHeartbeatRateSecond(rateStr)
	isFirst := updateHeartbeatInfo(sess.Swap(), time.Second*time.Duration(rateSecond))
	if isFirst && rateSecond == -1 {
		return drpc.NewStatus(drpc.CodeBadMessage, "invalid heartbeat rate", rateStr)
	}
	if rateSecond == 0 {
		internal.Infof("heart-pong: %s", sess.ID())
	} else {
		internal.Infof("heart-pong: %s, set rate: %ds", sess.ID(), rateSecond)
	}
	return nil
}

func parseHeartbeatRateSecond(s string) int {
	if len(s) == 0 {
		return 0
	}
	r, err := strconv.Atoi(s)
	if err != nil || r < 0 {
		return -1
	}
	return r
}

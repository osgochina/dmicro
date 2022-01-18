package heartbeat

import (
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/internal"
	"github.com/osgochina/dmicro/drpc/message"
	"github.com/osgochina/dmicro/utils/dgpool"
	"strconv"
	"sync"
	"time"
)

const (
	heartbeatServiceMethod = "/heartbeat"
	heartbeatMetaKey       = "hb_"
)

// Ping  插件接口
type Ping interface {
	// SetRate  设置心跳的频率
	SetRate(rateSecond int)
	// UseCall 使用CALL方法执行ping命令
	UseCall()
	// UsePush 使用PUSH方法执行ping命令
	UsePush()
	// Name 插件的名字
	Name() string
	// AfterNewEndpoint runs ping woker.
	AfterNewEndpoint(peer drpc.EarlyEndpoint) error
	// AfterDial initializes heartbeat information.
	AfterDial(sess drpc.EarlySession, isRedial bool) *drpc.Status
	// AfterAccept initializes heartbeat information.
	AfterAccept(sess drpc.EarlySession) *drpc.Status
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
	_ drpc.AfterNewEndpointPlugin    = Ping(nil)
	_ drpc.AfterDialPlugin           = Ping(nil)
	_ drpc.AfterAcceptPlugin         = Ping(nil)
	_ drpc.AfterWriteCallPlugin      = Ping(nil)
	_ drpc.AfterWritePushPlugin      = Ping(nil)
	_ drpc.AfterReadCallHeaderPlugin = Ping(nil)
	_ drpc.AfterReadPushHeaderPlugin = Ping(nil)
)

// NewPing 创建插件
func NewPing(rateSecond int, useCall bool) Ping {
	p := new(heartPing)
	p.useCall = useCall
	p.SetRate(rateSecond)
	return p
}

type heartPing struct {
	//端点对象
	endpoint drpc.Endpoint
	//ping的频率
	pingRate time.Duration
	//锁
	mu   sync.RWMutex
	once sync.Once

	//ping的频率秒数，需要作为参数传给对端
	pingRateSecond string
	//是否使用响应是心跳，true使用响应式，false使用push
	useCall bool
}

// SetRate 设置心跳频率
func (that *heartPing) SetRate(rateSecond int) {
	if rateSecond < minRateSecond {
		rateSecond = minRateSecond
	}
	that.mu.Lock()
	that.pingRate = time.Second * time.Duration(rateSecond)
	that.pingRateSecond = strconv.Itoa(rateSecond)
	that.mu.Unlock()
	internal.Infof("set heartbeat rate: %ds", rateSecond)
}

//获取心跳频率
func (that *heartPing) getRate() time.Duration {
	that.mu.RLock()
	defer that.mu.RUnlock()
	return that.pingRate
}

// 获取心跳频率具体数字
func (that *heartPing) getPingRateSecond() string {
	that.mu.RLock()
	defer that.mu.RUnlock()
	return that.pingRateSecond
}

// UseCall 使用响应式心跳，有问有答
func (that *heartPing) UseCall() {
	that.mu.Lock()
	that.useCall = true
	that.mu.Unlock()
}

// UsePush 使用通知式心跳，PUSH方法发送
func (that *heartPing) UsePush() {
	that.mu.Lock()
	that.useCall = false
	that.mu.Unlock()
}

//是否使用响应式心跳
func (that *heartPing) isCall() bool {
	that.mu.RLock()
	defer that.mu.RUnlock()
	return that.useCall
}

// Name 插件名字
func (that *heartPing) Name() string {
	return "heart-ping"
}

// AfterNewEndpoint 端点服务启动后，开启一个单独的协程来遍历会话，并发送心跳
func (that *heartPing) AfterNewEndpoint(endpoint drpc.EarlyEndpoint) error {
	rangeSession := endpoint.RangeSession

	go func() {
		var isCall bool
		for {
			time.Sleep(that.getRate())
			isCall = that.isCall()
			rangeSession(func(sess drpc.Session) bool {
				if !sess.Health() {
					return true
				}
				info, ok := getHeartbeatInfo(sess.Swap())
				if !ok {
					return true
				}
				//判断节点的最后心跳时间+心跳频率是否大于当前时间
				cp := info.elemCopy()
				if cp.last.Add(cp.rate).After(time.Now()) {
					return true
				}
				if isCall {
					that.goCall(sess)
				} else {
					that.goPush(sess)
				}
				return true
			})
		}
	}()
	return nil
}

// AfterDial 链接到远程服务端成功以后，触发该事件，并返回状态
func (that *heartPing) AfterDial(sess drpc.EarlySession, _ bool) *drpc.Status {
	return that.AfterAccept(sess)
}

// AfterAccept 接收到accept后，执行该事件，并返回status
func (that *heartPing) AfterAccept(sess drpc.EarlySession) *drpc.Status {
	rate := that.getRate()
	initHeartbeatInfo(sess.Swap(), rate)
	return nil
}

// AfterWriteCall 写入CALL消息之后执行该事件
func (that *heartPing) AfterWriteCall(ctx drpc.WriteCtx) *drpc.Status {
	return that.AfterWritePush(ctx)
}

// AfterWritePush 写入PUSH消息之后执行该事件
func (that *heartPing) AfterWritePush(ctx drpc.WriteCtx) *drpc.Status {
	that.update(ctx)
	return nil
}

// AfterReadCallHeader 读取CALL消息头后触发该事件
func (that *heartPing) AfterReadCallHeader(ctx drpc.ReadCtx) *drpc.Status {
	return that.AfterReadPushHeader(ctx)
}

// AfterReadPushHeader 读取PUSH消息头后触发该事件
func (that *heartPing) AfterReadPushHeader(ctx drpc.ReadCtx) *drpc.Status {
	that.update(ctx)
	return nil
}

//发送call类型的心跳
func (that *heartPing) goCall(sess drpc.Session) {
	dgpool.FILOGo(func() {
		//发送call方法，如果没有发送成功，则关闭会话
		stat := sess.Call(heartbeatServiceMethod, nil, nil, message.WithSetMeta(heartbeatMetaKey, that.getPingRateSecond())).Status()
		if stat != nil {
			err := sess.Close()
			if err != nil {
				return
			}
		}
	})
}

//发送push类型的心跳
func (that *heartPing) goPush(sess drpc.Session) {
	dgpool.FILOGo(func() {
		//发送push类型的心跳，如果没有发送成功，则关闭会话
		stat := sess.Push(heartbeatServiceMethod, nil, message.WithSetMeta(heartbeatMetaKey, that.getPingRateSecond()))
		if stat != nil {
			err := sess.Close()
			if err != nil {
				return
			}
		}
	})
}

//更新心跳源信息
func (that *heartPing) update(ctx drpc.EarlyCtx) {
	sess := ctx.Session()
	if !sess.Health() {
		return
	}
	updateHeartbeatInfo(sess.Swap(), that.getRate())
}

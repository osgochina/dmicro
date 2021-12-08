package drpc

import (
	"crypto/tls"
	"errors"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/os/grpool"
	"github.com/gogf/gf/util/grand"
	"github.com/osgochina/dmicro/drpc/codec"
	"github.com/osgochina/dmicro/drpc/proto"
	"github.com/osgochina/dmicro/drpc/socket"
	"github.com/osgochina/dmicro/drpc/status"
	"github.com/osgochina/dmicro/eventbus"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/utils/dgpool"
	errors2 "github.com/osgochina/dmicro/utils/errors"
	"github.com/osgochina/dmicro/utils/graceful"
	"github.com/osgochina/dmicro/utils/inherit"
	"net"
	"sync"
	"time"
)

type BaseEndpoint interface {

	// Close 关闭该端点
	Close() (err error)

	// CountSession 统计该端点上的session数量
	CountSession() int

	// GetSession 获取指定ID的 session
	GetSession(sessionID string) (Session, bool)

	// RangeSession 循环迭代session
	RangeSession(fn func(sess Session) bool)

	// SetTLSConfig 设置证书配置
	SetTLSConfig(tlsConfig *tls.Config)

	// SetTLSConfigFromFile 从文件中读取证书并设置证书配置
	SetTLSConfigFromFile(tlsCertFile, tlsKeyFile string, insecureSkipVerifyForClient ...bool) error

	// TLSConfig tls配置对象
	TLSConfig() *tls.Config

	// PluginContainer 插件容器对象
	PluginContainer() *PluginContainer

	// EventBus 消息总线
	EventBus() *eventbus.EventBus
}

type EarlyEndpoint interface {
	BaseEndpoint

	// Router 获取路由对象
	Router() *Router

	// SubRoute 获取分组路由对象
	SubRoute(pathPrefix string, plugin ...Plugin) *SubRouter

	// RouteCall 通过struct注册CALL类型的处理程序，并且返回注册的路径列表
	RouteCall(ctrlStruct interface{}, plugin ...Plugin) []string
	// RouteCallFunc 通过func注册CALL类型的处理程序，并且返回单个注册路径
	RouteCallFunc(callHandleFunc interface{}, plugin ...Plugin) string
	// RoutePush 通过struct注册PUSH类型的处理程序，并且返回注册的路径列表
	RoutePush(ctrlStruct interface{}, plugin ...Plugin) []string
	// RoutePushFunc 通过func注册PUSH类型的处理程序，并且返回单个注册路径
	RoutePushFunc(pushHandleFunc interface{}, plugin ...Plugin) string
	// SetUnknownCall 设置默认处理程序，当没有找到CALL的处理程序时将调用该处理程序。
	SetUnknownCall(fn func(UnknownCallCtx) (interface{}, *status.Status), plugin ...Plugin)
	// SetUnknownPush 设置默认处理程序，当没有找到PUSH的处理程序时将调用该处理程序。
	SetUnknownPush(fn func(UnknownPushCtx) *status.Status, plugin ...Plugin)
}

type Endpoint interface {
	EarlyEndpoint

	// ListenAndServe 打开服务监听
	ListenAndServe(protoFunc ...socket.ProtoFunc) error

	// Dial 作为客户端链接到指定的服务
	Dial(addr string, protoFunc ...socket.ProtoFunc) (Session, *status.Status)

	// ServeConn 传入指定的conn，生成session
	// 提示：
	// 1. 不支持断开链接后自动重拨
	// 2. 不检查TLS
	// 3. 执行 PostAcceptPlugin 插件
	ServeConn(conn net.Conn, protoFunc ...socket.ProtoFunc) (Session, *status.Status)
}

var (
	_ BaseEndpoint  = new(endpoint)
	_ EarlyEndpoint = new(endpoint)
	_ Endpoint      = new(endpoint)
)

type endpoint struct {
	router            *Router
	pluginContainer   *PluginContainer
	sessHub           *SessionHub
	eventbus          *eventbus.EventBus
	closeCh           chan struct{}
	defaultSessionAge time.Duration
	defaultContextAge time.Duration
	tlsConfig         *tls.Config
	slowCometDuration time.Duration
	timeNow           func() int64
	mu                sync.Mutex
	network           string
	defaultBodyCodec  byte
	printDetail       bool
	countTime         bool

	//只有作为server角色时候才有该对象
	listerAddr net.Addr
	listeners  map[net.Listener]struct{}

	//只有作为client角色时候才有该对象
	dialer *Dialer
}

func NewEndpoint(cfg EndpointConfig, globalLeftPlugin ...Plugin) Endpoint {

	//创建插件容器，并把全局插件加入到容器的最左，方便后续添加插件保存执行顺序
	pluginContainer := newPluginContainer()
	pluginContainer.AppendLeft(globalLeftPlugin...)
	//触发事件
	pluginContainer.beforeNewEndpoint(&cfg)

	// 检查配置项是否正确
	if err := cfg.check(); err != nil {
		logger.Fatalf("%v", err)
	}

	var e = &endpoint{
		router:            newRouter(pluginContainer),
		pluginContainer:   pluginContainer,
		sessHub:           newSessionHub(),
		eventbus:          eventbus.New(grand.S(8)),
		defaultSessionAge: cfg.DefaultSessionAge,
		defaultContextAge: cfg.DefaultContextAge,
		closeCh:           make(chan struct{}),
		slowCometDuration: cfg.slowCometDuration,
		network:           cfg.Network,
		listerAddr:        cfg.listenAddr,
		printDetail:       cfg.PrintDetail,
		countTime:         cfg.CountTime,
		listeners:         make(map[net.Listener]struct{}),
		dialer: &Dialer{
			network:        cfg.Network,
			dialTimeout:    cfg.DialTimeout,
			localAddr:      cfg.localAddr,
			redialInterval: cfg.RedialInterval,
			redialTimes:    cfg.RedialTimes,
		},
	}
	//默认的消息体编码格式
	if c, err := codec.GetByName(cfg.DefaultBodyCodec); err != nil {
		logger.Fatalf("%v", err)
	} else {
		e.defaultBodyCodec = c.ID()
	}
	//是否统计时间
	if e.countTime {
		e.timeNow = func() int64 {
			return time.Now().UnixNano()
		}
	} else {
		e.timeNow = func() int64 {
			return 0
		}
	}

	//addEndpoint(e)
	// 平滑重启添加endpoint
	graceful.Graceful().AddEndpoint(e)
	//触发事件
	e.pluginContainer.afterNewEndpoint(e)
	return e
}

// PluginContainer 获取端点的插件容器
func (that *endpoint) PluginContainer() *PluginContainer {
	return that.pluginContainer
}
func (that *endpoint) EventBus() *eventbus.EventBus {
	return that.eventbus
}

// TLSConfig 获取该端点的证书信息
func (that *endpoint) TLSConfig() *tls.Config {
	return that.tlsConfig
}

// SetTLSConfig 设置该端点的证书信息
func (that *endpoint) SetTLSConfig(tlsConfig *tls.Config) {
	that.tlsConfig = tlsConfig
	that.dialer.tlsConfig = tlsConfig
}

// SetTLSConfigFromFile 通过文件生成端点的证书信息
func (that *endpoint) SetTLSConfigFromFile(tlsCertFile, tlsKeyFile string, insecureSkipVerifyForClient ...bool) error {
	tlsConfig, err := inherit.NewTLSConfigFromFile(tlsCertFile, tlsKeyFile, insecureSkipVerifyForClient...)
	if err == nil {
		that.SetTLSConfig(tlsConfig)
	}
	return err
}

// GetSession 获取session
func (that *endpoint) GetSession(sessionID string) (Session, bool) {
	return that.sessHub.get(sessionID)
}

// RangeSession 遍历session
func (that *endpoint) RangeSession(fn func(sess Session) bool) {
	that.sessHub.sessions.Range(func(_, value interface{}) bool {
		return fn(value.(*session))
	})
}

// CountSession 统计当前session个数
func (that *endpoint) CountSession() int {
	return that.sessHub.len()
}

//从池子中获取会话上下文对象
func (that *endpoint) getHandleCtx(s *session, withWg bool) *handlerCtx {
	if withWg {
		// 优雅控制器增加1
		s.graceCtxWaitGroup.Add(1)
	}
	ctx := handlerCtxPool.Get().(*handlerCtx)
	ctx.clean()
	ctx.reInit(s)
	return ctx
}

//归还会话上下文对象到池子
func (that *endpoint) putHandleCtx(ctx *handlerCtx, withWg bool) {
	if withWg {
		// 处理成功，则优雅控制器done
		ctx.sess.graceCtxWaitGroup.Done()
	}
	handlerCtxPool.Put(ctx)
}

// Router 返回路由对象
func (that *endpoint) Router() *Router {
	return that.router
}

// SubRoute 设置路由分组
func (that *endpoint) SubRoute(pathPrefix string, plugin ...Plugin) *SubRouter {
	return that.router.SubRoute(pathPrefix, plugin...)
}

// RouteCall 通过结构体对象注册CALL命令路由
func (that *endpoint) RouteCall(callCtrlStruct interface{}, plugin ...Plugin) []string {
	return that.router.RouteCall(callCtrlStruct, plugin...)
}

// RouteCallFunc 通过对象的方法注册CALL命令路由
func (that *endpoint) RouteCallFunc(callHandleFunc interface{}, plugin ...Plugin) string {
	return that.router.RouteCallFunc(callHandleFunc, plugin...)
}

// RoutePush 通过结构体对象注册PUSH命令的路由
func (that *endpoint) RoutePush(pushCtrlStruct interface{}, plugin ...Plugin) []string {
	return that.router.RoutePush(pushCtrlStruct, plugin...)
}

// RoutePushFunc 通过对象的方法注册PUSH命令的路由
func (that *endpoint) RoutePushFunc(pushHandleFunc interface{}, plugin ...Plugin) string {
	return that.router.RoutePushFunc(pushHandleFunc, plugin...)
}

// SetUnknownCall 设置CALL命令的默认路由
func (that *endpoint) SetUnknownCall(fn func(UnknownCallCtx) (interface{}, *Status), plugin ...Plugin) {
	that.router.SetUnknownCall(fn, plugin...)
}

// SetUnknownPush 设置PUSH命令的默认路由
func (that *endpoint) SetUnknownPush(fn func(UnknownPushCtx) *Status, plugin ...Plugin) {
	that.router.SetUnknownPush(fn, plugin...)
}

//通过注册的路由地址，返回CALL处理方法
func (that *endpoint) getCallHandler(uriPath string) (*Handler, bool) {
	return that.router.subRouter.getCall(uriPath)
}

//通过注册的路由地址，返回PUSH处理方法
func (that *endpoint) getPushHandler(uriPath string) (*Handler, bool) {
	return that.router.subRouter.getPush(uriPath)
}

// Dial 拨号链接远端
func (that *endpoint) Dial(addr string, protoFunc ...proto.ProtoFunc) (Session, *Status) {

	var sess = newSession(that, nil, protoFunc)
	//连接到服务端之前，触发事件
	stat := that.pluginContainer.beforeDial(addr, false)
	if !stat.OK() {
		return nil, stat
	}
	//链接远端服务器，链接如果不成功，会重试
	_, err := that.dialer.dialWithRetry(addr, "", func(conn net.Conn) error {
		sess.socket.Reset(conn, protoFunc...)
		sess.socket.SetID(sess.LocalAddr().String())
		// 链接成功之后触发事件
		stat = that.pluginContainer.afterDial(sess, false)
		if !stat.OK() {
			_ = conn.Close()
			return stat.Cause()
		}
		return nil
	})
	if err != nil {
		//链接失败触发事件
		_ = that.pluginContainer.afterDialFail(sess, err, false)
		return nil, statDialFailed.Copy(err)
	}

	//如果重试次数不为0，则设置重试方法
	if that.dialer.RedialTimes() != 0 {
		sess.redialForClientLocked = func() bool {
			//获取链接的原始信息
			oldID := sess.ID()
			oldIP := sess.LocalAddr().String()
			oldConn := sess.getConn()
			//连接到服务端之前，触发事件
			_ = that.pluginContainer.beforeDial(addr, true)
			//重新链接服务器端
			_, err = that.dialer.dialWithRetry(addr, oldID, func(conn net.Conn) error {
				sess.socket.Reset(conn, protoFunc...)
				//如果原始的session id是使用的本地地址作为id值，则继续使用本地地址作为id值，
				//如果是使用的其他自定义的id值，则使用自定义的id值
				if oldID == oldIP {
					sess.socket.SetID(sess.LocalAddr().String())
				} else {
					sess.SetID(oldID)
				}
				//更改会话状态为初始状态
				sess.changeStatus(statusPreparing)
				//执行事件
				stat = that.pluginContainer.afterDial(sess, true)
				//如果执行事件返回的状态不是ok，则把当前链接关闭，把状态修改为重试中，继续重试。
				if !stat.OK() {
					_ = conn.Close()
					sess.changeStatus(statusRedialing)
					return stat.Cause()
				}
				return nil
			})
			//如果到了最大重试次数还没有链接成功，则把会话关闭，
			if err != nil {
				//链接失败触发事件
				_ = that.pluginContainer.afterDialFail(sess, err, true)
				_ = sess.closeLocked()
				//防止状态没有修改成功，再次尝试修改状态
				sess.tryChangeStatus(statusRedialFailed, statusRedialing)
				logger.Errorf("redial fail (network:%s, addr:%s, id:%s): %s", that.network, addr, oldID, err.Error())
				return false
			}
			//原始链接如果存在，则关闭
			if oldConn != nil {
				_ = oldConn.Close()
			}
			//修改会话状态为就绪，并且执行会话消息读取监听
			sess.changeStatus(statusOk)
			err = grpool.Add(sess.startReadAndHandle)
			if err != nil {
				logger.Errorf("redial fail (network:%s, addr:%s, id:%s): %s", that.network, addr, oldID, err.Error())
				return false
			}
			//把当前会话加入会话池
			that.sessHub.set(sess)
			logger.Infof("redial ok (network:%s, addr:%s, id:%s)", that.network, addr, sess.ID())
			return true
		}
	}
	logger.Infof("dial ok (network:%s, addr:%s, id:%s)", that.network, addr, sess.ID())
	//修改会话状态，并且启动响应监听
	sess.changeStatus(statusOk)
	err = grpool.Add(sess.startReadAndHandle)
	if err != nil {
		return nil, statDialFailed.Copy(err)
	}
	//把当前会话加入会话池
	that.sessHub.set(sess)
	return sess, nil
}

// ServeConn 通过提供的链接，生成会话并返回
// 1. 不支持自动重连
// 2. 不检查是否是 TLS链接
// 3. 会执行AfterAcceptPlugin 事件
func (that *endpoint) ServeConn(conn net.Conn, protoFunc ...proto.ProtoFunc) (Session, *Status) {
	network := conn.LocalAddr().Network()
	//if asQUIC(network) != "" {
	//	if _, ok := conn.(*quic.Conn); !ok {
	//		return nil, NewStatus(CodeWrongConn, "not support "+network, "network must be one of the following: tcp, tcp4, tcp6, unix, unixpacket, kcp or quic")
	//	}
	//	network = "quic"
	//} else if asKCP(network) != "" {
	//	if _, ok := conn.(*kcp.UDPSession); !ok {
	//		return nil, NewStatus(CodeWrongConn, "not support "+network, "network must be one of the following: tcp, tcp4, tcp6, unix, unixpacket, kcp or quic")
	//	}
	//	network = "kcp"
	//}

	var sess = newSession(that, conn, protoFunc)

	stat := that.pluginContainer.afterAccept(sess)
	if !stat.OK() {
		_ = sess.Close()
		return nil, stat
	}
	logger.Infof("serve ok (network:%s, addr:%s, id:%s)", network, sess.RemoteAddr().String(), sess.ID())
	sess.changeStatus(statusOk)
	err := grpool.Add(sess.startReadAndHandle)
	if err != nil {
		return nil, statUnknownError.Copy(err)
	}
	//把当前会话加入会话池
	that.sessHub.set(sess)
	return sess, nil
}

var ErrListenClosed = errors.New("listener is closed")

// 启动服务，并侦听指定地址
func (that *endpoint) serveListener(lis net.Listener, protoFunc ...proto.ProtoFunc) error {
	defer func() {
		_ = lis.Close()
	}()
	that.listeners[lis] = struct{}{}

	network := lis.Addr().Network()
	//switch lis.(type) {
	//case *quic.Listener:
	//	network = "quic"
	//case *kcp.Listener:
	//	network = "kcp"
	//}

	addr := lis.Addr().String()
	logger.Printf("启动监听并提供服务：(network:%s, addr:%s)", network, addr)
	that.pluginContainer.afterListen(lis.Addr())

	var tempDelay time.Duration
	var closeCh = that.closeCh

	for {
		conn, err := lis.Accept()
		if err != nil {
			//如果当前端点已关闭，则停止并返回
			select {
			case <-closeCh:
				return ErrListenClosed
			default:
			}
			//如果错误是网络错误，并且该错误是暂时的，则等待一段时间后继续
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				//最多1秒
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}

				logger.Infof("accept error: %s; retrying in %v", err.Error(), tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return err
		}
		//重新开始计算 暂时网络错误的等待时间
		tempDelay = 0

		dgpool.FILOAnywayGo(func() {
			//如果链接是tls加密链接，则设置超时时间，并进行握手
			if c, ok := conn.(*tls.Conn); ok {
				if that.defaultSessionAge > 0 {
					_ = c.SetReadDeadline(time.Now().Add(that.defaultSessionAge))
				}
				if that.defaultContextAge > 0 {
					_ = c.SetReadDeadline(time.Now().Add(that.defaultContextAge))
				}
				if err = c.Handshake(); err != nil {
					logger.Errorf("TLS handshake error from %s: %s", c.RemoteAddr(), err.Error())
					return
				}
			}
			//为当前请求创建会话
			var sess = newSession(that, conn, protoFunc)
			//触发accept事件
			if stat := that.pluginContainer.afterAccept(sess); !stat.OK() {
				_ = sess.Close()
				return
			}

			logger.Infof("accept ok (network:%s, addr:%s, id:%s)", network, sess.RemoteAddr().String(), sess.ID())
			that.sessHub.set(sess)
			sess.changeStatus(statusOk)
			// 启动消息侦听
			sess.startReadAndHandle()
		})
	}

}

// ListenAndServe 端点启动并监听，对外提供服务
func (that *endpoint) ListenAndServe(protoFunc ...proto.ProtoFunc) error {
	lis, err := NewInheritedListener(that.listerAddr, that.tlsConfig)
	if err != nil {
		logger.Fatalf("%v", err)
	}
	return that.serveListener(lis, protoFunc...)
}

// Close 关闭端点
func (that *endpoint) Close() (err error) {

	defer func() {
		if p := recover(); p != nil {
			err = gerror.NewSkipf(2, "panic:%v\n", p)
		}
	}()
	//关闭endpoint前，执行该事件
	err = that.pluginContainer.beforeCloseEndpoint(that)

	close(that.closeCh)
	for lis := range that.listeners {
		_ = lis.Close()
	}
	//deleteEndpoint(that)
	// 平滑重启移除endpoint
	graceful.Graceful().DeleteEndpoint(that)
	var (
		count int
		errCh = make(chan error, 1024)
	)
	that.sessHub.rangeCallback(func(s *session) bool {
		count++
		dgpool.FILOAnywayGo(func() {
			e := s.Close()
			if e != nil {
				logger.Error(e)
			}
			errCh <- e
		})
		return true
	})
	for i := 0; i < count; i++ {
		err = errors2.Merge(err, <-errCh)
	}
	close(errCh)

	//关闭endpoint后执行该事件
	that.pluginContainer.afterCloseEndpoint(that, err)
	return err
}

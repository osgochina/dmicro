package drpc

import (
	"context"
	"fmt"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/os/grpool"
	"github.com/osgochina/dmicro/drpc/codec"
	"github.com/osgochina/dmicro/drpc/message"
	"github.com/osgochina/dmicro/drpc/proto"
	"github.com/osgochina/dmicro/drpc/socket"
	"github.com/osgochina/dmicro/drpc/status"
	"github.com/osgochina/dmicro/logger"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// EarlySession 尚未启动 goroutine 读取数据的链接会话
type EarlySession interface {
	Endpoint() Endpoint

	// LocalAddr 本地地址
	LocalAddr() net.Addr

	// RemoteAddr 远端地址
	RemoteAddr() net.Addr

	// Swap 临时存储区内容
	Swap() *gmap.Map

	// SetID 设置session id
	SetID(newID string)

	// ControlFD 原始链接的fd
	ControlFD(f func(fd uintptr)) error

	// ModifySocket 修改session的底层socket
	ModifySocket(fn func(conn net.Conn) (modifiedConn net.Conn, newProtoFunc proto.ProtoFunc))

	// GetProtoFunc 获取协议方法
	GetProtoFunc() proto.ProtoFunc

	// EarlySend 在会话刚建立的时候临时发送消息，不执行任何中间件
	EarlySend(mType byte, serviceMethod string, body interface{}, stat *status.Status, setting ...message.MsgSetting) (opStat *status.Status)

	// EarlyReceive 在会话刚建立的时候临时接受信息，不执行任何中间件
	EarlyReceive(newArgs message.NewBodyFunc, ctx ...context.Context) (input message.Message)

	// EarlyCall 在会话刚建立的时候临时调用call发送和接受消息，不执行任何中间件
	EarlyCall(serviceMethod string, args, reply interface{}, callSetting ...message.MsgSetting) (opStat *status.Status)

	// EarlyReply 在会话刚建立的时候临时回复消息，不执行任何中间件
	EarlyReply(req message.Message, body interface{}, stat *status.Status, setting ...message.MsgSetting) (opStat *status.Status)

	// RawPush 发送原始push消息，不执行任何中间件
	RawPush(serviceMethod string, args interface{}, setting ...message.MsgSetting) (opStat *status.Status)

	// SessionAge 获取session最大的生存周期
	SessionAge() time.Duration

	// ContextAge 获取 CALL 和 PUSH 消息的最大生存周期
	ContextAge() time.Duration

	// SetSessionAge 设置session的最大生存周期
	SetSessionAge(duration time.Duration)

	// SetContextAge 设置单个 CALL 和 PUSH 消息的最大生存周期
	SetContextAge(duration time.Duration)
}

// BaseSession 基础的session
type BaseSession interface {
	Endpoint() Endpoint

	// ID 获取id
	ID() string

	// LocalAddr 本地地址
	LocalAddr() net.Addr

	// RemoteAddr 远端地址
	RemoteAddr() net.Addr

	// Swap 返回交换区的内容
	Swap() *gmap.Map
}

// CtxSession 在处理程序上下文中传递的会话对象
type CtxSession interface {

	// ID 获取id
	ID() string

	// LocalAddr 本地地址
	LocalAddr() net.Addr

	// RemoteAddr 远端地址
	RemoteAddr() net.Addr

	// Swap 返回交换区的内容
	Swap() *gmap.Map

	// CloseNotify 返回该链接被关闭时候的通知
	CloseNotify() <-chan struct{}

	// Health 检查该session是否健康
	Health() bool

	// AsyncCall 发送消息，并异步接收响应
	AsyncCall(serviceMethod string, args interface{}, result interface{}, callCmdChan chan<- CallCmd, setting ...message.MsgSetting) CallCmd

	// Call 发送消息并获得响应值
	Call(serviceMethod string, args interface{}, result interface{}, setting ...message.MsgSetting) CallCmd

	// Push 发送消息，不接收响应，只返回发送状态
	Push(serviceMethod string, args interface{}, setting ...message.MsgSetting) *status.Status

	// SessionAge 获取session最大的生存周期
	SessionAge() time.Duration

	// ContextAge 获取 CALL 和 PUSH 消息的最大生存周期
	ContextAge() time.Duration
}

type Session interface {
	Endpoint() Endpoint

	// SetID 设置session id
	SetID(newID string)

	// Close 关闭session
	Close() error

	CtxSession
}

// 会话的状态
//不能改变枚举值的顺序
const (
	statusPreparing      int32 = iota // 会话准备阶段 1
	statusOk                          // 会话就绪 2
	statusActiveClosing               // 会话主动关闭中 3
	statusActiveClosed                // 会话已经主动关闭 4
	statusPassiveClosing              // 会话被动关闭中 5
	statusPassiveClosed               // 会话被动关闭 6
	statusRedialing                   // 会话重建中 7
	statusRedialFailed                // 会话重建失败 8
)

type session struct {
	endpoint              *endpoint
	getCallHandler        func(serviceMethodPath string) (*Handler, bool)
	getPushHandler        func(serviceMethodPath string) (*Handler, bool)
	timeNow               func() int64
	callCmdMap            *gmap.Map
	protoFuncList         []proto.ProtoFunc
	socket                socket.Socket
	closeNotifyCh         chan struct{}
	writeLock             sync.Mutex
	graceCtxWaitGroup     sync.WaitGroup // 当前会话中的处理程序计数器
	graceCtxMutex         sync.Mutex
	graceCallCmdWaitGroup sync.WaitGroup // call cmd 方法等待组
	sessionAge            time.Duration
	contextAge            time.Duration
	sessionAgeLock        sync.RWMutex
	contextAgeLock        sync.RWMutex
	lock                  sync.RWMutex
	seq                   int32
	status                int32
	didCloseNotify        int32

	//链接如果断开，重新拨号，只有作为客户端角色的时候才有效果
	redialForClientLocked func() bool
}

var (
	_ EarlySession = new(session)
	_ BaseSession  = new(session)
	_ CtxSession   = new(session)
	_ Session      = new(session)
)

func newSession(e *endpoint, conn net.Conn, protoFunc []proto.ProtoFunc) *session {
	var s = &session{
		endpoint:       e,
		getCallHandler: e.router.subRouter.getCall,
		getPushHandler: e.router.subRouter.getPush,
		timeNow:        e.timeNow,
		protoFuncList:  protoFunc,
		status:         statusPreparing,
		socket:         socket.NewSocket(conn, protoFunc...),
		closeNotifyCh:  make(chan struct{}),
		callCmdMap:     gmap.New(true),
		sessionAge:     e.defaultSessionAge,
		contextAge:     e.defaultContextAge,
	}
	return s
}

//原子性修改session的状态
func (that *session) changeStatus(stat int32) {
	atomic.StoreInt32(&that.status, stat)
}

//原子性尝试修改session的状态
// 从fromList的多个状态修改成to的状态，修改成功则返回true，失败返回false
func (that *session) tryChangeStatus(to int32, fromList ...int32) (changed bool) {
	for _, from := range fromList {
		if atomic.CompareAndSwapInt32(&that.status, from, to) {
			return true
		}
	}
	return false
}

//判断session的状态是否在checkList中，如果在checkList中，则返回true，否则返回false
func (that *session) checkStatus(checkList ...int32) bool {
	stat := atomic.LoadInt32(&that.status)
	for _, v := range checkList {
		if v == stat {
			return true
		}
	}
	return false
}

//原子性的获取当前session的状态
func (that *session) getStatus() int32 {
	return atomic.LoadInt32(&that.status)
}

//判断session的状态是否是开始或者将要结束
func (that *session) goonRead() bool {
	return that.checkStatus(statusOk, statusActiveClosing)
}

// 通知session准备关闭
func (that *session) notifyClosed() {
	if atomic.CompareAndSwapInt32(&that.didCloseNotify, 0, 1) {
		close(that.closeNotifyCh)
	}
}

// CloseNotify session将要关闭的通知
func (that *session) CloseNotify() <-chan struct{} {
	return that.closeNotifyCh
}

// IsActiveClosed 判断链接是否处已经关闭，并且并且是主动关闭的
func (that *session) IsActiveClosed() bool {
	return that.checkStatus(statusActiveClosed)
}

// IsPassiveClosed 判断链接是否已经关闭，并且是被动关闭的
func (that *session) IsPassiveClosed() bool {
	return that.checkStatus(statusPassiveClosed)
}

// Health 判断session会话是否可用
func (that *session) Health() bool {
	s := that.getStatus()
	if s == statusOk {
		return true
	}
	if that.redialForClientLocked == nil {
		return false
	}
	if s == statusPassiveClosed {
		return true
	}
	return false
}

// Endpoint 获取当前session属于那个Endpoint
func (that *session) Endpoint() Endpoint {
	return that.endpoint
}

// ID 获取session的id
func (that *session) ID() string {
	return that.socket.ID()
}

// SetID 修改id
func (that *session) SetID(newID string) {
	oldID := that.ID()
	if oldID == newID {
		return
	}
	that.socket.SetID(newID)
	hub := that.endpoint.sessHub
	hub.set(that)
	hub.delete(oldID)
	logger.Infof("session changes id: %s -> %s", oldID, newID)
}

// ControlFD 处理底层fd
func (that *session) ControlFD(f func(fd uintptr)) error {
	that.lock.RLock()
	defer that.lock.RUnlock()
	return that.socket.ControlFD(f)
}

//获取会话的原始链接
func (that *session) getConn() net.Conn {
	return that.socket.Raw()
}

// ModifySocket 替换底层的socket套接字
// NOTE:
// 连接fd不允许更改!
// 继承以前的session id和Swap;
// 如果 modifiedConn!=nil,重置 net.Conn
// 如果 newProtoFunc!=nil, 重置 ProtoFunc.
func (that *session) ModifySocket(fn func(conn net.Conn) (modifiedConn net.Conn, newProtoFunc proto.ProtoFunc)) {
	conn := that.getConn()
	modifiedConn, newProtoFunc := fn(conn)
	isModifiedConn := modifiedConn != nil
	isNewProtoFunc := newProtoFunc != nil
	if isNewProtoFunc {
		that.protoFuncList = that.protoFuncList[:0]
		that.protoFuncList = append(that.protoFuncList, newProtoFunc)
	}
	if !isModifiedConn && !isNewProtoFunc {
		return
	}
	var pub *gmap.Map
	if that.socket.SwapLen() > 0 {
		pub = that.socket.Swap()
	}
	id := that.ID()
	that.socket.Reset(modifiedConn, that.protoFuncList...)
	that.socket.Swap(pub) // set the old swap
	that.socket.SetID(id)
}

// GetProtoFunc 获取协议方法
func (that *session) GetProtoFunc() proto.ProtoFunc {
	if len(that.protoFuncList) > 0 && that.protoFuncList[0] != nil {
		return that.protoFuncList[0]
	}
	return socket.DefaultProtoFunc()
}

// LocalAddr 获取本地监听地址
func (that *session) LocalAddr() net.Addr {
	return that.socket.LocalAddr()
}

// RemoteAddr 获取远程链接的地址
func (that *session) RemoteAddr() net.Addr {
	return that.socket.RemoteAddr()
}

// SessionAge 获取session的生存周期
func (that *session) SessionAge() time.Duration {
	that.sessionAgeLock.RLock()
	age := that.sessionAge
	that.sessionAgeLock.RUnlock()
	return age
}

// SetSessionAge 设置会话的最大生命周期
func (that *session) SetSessionAge(duration time.Duration) {
	that.sessionAgeLock.Lock()
	that.sessionAge = duration
	if duration > 0 {
		_ = that.socket.SetReadDeadline(time.Now().Add(duration))
	} else {
		_ = that.socket.SetReadDeadline(time.Time{})
	}
	that.sessionAgeLock.Unlock()
}

// ContextAge 获取CALL 或者PUSH上下文的最大生命周期
func (that *session) ContextAge() time.Duration {
	that.contextAgeLock.RLock()
	age := that.contextAge
	that.contextAgeLock.RUnlock()
	return age
}

// SetContextAge 设置CALL 或者PUSH上下文的最大生命周期
func (that *session) SetContextAge(duration time.Duration) {
	that.contextAgeLock.Lock()
	that.contextAge = duration
	that.contextAgeLock.Unlock()
}

// Close 关闭当前session
func (that *session) Close() error {
	that.lock.Lock()
	defer that.lock.Unlock()
	return that.closeLocked()
}

// Swap 获取session的临时交换区数据
func (that *session) Swap() *gmap.Map {
	return that.socket.Swap()
}

//发送消息
func (that *session) send(
	mType byte,
	seq int32,
	serviceMethod string,
	body interface{},
	stat *status.Status,
	setting []message.MsgSetting) (message.Message, *status.Status) {

	//生成output消息对象
	output := message.GetMessage(setting...)
	//设置消息类型
	output.SetMType(mType)
	if seq == 0 {
		seq = atomic.AddInt32(&that.seq, 1)
	}
	//设置序列号
	output.SetSeq(seq)
	//设置消息内容编码格式
	if output.BodyCodec() == codec.NilCodecID {
		output.SetBodyCodec(that.endpoint.defaultBodyCodec)
	}
	//设置请求名
	if len(serviceMethod) > 0 {
		output.SetServiceMethod(serviceMethod)
	}
	//设置消息内容
	if body != nil {
		output.SetBody(body)
	}
	//判断状态是否是ok，不是的话则要设置
	if !stat.OK() {
		output.SetStatus(stat)
	}
	return output, that.doSend(output)
}

//发送消息
func (that *session) doSend(output message.Message) *status.Status {

	//如果设置了上下文生存时间，则应该对会话设置生存时间
	if age := that.ContextAge(); age > 0 {
		ctxTimout, cancel := context.WithTimeout(output.Context(), age)
		defer cancel()
		message.WithContext(ctxTimout)(output)
	}
	that.writeLock.Lock()
	defer that.writeLock.Unlock()

	ctx := output.Context()
	select {
	case <-ctx.Done():
		return statWriteFailed.Copy(ctx.Err())
	default:
		deadline, _ := ctx.Deadline()              //获取上下文超时时间
		_ = that.socket.SetWriteDeadline(deadline) //设置链接写入的超时时间
		err := that.socket.WriteMessage(output)    // 写入消息
		if err == nil {
			return nil
		}
		if err == io.EOF || err == socket.ErrProactivelyCloseSocket {
			return statConnClosed
		}
		logger.Debugf("write error: %s", err.Error())
		return statWriteFailed.Copy(err)
	}
}

// EarlySend 当会话刚刚建立的时候，可以使用 EarlySend 发送临时消息，它不会经过插件
func (that *session) EarlySend(mType byte, serviceMethod string, body interface{}, stat *Status, setting ...message.MsgSetting) (opStat *Status) {
	if !that.checkStatus(statusPreparing) {
		return statUnpreparedError
	}
	var output message.Message
	defer func() {
		if output != nil {
			message.PutMessage(output)
		}
		if p := recover(); p != nil {
			opStat = statBadMessage.Copy(p, 3)
		}
	}()
	output, opStat = that.send(mType, 0, serviceMethod, body, stat, setting)
	return opStat
}

// EarlyReceive 当会话刚刚建立的时候可以使用EarlyReceive接收临时消息，它不会调用其他插件
func (that *session) EarlyReceive(newArgs message.NewBodyFunc, ctx ...context.Context) (input message.Message) {
	if len(ctx) > 0 {
		input = message.GetMessage(message.WithContext(ctx[0]))
	} else {
		input = message.GetMessage()
	}
	//判断会话当前状态是不是处于准备阶段
	if !that.checkStatus(statusPreparing) {
		input.SetStatus(statUnpreparedError)
		return input
	}

	input.SetNewBody(newArgs)
	// 如果处理失败了，则把错误赋值给status
	defer func() {
		if p := recover(); p != nil {
			input.SetStatus(statBadMessage.Copy(p, 3))
		}
	}()

	//给消息处理上下文增加超时时间
	if age := that.ContextAge(); age > 0 {
		ctxTimeout, cancel := context.WithTimeout(input.Context(), age)
		defer cancel()
		message.WithContext(ctxTimeout)(input)
	}
	//给链接增加超时时间
	deadline, _ := input.Context().Deadline()
	_ = that.socket.SetDeadline(deadline)

	//读取消息
	if err := that.socket.ReadMessage(input); err != nil {
		input.SetStatus(statConnClosed.Copy(err))
	}
	return input
}

// EarlyCall 会话准备阶段调用EarlyCall，不会调用其他的插件
func (that *session) EarlyCall(serviceMethod string, args, replay interface{}, callSetting ...message.MsgSetting) (opStat *Status) {
	//判断会话当前状态是不是处于准备阶段
	if !that.checkStatus(statusPreparing) {
		return statUnpreparedError
	}
	//异常捕获
	defer func() {
		if p := recover(); p != nil {
			opStat = statBadMessage.Copy(p, 3)
		}
	}()
	//发送消息
	var output message.Message
	output, opStat = that.send(message.TypeCall, 0, serviceMethod, args, nil, callSetting)
	if !opStat.OK() {
		message.PutMessage(output)
		return opStat
	}
	//接收消息
	ctx := output.Context()
	return that.EarlyReceive(func(_ message.Header) interface{} {
		return replay
	}, ctx).Status()
}

// EarlyReply 会话在准备阶段调用EarlyReply回复
func (that *session) EarlyReply(req message.Message, body interface{}, stat *status.Status, setting ...message.MsgSetting) (opStat *status.Status) {
	if !that.checkStatus(statusPreparing) {
		return statUnpreparedError
	}
	var output message.Message
	defer func() {
		if output != nil {
			message.PutMessage(output)
		}
		if p := recover(); p != nil {
			opStat = statBadMessage.Copy(p, 3)
		}
	}()
	output, opStat = that.send(message.TypeReply, req.Seq(), req.ServiceMethod(), body, stat, setting)
	return opStat
}

//RawPush 发送原始的TypePush消息，不执行任何插件
// 不适用外部设置的seq，也不支持断开后自动重拨
func (that *session) RawPush(serviceMethod string, args interface{}, setting ...message.MsgSetting) (opStat *Status) {
	var output message.Message
	defer func() {
		if output != nil {
			message.PutMessage(output)
		}
		if p := recover(); p != nil {
			opStat = statBadMessage.Copy(p, 3)
		}
	}()
	output, opStat = that.send(message.TypePush, 0, serviceMethod, args, nil, setting)
	return opStat
}

// Push 发送PUSH消息，不会等待响应
func (that *session) Push(serviceMethod string, args interface{}, setting ...message.MsgSetting) (opStat *Status) {
	//从池子中获取上下文对象
	ctx := that.endpoint.getHandleCtx(that, true)

	defer func() {
		//把上下文对象放回池子
		that.endpoint.putHandleCtx(ctx, true)
		if p := recover(); p != nil {
			logger.Errorf("panic:%v\n%s", p, status.PanicStackTrace())
		}
	}()

	ctx.start = that.timeNow()
	output := ctx.output
	output.SetMType(message.TypePush)
	output.SetServiceMethod(serviceMethod)
	output.SetBody(args)
	//针对输出消息执行方法
	for _, fn := range setting {
		if fn != nil {
			fn(output)
		}
	}

	output.SetSeq(atomic.AddInt32(&that.seq, 1))

	if output.BodyCodec() == codec.NilCodecID {
		output.SetBodyCodec(that.endpoint.defaultBodyCodec)
	}
	//设置上下文存活时间
	if age := that.ContextAge(); age > 0 {
		ctxTimout, cancel := context.WithTimeout(output.Context(), age)
		defer cancel()
		message.WithContext(ctxTimout)(output)
	}
	//在push write之前执行插件
	stat := that.endpoint.pluginContainer.beforeWritePush(ctx)
	if !stat.OK() {
		return stat
	}
	var usedConn net.Conn

W:
	if usedConn, stat = that.write(output); !stat.OK() {
		if stat == statConnClosed && that.redialForClient(usedConn) {
			goto W
		}
		return stat
	}
	if enablePrintRunLog() {
		that.printRunLog("", time.Duration(that.timeNow()-ctx.start), nil, output, typePushLaunch)
	}
	//在push write 之后执行插件
	that.endpoint.pluginContainer.afterWritePush(ctx)
	return nil
}

// AsyncCall 异步发送call消息
func (that *session) AsyncCall(serviceMethod string, args interface{}, result interface{}, callCmdChan chan<- CallCmd, setting ...message.MsgSetting) CallCmd {
	if callCmdChan == nil {
		callCmdChan = make(chan CallCmd, 10)
	} else {
		if cap(callCmdChan) == 0 {
			logger.Panicf("*session.AsyncCall(): callCmdChan channel is unbuffered")
		}
	}
	output := message.NewMessage()
	output.SetServiceMethod(serviceMethod)
	output.SetBody(args)
	output.SetMType(message.TypeCall)
	for _, fn := range setting {
		if fn != nil {
			fn(output)
		}
	}
	seq := atomic.AddInt32(&that.seq, 1)
	output.SetSeq(seq)
	if output.BodyCodec() == codec.NilCodecID {
		output.SetBodyCodec(that.endpoint.defaultBodyCodec)
	}
	if age := that.ContextAge(); age > 0 {
		ctxTimout, cancel := context.WithTimeout(output.Context(), age)
		defer cancel()
		message.WithContext(ctxTimout)(output)
	}
	cmd := &callCmd{
		sess:        that,
		output:      output,
		result:      result,
		callCmdChan: callCmdChan,
		doneChan:    make(chan struct{}),
		start:       that.timeNow(),
		swap:        gmap.New(true),
	}
	// 计数 call cmd
	that.graceCallCmdWaitGroup.Add(1)

	if that.socket.SwapLen() > 0 {
		cmd.swap = that.socket.Swap().Clone(true)
	}
	cmd.mu.Lock()
	defer cmd.mu.Unlock()

	that.callCmdMap.Set(seq, cmd)
	defer func() {
		if p := recover(); p != nil {
			logger.Errorf("panic:%v\n%s", p, status.PanicStackTrace())
		}
	}()
	//write call消息写入之前执行插件
	cmd.stat = that.endpoint.pluginContainer.beforeWriteCall(cmd)
	if !cmd.stat.OK() {
		cmd.done()
		return cmd
	}
	var usedConn net.Conn
W:
	if usedConn, cmd.stat = that.write(output); !cmd.stat.OK() {
		if cmd.stat == statConnClosed && that.redialForClient(usedConn) {
			goto W
		}
		cmd.done()
		return cmd
	}
	//发送call消息之后，执行插件
	that.endpoint.pluginContainer.afterWriteCall(cmd)
	return cmd
}

// Call 发送call消息，并且同步返回结果
func (that *session) Call(serviceMethod string, args interface{}, result interface{}, setting ...message.MsgSetting) CallCmd {
	cCmd := that.AsyncCall(serviceMethod, args, result, make(chan CallCmd, 1), setting...)
	<-cCmd.Done()
	return cCmd
}

//写入消息
func (that *session) write(msg message.Message) (net.Conn, *Status) {
	usedConn := that.getConn()
	sta := that.getStatus()

	//当前会话的状态必须 可用状态 或者 主动关闭中的情况下，需要发送的消息是回复消息
	if !(sta == statusOk || (sta == statusActiveClosing && msg.MType() == message.TypeReply)) {
		return usedConn, statConnClosed
	}
	var (
		err         error
		ctx         = msg.Context()
		deadline, _ = ctx.Deadline()
	)
	//判断当前消息的上下文是否可用
	select {
	case <-ctx.Done():
		err = ctx.Err()
		goto ERR
	default:
	}
	//加锁
	that.writeLock.Lock()
	defer that.writeLock.Unlock()

	select {
	case <-ctx.Done():
		err = ctx.Err()
		goto ERR
	default:
		//设置写入超时时间
		_ = that.socket.SetWriteDeadline(deadline)
		//写入消息
		err = that.socket.WriteMessage(msg)
	}
	//如果写入成功，则返回
	if err == nil {
		return usedConn, nil
	}
	//链接已断开，或者链接已经主动关闭
	if err == io.EOF || err == socket.ErrProactivelyCloseSocket {
		return usedConn, statConnClosed
	}
	logger.Debugf("write error: %s", err.Error())
ERR:
	return usedConn, statWriteFailed.Copy(err)
}

//重新链接
func (that *session) redialForClient(oldConn net.Conn) bool {
	if that.redialForClientLocked == nil {
		return false
	}
	that.lock.Lock()
	defer that.lock.Unlock()

	//避免从write和readDisconnected方法重复调用
	if oldConn != that.getConn() {
		return true
	}

	if that.tryChangeStatus(statusRedialing, statusOk, statusPassiveClosing, statusPassiveClosed, statusRedialFailed) {
		return that.redialForClientLocked()
	}
	return false
}

//安全关闭会话
func (that *session) closeLocked() error {
	// 把会话从 statusOk 和 statusPreparing 状态转换成 statusActiveClosing
	//如果失败，则返回nil
	if !that.tryChangeStatus(statusActiveClosing, statusOk, statusPreparing) {
		return nil
	}
	//从端点的会话池中删除自己
	that.endpoint.sessHub.delete(that.ID())
	//发送会话准备关闭通知
	that.notifyClosed()
	// 优雅的结束会话
	that.graceCtxWait()
	// 优雅的等待会话中的链接关闭
	that.graceCallCmdWaitGroup.Wait()
	//更改会话的状态为主动关闭
	that.changeStatus(statusActiveClosed)
	//关闭链接
	err := that.socket.Close()
	// 执行链接关闭插件
	that.endpoint.pluginContainer.afterDisconnect(that)
	return err
}

//读取消息，并根据消息内容，执行对应的操作
func (that *session) startReadAndHandle() {
	var withContext message.MsgSetting
	//如果设置了会话的生命期，则针对该上下文设置生命周期
	if readTimeout := that.SessionAge(); readTimeout > 0 {
		_ = that.socket.SetReadDeadline(time.Now().Add(readTimeout))
		ctxTimout, cancel := context.WithTimeout(context.Background(), readTimeout)
		defer cancel()
		withContext = message.WithContext(ctxTimout)
	} else {
		_ = that.socket.SetReadDeadline(time.Time{})
		withContext = message.WithContext(nil)
	}
	var (
		err      error
		usedConn = that.getConn()
	)

	defer func() {
		//捕获错误，并且继续执行
		if p := recover(); p != nil {
			err = fmt.Errorf("panic:%v\n%s", p, status.PanicStackTrace())
		}
		that.readDisconnected(usedConn, err)
	}()
	//判断该会话是否可以继续执行
	for that.goonRead() {
		var ctx = that.endpoint.getHandleCtx(that, false)
		//给input消息设置上下文信息，主要是生存周期
		withContext(ctx.input)
		//读取消息前，执行插件
		if that.endpoint.pluginContainer.beforeReadHeader(ctx) != nil {
			that.endpoint.putHandleCtx(ctx, false)
			return
		}
		err = that.socket.ReadMessage(ctx.input)
		//读取消息失败，或者当前会话状态不能继续执行
		if (err != nil && ctx.GetBodyCodec() == codec.NilCodecID) || !that.goonRead() {
			that.endpoint.putHandleCtx(ctx, false)
			return
		}
		//如果读取消息有错误，则把错误赋值给处理器上下文
		if err != nil {
			ctx.stat = statBadMessage.Copy(err)
		}
		// 给优雅处理器添加一次记录,优雅的结束会话之前，需要等待改协程处理完毕
		that.graceCtxWaitGroup.Add(1)

		if err = grpool.Add(func() {
			defer that.endpoint.putHandleCtx(ctx, true)
			ctx.handle()
		}); err != nil {
			that.endpoint.putHandleCtx(ctx, true)
		}
	}
}

// 处理已经断开的链接
func (that *session) readDisconnected(oldConn net.Conn, err error) {
	stat := that.getStatus()
	switch stat {
	//如果当前会话状态是 主动关闭，被动关闭，被动关闭中，则什么都不做，直接返回
	case statusPassiveClosed, statusActiveClosed, statusPassiveClosing:
		return
	//如果是主动关闭中，则不做什么，继续执行逻辑
	case statusActiveClosing:
	//所有其他逻辑，都把该会话设置为被动关闭中
	default:
		that.changeStatus(statusPassiveClosing)
	}
	//删除端点中会话池中的自己
	that.endpoint.sessHub.delete(that.ID())

	var reason string
	//如果错误不是主动关闭会话
	if err != nil && err != socket.ErrProactivelyCloseSocket {
		//如果错误不是链接断开
		if errStr := err.Error(); errStr != "EOF" {
			reason = errStr
			//记录会话关闭原因
			logger.Warningf("disconnect when reading: %T %s", err, errStr)
		}
	}
	//优化的等待所有处理程序结束
	that.graceCtxWait()
	// 循环处理该会话中的各个请求
	for _, v := range that.callCmdMap.Values() {
		cCmd := v.(*callCmd)
		cCmd.mu.Lock()
		//如果该请求不是回复，并且该请求当前状态是ok，则主动取消它
		if !cCmd.hasReply() && cCmd.stat.OK() {
			cCmd.cancel(reason)
		}
		cCmd.mu.Unlock()
	}
	//如果当前会话为主动关闭
	if stat == statusActiveClosing {
		return
	}
	if that.socket != nil {
		//关闭链接
		_ = that.socket.Close()
	}

	//重新链接失败
	if !that.redialForClient(oldConn) {
		//设置当前会话为被动关闭
		that.changeStatus(statusPassiveClosed)
		//发送会话关闭通知
		that.notifyClosed()
		//执行会话关闭事件
		that.endpoint.pluginContainer.afterDisconnect(that)
	}
}

func (that *session) graceCtxWait() {
	that.graceCtxMutex.Lock()
	that.graceCtxWaitGroup.Wait()
	that.graceCtxMutex.Unlock()
}

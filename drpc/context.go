package drpc

import (
	"context"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/util/gconv"
	"github.com/osgochina/dmicro/drpc/codec"
	"github.com/osgochina/dmicro/drpc/message"
	"github.com/osgochina/dmicro/drpc/status"
	"github.com/osgochina/dmicro/logger"
	"reflect"
	"sync"
	"time"
)

// EarlyCtx 基础上下文
type EarlyCtx interface {

	// Endpoint 获取当前Endpoint
	Endpoint() Endpoint

	// Session 返回当前的session
	Session() CtxSession

	// IP 返回远端ip
	IP() string

	// RealIP 返回远端真实ip
	RealIP() string

	// Swap 返回自定义交换区数据
	Swap() *gmap.Map

	// Context 获取上下文
	Context() context.Context
}

// WriteCtx 写消息时使用的上下文方法
type WriteCtx interface {
	EarlyCtx

	// Output 将要发送的消息对象
	Output() message.Message

	// StatusOK 状态是否ok
	StatusOK() bool

	// Status 当前步骤的状态
	Status() *status.Status
}

// InputCtx 该上下文是一个公共上下文
type inputCtx interface {
	EarlyCtx

	// Seq 获取消息的序列号
	Seq() int32

	// PeekMeta 窥视消息的元数据
	PeekMeta(key string) interface{}

	// VisitMeta 浏览消息的元数据
	VisitMeta(f func(key, value interface{}) bool)

	// CopyMeta 赋值消息的元数据
	CopyMeta() *gmap.Map

	// ServiceMethod 该消息需要访问的服务名
	ServiceMethod() string

	// ResetServiceMethod 重置该消息将要访问的服务名
	ResetServiceMethod(string)
}

// ReadCtx 读取消息使用的上下文
type ReadCtx interface {
	inputCtx

	// Input 获取传入的消息
	Input() message.Message

	// StatusOK 状态是否ok
	StatusOK() bool

	// Status 当前步骤的状态
	Status() *status.Status
}

// PushCtx push消息使用的上下文
type PushCtx interface {
	inputCtx

	// GetBodyCodec 获取当前消息的编码格式
	GetBodyCodec() byte
}

// CallCtx call消息使用的上下文
type CallCtx interface {
	inputCtx

	// Input 获取传入的消息
	Input() message.Message

	// GetBodyCodec 获取当前消息的编码格式
	GetBodyCodec() byte

	// Output 将要发送的消息对象
	Output() message.Message

	// ReplyBodyCodec 获取响应消息的编码格式
	ReplyBodyCodec() byte

	// SetBodyCodec 设置响应消息的编码格式
	SetBodyCodec(byte)

	// SetMeta 设置指定key的值
	SetMeta(key, value string)

	// AddTFilterId 设置回复消息传输层的编码过滤方法id
	AddTFilterId(filterID ...byte)
}

// UnknownPushCtx 未知push消息的上下文
type UnknownPushCtx interface {
	inputCtx

	// GetBodyCodec 获取当前消息的编码格式
	GetBodyCodec() byte

	// InputBodyBytes 传入消息体
	InputBodyBytes() []byte

	// BuildBody 如果push消息是未知的消息，则使用v对象解析消息内容
	BuildBody(v interface{}) (bodyCodec byte, err error)
}

// UnknownCallCtx 未知call消息的上下文
type UnknownCallCtx interface {
	inputCtx

	// GetBodyCodec 获取当前消息的编码格式
	GetBodyCodec() byte

	// InputBodyBytes 传入消息体
	InputBodyBytes() []byte

	// BuildBody 如果push消息是未知的消息，则使用v对象解析消息内容
	BuildBody(v interface{}) (bodyCodec byte, err error)

	// SetBodyCodec 设置回复消息的编码格式
	SetBodyCodec(byte)

	// SetMeta 设置指定key的值
	SetMeta(key, value string)

	// AddTFilterId 设置回复消息传输层的编码过滤方法id
	AddTFilterId(filterID ...byte)
}

var (
	_ EarlyCtx       = new(handlerCtx)
	_ inputCtx       = new(handlerCtx)
	_ WriteCtx       = new(handlerCtx)
	_ ReadCtx        = new(handlerCtx)
	_ PushCtx        = new(handlerCtx)
	_ CallCtx        = new(handlerCtx)
	_ UnknownPushCtx = new(handlerCtx)
	_ UnknownCallCtx = new(handlerCtx)
)

var emptyValue = reflect.Value{}

// handlerCtx是 PushCtx 和 CallCtx 的底层公共实例
type handlerCtx struct {
	sess            *session
	input           message.Message
	output          message.Message
	handler         *Handler
	arg             reflect.Value // 消息传入的参数
	callCmd         *callCmd
	swap            *gmap.Map
	start           int64
	cost            time.Duration
	pluginContainer *PluginContainer
	stat            *status.Status
	context         context.Context
}

//newReadHandleCtx 创建一个给request/response或push使用的上下文
func newReadHandleCtx() *handlerCtx {
	c := new(handlerCtx)
	c.input = message.NewMessage()
	c.input.SetNewBody(c.buildingBody)
	c.output = message.NewMessage()
	return c
}

//会话上下文生成池
var handlerCtxPool = sync.Pool{
	New: func() interface{} {
		return newReadHandleCtx()
	},
}

//重新初始化
func (that *handlerCtx) reInit(s *session) {
	that.sess = s
	that.swap = s.socket.Swap().Clone(true)
}

//清除上下文
func (that *handlerCtx) clean() {
	that.sess = nil
	that.handler = nil
	that.arg = emptyValue
	that.swap = nil
	that.callCmd = nil
	that.cost = 0
	that.pluginContainer = nil
	that.stat = nil
	that.context = nil
	that.input.Reset(message.WithNewBody(that.buildingBody))
	that.output.Reset()
}

// Endpoint 获取当前端点
func (that *handlerCtx) Endpoint() Endpoint {
	return that.sess.Endpoint()
}

// Session 获取当前的会话
func (that *handlerCtx) Session() CtxSession {
	return that.sess
}

// Input 获取输入消息
func (that *handlerCtx) Input() message.Message {
	return that.input
}

// Output 获取输出消息
func (that *handlerCtx) Output() message.Message {
	return that.output
}

// Swap 获取交换区数据
func (that *handlerCtx) Swap() *gmap.Map {
	return that.swap
}

// Seq 获取请求序列号
func (that *handlerCtx) Seq() int32 {
	return that.input.Seq()
}

// ServiceMethod 获取请求的服务名
func (that *handlerCtx) ServiceMethod() string {
	return that.input.ServiceMethod()
}

// ResetServiceMethod 重写请求的服务名
func (that *handlerCtx) ResetServiceMethod(serviceMethod string) {
	that.input.SetServiceMethod(serviceMethod)
}

// PeekMeta 查看请求消息的元数据
func (that *handlerCtx) PeekMeta(key string) interface{} {
	return that.input.Meta().Get(key)
}

// VisitMeta 迭代请求消息的元数据
func (that *handlerCtx) VisitMeta(f func(key, value interface{}) bool) {
	that.input.Meta().Iterator(f)
}

// CopyMeta 复制请求消息的元数据
func (that *handlerCtx) CopyMeta() *gmap.Map {
	return that.input.Meta().Clone(true)
}

// SetMeta 设置发送消息的元数据
func (that *handlerCtx) SetMeta(key, value string) {
	that.output.Meta().Set(key, value)
}

// GetBodyCodec 获取请求消息的编码格式
func (that *handlerCtx) GetBodyCodec() byte {
	return that.input.BodyCodec()
}

// SetBodyCodec 设置发送消息的编码格式
func (that *handlerCtx) SetBodyCodec(bodyCodec byte) {
	that.output.SetBodyCodec(bodyCodec)
}

// AddTFilterId 针对发送消息添加处理管道
func (that *handlerCtx) AddTFilterId(filterID ...byte) {
	_ = that.output.PipeTFilter().Append(filterID...)
}

// IP 获取远程服务的ip
func (that *handlerCtx) IP() string {
	return that.sess.RemoteAddr().String()
}

// RealIP 获取远程服务的真实ip
func (that *handlerCtx) RealIP() string {
	realIP := gconv.String(that.PeekMeta(message.MetaRealIP))
	if len(realIP) > 0 {
		return realIP
	}
	return that.sess.RemoteAddr().String()
}

// Context 获取当前上下文
func (that *handlerCtx) Context() context.Context {
	if that.context == nil {
		return that.input.Context()
	}
	return that.context
}

// 设置上下文
func (that *handlerCtx) setContext(ctx context.Context) {
	that.context = ctx
}

// StatusOK 判断该上下文的状态是否是ok
func (that *handlerCtx) StatusOK() bool {
	return that.stat.OK()
}

// Status 获取当前上下问的状态
func (that *handlerCtx) Status() *status.Status {
	return that.stat
}

// InputBodyBytes 获取接收消息的消息体
func (that *handlerCtx) InputBodyBytes() []byte {
	b, ok := that.input.Body().(*[]byte)
	if !ok {
		return nil
	}
	return *b
}

// BuildBody 把原始接收到的消息，类型转换成 传入的v类型，并且进行解码，最后该消息获得的body就是 v类型
func (that *handlerCtx) BuildBody(v interface{}) (byte, error) {
	b := that.InputBodyBytes()
	if b == nil {
		return codec.NilCodecID, nil
	}
	that.input.SetBody(v)
	err := that.input.UnmarshalBody(b)
	return that.input.BodyCodec(), err
}

// ReplyBodyCodec 设置reply消息的正文编解码器，并返回
func (that *handlerCtx) ReplyBodyCodec() byte {
	id := that.output.BodyCodec()
	if id != codec.NilCodecID {
		return id
	}
	id, ok := message.GetAcceptBodyCodec(that.input.Meta())
	if ok {
		if _, err := codec.Get(id); err == nil {
			that.output.SetBodyCodec(id)
			return id
		}
	}
	id = that.input.BodyCodec()
	that.output.SetBodyCodec(id)
	return id
}

//设置响应消息的编解码器，如果产生了错误，则不设置
func (that *handlerCtx) setReplyBodyCodec(hasError bool) {
	if hasError {
		return
	}
	that.ReplyBodyCodec()
}

//通过消息头，构建消息体
func (that *handlerCtx) buildingBody(header message.Header) (body interface{}) {
	that.start = that.sess.timeNow()
	that.pluginContainer = that.sess.endpoint.pluginContainer
	switch header.MType() {
	case message.TypeReply:
		return that.buildReplyBody(header)
	case message.TypePush:
		return that.buildPushBody(header)
	case message.TypeCall:
		return that.buildCallBody(header)
	default:
		that.stat = statCodeMTypeNotAllowed
		return nil
	}
}

// 根据消息头正确打包CALL消息体
func (that *handlerCtx) buildCallBody(header message.Header) interface{} {

	//执行事件
	that.stat = that.pluginContainer.afterReadCallHeader(that)
	if !that.stat.OK() {
		return nil
	}
	//传入的消息如果没有服务方法
	if len(header.ServiceMethod()) == 0 {
		that.stat = statBadMessage.Copy("invalid service method for message")
		return nil
	}
	//如果请求消息的服务名没有命中处理方法，
	//注意，这里会调用路由器匹配，如果路由器匹配不上，但是设置了默认处理方法，也是会返回默认处理方法的
	var ok bool
	that.handler, ok = that.sess.getCallHandler(header.ServiceMethod())
	if !ok {
		that.stat = statNotFound
		return nil
	}

	//使用消息处理方法的插件容器来处理消息
	that.pluginContainer = that.handler.pluginContainer

	//如果处理程序是默认处理未知命令的方法，则把消息内容设置为空byte数组
	if that.handler.IsUnknown() {
		that.input.SetBody(new([]byte))
	} else {
		//如果是正常的处理方法，则正常处理
		that.arg = that.handler.NewArgValue()
		that.input.SetBody(that.arg.Interface())
	}
	//产生读取消息体事件
	that.stat = that.pluginContainer.beforeReadCallBody(that)
	if !that.stat.OK() {
		return nil
	}
	return that.input.Body()
}

// 根据消息头正确读取push消息体
func (that *handlerCtx) buildPushBody(header message.Header) interface{} {

	//触发事件
	that.stat = that.pluginContainer.afterReadPushHeader(that)
	if !that.stat.OK() {
		return nil
	}
	//传入的消息如果没有服务方法
	if len(header.ServiceMethod()) == 0 {
		that.stat = statBadMessage.Copy("invalid service method for message")
		return nil
	}

	//如果请求消息的服务名没有命中处理方法，
	//注意，这里会调用路由器匹配，如果路由器匹配不上，但是设置了默认处理方法，也是会返回默认处理方法的
	var ok bool
	that.handler, ok = that.sess.getPushHandler(header.ServiceMethod())
	if !ok {
		that.stat = statNotFound
		return nil
	}

	//使用消息处理方法的插件容器来处理消息
	that.pluginContainer = that.handler.pluginContainer

	//设置body的对象格式
	that.arg = that.handler.NewArgValue()
	that.input.SetBody(that.arg.Interface())
	that.stat = that.pluginContainer.beforeReadPushBody(that)
	if !that.stat.OK() {
		return nil
	}
	return that.input.Body()
}

//根据消息头构建Reply消息体
func (that *handlerCtx) buildReplyBody(header message.Header) interface{} {

	//从call消息暂存池获取改消息对象
	_callCmd, ok := that.sess.callCmdMap.Search(header.Seq())
	if !ok {
		logger.Warningf("not found call cmd: %v", that.input)
	}
	that.callCmd = _callCmd.(*callCmd)

	// 在handleReply方法中解锁
	that.callCmd.mu.Lock()
	//把收到的回复消息中待的服务名赋值给input消息对象，记录日志使用
	that.input.SetServiceMethod(that.callCmd.output.ServiceMethod())
	//回复中的交换数据
	that.swap = that.callCmd.swap
	//设置回复消息体编码格式
	that.callCmd.inputBodyCodec = that.GetBodyCodec()
	//把获取到的消息元数据赋值给调用方
	that.callCmd.inputMeta = that.input.Meta().Clone(true)
	//设置上下文
	that.setContext(that.callCmd.output.Context())
	//设置消息体的格式
	that.input.SetBody(that.callCmd.result)

	//读取reply消息头之后执行该事件
	stat := that.pluginContainer.afterReadReplyHeader(that)
	if !stat.OK() {
		that.callCmd.stat = stat
		return nil
	}
	//读取Reply消息体之前执行该事件
	stat = that.pluginContainer.beforeReadReplyBody(that)
	if !stat.OK() {
		that.callCmd.stat = stat
		return nil
	}
	return that.input.Body()
}

//处理回复消息
func (that *handlerCtx) handleReply() {
	// 如果callCmd解析失败
	if that.callCmd == nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			logger.Errorf("panic:%v\n%s", p, status.PanicStackTrace())
		}
		//把响应消息的消息体赋值给返回值
		that.callCmd.result = that.input.Body()
		that.stat = that.callCmd.stat
		//告诉调用方，此次调用已经完成
		that.callCmd.done()
		//计算调用事件
		that.callCmd.cost = time.Duration(that.sess.timeNow() - that.callCmd.start)
		if enablePrintRunLog() {
			that.sess.printRunLog(that.RealIP(), that.callCmd.cost, that.input, that.callCmd.output, typeCallLaunch)
		}
		// lock: bindReply
		that.callCmd.mu.Unlock()
	}()

	//判断处理状态是否成功，如果都很成功，则触发事件
	if that.callCmd.stat.OK() {
		stat := that.input.Status()
		if stat.OK() {
			stat = that.pluginContainer.afterReadReplyBody(that)
		}
		that.callCmd.stat = stat
	}
}

const logFormatDisconnected = "disconnected due to unsupported message type: %d %s %s %q RECV(%s)"

//收到消息，处理消息
func (that *handlerCtx) handle() {
	if that.stat.Code() == CodeMTypeNotAllowed {
		goto E
	}

	switch that.input.MType() {
	case message.TypeReply:
		// handles call reply
		that.handleReply()
		return

	case message.TypePush:
		//  handles push
		that.handlePush()
		return

	case message.TypeCall:
		// handles and replies call
		that.handleCall()
		return

	default:
	}
E:
	that.output.SetStatus(statCodeMTypeNotAllowed)
	logger.Errorf(logFormatDisconnected,
		that.input.MType(), that.IP(), that.input.ServiceMethod(), that.input.Seq(),
		messageLogBytes(that.input, that.sess.endpoint.printDetail))

	go func() {
		_ = that.sess.Close()
	}()
}

//处理push消息
func (that *handlerCtx) handlePush() {

	//判断push消息时候有处理时间限制
	age := that.sess.ContextAge()
	if age > 0 {
		ctxTimout, cancel := context.WithTimeout(context.Background(), age)
		defer cancel()
		that.setContext(ctxTimout)
	}
	defer func() {
		if p := recover(); p != nil {
			logger.Errorf("panic:%v\n%s", p, status.PanicStackTrace())
		}
		//计算该请求处理消耗时间
		that.recordCost()
		//打印处理log
		if enablePrintRunLog() {
			that.sess.printRunLog(that.RealIP(), that.cost, that.input, nil, typePushHandle)
		}
	}()
	//消息状态正确，且有注册的处理函数
	if that.stat.OK() && that.handler != nil {
		//执行处理事件
		if that.pluginContainer.afterReadPushBody(that) == nil {
			if that.handler.IsUnknown() {
				that.handler.unknownHandleFunc(that)
			} else {
				that.handler.handleFunc(that, that.arg)
			}
		}
	}
	if !that.stat.OK() {
		logger.Warningf("%s", that.stat.String())
	}
}

// 处理call请求
func (that *handlerCtx) handleCall() {
	var isWrite bool

	defer func() {
		if p := recover(); p != nil {
			logger.Errorf("panic:%v\n%s", p, status.PanicStackTrace())
			//报错的情况，如果没有写入响应，则再次写入响应
			if !isWrite {
				if that.stat.OK() {
					that.stat = statInternalServerError.Copy(p)
				}
				that.writeReply(that.stat)
			}
		}
		//计算消耗时间
		that.recordCost()
		//打印处理日志
		if enablePrintRunLog() {
			that.sess.printRunLog(that.RealIP(), that.cost, that.input, that.output, typeCallHandle)
		}
	}()

	//先设置好将要回复的消息参数
	that.output.SetMType(message.TypeReply)
	//把客户端传入的序列号原封不动的带回去
	that.output.SetSeq(that.input.Seq())
	//把请求消息的服务名带入到返回消息中,实际上并不会返回给客户端
	//记录响应消息日志的时候，需要用到这个服务名，所以只能这里设置该服务名
	that.output.SetServiceMethod(that.input.ServiceMethod())
	// 设置返回消息的管道处理器
	that.output.PipeTFilter().AppendFrom(that.input.PipeTFilter())

	age := that.sess.ContextAge()
	if age > 0 {
		ctxTimout, cancel := context.WithTimeout(that.input.Context(), age)
		defer cancel()
		//为自己和响应消息设置生存周期
		that.setContext(ctxTimout)
		message.WithContext(ctxTimout)(that.output)
	}
	if that.stat.OK() {
		that.stat = that.output.Status()
	}
	if that.stat.OK() {
		//触发事件
		that.stat = that.pluginContainer.afterReadCallBody(that)
		if that.stat.OK() {
			//处理
			if that.handler.isUnknown {
				that.handler.unknownHandleFunc(that)
			} else {
				that.handler.handleFunc(that, that.arg)
			}
		}
	}
	//响应
	that.setReplyBodyCodec(!that.stat.OK()) //设置响应正文的编解码器，默认使用请求消息的正文编解码器
	//触发事件
	that.pluginContainer.beforeWriteReply(that)
	//写入回复
	stat := that.writeReply(that.stat)

	//如果写入失败
	if !stat.OK() {
		if that.stat.OK() {
			that.stat = stat
		}
		//写入失败，但不是链接被关闭的状态，继续写入服务错误消息
		if stat.Code() != CodeConnClosed {
			that.writeReply(statInternalServerError.Copy(stat.Cause()))
		}
		return
	}
	isWrite = true
	that.pluginContainer.afterWriteReply(that)
}

//计算该请求处理消耗时间
func (that *handlerCtx) recordCost() {
	that.cost = time.Duration(that.sess.timeNow() - that.start)
}

// 写入reply回复
func (that *handlerCtx) writeReply(stat *Status) *Status {
	//如果处理失败
	if !stat.OK() {
		//把失败原因写入响应消息
		that.output.SetStatus(stat)
		//消息内容置空
		that.output.SetBody(nil)
		//消息体编解码器设置为原始
		that.output.SetBodyCodec(codec.NilCodecID)
	}

	serviceMethod := that.output.ServiceMethod()
	//发送消息的时候，把服务名设置为空，因为响应的时候本来就不应该有服务名，
	//但是记录响应消息日志的时候，又需要用到这个服务名，所以只能这里做一些妥协
	that.output.SetServiceMethod("")
	_, stat = that.sess.write(that.output)
	that.output.SetServiceMethod(serviceMethod)
	return stat
}

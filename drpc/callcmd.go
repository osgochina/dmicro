package drpc

import (
	"context"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/util/gconv"
	"github.com/osgochina/dmicro/drpc/message"
	"github.com/osgochina/dmicro/drpc/status"
	"sync"
	"time"
)

// CallCmd CALL 命令调用后，响应操作的命令
type CallCmd interface {
	TraceEndpoint() (e Endpoint, found bool)
	TraceSession() (sess Session, found bool)

	// Context 协程上下文
	Context() context.Context
	// Output 发送的消息
	Output() message.Message

	// StatusOK 状态是否是OK
	StatusOK() bool
	// Status 状态
	Status() *Status
	// Done 返回指示是否已经完毕的chan
	Done() <-chan struct{}
	// Reply 返回应答
	Reply() (interface{}, *Status)

	// InputBodyCodec 接收到的消息使用的编码
	InputBodyCodec() byte
	// InputMeta 接收到的消息传入的元数据
	InputMeta() *gmap.Map
	// CostTime 消耗的时间
	CostTime() time.Duration
}

type callCmd struct {
	start          int64
	cost           time.Duration
	sess           *session
	output         message.Message
	result         interface{}
	stat           *status.Status
	inputMeta      *gmap.Map
	swap           *gmap.Map
	mu             sync.Mutex
	callCmdChan    chan<- CallCmd
	doneChan       chan struct{}
	inputBodyCodec byte
}

var _ WriteCtx = new(callCmd)

func (that *callCmd) TraceEndpoint() (Endpoint, bool) {
	return that.Endpoint(), true
}

func (that *callCmd) Endpoint() Endpoint {
	return that.sess.Endpoint()
}

func (that *callCmd) TraceSession() (Session, bool) {
	return that.sess, true
}

func (that *callCmd) Session() CtxSession {
	return that.sess
}

// IP 远端请求的ip
func (that *callCmd) IP() string {
	return that.sess.RemoteAddr().String()
}

// RealIP 远端请求的真是ip
func (that *callCmd) RealIP() string {
	realIP, found := that.inputMeta.Search(message.MetaRealIP)
	if found {
		return gconv.String(realIP)
	}
	return that.sess.RemoteAddr().String()
}

// Swap 临时存储区
func (that *callCmd) Swap() *gmap.Map {
	return that.swap
}

// Output 输出消息
func (that *callCmd) Output() message.Message {
	return that.output
}

// Context 输出消息的上下文
func (that *callCmd) Context() context.Context {
	return that.output.Context()
}

// StatusOK 状态是否是ok
func (that *callCmd) StatusOK() bool {
	return that.stat.OK()
}

// Status 状态
func (that *callCmd) Status() *Status {
	return that.stat
}

// Done 是否处理完毕
func (that *callCmd) Done() <-chan struct{} {
	return that.doneChan
}

// Reply 获取远端返回的结果
func (that *callCmd) Reply() (interface{}, *Status) {
	<-that.Done()
	return that.result, that.stat
}

// InputBodyCodec 接收到的消息内容是什么编码
func (that *callCmd) InputBodyCodec() byte {
	<-that.Done()
	return that.inputBodyCodec
}

// InputMeta 接收到的消息元数据
func (that *callCmd) InputMeta() *gmap.Map {
	<-that.Done()
	return that.inputMeta
}

// CostTime 消耗时间
func (that *callCmd) CostTime() time.Duration {
	<-that.Done()
	return that.cost
}

// 处理完成
func (that *callCmd) done() {
	that.sess.callCmdMap.Remove(that.output.Seq())
	that.callCmdChan <- that
	close(that.doneChan)
	// free count call-launch
	that.sess.graceCallCmdWaitGroup.Done()
}

// 取消请求
func (that *callCmd) cancel(reason string) {
	that.sess.callCmdMap.Remove(that.output.Seq())
	if reason != "" {
		that.stat = statConnClosed.Copy(reason)
	} else {
		that.stat = statConnClosed
	}
	that.callCmdChan <- that
	close(that.doneChan)
	// free count call-launch
	that.sess.graceCallCmdWaitGroup.Done()
}

//是否是回复消息
func (that *callCmd) hasReply() bool {
	return that.inputMeta != nil
}

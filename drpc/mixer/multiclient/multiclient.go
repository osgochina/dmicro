package multiclient

import (
	"context"
	"github.com/gogf/gf/v2/container/gpool"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/internal"
	"github.com/osgochina/dmicro/drpc/message"
	"github.com/osgochina/dmicro/drpc/proto"
	"time"
)

// MultiClient 并发客户端，使用连接池的方式实现
type MultiClient struct {
	addr     string
	endpoint drpc.Endpoint
	pool     *gpool.Pool
}

// New 创建并发请求客户端
// addr: 要请求的服务端地址
// maxIdleDuration: 链接的最大闲置时间，超出该时间会自动关闭，0表示不限制,小于0表示使用后立即过期，大于0表示过期时间
// protoFunc: 要使用的协议
func New(endpoint drpc.Endpoint, addr string, maxIdleDuration time.Duration, protoFunc ...proto.ProtoFunc) *MultiClient {
	newSessionFunc := func() (interface{}, error) {
		session, stat := endpoint.Dial(addr, protoFunc...)
		return session, stat.Cause()
	}
	expireFunc := func(sess interface{}) {
		if s, ok := sess.(drpc.Session); ok {
			_ = s.Close()
		}
	}
	return &MultiClient{
		addr:     addr,
		endpoint: endpoint,
		pool:     gpool.New(maxIdleDuration, newSessionFunc, expireFunc),
	}
}

// Addr 请求地址
func (that *MultiClient) Addr() string {
	return that.addr
}

// Endpoint 返回请求终端
func (that *MultiClient) Endpoint() drpc.Endpoint {
	return that.endpoint
}

// Close 关闭客户端
func (that *MultiClient) Close() {
	that.pool.Close()
}

// Size 获取当前链接数
func (that *MultiClient) Size() int {
	return that.pool.Size()
}

// AsyncCall 异步请求
func (that *MultiClient) AsyncCall(uri string, arg interface{}, result interface{}, callCmdChan chan<- drpc.CallCmd, setting ...message.MsgSetting) drpc.CallCmd {
	_sess, err := that.pool.Get()
	if err != nil {
		callCmd := drpc.NewFakeCallCmd(uri, arg, result, drpc.NewStatusByCodeText(drpc.CodeWrongConn, err, false))
		if callCmdChan != nil && cap(callCmdChan) == 0 {
			internal.Panicf(context.TODO(), "*MultiClient.AsyncCall(): callCmdChan channel is unbuffered")
		}
		callCmdChan <- callCmd
		return callCmd
	}
	sess := _sess.(drpc.Session)
	defer func() {
		_ = that.pool.Put(sess)
	}()
	return sess.AsyncCall(uri, arg, result, callCmdChan, setting...)
}

// Call 阻塞请求
func (that *MultiClient) Call(uri string, arg interface{}, result interface{}, setting ...message.MsgSetting) drpc.CallCmd {
	callCmd := that.AsyncCall(uri, arg, result, make(chan drpc.CallCmd, 1), setting...)
	<-callCmd.Done()
	return callCmd
}

// Push 发送push消息
func (that *MultiClient) Push(uri string, arg interface{}, setting ...message.MsgSetting) *drpc.Status {
	_sess, err := that.pool.Get()
	if err != nil {
		return drpc.NewStatusByCodeText(drpc.CodeWrongConn, err, false)
	}
	sess := _sess.(drpc.Session)
	defer func() {
		_ = that.pool.Put(sess)
	}()
	return sess.Push(uri, arg, setting...)
}

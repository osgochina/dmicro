package client

import (
	"github.com/gogf/gf/util/guid"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/proto"
	"sync"
	"time"
)

// Pool 连接池对象
type Pool interface {
	// Close 关闭该连接池
	Close() error
	// Get 获取链接
	Get(addr string) (Conn, error)
	// Release  释放链接
	Release(c Conn, status *drpc.Status) error
}

// Conn 链接对象
type Conn interface {
	// Session session
	drpc.Session
	// Created 链接创建时间
	Created() time.Time

	Id() string
}

// NewPool 创建连接池
func NewPool(opts poolOptions) Pool {
	return newPool(opts)
}

type pool struct {
	size      int
	ttl       time.Duration
	endpoint  drpc.Endpoint
	protoFunc proto.ProtoFunc
	sync.Mutex
	conns map[string][]*poolConn
}

type poolConn struct {
	drpc.Session
	created time.Time
	id      string
}

func newPool(options poolOptions) *pool {
	return &pool{
		size:      options.Size,
		ttl:       options.TTL,
		endpoint:  options.Endpoint,
		protoFunc: options.ProtoFunc,
		conns:     make(map[string][]*poolConn),
	}
}

func (that *pool) Close() error {
	that.Lock()
	for k, c := range that.conns {
		for _, conn := range c {
			_ = conn.Close()
		}
		delete(that.conns, k)
	}
	that.Unlock()
	return nil
}

func (that *pool) Get(addr string) (Conn, error) {
	that.Lock()
	conns := that.conns[addr]
	for len(conns) > 0 {
		// 获取最后一个可用的链接
		conn := conns[len(conns)-1]
		conns = conns[:len(conns)-1]
		that.conns[addr] = conns

		// 判断该链接时候过了存活时间，超过该事件则关闭链接
		if d := time.Since(conn.Created()); d > that.ttl {
			_ = conn.Close()
			continue
		}
		that.Unlock()

		return conn, nil
	}

	that.Unlock()

	sess, stat := that.endpoint.Dial(addr, that.protoFunc)
	if !stat.OK() {
		return nil, stat.Cause()
	}

	return &poolConn{
		Session: sess,
		created: time.Now(),
		id:      guid.S(),
	}, nil
}

func (that *pool) Release(conn Conn, stat *drpc.Status) error {
	if !stat.OK() && drpc.IsConnError(stat) {
		return conn.(*poolConn).Close()
	}
	that.Lock()
	conns := that.conns[conn.RemoteAddr().String()]
	if len(conns) >= that.size {
		that.Unlock()
		return conn.(*poolConn).Close()
	}
	that.conns[conn.RemoteAddr().String()] = append(conns, conn.(*poolConn))
	that.Unlock()
	return nil
}

func (that *poolConn) Created() time.Time {
	return that.created
}

func (that *poolConn) Id() string {
	return that.id
}

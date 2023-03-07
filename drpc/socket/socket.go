package socket

import (
	"bufio"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/container/gtype"
	"github.com/osgochina/dmicro/drpc/message"
	"github.com/osgochina/dmicro/drpc/proto"
	"net"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type Message = message.Message

type Socket interface {

	// ControlFD 获取原始句柄
	ControlFD(f func(fd uintptr)) error

	// LocalAddr 获取socket本地的地址
	LocalAddr() net.Addr

	// RemoteAddr 获取socket远端的地址
	RemoteAddr() net.Addr

	// SetDeadline 设置超时时间
	SetDeadline(t time.Time) error

	// SetReadDeadline 设置读取数据的超时时间
	SetReadDeadline(t time.Time) error

	// SetWriteDeadline 设置发送数据的超时时间
	SetWriteDeadline(t time.Time) error

	// WriteMessage 往链接中写入消息
	WriteMessage(message Message) error

	// ReadMessage 从链接中读取消息头和消息体，并填充到消息对象中
	ReadMessage(message Message) error

	// Read 从链接中读取字符
	Read(b []byte) (n int, err error)

	// Write 写入字符到链接
	Write(b []byte) (n int, err error)

	// Close 关闭链接
	Close() error

	// Swap 链接的自定义数据，如果 newSwap不为空，则会替换内部数据，并返回
	Swap(newSwap ...*gmap.Map) *gmap.Map

	// SwapLen 返回链接中自定义数据的长度
	SwapLen() int

	// ID 返回链接的id
	ID() string

	// SetID 设置链接id
	SetID(string)

	// Reset 重置net.Conn
	Reset(netConn net.Conn, protoFunc ...ProtoFunc)

	// Raw 返回原始链接
	Raw() net.Conn
}

// UnsafeSocket 比socket接口多了更多不安全的方法
type UnsafeSocket interface {
	Socket
	// RawLocked returns the raw net.Conn,
	// can be called in ProtoFunc.
	// NOTE:
	//  Make sure the external is locked before calling
	RawLocked() net.Conn
}

//自定义链接
type socket struct {
	net.Conn
	readerWithBuffer *bufio.Reader
	protocol         proto.Proto
	id               *gtype.String
	idMutex          sync.RWMutex
	swap             *gmap.Map
	swapMutex        sync.RWMutex
	mu               sync.RWMutex
	curState         int32
	fromPool         bool
}

var (
	_ net.Conn     = Socket(nil)
	_ UnsafeSocket = new(socket)
)

var readerSize = 1024

// GetSocket 获取一个socket
func GetSocket(c net.Conn, protoFunc ...ProtoFunc) Socket {
	s := socketPool.Get().(*socket)
	s.Reset(c, protoFunc...)
	return s
}

// NewSocket 对外暴露创建链接的接口
func NewSocket(c net.Conn, protoFunc ...ProtoFunc) Socket {
	return newSocket(c, protoFunc)
}

//创建链接
func newSocket(c net.Conn, protoFuncList []ProtoFunc) *socket {
	var s = &socket{
		Conn:             c,
		readerWithBuffer: bufio.NewReaderSize(c, readerSize),
	}
	s.id = gtype.NewString()
	s.protocol = getProto(protoFuncList, s)
	s.initOptimize()
	return s
}

// Raw 获取原始链接
func (that *socket) Raw() net.Conn {
	that.mu.RLock()
	conn := that.Conn
	that.mu.RUnlock()
	return conn
}

// ControlFD 获取链接的原始句柄
func (that *socket) ControlFD(f func(fd uintptr)) error {
	syscallConn, ok := that.Raw().(syscall.Conn)
	if !ok {
		return syscall.EINVAL
	}
	ctrl, err := syscallConn.SyscallConn()
	if err != nil {
		return err
	}
	return ctrl.Control(f)
}

// RawLocked 获取原始链接
func (that *socket) RawLocked() net.Conn {
	return that.Conn
}

//读取链接中指定字节的数据
func (that *socket) Read(b []byte) (int, error) {
	return that.readerWithBuffer.Read(b)
}

// ReadMessage 读取数据到消息
func (that *socket) ReadMessage(message Message) error {
	that.mu.RLock()
	protocol := that.protocol
	that.mu.RUnlock()
	return protocol.Unpack(message)
}

// WriteMessage 写入消息
func (that *socket) WriteMessage(message Message) error {
	that.mu.RLock()
	protocol := that.protocol
	that.mu.RUnlock()
	err := protocol.Pack(message)
	if err != nil && that.isActiveClosed() {
		err = ErrProactivelyCloseSocket
	}
	return err
}

// Swap 链接的自定义数据，如果 newSwap不为空，则会替换内部数据，并返回
func (that *socket) Swap(newSwap ...*gmap.Map) *gmap.Map {
	that.swapMutex.Lock()
	defer that.swapMutex.Unlock()
	if len(newSwap) > 0 {
		that.swap = newSwap[0]
	} else if that.swap == nil {
		that.swap = gmap.New(true)
	}
	swap := that.swap
	return swap
}

// SwapLen 返回链接中自定义数据的长度
func (that *socket) SwapLen() int {
	that.swapMutex.RLock()
	defer that.swapMutex.RUnlock()
	if that.swap == nil {
		return 0
	}
	return that.swap.Size()
}

func (that *socket) ID() string {
	that.idMutex.RLock()
	defer that.idMutex.RUnlock()
	id := that.id
	if len(id.Val()) == 0 {
		id.Set(that.RemoteAddr().String())
	}
	return id.Val()
}

func (that *socket) SetID(id string) {
	that.idMutex.Lock()
	defer that.idMutex.Unlock()
	that.id.Set(id)
}

func (that *socket) Reset(netConn net.Conn, protoFunc ...ProtoFunc) {
	atomic.StoreInt32(&that.curState, activeClose)
	that.mu.Lock()
	defer that.mu.Unlock()
	that.Conn = netConn
	_, _ = that.readerWithBuffer.Discard(that.readerWithBuffer.Buffered())
	that.readerWithBuffer.Reset(netConn)
	that.protocol = getProto(protoFunc, that)
	that.SetID("")
	that.swapMutex.Lock()
	that.swap = nil
	that.swapMutex.Unlock()
	atomic.StoreInt32(&that.curState, normal)
	that.initOptimize()
}

func (that *socket) Close() error {
	if that.isActiveClosed() {
		return nil
	}
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.isActiveClosed() {
		return nil
	}
	atomic.StoreInt32(&that.curState, activeClose)

	var err error
	if that.Conn != nil {
		_ = that.Conn.Close()
	}
	if that.fromPool {
		that.Conn = nil
		that.swapMutex.Lock()
		that.swap = nil
		that.swapMutex.Unlock()
		that.protocol = nil
		socketPool.Put(that)
	}
	return err
}

//根据参数优化链接
func (that *socket) initOptimize() {
	TryOptimize(that.Conn)
}

//当前链接状态
func (that *socket) isActiveClosed() bool {
	return atomic.LoadInt32(&that.curState) == activeClose
}

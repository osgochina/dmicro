package socket

import (
	"errors"
	"fmt"
	"net"
	"time"
)

const (
	normal      int32 = 0 //链接正常
	activeClose int32 = 1 //链接已关闭
)

var ErrProactivelyCloseSocket = errors.New("socket is closed proactively")

var (
	writeBuffer     int           = -1
	readBuffer      int           = -1
	changeKeepAlive bool          = false
	keepAlive       bool          = true
	keepAlivePeriod time.Duration = -1
	noDelay         bool          = true
)

// 支持keepAlive的链接
type setKeepAliveInterface interface {
	// SetKeepAlive 开启关闭链接的keepAlive支持
	SetKeepAlive(keepalive bool) error
	// SetKeepAlivePeriod 设置keepAlive 心跳的间隔时间
	SetKeepAlivePeriod(d time.Duration) error
}

// 支持读写缓冲区的链接
type setBufferInterface interface {
	// SetReadBuffer 设置读缓冲区长度
	SetReadBuffer(bytes int) error
	// SetWriteBuffer 设置写缓冲区长度
	SetWriteBuffer(bytes int) error
}

//支持no daly开关
type setNoDelayInterface interface {
	// SetNoDelay 开启关闭tcp no delay算法，开启后TCP 连接发送数据时会关闭 Nagle 合并算法，立即发往对端 TCP 连接。
	//在某些场景下，如命令行终端，敲一个命令就需要立马发到服务器，可以提升响应速度，请自行 Google Nagle 算法。
	SetNoDelay(noDelay bool) error
}

func TryOptimize(conn net.Conn) {
	if c, ok := conn.(setKeepAliveInterface); ok {
		if changeKeepAlive {
			_ = c.SetKeepAlive(keepAlive)
		}
		if keepAlivePeriod >= 0 && keepAlive {
			_ = c.SetKeepAlivePeriod(keepAlivePeriod)
		}
	}
	if c, ok := conn.(setBufferInterface); ok {
		if readBuffer >= 0 {
			_ = c.SetReadBuffer(readBuffer)
		}
		if writeBuffer >= 0 {
			_ = c.SetWriteBuffer(writeBuffer)
		}
	}
	if c, ok := conn.(setNoDelayInterface); ok {
		if !noDelay {
			_ = c.SetNoDelay(noDelay)
		}
	}
}

// SetKeepAlive 开启链接保活
func SetKeepAlive(keepalive bool) {
	changeKeepAlive = true
	keepAlive = keepalive
}

// SetKeepAlivePeriod 链接保活间隔时间
func SetKeepAlivePeriod(d time.Duration) {
	if d >= 0 {
		keepAlivePeriod = d
	} else {
		fmt.Println("socket: SetKeepAlivePeriod: invalid keepAlivePeriod:", d)
	}
}

// ReadBuffer 获取链接读缓冲区长度
func ReadBuffer() (bytes int, isDefault bool) {
	return readBuffer, readBuffer == -1
}

// SetReadBuffer 设置链接读缓冲区长度
func SetReadBuffer(bytes int) {
	if bytes >= 0 {
		readBuffer = bytes
	} else {
		fmt.Println("socket: SetReadBuffer: invalid readBuffer size:", bytes)
	}
}

// WriteBuffer 获取链接写缓冲区长度
func WriteBuffer() (bytes int, isDefault bool) {
	return writeBuffer, writeBuffer == -1
}

// SetWriteBuffer 设置链接写缓冲区长度
func SetWriteBuffer(bytes int) {
	if bytes >= 0 {
		writeBuffer = bytes
	} else {
		fmt.Println("socket: SetWriteBuffer: invalid writeBuffer size:", bytes)
	}
}

// SetNoDelay 开启关闭no delay算法
func SetNoDelay(_noDelay bool) {
	noDelay = _noDelay
}

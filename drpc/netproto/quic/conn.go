package quic

import (
	"github.com/quic-go/quic-go"
	"net"
	"time"
)

type Conn struct {
	sess   quic.Connection
	stream quic.Stream
}

var _ net.Conn = new(Conn)

func (that *Conn) Read(b []byte) (n int, err error) {
	return that.stream.Read(b)
}

func (that *Conn) Write(b []byte) (n int, err error) {
	return that.stream.Write(b)
}

func (that *Conn) Close() error {
	err := that.stream.Close()
	if err != nil {
		_ = that.sess.CloseWithError(1, err.Error())
		return err
	}
	return that.sess.CloseWithError(0, "")
}

func (that *Conn) LocalAddr() net.Addr {
	return that.sess.LocalAddr()
}

func (that *Conn) RemoteAddr() net.Addr {
	return that.sess.RemoteAddr()
}

func (that *Conn) SetDeadline(t time.Time) error {
	return that.stream.SetDeadline(t)
}

func (that *Conn) SetReadDeadline(t time.Time) error {
	return that.stream.SetReadDeadline(t)
}

func (that *Conn) SetWriteDeadline(t time.Time) error {
	return that.stream.SetWriteDeadline(t)
}

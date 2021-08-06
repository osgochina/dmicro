package proto

import (
	"github.com/osgochina/dmicro/drpc/message"
	"io"
)

type Message message.Message

type Proto interface {
	Version() (byte, string)
	Pack(Message) error
	Unpack(Message) error
}

type IOWithReadBuffer interface {
	io.ReadWriter
}

type ProtoFunc func(IOWithReadBuffer) Proto

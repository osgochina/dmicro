package drpc_test

import (
	"bufio"
	"github.com/gogf/gf/test/gtest"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/codec"
	"github.com/osgochina/dmicro/drpc/proto/jsonproto"
	"github.com/osgochina/dmicro/drpc/proto/rawproto"
	"math"
	"testing"
	"time"
)

func TestListenAddr(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		endpointConfig := &drpc.EndpointConfig{
			Network:    "tcp",
			ListenIP:   "0.0.0.0",
			ListenPort: 9091,
			LocalIP:    "127.0.0.1",
			LocalPort:  9092,
		}
		addr := endpointConfig.ListenAddr()
		t.Assert(addr.Network(), "tcp")
		t.Assert(addr.String(), "0.0.0.0:9091")

		addr = endpointConfig.LocalAddr()
		t.Assert(addr.Network(), "tcp")
		t.Assert(addr.String(), "127.0.0.1:9092")
	})
}

func TestLocalAddr(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		endpointConfig := &drpc.EndpointConfig{
			LocalIP:   "127.0.0.1",
			LocalPort: 9092,
		}
		addr := endpointConfig.LocalAddr()
		t.Assert(addr.String(), "127.0.0.1:9092")
		//为传入ListenIP的时候，默认使用LocalIP的值
		addr = endpointConfig.ListenAddr()
		t.Assert(addr.String(), "127.0.0.1:0")
	})
}

func TestSlowCometDuration(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		endpointConfigV2 := &drpc.EndpointConfig{
			SlowCometDuration: 10 * time.Second,
		}
		//为了调用check方法而已
		_ = endpointConfigV2.LocalAddr()
		t.Assert(endpointConfigV2.SlowCometDuration, 10*time.Second)
	})
}

func TestDefaultBodyCodec(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		endpointConfig := &drpc.EndpointConfig{}
		//为了调用check方法而已
		_ = endpointConfig.LocalAddr()
		t.Assert(endpointConfig.DefaultBodyCodec, new(codec.JSONCodec).Name())
		err := drpc.SetDefaultBodyCodec(codec.XmlName)
		t.Assert(err, nil)
		t.Assert(drpc.DefaultBodyCodec().ID(), codec.XmlId)
		endpointConfigV2 := &drpc.EndpointConfig{}
		//为了调用check方法而已
		_ = endpointConfigV2.LocalAddr()
		t.Assert(endpointConfigV2.DefaultBodyCodec, new(codec.XMLCodec).Name())
	})
}

func TestDefaultProtoFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(drpc.DefaultProtoFunc(), rawproto.RawProtoFunc)
		drpc.SetDefaultProtoFunc(jsonproto.NewJSONProtoFunc())
		protoFunc := drpc.DefaultProtoFunc()
		id, name := protoFunc(bufio.ReadWriter{}).Version()
		id2, name2 := jsonproto.NewJSONProtoFunc()(bufio.ReadWriter{}).Version()
		t.Assert(id, id2)
		t.Assert(name, name2)
	})
}

func TestGetReadLimit(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(drpc.GetReadLimit(), math.MaxUint32)
		drpc.SetReadLimit(1024 * 1024 * 2)
		t.Assert(drpc.GetReadLimit(), 1024*1024*2)
	})
}

func TestSocketReadBuffer(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		size, isDefault := drpc.SocketReadBuffer()
		t.Assert(size, -1)
		t.Assert(isDefault, true)
		drpc.SetSocketReadBuffer(1024 * 1024 * 8)
		drpc.SetSocketReadBuffer(-1)
		size, isDefault = drpc.SocketReadBuffer()
		t.Assert(size, 1024*1024*8)
		t.Assert(isDefault, false)

	})
}

func TestSocketWriteBuffer(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		size, isDefault := drpc.SocketWriteBuffer()
		t.Assert(size, -1)
		t.Assert(isDefault, true)
		drpc.SetSocketWriteBuffer(1024 * 1024 * 8)
		drpc.SetSocketWriteBuffer(-1)
		size, isDefault = drpc.SocketWriteBuffer()
		t.Assert(size, 1024*1024*8)
		t.Assert(isDefault, false)

	})
}

func TestTryOptimize(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		drpc.SetSocketKeepAlive(true)
		drpc.SetSocketKeepAlivePeriod(-1)
		drpc.SetSocketKeepAlivePeriod(10 * time.Second)
		drpc.SetSocketNoDelay(true)
	})
}

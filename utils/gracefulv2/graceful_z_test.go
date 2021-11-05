package gracefulv2

import (
	"github.com/gogf/gf/test/gtest"
	"testing"
)

func TestGraceful_PushParentAddr(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		graceful := NewGraceful(GracefulChangeProcess)
		graceful.PushParentAddr("tcp", "127.0.0.1", "127.0.0.1:8399")
		graceful.PushParentAddr("tcp", "127.0.0.1", "127.0.0.1:0")
		graceful.SetParentListenAddrList()
	})
}

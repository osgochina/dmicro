package server

import (
	"fmt"
	"github.com/gogf/gf/v2/test/gtest"
	"testing"
)

func TestOptions_EndpointConfig(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := NewRpcServer("test_one",
			OptListenAddress("127.0.0.1:8199"),
		)
		cfg := s.Options().EndpointConfig()
		fmt.Println(cfg)
	})
}

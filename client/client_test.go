package client_test

import (
	"context"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/osgochina/dmicro/client"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/message"
	"github.com/osgochina/dmicro/drpc/plugin/ignorecase"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/registry"
	"github.com/osgochina/dmicro/server"
	"sync"
	"testing"
	"time"
)

type Math struct {
	drpc.CallCtx
}

func (that *Math) Add(arg *[]int) (int, *drpc.Status) {
	logger.Infof(context.TODO(), "author: %s", that.PeekMeta("author"))
	var r int
	for _, a := range *arg {
		r += a
	}
	return r, nil
}

var once sync.Once

func getServer(serverName string, addr string) {
	once.Do(func() {
		svr := server.NewRpcServer(serverName,
			server.OptListenAddress(addr),
			server.OptPrintDetail(true),
			server.OptEnableHeartbeat(true),
			server.OptGlobalPlugin(ignorecase.NewIgnoreCase()),
		)
		svr.RouteCall(new(Math))
		_ = svr.ListenAndServe()
	})
}

func TestNewRpcClientDefault(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		serviceName := "testregistry"
		addr := "127.0.0.1:9091"
		go func() {
			getServer(serviceName, addr)
		}()
		time.Sleep(1 * time.Second)
		cli := client.NewRpcClient(serviceName,
			client.OptHeartbeatTime(3*time.Second),
		)
		var result int
		stat := cli.Call("/math/Add", []int{1, 2, 3, 4, 5}, &result,
			message.WithSetMeta("author", "clownfish"),
		).Status()
		if !stat.OK() {
			logger.Fatalf(context.TODO(), "%v", stat)
		}
		t.Assert(stat.OK(), true)
		t.Assert(result, 15)
	})
}

func TestNewRpcClientMDNS(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		serviceName := "testregistry"
		addr := "127.0.0.1:9091"
		go func() {
			getServer(serviceName, addr)
		}()
		time.Sleep(1 * time.Second)
		cli := client.NewRpcClient(serviceName, client.OptRegistry(registry.DefaultRegistry))
		//cli := client.NewRpcClient(serviceName, client.OptRegistry(registry.NewRegistry()))
		var result int
		stat := cli.Call("/math/add", []int{1, 2, 3, 4, 5}, &result,
			message.WithSetMeta("author", "clownfish"),
		).Status()
		if !stat.OK() {
			logger.Fatalf(context.TODO(), "%v", stat)
		}
		t.Assert(stat.OK(), true)
		t.Assert(result, 15)
	})
}

func TestNewRpcClientMemory(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		serviceName := "testregistry"
		addr := "127.0.0.1:9091"
		go func() {
			getServer(serviceName, addr)
		}()
		time.Sleep(1 * time.Second)
		s := &registry.Service{
			Nodes: []*registry.Node{
				{
					Address: "127.0.0.1:9091",
				},
			},
		}
		cli := client.NewRpcClient(serviceName,
			client.OptServiceVersion("1.0.1"),
			client.OptCustomService(s),
		)
		var result int
		stat := cli.Call("/math/add", []int{1, 2, 3, 4, 5}, &result,
			message.WithSetMeta("author", "clownfish"),
		).Status()
		if !stat.OK() {
			logger.Fatalf(context.TODO(), "%v", stat)
		}
		t.Assert(stat.OK(), true)
		t.Assert(result, 15)
	})
}

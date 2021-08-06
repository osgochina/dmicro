package heartbeat_test

import (
	"github.com/gogf/gf/test/gtest"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/plugin/heartbeat"
	"testing"
	"time"
)

func TestHeartbeatCALl(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		svr := drpc.NewEndpoint(drpc.EndpointConfig{
			ListenPort:  9090,
			PrintDetail: true,
		}, heartbeat.NewPong())
		go svr.ListenAndServe()

		time.Sleep(time.Second)

		cli := drpc.NewEndpoint(drpc.EndpointConfig{PrintDetail: true}, heartbeat.NewPing(3, true))
		cli.Dial(":9090")
		time.Sleep(time.Second * 20)

	})
}

func TestHeartbeatCALl2(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		svr := drpc.NewEndpoint(drpc.EndpointConfig{
			ListenPort:  9090,
			PrintDetail: true,
		}, heartbeat.NewPong())
		go svr.ListenAndServe()

		time.Sleep(time.Second)

		cli := drpc.NewEndpoint(drpc.EndpointConfig{PrintDetail: true}, heartbeat.NewPing(3, true))
		sess, _ := cli.Dial(":9090")
		for i := 0; i < 8; i++ {
			sess.Call("/", nil, nil).Status()
			time.Sleep(time.Second)
		}
		time.Sleep(time.Second * 10)
	})
}

func TestHeartbeatPush1(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		srv := drpc.NewEndpoint(
			drpc.EndpointConfig{ListenPort: 9090, PrintDetail: true},
			heartbeat.NewPing(3, false),
		)
		go srv.ListenAndServe()
		time.Sleep(time.Second)

		cli := drpc.NewEndpoint(
			drpc.EndpointConfig{PrintDetail: true},
			heartbeat.NewPong(),
		)
		cli.Dial(":9090")
		time.Sleep(time.Second * 10)
	})

}

func TestHeartbeatPush2(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		srv := drpc.NewEndpoint(
			drpc.EndpointConfig{ListenPort: 9090, PrintDetail: true},
			heartbeat.NewPing(3, false),
		)
		go srv.ListenAndServe()
		time.Sleep(time.Second)

		cli := drpc.NewEndpoint(
			drpc.EndpointConfig{PrintDetail: true},
			heartbeat.NewPong(),
		)
		sess, _ := cli.Dial(":9090")
		for i := 0; i < 8; i++ {
			sess.Push("/", nil)
			time.Sleep(time.Second)
		}
		time.Sleep(time.Second * 5)
	})

}

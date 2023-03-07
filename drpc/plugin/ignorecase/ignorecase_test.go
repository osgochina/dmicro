package ignorecase_test

import (
	"context"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/internal"
	"github.com/osgochina/dmicro/drpc/message"
	"github.com/osgochina/dmicro/drpc/plugin/ignorecase"
	"testing"
	"time"
)

type Home struct {
	drpc.CallCtx
}

func (that *Home) Test(arg *map[string]string) (map[string]interface{}, *drpc.Status) {
	that.Session().Push("/push/test", map[string]string{
		"your_id": gconv.String(that.PeekMeta("peer_id")),
	})
	return map[string]interface{}{
		"arg": *arg,
	}, nil
}

func TestIgnoreCase(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		// Server
		srv := drpc.NewEndpoint(drpc.EndpointConfig{ListenPort: 9087, Network: "tcp"}, ignorecase.NewIgnoreCase())
		srv.RouteCall(new(Home))
		go srv.ListenAndServe()
		time.Sleep(1e9)

		// Client
		cli := drpc.NewEndpoint(drpc.EndpointConfig{Network: "tcp"}, ignorecase.NewIgnoreCase())
		cli.RoutePush(new(Push))
		sess, stat := cli.Dial(":9087")
		if !stat.OK() {
			t.Fatal(stat)
		}
		var result interface{}
		stat = sess.Call("/home/TesT",
			map[string]string{
				"author": "clownfish",
			},
			&result,
			message.WithSetMeta("peer_id", "110"),
		).Status()
		if !stat.OK() {
			t.Error(stat)
		}
		t.Logf("result:%v", result)
		time.Sleep(3e9)
	})

}

type Push struct {
	drpc.PushCtx
}

func (p *Push) Test(arg *map[string]string) *drpc.Status {
	internal.Infof(context.TODO(), "receive push(%s):\narg: %#v\n", p.IP(), arg)
	return nil
}

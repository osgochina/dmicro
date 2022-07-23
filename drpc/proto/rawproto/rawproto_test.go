package rawproto_test

import (
	"github.com/gogf/gf/util/gconv"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/internal"
	"github.com/osgochina/dmicro/drpc/tfilter"
	"testing"
	"time"
)

func TestRawProto(t *testing.T) {
	tfilter.Reg('g', "gizp-5", 5)

	// server
	srv := drpc.NewEndpoint(drpc.EndpointConfig{ListenPort: 9090})
	srv.RouteCall(new(Home))
	go srv.ListenAndServe()
	time.Sleep(1e9)

	// client
	cli := drpc.NewEndpoint(drpc.EndpointConfig{})
	cli.RoutePush(new(Push))
	sess, stat := cli.Dial(":9090")
	if !stat.OK() {
		t.Fatal(stat)
	}
	var result interface{}
	stat = sess.Call("/home/test",
		map[string]string{
			"author": "osgochina@gmail.com",
		},
		&result,
		drpc.WithSetMeta("endpoint_id", "110"),
		drpc.WithXFerPipe('g'),
	).Status()
	if !stat.OK() {
		t.Error(stat)
	}
	t.Logf("result:%v", result)
	time.Sleep(3e9)
}

type Home struct {
	drpc.CallCtx
}

func (h *Home) Test(arg *map[string]string) (map[string]interface{}, *drpc.Status) {
	h.Session().Push("/push/test", map[string]string{
		"your_id": gconv.String(h.PeekMeta("endpoint_id")),
	})
	return map[string]interface{}{
		"arg": *arg,
	}, nil
}

type Push struct {
	drpc.PushCtx
}

func (p *Push) Test(arg *map[string]string) *drpc.Status {
	internal.Infof("receive push(%s):\narg: %#v\n", p.IP(), arg)
	return nil
}

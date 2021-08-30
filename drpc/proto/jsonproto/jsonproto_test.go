package jsonproto_test

import (
	"github.com/gogf/gf/util/gconv"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/message"
	"github.com/osgochina/dmicro/drpc/proto/jsonproto"
	"github.com/osgochina/dmicro/drpc/tfilter/gzip"
	"github.com/osgochina/dmicro/logger"
	"testing"
	"time"
)

type Home struct {
	drpc.CallCtx
}

func (h *Home) Test(arg *map[string]string) (map[string]interface{}, *drpc.Status) {
	h.Session().Push("/push/test", map[string]string{
		"your_id": gconv.String(h.PeekMeta("peer_id")),
	})
	return map[string]interface{}{
		"arg": *arg,
	}, nil
}

func TestJSONProto(t *testing.T) {
	gzip.Reg('g', "gizp-5", 5)

	// Server
	srv := drpc.NewEndpoint(drpc.EndpointConfig{ListenPort: 9090})
	srv.RouteCall(new(Home))
	go srv.ListenAndServe(jsonproto.NewJSONProtoFunc())
	time.Sleep(1e9)

	// Client
	cli := drpc.NewEndpoint(drpc.EndpointConfig{})
	cli.RoutePush(new(Push))
	sess, stat := cli.Dial(":9090", jsonproto.NewJSONProtoFunc())
	if !stat.OK() {
		t.Fatal(stat)
	}
	var result interface{}
	stat = sess.Call("/home/test",
		map[string]string{
			"author": "henrylee2cn",
		},
		&result,
		message.WithSetMeta("peer_id", "110"),
		message.WithXFerPipe('g'),
	).Status()
	if !stat.OK() {
		t.Error(stat)
	}
	t.Logf("result:%v", result)
	time.Sleep(3e9)
}

type Push struct {
	drpc.PushCtx
}

func (p *Push) Test(arg *map[string]string) *drpc.Status {
	logger.Infof("receive push(%s):\narg: %#v\n", p.IP(), arg)
	return nil
}

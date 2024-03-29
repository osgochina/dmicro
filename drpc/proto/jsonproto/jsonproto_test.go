package jsonproto_test

import (
	"context"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/internal"
	"github.com/osgochina/dmicro/drpc/message"
	"github.com/osgochina/dmicro/drpc/proto/jsonproto"
	"github.com/osgochina/dmicro/drpc/tfilter"
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

func (h *Home) RetString(arg *map[string]string) (string, *drpc.Status) {
	return "RetString", nil
}

func TestJSONProto(t *testing.T) {
	tfilter.RegGzip(5)

	// Server
	srv := drpc.NewEndpoint(drpc.EndpointConfig{ListenPort: 9092})
	srv.RouteCall(new(Home))
	go srv.ListenAndServe(jsonproto.NewJSONProtoFunc())
	time.Sleep(1e9)

	// Client
	cli := drpc.NewEndpoint(drpc.EndpointConfig{})
	cli.RoutePush(new(Push))
	sess, stat := cli.Dial(":9092", jsonproto.NewJSONProtoFunc())
	if !stat.OK() {
		t.Fatal(stat)
	}
	var result interface{}
	stat = sess.Call("/home/test",
		map[string]string{
			"author": "osgochina@gmail.com",
		},
		&result,
		message.WithSetMeta("endpoint_id", "110"),
		message.WithTFilterPipe(tfilter.GzipId),
	).Status()
	if !stat.OK() {
		t.Error(stat)
	}
	t.Logf("result:%v", result)
	time.Sleep(3e9)

	var resultString string
	stat = sess.Call("/home/ret_string",
		map[string]string{
			"author": "osgochina@gmail.com",
		},
		&resultString,
	).Status()
	if !stat.OK() {
		t.Error(stat)
	}
	t.Logf("resultString:%v", resultString)
	time.Sleep(3e9)
}

type Push struct {
	drpc.PushCtx
}

func (p *Push) Test(arg *map[string]string) *drpc.Status {
	internal.Infof(context.TODO(), "receive push(%s):\narg: %#v\n", p.IP(), arg)
	return nil
}

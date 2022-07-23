package pbproto_test

import (
	"github.com/gogf/gf/util/gconv"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/codec"
	"github.com/osgochina/dmicro/drpc/internal"
	"github.com/osgochina/dmicro/drpc/proto/pbproto"
	"github.com/osgochina/dmicro/drpc/proto/pbproto/pb_test"
	"github.com/osgochina/dmicro/drpc/tfilter"
	"testing"
	"time"
)

type Home struct {
	drpc.CallCtx
}

func (that *Home) Test(arg *pb_test.Request) (*pb_test.Response, *drpc.Status) {
	that.Session().Push("/push/test", &pb_test.Push{
		PeerId: gconv.Int32(that.PeekMeta("peer_id")),
	})
	return &pb_test.Response{
		Author: arg.GetAuthor(),
		Uid:    arg.GetUid(),
		Email:  arg.GetEmail(),
		Phone:  arg.GetPhone(),
	}, nil
}

type Push struct {
	drpc.PushCtx
}

func (that *Push) Test(arg *pb_test.Push) *drpc.Status {
	internal.Infof("receive push(%s):\narg: %#v\n", that.IP(), arg.PeerId)
	return nil
}

func TestPbProto(t *testing.T) {
	tfilter.Reg('g', "gizp-5", 5)

	// server
	srv := drpc.NewEndpoint(drpc.EndpointConfig{ListenPort: 9090, DefaultBodyCodec: codec.NameProtobuf})
	srv.RouteCall(new(Home))
	go srv.ListenAndServe(pbproto.NewPbProtoFunc())
	time.Sleep(1e9)
	// client
	cli := drpc.NewEndpoint(drpc.EndpointConfig{DefaultBodyCodec: codec.NameProtobuf})
	cli.RoutePush(new(Push))

	sess, stat := cli.Dial(":9090", pbproto.NewPbProtoFunc())
	if !stat.OK() {
		t.Fatal(stat)
	}
	var result pb_test.Response
	stat = sess.Call("/home/test",
		&pb_test.Request{
			Author: "liuzhiming",
			Uid:    100,
		},
		&result,
		drpc.WithSetMeta("peer_id", "110"),
		drpc.WithXFerPipe('g'),
	).Status()
	if !stat.OK() {
		t.Error(stat)
	}
	t.Logf("result:%v", &result)
	time.Sleep(3e9)
}

### 概述

pbproto实现了`ProtoBuf`协议。

### 消息的格式说明

`{length bytes}` `{xferPipe length byte}` `{xferPipe bytes}` `{protobuf bytes}`

- `{length bytes}`: uint32, 4 bytes, big endian
- `{xferPipe length byte}`: 1 byte
- `{xferPipe bytes}`: one byte one xfer
- `{protobuf bytes}`: {"seq":%d,"mtype":%d,"serviceMethod":%q,"status":%q,"meta":%q,"bodyCodec":%d,"body":"%s"}

### 如何引用

`import "github.com/osgochina/dmicro/drpc/proto/pbproto"`

#### 使用示例

```go
package pbproto_test

import (
	"github.com/gogf/gf/util/gconv"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/codec"
	"github.com/osgochina/dmicro/drpc/proto/pbproto"
	"github.com/osgochina/dmicro/drpc/proto/pbproto/pb_test"
	"github.com/osgochina/dmicro/drpc/tfilter/gzip"
	"github.com/osgochina/dmicro/logger"
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
	logger.Infof("receive push(%s):\narg: %#v\n", that.IP(), arg.PeerId)
	return nil
}

func TestPbProto(t *testing.T) {
	gzip.Reg('g', "gizp-5", 5)

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

```


执行以下命令测试:

```sh
$ cd dmicro/drpc/proto/pbproto
$ go test -v -run=TestPbProto
```
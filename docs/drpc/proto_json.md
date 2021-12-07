### 概述

jsonproto是实现JSON套接字通信协议。

### 消息的格式说明

`{length bytes}` `{xferPipe length byte}` `{xferPipe bytes}` `{JSON bytes}`

- `{length bytes}`: uint32, 4 bytes, big endian
- `{xferPipe length byte}`: 1 byte
- `{xferPipe bytes}`: one byte one xfer
- `{JSON bytes}`: {"seq":%d,"mtype":%d,"serviceMethod":%q,"status":%q,"meta":%q,"bodyCodec":%d,"body":"%s"}

### 如何引用

`import "github.com/osgochina/dmicro/drpc/proto/jsonproto"`

#### 使用示例

```go
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
			"author": "osgochina@gmail.com",
		},
		&result,
		message.WithSetMeta("endpoint_id", "110"),
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


```

执行以下命令测试:

```sh
$ cd dmicro/drpc/proto/jsonproto
$ go test -v -run=TestJSONProto
```

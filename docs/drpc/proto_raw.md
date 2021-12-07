### 概述

Rawproto 是默认的帧拼接格式协议。

### 协议发送格式

raw protocol format(Big Endian):

```sh
{4 bytes message length}
{1 byte protocol version} # 6
{1 byte transfer pipe length}
{transfer pipe IDs}
# The following is handled data by transfer pipe
{1 bytes sequence length}
{sequence (HEX 36 string of int32)}
{1 byte message type} # e.g. CALL:1; REPLY:2; PUSH:3
{1 bytes service method length}
{service method}
{2 bytes status length}
{status(urlencoded)}
{2 bytes metadata length}
{metadata(urlencoded)}
{1 byte body codec id}
{body}
```

### 如何引入

`import "github.com/osgochina/dmicro/drpc/proto/rawproto"`

#### 使用示例

```go
package rawproto_test

import (
	"github.com/gogf/gf/util/gconv"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/tfilter/gzip"
	"github.com/osgochina/dmicro/logger"
	"testing"
	"time"
)

func TestRawProto(t *testing.T) {
	gzip.Reg('g', "gizp-5", 5)

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
	logger.Infof("receive push(%s):\narg: %#v\n", p.IP(), arg)
	return nil
}

```

执行以下命令测试:

```sh
$ cd dmicro/drpc/proto/rawproto
$ go test -v -run=TestRawProto
```

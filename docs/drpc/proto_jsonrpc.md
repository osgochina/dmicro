### 概述

jsonrpc协议是实现`JSONRPC`标准的套接字通信协议。
注意，因为`JSONRPC`标准中并未有push逻辑，所以使用`JSONRPC`只能使用call，reply模式，也就是`请求-应答`模式。

### 消息的格式说明

参考[jsonrpc标准文档](https://www.w3cschool.cn/ycuott/z7er3ozt.html)

### 如何引用

`import "github.com/osgochina/dmicro/drpc/proto/jsonproto"`

#### 使用示例

```go
package jsonrpcproto_test

import (
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/message"
	"github.com/osgochina/dmicro/drpc/proto/jsonrpcproto"
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

func TestJSONRPCProto(t *testing.T) {
	//gzip.Reg('g', "gizp-5", 5)

	// Server
	srv := drpc.NewEndpoint(drpc.EndpointConfig{ListenPort: 9090})
	srv.RouteCall(new(Home))
	go srv.ListenAndServe(jsonrpcproto.NewJSONRPCProtoFunc())
	time.Sleep(1e9)

	// Client
	cli := drpc.NewEndpoint(drpc.EndpointConfig{})
	cli.RoutePush(new(Push))
	sess, stat := cli.Dial(":9090", jsonrpcproto.NewJSONRPCProtoFunc())
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
	logger.Infof(context.TODO(),"receive push(%s):\narg: %#v\n", p.IP(), arg)
	return nil
}


```

执行以下命令测试:

```sh
$ cd dmicro/drpc/proto/jsonproto
$ go test -v -run=TestJSONProto
```

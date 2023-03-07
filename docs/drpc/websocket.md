# WebSocket支持

`DRpc` 不只是支持服务器之间的通讯，还支持浏览器与服务器之间通过`websocket`协议进行`rpc`通讯.

并且默认支持`json`，`protobuf`协议。

## 启动`WebSocket`的支持

简单示例：
```go
package main

import (
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/mixer/websocket"
	"github.com/osgochina/dmicro/logger"
	"time"
)

type Arg struct {
	A int
	B int `param:"<range:1:>"`
}

type P struct {
	drpc.CallCtx
}

func (p *P) Divide(arg *Arg) (int, *drpc.Status) {
	return arg.A / arg.B, nil
}

func main() {
	srv := websocket.NewServer("/", drpc.EndpointConfig{ListenPort: 9090})
	srv.RouteCall(new(P))
	go srv.ListenAndServe()
	time.Sleep(time.Second * 1)
	cli := websocket.NewClient("/", drpc.EndpointConfig{})
	sess, stat := cli.Dial(":9090")
	if !stat.OK() {
		logger.Fatal(context.TODO(),stat)
	}
	var result int
	stat = sess.Call("/p/divide", &Arg{
		A: 10,
		B: 2,
	}, &result,
	).Status()
	if !stat.OK() {
		logger.Fatal(context.TODO(),stat)
	}
	logger.Println(context.TODO(),"10/2=%d", result)
	time.Sleep(time.Second)
}
```

还支持启动http服务后附加到`drpc服务`:

```go
package main

import (
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/mixer/websocket"
	"github.com/osgochina/dmicro/drpc/mixer/websocket/jsonSubProto"
	"github.com/osgochina/dmicro/logger"
	"net/http"
	"time"
)

type Arg struct {
	A int
	B int `param:"<range:1:>"`
}

type P struct {
	drpc.CallCtx
}

func (p *P) Divide(arg *Arg) (int, *drpc.Status) {
	return arg.A / arg.B, nil
}

func main() {
	srv := drpc.NewEndpoint(drpc.EndpointConfig{})
	http.Handle("/token", websocket.NewJSONServeHandler(srv, nil))
	go http.ListenAndServe(":9094", nil)
	srv.RouteCall(new(P))
	time.Sleep(time.Millisecond * 200)

	// example in Browser: ws://localhost/token?access_token=clientAuthInfo
	cli := drpc.NewEndpoint(drpc.EndpointConfig{}, websocket.NewDialPlugin("/token"))
	sess, stat := cli.Dial(":9094", jsonSubProto.NewJSONSubProtoFunc())
	if !stat.OK() {
		logger.Fatal(context.TODO(),stat)
	}
	var result int
	stat = sess.Call("/p/divide", &Arg{
		A: 10,
		B: 2,
	}, &result,
	).Status()
	if !stat.OK() {
		logger.Fatal(context.TODO(),stat)
	}
	logger.Println(context.TODO(),"10/2=%d", result)
	time.Sleep(time.Millisecond * 200)
}

```


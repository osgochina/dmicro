# 并发请求客户端

使用连接池创建session，满足并发请求的场景。


```go
package main

import (
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/mixer/multiclient"
	"github.com/osgochina/dmicro/logger"
	"time"
)

type Arg struct {
	A int
	B int `param:"<range:1:>"`
}

type P struct{ drpc.CallCtx }

func (p *P) Divide(arg *Arg) (int, *drpc.Status) {
	return arg.A / arg.B, nil
}

func main() {
	srv := drpc.NewEndpoint(drpc.EndpointConfig{
		ListenPort: 9090,
	})
	srv.RouteCall(new(P))
	go srv.ListenAndServe()
	time.Sleep(time.Second)

	cli := multiclient.New(
		drpc.NewEndpoint(drpc.EndpointConfig{}),
		":9090",
		time.Second*5,
	)
	go func() {
		for {
			logger.Printf("%d", cli.Size())
			time.Sleep(time.Millisecond * 500)
		}
	}()
	go func() {
		var result int
		for i := 0; ; i++ {
			stat := cli.Call("/p/divide", &Arg{
				A: i,
				B: 2,
			}, &result).Status()
			if !stat.OK() {
				logger.Print(stat)
			} else {
				logger.Printf("%d/2=%v", i, result)
			}
			time.Sleep(time.Millisecond * 500)
		}
	}()
	time.Sleep(time.Second * 6)
	cli.Close()
	time.Sleep(time.Second * 3)
}

```
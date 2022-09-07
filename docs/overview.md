# 快速开始

## rpc服务端
如何快速的通过简单的代码创建一个真正的rpc服务。
以下就是示例代码：
```go
package main

import (
	"fmt"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/dserver"
	"github.com/osgochina/dmicro/logger"
)
// DRpcSandBox  默认的服务
type DRpcSandBox struct {
	dserver.BaseSandbox
	endpoint drpc.Endpoint
}

func (that *DRpcSandBox) Name() string {
	return "DRpcSandBox"
}

func (that *DRpcSandBox) Setup() error {
	fmt.Println("DRpcSandBox Setup")
	cfg := that.Config.EndpointConfig(that.Name())
	cfg.ListenPort = 9091
	cfg.PrintDetail = true
	that.endpoint = drpc.NewEndpoint(cfg)
	that.endpoint.RouteCall(new(Math))
	return that.endpoint.ListenAndServe()
}

func (that *DRpcSandBox) Shutdown() error {
	fmt.Println("DRpcSandBox Shutdown")
	return that.endpoint.Close()
}


// Math rpc请求的最终处理器，必须集成drpc.CallCtx
type Math struct {
	drpc.CallCtx
}

func (m *Math) Add(arg *[]int) (int, *drpc.Status) {
	// test meta
	logger.Infof("author: %s", m.PeekMeta("author"))
	// add
	var r int
	for _, a := range *arg {
		r += a
	}
	// response
	return r, nil
}

func main() {
	dserver.Authors = "osgochina@gmail.com"
	dserver.SetName("DMicro_drpc")
	dserver.Setup(func(svr *dserver.DServer) {
		err := svr.AddSandBox(new(DRpcSandBox))
		if err != nil {
			logger.Fatal(err)
		}
	})
}

```


# rpc客户端

服务已经建立完毕，如何通过client链接它呢？

```go

package main

import (
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/message"
	"github.com/osgochina/dmicro/logger"
	"time"
)

func main() {

	cli := drpc.NewEndpoint(drpc.EndpointConfig{PrintDetail: true, RedialTimes: -1, RedialInterval: time.Second})
	defer cli.Close()
	

	sess, stat := cli.Dial("127.0.0.1:9091")
	if !stat.OK() {
		logger.Fatalf("%v", stat)
	}
    var result int
    stat = sess.Call("/math/add",
        []int{1, 2, 3, 4, 5},
        &result,
        message.WithSetMeta("author", "liuzhiming"),
    ).Status()
    if !stat.OK() {
        logger.Fatalf("%v", stat)
    }
    logger.Printf("result: %d", result)
}
```
通过以上的代码事例，大家基本可以了解`drpc`框架是怎么使用。

## rpc服务

如何快速的通过简单的代码创建一个真正的rpc服务。
以下就是示例代码：
```go
package main

import (
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/logger"
)

func main() {
	//开启信号监听
	go drpc.GraceSignal()

	// 创建一个rpc服务
	svr := drpc.NewEndpoint(drpc.EndpointConfig{
		CountTime:   true,
		LocalIP:     "127.0.0.1",
		ListenPort:  9091,
		PrintDetail: true,
	})
	//注册处理方法
	svr.RouteCall(new(Math))

	err := svr.ListenAndServe()
	logger.Error(err)
}

type Math struct {
	drpc.CallCtx
}

// Add 数据方法，把传入的参数累加，把结果返回
func (m *Math) Add(arg *[]int) (int, *drpc.Status) {
	// 查看传入的元数据
	logger.Infof("author: %s", m.PeekMeta("author"))
	// add
	var r int
	for _, a := range *arg {
		r += a
	}
	// response
	return r, nil
}
```

## rpc客户端

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

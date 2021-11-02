# dmicro [![GitHub release](https://img.shields.io/github/v/release/osgochina/dmicro.svg?style=flat-square)](https://github.com/osgochina/dmicro/releases) [![report card](https://goreportcard.com/badge/github.com/osgochina/dmicro?style=flat-square)](http://goreportcard.com/report/osgochina/dmicro) [![github issues](https://img.shields.io/github/issues/osgochina/dmicro.svg?style=flat-square)](https://github.com/osgochina/dmicro/issues?q=is%3Aopen+is%3Aissue) [![github closed issues](https://img.shields.io/github/issues-closed-raw/osgochina/dmicro.svg?style=flat-square)](https://github.com/osgochina/dmicro/issues?q=is%3Aissue+is%3Aclosed) [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/osgochina/dmicro) [![view examples](https://img.shields.io/badge/learn%20by-examples-00BCD4.svg?style=flat-square)](https://github.com/osgochina/dmicro/tree/main/.examples)
## dmicro简介

> dmicro是一个高效、可扩展且简单易用的微服务框架。包含drpc,easyserver等组件。

该项目的诞生离不开`erpc`和`GoFrame`两个优秀的项目。

其中`drpc`组件参考`erpc`项目的架构思想，依赖的基础库是`GoFrame`。

* [erpc](https://gitee.com/henrylee/erpc)
* [GoFrame](https://gitee.com/johng/gf)

## 详细文档

[中文文档](http://dmicro.clownfish.site/)

## 安装

```html
go get -u -v github.com/osgochina/dmicro
```
推荐使用 `go.mod`:
```
require github.com/osgochina/dmicro latest
```

* import

```go
import "github.com/osgochina/dmicro"
```

## 限制
```shell
golang版本 >= 1.13
```

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
	//启动监听
	err := svr.ListenAndServe()
	logger.Warning(err)
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
通过以上的代码事例，大家基本可以了解`drpc`框架是怎么使用。


# 快速开始

## 最简单的应用
跟随我，创建一个全新的应用吧！

```go
package main

import (
	"fmt"
	"github.com/osgochina/dmicro/dserver"
)
func main() {
    dserver.SetName("DMicro_simple")
    dserver.Setup(func(svr *dserver.DServer) {
        fmt.Println("start success!")
    })
}
```
以上代码可以创建出最简单的应用。

1. `dserver.SetName("DMicro_simple")`

    设置应用名,建议设置独特个性化的引用名，因为管理链接，日志目录等地方会用到它。
    如果不设置，默认是"DServer_xxx",启动xxx为二进制名
2. `dserver.Setup()`
    
    整个`dserver`项目的入口方法,所有的启动逻辑应该写在该方法内。

编译出你的第一个应用吧！
```shell
$ go build main.go
```
编译后你可以通过`main version`,`main help` 来获取使用帮助.
```shell
$ ./main version

  ____    ____                                      
 |  _ \  / ___|    ___   _ __  __   __   ___   _ __ 
 | | | | \___ \   / _ \ | '__| \ \ / /  / _ \ | '__|
 | |_| |  ___) | |  __/ | |     \ V /  |  __/ | |   
 |____/  |____/   \___| |_|      \_/    \___| |_|  
Version:         60c373c
Go Version:      go version go1.18.3 linux/amd64
DMicro Version:  v1.0.0
GF Version:      v1.16.9
Git Commit:      60c373c
Build Time:      2022-07-16 15:42:26
Authors:         osgochina@gmail.com
Install Path:    /path/to/dmicro/examples/dserver/main

```

```shell
$ ./main help
USAGE
  /path/to/server start|stop|reload|quit [OPTION] [sandboxName1|sandboxName2...] 
OPTION
  -c,--config     指定要载入的配置文件，该参数与gf.gcfg.file参数二选一，建议使用该参数
  -d,--daemon     使用守护进程模式启动
  --env           环境变量，表示当前启动所在的环境,有[dev,test,product]这三种，默认是product
  --debug         是否开启debug 默认debug=false
  --pid           设置pid文件的地址，默认是/tmp/[server].pid
  -h,--help       获取帮助信息
  -v,--version    获取编译版本信息
  -m,--model      进程模型，0表示单进程模型，1表示多进程模型    
EXAMPLES
  /path/to/server 
  /path/to/server start --env=dev --debug=true --pid=/tmp/server.pid
  /path/to/server start --gf.gcfg.file=config.product.toml
  /path/to/server start -c=config.product.toml
  /path/to/server start --config=config.product.toml user admin
  /path/to/server start user
  /path/to/server stop
  /path/to/server quit
  /path/to/server reload
  /path/to/server version
  /path/to/server help

```

## 第一个Sandbox

前面我们建立了一个最基础的可以运行`DServer`，让大家能够快速的运行起来，
接下来我们真正建立起可以运行业务的应用。

```go
// main.go
package main

import (
	"fmt"
	"github.com/osgochina/dmicro/dserver"
	"github.com/osgochina/dmicro/logger"
)

// DefaultSandBox  默认的服务
type DefaultSandBox struct {
	dserver.BaseSandbox
}

func (that *DefaultSandBox) Name() string {
	return "DefaultSandBox"
}

func (that *DefaultSandBox) Setup() error {
	fmt.Println("DefaultSandBox Setup")
	return nil
}

func (that *DefaultSandBox) Shutdown() error {
	fmt.Println("DefaultSandBox Shutdown")
	return nil
}

func main() {
	dserver.Authors = "osgochina@gmail.com"
	dserver.SetName("DMicro_foo")
	dserver.Setup(func(svr *dserver.DServer) {
		err := svr.AddSandBox(new(DefaultSandBox))
		if err != nil {
			logger.Fatal(err)
		}
	})
}

```

1. `dserver.Authors` 
    
    设置该应用的开发者，编译后，运行version的时候会显示该信息
2. `err := svr.AddSandBox(new(DefaultSandBox))`

    把`sandbox`对象添加到`DServer`中管理。

3. `Name()` sandbox的name必须在应用中唯一。可以设定常量字符串。
4. `Setup() error ` 应用初始化完毕，会新开一个协程调用该方法，你可以在这里启动`drpc`,`ghttp`等服务.可以阻塞。
5. `Shutdown() error` 进程结束，平滑重启,手动关闭sandbox的时候，会调用该方法，你可以在这里进行业务的收尾动作。

`Sandbox`讲解

后面的文档中有`sandbox`原理的详细讲解，我们先不去管它，只需要知道怎么用就好。
我们自己创建的`sandbox`结构必须实现`dserver.ISandBox`接口。需要实现`Name()`,`Setup() error`,`Shutdown() error`三个方法。
并且必须继承`dserver.BaseSandbox`.

编译应用
```shell
$ go build main.go
```

启动应用
```shell
$ ./main start
```

## `Http`服务

```go
// http.go
package main

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/osgochina/dmicro/dserver"
	"github.com/osgochina/dmicro/logger"
	"os"
)

// HttpSandBox  默认的服务
type HttpSandBox struct {
	dserver.BaseSandbox
	http *ghttp.Server
}

func (that *HttpSandBox) Name() string {
	return "HttpSandBox"
}

func (that *HttpSandBox) Setup() error {
	fmt.Println("HttpSandBox Setup")
	that.http = g.Server("ghttp")
	that.http.BindHandler("/", func(r *ghttp.Request) {
		r.Response.WriteExit("hello world!", "pid:"+gconv.String(os.Getpid()),"\n")
	})
	that.http.SetPort(8080)
	return that.http.Start()
}

func (that *HttpSandBox) Shutdown() error {
	fmt.Println("HttpSandBox Shutdown")
	return that.http.Shutdown()
}

func main() {
	dserver.Authors = "osgochina@gmail.com"
	dserver.SetName("DMicro_http")
	dserver.Setup(func(svr *dserver.DServer) {
		err := svr.AddSandBox(new(HttpSandBox))
		if err != nil {
			logger.Fatal(err)
		}
	})
}

```

1. `that.http = g.Server("ghttp")` 创建http服务对象。
2. `that.http.BindHandler` 绑定指定路径的处理方法。
3. `that.http.SetPort(8080)` 监听8080方法并启动。

编译应用
```shell
$ go build http.go
```

启动应用
```shell
$ ./http start
```
访问服务
```shell
$ curl http://127.0.0.1:8080/

hello world!pid:775581
```


## `DRPC`服务

```go
// drpc.go
package main

import (
	"fmt"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/dserver"
	"github.com/osgochina/dmicro/logger"
)

func Home(ctx drpc.CallCtx, args *struct{}) (string, *drpc.Status) {
	return "home", nil
}

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
	cfg.ListenPort = 8199
	that.endpoint = drpc.NewEndpoint(cfg)
	that.endpoint.SubRoute("/app").RouteCallFunc(Home)
	return that.endpoint.ListenAndServe()
}

func (that *DRpcSandBox) Shutdown() error {
	fmt.Println("DRpcSandBox Shutdown")
	return that.endpoint.Close()
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

1. `cfg := that.Config.EndpointConfig(that.Name())` 

    通过sandbox对象中的`Config`属性，能够获取当前服务的配置对象。
    注意，`Config`对象是通过反射，由`DServer`注入进去的。
2. `return that.endpoint.ListenAndServe()` 
    
    可以看到`ListenAndServe()`方法是阻塞监听的，放这里没有问题。因为是新开了单独的协程调用`Setup()`.
3. `return that.endpoint.Close()`

    这里是sandbox关闭的时候，把它打开的`endpoint`也关闭。

编译应用
```shell
$ go build drpc.go
```

启动应用
```shell
$ ./drpc start
```

## 多个Sandbox

```go
// drpc_http.go

func main() {
	dserver.Authors = "osgochina@gmail.com"
	dserver.SetName("DMicro_drpc")
	dserver.Setup(func(svr *dserver.DServer) {
		err := svr.AddSandBox(new(DRpcSandBox))
		if err != nil {
			logger.Fatal(err)
		}
		err = svr.AddSandBox(new(HttpSandBox))
		if err != nil {
			logger.Fatal(err)
		}
	})
}
```

以上代码展示了同时载入`DRpcSandBox`和`HttpSandBox`的用法。

`svr.AddSandBox`方法可以载入任意数量的sandbox。

## Sandbox之间调用

```go
func (that *DRpcSandBox) Setup() error {
	sandbox, found := that.Service.SearchSandBox("HttpSandBox")
	if !found {
		return fmt.Errorf("not found HttpSandBox")
	}
	httpSandBox := sandbox.(*HttpSandBox)
	fmt.Println(httpSandBox.Name())
	return nil
}
```

通过`SearchSandBox(name string)`搜索同一个service中的其他sandbox.

## 多个Service


```go
// drpc_http.go

func main() {
	dserver.Authors = "osgochina@gmail.com"
	dserver.SetName("DMicro_drpc")
	dserver.Setup(func(svr *dserver.DServer) {
		err := svr.AddSandBox(new(DRpcSandBox),svr.NewService("rpc"))
		if err != nil {
			logger.Fatal(err)
		}
		err = svr.AddSandBox(new(HttpSandBox),svr.NewService("http"))
		if err != nil {
			logger.Fatal(err)
		}
	})
}
```

调用`svr.AddSandBox`方法的时候，第二个可选参数，支持传入新的`Service`对象。

通过`svr.NewService(name string)` 方法可以创建一个全新的对象。
- 在单进程模式下可以隔离`sandbox`。
- 多进程模式下，每个`Service`为一个独立的进程。

## 开启多进程模式

进程模式有两种开启方法。
1. 通过在代码内调用`svr.ProcessModel(dserver.ProcessModelMulti)`开启多进程模式。
2. 启动应用进程的时候传入`--model=1`开启多进程模式，`--model=0`开启单进程模式。

注意，第二种启动命令`--model`的优先级更高。

## 增加自定义命令

如果你的应用需要增加自己的命令行参数，可以调用`dserver.GrumbleApp()`命令获取app对象，
更多用法可以参考[grumble](github.com/desertbit/grumble)

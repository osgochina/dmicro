# 创建服务

`easyService` 让你能专注于业务，快速的创建服务，非常简单的就拥有功能强大的启动命令行。


## 快速开始

```go
package main

import (
	"fmt"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/osgochina/dmicro/examples/easyservice/sandbox"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/easyservice"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/utils"
	"github.com/osgochina/dmicro/utils/graceful"
)

func Home(ctx drpc.CallCtx, args *struct{}) (string, *drpc.Status) {
	return "home", nil
}

func main() {
	logger.SetDebug(true)
	easyservice.Setup(func(svr *easyservice.EasyService) {
		//注册服务停止时要执行法方法
		svr.BeforeStop(func(service *easyservice.EasyService) bool {
			fmt.Println("server stop")
			return true
		})
		// 获取默认配置文件
		var cfg = easyservice.DefaultBoxConf(svr.CmdParser(), svr.Config())
		// 创建rpc服务沙盒
		rpc := sandbox.NewDefaultSandBox(cfg)
		rpc.Endpoint().SubRoute("/app").RouteCallFunc(Home)
		//把服务沙盒加入service容器
		svr.AddSandBox(rpc)
		
		// 启动一个http服务沙盒
		http := sandbox.NewHttpSandBox(svr)
		svr.AddSandBox(http)
	})
}
```

## 名词解释

* ```EasyService```

    `服务容器`，一个进程可以创建多个服务容器，但是建议使用默认创建的那个。
* ```SandBox```

    `服务沙盒`,每个`服务容器`可以添加多个沙盒，通常这里是按业务区分，比如`admin管理后台`，`api服务`，`rpc服务`，可以放在一个进程中同时启动
    也可以编译到一个进程中，通过命令行参数启动。可以在某个进程中启动一个或多个。

* ```BoxConf```
    
    `沙盒配置`,通过配置文件中的名称，区分不同的沙盒，获取配置。


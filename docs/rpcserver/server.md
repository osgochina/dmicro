# RPC Server

`drpc`原生就支持对外提供服务，为什么还需要`RPC Server`呢？ 通过`drpc`的`endpoint`创建一个可用的的`rpc server`是非常简单的。

但是你如果要使用组件，如`Registry`,`Selector`,`Broker`等等组件，就不太友好了，组合起来会显得特别麻烦。

`RPC Server`作为`DMicro`核心组件的价值就发挥出来了，它起到一个糅合`endpoint`与各个组件的作用。
暴露简单的`Api`,让开发者使用起来简单方便。

## 快速开始

```go
func main() {
	serviceName := "foo"
	serviceVersion := "1.0.0"
	svr := server.NewRpcServer(serviceName,
		server.OptListenAddress("127.0.0.1:9091"),
		server.OptCountTime(true),
		server.OptPrintDetail(true),
		server.OptServiceVersion(serviceVersion),
		server.OptRegistry(reg),
	)
	svr.RouteCall(new(Math))
	_ = svr.ListenAndServe()
}
```

## 支持的配置方法

## 使用组件

# 服务注册中心 (Registry)

微服务架构决定了要运行`业务服务`的数量会非常多，且每个`业务服务`需要多个`Node`来提供高可用的服务。
为了保障`业务服务`的正常运行，需要引入分布式的中间件`服务注册中心`来进行服务的注册，更新。

`服务注册中心`需要要保存服务的名称，版本，节点地址，节点数量，元数据等信息。
当某个`业务服务`状态有变化的时候，需要上报到该`中心`，并由它来进行服务状态更新，通知到服务的调用方。

在`DMicro`中，要使用`Registry`是非常方便的，你既可以使用内置的组件，也可以自己实现`Registry interface`来接入你自己的服务注册中心。

目前`DMicro`已内置`etcd`,`mdns`,`memory`三种服务注册中心。

## 快速开始

### 使用默认的`MDNS`服务注册

```go
func main() {
	serviceName := "testregistry"
	serviceVersion := "1.0.0"
	reg := registry.DefaultRegistry
	err := reg.Init(registry.ServiceName(serviceName), registry.ServiceVersion(serviceVersion))
	if err != nil {
		logger.Fatal(err)
	}
	svr := server.NewRpcServer(serviceName,
		server.OptListenAddress("127.0.0.1:9091"),
		server.OptCountTime(true),
		server.OptPrintDetail(true),
		server.OptRegistry(reg),
	)
	svr.RouteCall(new(Math))
	_ = svr.ListenAndServe()
}
```
`registry.DefaultRegistry` 默认的服务注册中心使用的是`mdns`,这样你就可以在本地局域网直接通过`rpc client`请求到该服务。

服务注册的是要需要传入`serviceName`,`serviceVersion`来区分你的服务。
节点信息`RPC Server`会自动补全。

当然，你也可以不显示的传入要使用的`Registry`,通过`RPC Server`自动生成，这样你写的服务就更简洁了。

```go
func main() {
	serviceName := "testregistry"
	serviceVersion := "1.0.0"
	svr := server.NewRpcServer(serviceName,
		server.OptListenAddress("127.0.0.1:9091"),
		server.OptCountTime(true),
		server.OptPrintDetail(true),
		server.OptServiceVersion(serviceVersion),
	)
	svr.RouteCall(new(Math))
	_ = svr.ListenAndServe()
}
```

### 使用`ETCD`作为服务注册中心

```go
reg := etcd.NewRegistry(
    registry.ServiceName(serviceName),
    registry.ServiceVersion(serviceVersion),
    registry.AddrList("127.0.0.1:12379", "127.0.0.1:22379", "127.0.0.1:32379"),
    registry.LeasesInterval(10*time.Second),
    etcd.RegisterTTL(20*time.Second),
)
svr := server.NewRpcServer(serviceName,
    server.OptServiceVersion(serviceVersion),
    server.OptRegistry(reg),
    server.OptListenAddress("127.0.0.1:9091"),
    server.OptCountTime(true),
    server.OptPrintDetail(true),
)
svr.RouteCall(new(Math))
_ = svr.ListenAndServe()
```

`etcd.NewRegistry` 创建etcd的注册中心，参数是创建etcd注册中心所需要的参数，返回对象本身，
再把它作为参数传给`RPC Server`,启动服务后，会自动把当前服务的节点信息注册到`etcd`,并且会定时上报服务状态。


## 实现原理

只需要实现`Registry interface`就能定义自己的服务注册中心。

```go
type Registry interface {
	Init(...Option) error
	Options() Options
	Register(*Service, ...RegisterOption) error
	Deregister(*Service, ...DeregisterOption) error
	GetService(string, ...GetOption) ([]*Service, error)
	ListServices(...ListOption) ([]*Service, error)
	Watch(...WatchOption) (Watcher, error)
	String() string
}
```

有几个关键的接口

* `Register` 注册服务
* `Deregister` 注销服务
* `GetService` 获取指定服务名的服务列表，区分版本
* `ListServices` 展示所有服务列表
* `Watch` 监听服务信息变化，包括节点上下线，服务版本修改，服务元数据修改
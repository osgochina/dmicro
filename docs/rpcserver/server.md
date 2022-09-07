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
		server.OptPrintDetail(true),
		server.OptServiceVersion(serviceVersion),
		server.OptRegistry(reg),
	)
	svr.RouteCall(new(Math))
	_ = svr.ListenAndServe()
}
```

## 支持的配置方法

* `OptServiceName(name string)` 设置服务名称
* `OptServiceVersion(version string) ` 当前服务版本
* `OptRegistry(r registry.Registry) ` 设置服务注册中心
* `OptGlobalPlugin(plugin ...drpc.Plugin)` 设置插件
* `OptEnableHeartbeat(t bool)` 是否开启心跳检测
* `OptTlsFile(tlsCertFile string, tlsKeyFile string)` 设置证书内容
* `OptTlsConfig(config *tls.Config)` 设置证书对象
* `OptProtoFunc(pf proto.ProtoFunc)` 设置协议
* `OptSessionAge(n time.Duration)` 设置会话生命周期
* `OptContextAge(n time.Duration)` 设置单次请求生命周期
* `OptSlowCometDuration(n time.Duration)` 设置慢请求的定义时间
* `OptBodyCodec(c string)` 设置消息内容编解码器
* `OptPrintDetail(c bool)` 是否打印消息详情
* `OptNetwork(net string)` 设置网络类型
* `OptNetwork(net string)` 设置网络类型
* `OptListenAddress(addr string)` 设置监听的网络地址
* `OptListenAddress(addr string)` 设置监听的网络地址
* `OptMetrics(m metrics.Metrics)` 设置指标统计对象

## 使用组件

### 使用 `Registry` 组件，请参考文档  [Registry](/component/registry.md)
### 使用 `Metrics` 组件，请参考文档  [Metrics](/component/metrics.md)
# RPC 客户端

直接使用`drpc`的`Endpoint`就能创建一个可用的rpc客户端。
那为什么还需要`rpc client`呢？

前面的文档已经阐述了如何创建一个原始的`drpc.Endpoint`,链接到服务端。
但是一个成熟的微服务框架的客户端不仅仅是一个单独的链接就可行的。

它需要包含`服务发现`,`链路追踪`,`限流`,`指标上报`等等服务治理的功能。这些组件与`drpc.Endpoint`如何有机的结合起来，就需要`rpc client`组件来实现了。

## 快速开始

```go
func main() {
	serviceName := "foo"
	cli := client.NewRpcClient(serviceName,
            client.OptPrintDetail(true),
            client.OptLocalIP("127.0.0.1"),
            client.OptMetrics(prometheus.NewPromMetrics(
            metrics.OptHost("0.0.0.0"),
            metrics.OptPort(9102),
	    )),
	)
	var result int
	stat := cli.Call("/math/add",
		[]int{1, 2, 3, 4, 5},
		&result,
		message.WithSetMeta("author", "clownfish"),
	).Status()
	if !stat.OK() {
		logger.Fatalf(context.TODO(),"%v", stat)
	}
	logger.Printf(context.TODO(),"result: %d", result)
}

```

## 支持的配置方法

* `OptServiceName(name string)` 设置服务名称
* `OptServiceVersion(version string) ` 当前服务版本
* `OptRegistry(r registry.Registry) ` 设置服务注册中心
* `OptSelector(s selector.Selector)` 设置服务选择器
* `OptGlobalPlugin(plugin ...drpc.Plugin)` 设置插件
* `OptHeartbeatTime(t time.Duration)` 设置心跳包时间
* `OptTlsFile(tlsCertFile string, tlsKeyFile string)` 设置证书内容
* `OptTlsConfig(config *tls.Config) ` 设置证书对象
* `OptProtoFunc(pf proto.ProtoFunc) ` 设置协议方法
* `OptRetryTimes(n int) ` 设置重试次数
* `OptSessionAge(n time.Duration) ` 设置会话生命周期
* `OptContextAge(n time.Duration)` 设置单次请求生命周期
* `OptSlowCometDuration(n time.Duration)` 设置慢请求的定义时间
* `OptBodyCodec(c string)` 设置消息内容编解码器
* `OptPrintDetail(c bool)` 是否打印消息详情
* `OptNetwork(net string)` 设置网络类型
* `OptLocalIP(addr string)` 设置本地监听的地址
* `OptCustomService(service *registry.Service)` 设置自定义service
* `OptMetrics(m metrics.Metrics)` 设置统计数据接口

## 使用组件


### 使用 `Registry` 组件，请参考文档  [Registry](/component/registry.md)
### 使用 `Selector` 组件，请参考文档  [Selector](/component/selector.md)
### 使用 `Metrics` 组件，请参考文档  [Metrics](/component/metrics.md)
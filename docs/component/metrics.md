# 指标 (Metrics)


## 概念

## 功能特性

* 抽象出 `Metrics` 接口,方便实现多种指标组件
* 内置 `Prometheus` 组件。

## 快速使用

### 在 `rpc server` 中使用 `Metrics` 组件

```go
	server.NewRpcServer("test_one",
		server.OptEnableHeartbeat(true),
		server.OptListenAddress("127.0.0.1:8199"),
		server.OptMetrics(prometheus.NewPromMetrics(
			metrics.OptHost("0.0.0.0"),
			metrics.OptPort(9101),
			metrics.OptPath("/metrics"),
			metrics.OptServiceName("test_one"),
		)),
	)
```
使用 `server.OptMetrics()` 方法传入要使用的 `Metrics`组件。以上代码示例是使用 `Prometheus` 作为指标组件。

使用 `prometheus.NewPromMetrics()` 创建对象的时候需要传入参数。

```go
type Options struct {
	Host        string  // Prometheus监听的地址
	Port        int     // 监听的端口
	Path        string  // 监听的路径
	ServiceName string  // 当前服务的名称
	Plugins     []drpc.Plugin // 需要使用的插件列表
}
```

* `metrics.OptHost("127.0.0.1")` Prometheus 监听的地址
* `metrics.OptPort(9101)` Prometheus 监听的端口
* `metrics.OptPath("/metrics")` http请求的路径

以上配置后 `Prometheus` 指标的请求地址为: `http://127.0.0.1:9101/metrics`.

* `metrics.OptServiceName("test_one")` 自定义的服务名称，如果不传入该参数，则默认使用 server name。
* `metrics.OptPlugin(plugin)` 框架默认已经注册了一个统计指标的插件，如果有需要更多的指标统计，可以注册自定义的插件。

### 在 `rpc client` 中使用 `Metrics` 组件

```go
c := client.NewRpcClient("test_one",
		client.OptMetrics(prometheus.NewPromMetrics(
			metrics.OptHost("0.0.0.0"),
			metrics.OptPort(9102),
            metrics.OptPath("/metrics"),
            metrics.OptServiceName("test_one"),
		)),
	)
```
使用 `client.OptMetrics()` 方法传入要使用的 `Metrics`组件。以上代码示例是使用 `Prometheus` 作为指标组件。

其他使用方法同 `rpc server`。

## 框架内支持的指标

作为 server 的指标
* `rpc_server_reply_code_total counter` 统计 `call` 请求的响应 `code` 值。
* `rpc_server_reply_duration_ms histogram` 统计 `call` 请求的处理耗时(仅表示处理时间，不包含网络通讯时间)。

作为 client 的指标
* `rpc_server_call_code_total counter` 统计 `call` 请求的响应 `code` 值。
* `rpc_server_call_duration_ms histogram` 统计 `call` 请求的响应总耗时(包含网络通讯时间)。
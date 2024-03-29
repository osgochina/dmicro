### 心跳

#### 如何使用

```go
func main() {
	svr := drpc.NewEndpoint(drpc.EndpointConfig{
		ListenPort:  9090,
		PrintDetail: true,
	},
		heartbeat.NewPong(), //使用pong插件，响应客户端的ping请求
	)
	go svr.ListenAndServe()

	time.Sleep(time.Second)

	cli := drpc.NewEndpoint(
		drpc.EndpointConfig{PrintDetail: false},
		heartbeat.NewPing(3, true), // 使用ping插件，发送ping请求给服务端
	)
	cli.Dial(":9090")
	time.Sleep(time.Second * 20)
}
```

#### 原理

- 客户端端点服务启动后`AfterNewEndpoint`，开启一个单独的协程来遍历会话，并发送心跳.
- 服务端端点启动的时候`AfterNewEndpoint`，开启一个协程，遍历会话，检查超时时间.
- 服务端注册`ping`处理方法，更新心跳数据。
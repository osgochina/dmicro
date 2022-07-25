### 忽略请求方法名的大小写

#### 如何使用
```go
func main() {
	svr := drpc.NewEndpoint(drpc.EndpointConfig{
		ListenPort:  9090,
		PrintDetail: true,
	},
		ignorecase.NewIgnoreCase(), //使用IgnoreCase插件，活跃serviceName的大小写
	)
	go svr.ListenAndServe()

	time.Sleep(time.Second)

	cli := drpc.NewEndpoint(
		drpc.EndpointConfig{PrintDetail: false},
		ignorecase.NewIgnoreCase(), // 使用IgnoreCase插件，活跃serviceName的大小写
	)
	cli.Dial(":9090")
	time.Sleep(time.Second * 20)
}
```

#### 原理

只需要在`AfterReadCallHeader`和`AfterReadPushHeader`两个`hook点`执行`ResetServiceMethod`方法，重置`ServiceName`就可以。

```go
func (that *ignoreCase) AfterReadCallHeader(ctx drpc.ReadCtx) *drpc.Status {
	// Dynamic transformation path is lowercase
	ctx.ResetServiceMethod(strings.ToLower(ctx.ServiceMethod()))
	return nil
}

func (that *ignoreCase) AfterReadPushHeader(ctx drpc.ReadCtx) *drpc.Status {
	return that.AfterReadCallHeader(ctx)
}
```

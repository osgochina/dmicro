### 概念

事件机制是一种经过充分测试的可靠机制，是一种非常适用与解耦的机制。
在`drpc`中，事件贯穿与整个`Endpoint`的生命周期，是它不可或缺的重要一环。

插件是使用多个`事件`有逻辑的组合形成的实现特定功能的一种机制，整个`drpc`的插件都离不开事件机制。

比如`ignoreCase`插件，就使用`AfterReadCallHeader`,`AfterReadPushHeader`两个事件，来修改请求的`ServiceMethod`，把它转换成全小写,
从而达到忽略大小写的效果。

### 事件与插件的关系

事件架构图待完成

### Endpoint生命周期中会触发的事件列表

#### BeforeNewEndpoint

 创建Endpoint之前触发该事件

在创建endpoint之前触发该事件。参数为`EndpointConfig`可以在该事件正查看，修改配置信息。

```go
BeforeNewEndpoint(*EndpointConfig, *PluginContainer) error
```

#### AfterNewEndpoint

 创建Endpoint之后触发该事件
```go
AfterNewEndpoint(EarlyEndpoint) error
```

创建Endpoint之后触发该事件。参数为`EarlyEndpoint`接口，该接口具体定义可以参考源码。

#### BeforeCloseEndpoint

 关闭Endpoint之前触发该事件

```go
BeforeCloseEndpoint(Endpoint) error
```

#### AfterCloseEndpoint

关闭Endpoint之后触发该事件

```go
AfterCloseEndpoint(Endpoint, error) error
```

#### AfterRegRouter

路由注册成功触发该事件

```go
AfterRegRouter(*Handler) error
```

#### AfterRegRouter

服务端监听以后触发该事件

```go
AfterListen(net.Addr) error
```

#### BeforeDial

作为客户端链接到服务端之前调用该事件

```go
BeforeDial(sess EarlySession, isRedial bool) *Status
```

#### AfterDial

作为客户端链接到服务端成功以后触发该事件

```go
AfterDial(sess EarlySession, isRedial bool) *Status
```

#### AfterDialFail

作为客户端链接到服务端失败以后触发该事件

```go
	AfterDialFail(sess EarlySession, err error, isRedial bool) *Status
```

#### AfterAccept

作为服务端，接收到客户端的链接后触发该事件

```go
AfterAccept(EarlySession) *Status
```


#### BeforeWriteCall

写入CALL消息之前触发该事件

```go
BeforeWriteCall(WriteCtx) *Status
```

#### AfterWriteCall

写入CALL消息成功之后触发该事件

```go
AfterWriteCall(WriteCtx) *Status
```


#### BeforeWriteReply

写入Reply消息之前触发该事件

```go
BeforeWriteReply(WriteCtx) *Status
```

#### AfterWriteReply

写入Reply消息成功之后触发该事件

```go
AfterWriteReply(WriteCtx) *Status
```

#### BeforeWritePush

写入PUSH消息之前触发该事件

```go
BeforeWritePush(WriteCtx) *Status
```

#### AfterWritePush

写入PUSH消息成功之后触发该事件

```go
	AfterWritePush(WriteCtx) *Status
```

#### BeforeReadHeader

执行读取Header之前触发该事件

```go
BeforeReadHeader(EarlyCtx) error
```

#### AfterReadCallHeader

读取CALL消息的Header之后触发该事件

```go
AfterReadCallHeader(ReadCtx) *Status
```


#### BeforeReadCallBody

读取CALL消息的body之前触发该事件

```go
BeforeReadCallBody(ReadCtx) *Status
```

#### AfterReadCallBody

读取CALL消息的body之后触发该事件

```go
AfterReadCallBody(ReadCtx) *Status
```

#### AfterReadPushHeader

读取PUSH消息Header之后触发该事件

```go
AfterReadPushHeader(ReadCtx) *Status
```

#### BeforeReadPushBody

读取PUSH消息body之前触发该事件

```go
BeforeReadPushBody(ReadCtx) *Status
```

#### AfterReadPushBody

读取PUSH消息body之后触发该事件

```go
AfterReadPushBody(ReadCtx) *Status
```

#### AfterReadReplyHeader

读取REPLY消息Header之前触发该事件

```go
AfterReadReplyHeader(ReadCtx) *Status
```

#### BeforeReadReplyBody

读取REPLY消息body之前触发该事件

```go
BeforeReadReplyBody(ReadCtx) *Status
```

#### AfterReadReplyBody

读取REPLY消息body之后触发该事件

```go
AfterReadReplyBody(ReadCtx) *Status
```

#### AfterDisconnect

断开会话以后触发该事件

```go
AfterDisconnect(BaseSession) *Status
```




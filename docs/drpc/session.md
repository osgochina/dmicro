# Session 会话

`client`与`server`之间通讯的时候，因为网络协议的多样性，我们抽象出了`Socket`接口，来屏蔽网络协议的复杂性,
有了`Socket`接口后，我们能进行通讯了，但是在开发的过程中，如果直接对`Socket`进行操作，缺少灵活性。
在`Socket`层之上，我们有抽象出了`Session`层。

`Session` 层的意义就是屏蔽底层的复杂性，并保证了`client`与`server`之间通讯链接的唯一性，能在其上增加很多特殊的逻辑，利于扩展。

## Session的生命周期及状态

一个会话从创建到关闭，会经历多个阶段，在不同的阶段有不同的状态及可用接口。

以下4个`Interface`,提供不同的方法，适用与不同的生命周期。
* `EarlySession`
> 会话刚被创建，还未启动`goroutine`读取数据。

    该对象主要提供给`event`使用。
    1. `AfterDial`
    2. `AfterDialFail`
    3. `AfterAccept`

* `BaseSession`
> 最基础的会话信息接口。被用在`AfterDisconnect`事件中。

* `CtxSession`
> 在处理器处理逻辑的时候，在上下文中传递。

* `Session`

> 会话的完全体，提供了session的所有功能。

会话不同的生命周期，对应了不同的状态，状态列表如下所示：

* 会话准备阶段
* 会话就绪
* 会话主动关闭中
* 会话已经主动关闭
* 会话被动关闭中
* 会话被动关闭
* 会话重建中
* 会话重建失败

## Session的配置

在`EndpointConfig`配置项中，对`Session`其效果的只有两项，分别是：
### DefaultSessionAge
> 默认session会话生命期(超时时间)。
### DefaultContextAge
> 默认单次(请求/响应)的超时时间.


更详细的配置介绍，请参考 [config](drpc/config.md) 章节。

## Session的方法

### 当前Session属于那个Endpoint
* `Endpoint() Endpoint`

实现它的`Interface`有：`EarlySession`,`BaseSession`,`Session`

### 获取本地监听地址
* `LocalAddr() net.Addr`

所有`Interface`都实现。

### 获取远端地址
* `RemoteAddr() net.Addr`

所有`Interface`都实现。

### 临时存储区对象
* `Swap() *gmap.Map`

所有`Interface`都实现。

### 设置会话ID
* `SetID(newID string)`

实现它的`Interface`有：`EarlySession`,`Session`

### 处理原始链接的fd
* `ControlFD(f func(fd uintptr)) error`

实现它的`Interface`有：`EarlySession`

### 修改会话的底层socket
* `ModifySocket(fn func(conn net.Conn) (modifiedConn net.Conn, newProtoFunc proto.ProtoFunc))`

实现它的`Interface`有：`EarlySession`

### 获取当前会话适用的协议
* `GetProtoFunc() proto.ProtoFunc`

实现它的`Interface`有：`EarlySession`

### 会话刚建立时临时发送消息
* `EarlySend(mType byte, serviceMethod string, body interface{}, stat *status.Status, setting ...message.MsgSetting) (opStat *status.Status)`

不执行任何插件。

实现它的`Interface`有：`EarlySession`

### 在会话刚建立时临时接受信息
* `EarlyReceive(newArgs message.NewBodyFunc, ctx ...context.Context) (input message.Message)`

不执行任何插件。

实现它的`Interface`有：`EarlySession`

### 在会话刚建立时临时调用call发送和接收消息
* `EarlyCall(serviceMethod string, args, reply interface{}, callSetting ...message.MsgSetting) (opStat *status.Status)`

不执行任何插件。

实现它的`Interface`有：`EarlySession`

### 在会话刚建立时临时回复消息
* `EarlyReply(req message.Message, body interface{}, stat *status.Status, setting ...message.MsgSetting) (opStat *status.Status)`

不执行任何插件。

实现它的`Interface`有：`EarlySession`

### 在会话刚建立时发送原始push消息
* `RawPush(serviceMethod string, args interface{}, setting ...message.MsgSetting) (opStat *status.Status)`

不执行任何插件。

实现它的`Interface`有：`EarlySession`

### 获取会话的最大的生存周期
* `SessionAge() time.Duration`

实现它的`Interface`有：`EarlySession`,`CtxSession`

### 获取CALL和PUSH消息的最大生存周期
* `ContextAge() time.Duration`

实现它的`Interface`有：`EarlySession`,`CtxSession`

### 设置会话的最大生存周期
* `SetSessionAge(duration time.Duration)`

实现它的`Interface`有：`EarlySession`

### 设置单个CALL和PUSH消息的最大生存周期
* `SetContextAge(duration time.Duration)`

实现它的`Interface`有：`EarlySession`

### 获取会话id
* `ID() string`

实现它的`Interface`有：`BaseSession`,`CtxSession`,`Session`.

### 返回该会话被关闭时候的通知
* `CloseNotify() <-chan struct{}`

实现它的`Interface`有：`CtxSession`,`Session`.

### 检查该会话是否健康
* `Health() bool`

实现它的`Interface`有：`CtxSession`,`Session`.

### 发送CALL消息，并异步接收响应
* `AsyncCall(serviceMethod string, args interface{}, result interface{}, callCmdChan chan<- CallCmd, setting ...message.MsgSetting) CallCmd`

实现它的`Interface`有：`CtxSession`,`Session`.

### 发送CALL消息并阻塞获取响应值
* `Call(serviceMethod string, args interface{}, result interface{}, setting ...message.MsgSetting) CallCmd`

实现它的`Interface`有：`CtxSession`,`Session`.

### 发送PUSH消息
* `Push(serviceMethod string, args interface{}, setting ...message.MsgSetting) *status.Status`

不接收响应，只返回发送状态。

实现它的`Interface`有：`CtxSession`,`Session`.

### 关闭会话
* `Close() error`

实现它的`Interface`有：`Session`.



# Endpoint

## 介绍

`Endpoint`是`drpc`的核心入口，整个框架的主对象，它可以支持不同的角色，既可以做客户端，也可以做服务端。
为`server`和`client`提供相同的`API`封装，方便开发者使用。

一个进程中可以支持多个`Endpoint`，它们相互之间是独立的，彼此没有影响。

## 生命周期

### 框架生命周期

从`NewEndpoint`生成一个对象开始，根据执行不同的方法，会转变成不同的角色。

* 执行`ListenAndServe`方法，则会变成服务端，监听指定的地址端口，处理客户端请求。
* 执行`Dial` 方法，则会链接到远程服务端，可以请求服务端接口。

整个框架的生命周期可以参考 [event](drpc/event.md) 模块。

### 请求的生命周期

`drpc`在处理每个连接时，会默认创建一个协程去处理.
在接收到一个链接的`Accept`后，会新开一个协程，阻塞等待数据的读取.
读取号数据后，经过解析，执行对应的方法。 可以理解为每个请求都是一个协程。


### Endpoint的使用

在快速开始的示例中，我们演示了`Endpoint`的简单用法，在这里我们详细介绍它的高级用法。

`Endpoint` 提供了三个对象`BaseEndpoint`,`EarlyEndpoint`,`Endpoint`分别对应不同的生命周期使用。

* `BaseEndpoint` 基础对象，继承使用。
* `EarlyEndpoint` 在事件`AfterNewEndpoint`中作为参数传入，因为此时Endpoint刚刚创立，很多`Endpoint`提供的方法还不能使用。
* `Endpoint` 完全对象，能使用所有的方法。

### Endpoint的方法

#### `NewEndpoint()`

* 新建Endpoint对象。

```go

var cfg = drpc.EndpointConfig{}
var plugins  []drpc.Plugin
plugins = append(plugins,ignorecase.NewIgnoreCase())
drpc.NewEndpoint(cfg,plugins...)

```

以上代码就是创建一个简单的`Endpoint`的用法，需要注意的是两个参数。
cfg 是配置对象，具体配置信息可以参考 [config](drpc/config.md).

重点来介绍下`plugins`,这个参数是可选参数，并且是一个数组，支持传入多个插件，并按传入的顺序执行。

插件是由多个事件组合，并且根据不同的逻辑形成的特定组件。
具体信息参考 [插件](drpc/plugin.md)。

#### `Router()`

* 返回路由对象

#### `SubRoute()`

* 设置路由分组

#### `RouteCall()`

* 通过结构体对象注册CALL命令路由

```go
type Math struct {
    drpc.CallCtx
}

func (m *Math) Add(arg *[]int) (int, *drpc.Status) {
    var r int
    for _, a := range *arg {
    r += a
    }
    return r, nil
}
endpointSvr.RouteCall(new(Math))
```

#### `RouteCallFunc()`

* 通过对象的方法注册CALL命令路由

```go
type Math struct {
    drpc.CallCtx
}

func (m *Math) Add(arg *[]int) (int, *drpc.Status) {
    var r int
    for _, a := range *arg {
    r += a
    }
    return r, nil
}
endpoint.RouteCallFunc((*Math).Add)
```

#### `RoutePush()`

* 通过结构体对象注册PUSH命令的路由

```go
type MathPush struct {
    drpc.PushCtx
}

func (m *MathPush) Add(arg *[]int) *drpc.Status {
    var r int
    for _, a := range *arg {
        r += a
    }
    return nil
}
endpoint.RoutePush(new(MathPush))
```

#### `RoutePushFunc()`

* 通过对象的方法注册PUSH命令的路由

```go
type MathPush struct {
    drpc.PushCtx
}

func (m *MathPush) Add(arg *[]int) *drpc.Status {
    var r int
    for _, a := range *arg {
        r += a
    }
    return nil
}
endpoint.RoutePush((*MathPush).Add)
```

#### `SetUnknownCall()`

* 设置CALL命令的默认路由

```go
endpoint.SetUnknownCall(func(ctx drpc.UnknownCallCtx) (interface{}, *drpc.Status){
    return nil, nil
},ignorecase.NewIgnoreCase())
```

#### `SetUnknownPush()`

* 设置PUSH命令的默认路由

```go
endpoint.SetUnknownCall(ctx drpc.UnknownPushCtx) *drpc.Status{
    return nil
},ignorecase.NewIgnoreCase())
```

#### `Dial()`

* 拨号链接远端

```go
sess, stat := cli.Dial("127.0.0.1:9091",jsonproto.NewJSONProtoFunc())
if !stat.OK() {
    logger.Fatalf("%v", stat)
}
```

大家注意第二个参数，这里是可以自定义协议的的，也就是说一个Endpoint，链接多个远端服务时，可以使用不同的协议。

#### `ListenAndServe()`

* 端点启动并监听，对外提供服务

```go
_ = svr.ListenAndServe(jsonproto.NewJSONProtoFunc())
```

监听的时候也可以指定协议。

#### `Close()`

* 关闭端点

#### `GetSession()`

* 根据sessionid获取session详情。

#### `RangeSession()`

* 遍历所有session。

#### `CountSession()`

* 统计当前session个数。

#### `ServeConn()`

* 通过提供的链接，生成会话并返回。

1. 不支持自动重连
2. 不检查是否是 TLS链接
3. 会执行AfterAcceptPlugin 事件

#### `PluginContainer()`

* 获取端点的插件容器，可以通过它对插件进行操作。

#### `TLSConfig()`

* 获取该端点的证书信息。

#### `SetTLSConfig()`

* 设置该端点的证书信息。

#### `SetTLSConfigFromFile()`

* 通过文件生成端点的证书信息。
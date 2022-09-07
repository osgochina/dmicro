## 配置详情

在创建`Endpoint`的时候，需要传入`EndpointConfig`这个struct作为参数，那么这个配置项的具体含义就需要了解。

### EndpointConfig 定义

在源码 [drpc/config.go](https://github.com/osgochina/dmicro/blob/main/drpc/config.go) 中可以查到具体定义。
```go
package drpc
import (
	"net"
	"time"
)
type EndpointConfig struct {

    Network string
    ListenIP string
    ListenPort uint16
    LocalIP string
    LocalPort uint16
    DefaultBodyCodec string
    DefaultSessionAge time.Duration
    DefaultContextAge time.Duration
    SlowCometDuration time.Duration
    PrintDetail bool
    DialTimeout time.Duration
    RedialTimes int
    RedialInterval time.Duration
    
    listenAddr net.Addr
    localAddr net.Addr
    slowCometDuration time.Duration
    checked bool
}
```

我们可以把 `EndpointConfig`作为参数，在`NewEndpoint`的时候传入。
> 注意`EndpointConfig`的参数会根据`Endpoint`的角色不同，而起到不同的作用。详细介绍中会有具体介绍。

```go
config := rpc.EndpointConfig{
    CountTime:   true,
    LocalIP:     "127.0.0.1",
    ListenPort:  9091,
    PrintDetail: true,
}
svr := drpc.NewEndpoint(config)
```

### 参数详解

#### Network

当前服务要使用的网络类型,可以是`tcp`, `tcp4`, `tcp6`, `unix`, `unixpacket`,`kcp`,`quic`。

* 如果角色是服务端，则会使用对应的网络类型监听具体的端口或地址。
* 如果角色是客户端，则会使用对应的网络类型链接到服务端。
* 注意，框架已经支持`kcp`,`quic`协议，具体的用法可以参考
* [kcp](https://github.com/osgochina/dmicro/tree/main/examples/kcp)
* [quic](https://github.com/osgochina/dmicro/tree/main/examples/quic)

#### ListenIP

作为服务端角色时，要监听的服务器本地IP或unix socket地址。

#### ListenPort

作为服务端角色时，需要监听的本地端口号。如果不传入，则表示随机监听本地端口。

#### LocalIP

作为客户端角色时,请求服务端时候，本地使用的地址。当`ListenIP`没有传入的时候，会默认把`LocalIP`的值赋值给`ListenIP`。

#### LocalPort

作为客户端角色时,请求服务端时候，本地使用的地址端口号。如果不传入，则表示随机使用端口去链接服务端。

#### DefaultBodyCodec

RPC请求响应内容的编码格式，默认使用json作为消息内容的编码格式。
目前框架支持的编码格式有：

* json
* form
* xml
* plain

使用方法如下：

* 获取默认`Endpoint`的编码格式.

```go
 drpc.DefaultBodyCodec()
```

* 设置编码格式

```go
drpc.SetDefaultBodyCodec(codec.IdXml)
drpc.SetDefaultBodyCodec(codec.IdPlain)
drpc.SetDefaultBodyCodec(codec.IdJson)
drpc.SetDefaultBodyCodec(codec.IdForm)
```

注意，如果需要更改默认的编码格式，需要在`NewEndpoint`之前设置，不然不起效果。

#### DefaultSessionAge

默认session会话生命期(超时时间)。

当`Endpoint`作为服务端或者客户端，它们之间的链接建立以后，会生成会话`session`，如果设置了`DefaultSessionAge`,则表示生成的session有效期为设置的值。
如果不设置该参数，则表示`session`不过期。

#### DefaultContextAge

默认单次(请求/响应)的超时时间.

生成了会话以后，会有请求及响应，`DefaultContextAge`表示单次请求/响应的超时时间。

#### SlowCometDuration

设置慢处理的定义。

定义请求处理多少时间是慢请求。默认是`math.MaxInt64`,表示不记录慢处理。

#### PrintDetail

打印处理日志的时候，是否需要打印出详细的`body`及`metadata`。默认是`false`.

#### CountTime

是否统计每个请求消耗的时间。默认是`false`.

#### DialTimeout

作为客户端角色时，链接服务端的超时时间。

#### RedialTimes

仅限客户端角色使用,链接中断时候，试图重新链接服务端的最大重试次数。

#### RedialInterval

仅限客户端角色使用 试图链接服务端时候，每次重试之间的时间间隔.
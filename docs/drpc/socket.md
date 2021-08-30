
### 设置socket的属性

#### GetReadLimit

> 获取消息的最大长度设置信息。

在`rpc`调用的时候，传输的消息有长度限制，默认是`math.MaxUint32`,约等于不限制长度。

```go
limit := drpc.GetReadLimit()
```

#### SetReadLimit

> 设置消息的最大长度。

设置完该值以后，消息如果超过改长度，则会报错。

```go
drpc.SetReadLimit(limit)
```

#### SetSocketKeepAlive

> 开启死链检测

在`TCP`中有一个`Keep-Alive`的机制可以检测死连接，应用层如果对于死链接周期不敏感或者没有实现心跳机制，可以使用操作系统提供的 keepalive 机制来踢掉死链接。

#### SetSocketKeepAlivePeriod

> 死链检测间隔时间。

#### SocketReadBuffer

> 获取接收输入缓存区内存尺寸

#### SetSocketReadBuffer

> 设置接收输入缓存区内存尺寸

#### SocketWriteBuffer

> 获取发送输出缓存区内存尺寸

#### SetSocketWriteBuffer

> 获取发送输出缓存区内存尺寸


#### SetSocketNoDelay

> 开启关闭`tcp no delay`算法,默认是`true`.

开启后TCP 连接发送数据时会关闭 Nagle 合并算法，立即发往对端 TCP 连接。
在某些场景下，如命令行终端，敲一个命令就需要立马发到服务器，可以提升响应速度，请自行 Google Nagle 算法。

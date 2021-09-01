# Socket 套接字

`client`与`server`之间通讯的时候，因为网络协议的多样性，我们抽象出了`Socket`接口，来屏蔽网络协议的复杂性。
支持多种网络协议。

### Socket的方法

#### 创建Socket

* `NewSocket(c net.Conn, protoFunc ...ProtoFunc) Socket `

#### 获取原始句柄

* `ControlFD(f func(fd uintptr)) error `

#### 获取socket本地的地址

* `LocalAddr() net.Addr`


#### 获取socket远端的地址

* `RemoteAddr() net.Addr`

#### 设置超时时间

* `SetDeadline(t time.Time) error`

#### 设置读取数据的超时时间

* `SetReadDeadline(t time.Time) error`

#### 设置发送数据的超时时间

* `SetWriteDeadline(t time.Time) error`

#### 发送消息

* `WriteMessage(message Message) error`

#### 读取消息

* `ReadMessage(message Message) error`

#### 读取指定的字符

* `Read(b []byte) (n int, err error)`

#### 写入指定的字符

* `Write(b []byte) (n int, err error)`

#### 关闭链接

* `Close() error`

#### 设置临时交换区数据

* `Swap(newSwap ...*gmap.Map) *gmap.Map`

#### 返回临时交换区长度

* `SwapLen() int`

#### 返回套接字的id

* `ID() string`

#### 设置套接字的id

* `SetID(string)`

#### 把套接字的 net.Conn重置为新的

* `Reset(netConn net.Conn, protoFunc ...ProtoFunc)`

#### 返回原始 net.Conn

* `Raw() net.Conn`

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

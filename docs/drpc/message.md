# Message 消息

从客户端链接到服务端，发起一个请求，该请求会被封装成`消息`，经过`codec`的编解码，用`TFilter`的处理，使用`proto`拼装，最后通过`socket`发送到服务端。

从服务端处理完成以后，给客户端的响应，也会被封装成消息，进过以上的相反步骤，返回给客户端。

通过上面的介绍，大家对`消息`的作用应该有个基本得认识。总的来说，消息就是请求及响应的封装。

## 消息的类型

从角色的角度来看，消息可以分为`Input`，`Output`两种类型。

* 客户端角色：发送的消息是`Output`，接收的消息是`Input`。
* 服务端角色：接收的消息是`Input`，发送的消息是`Output`。

从功能的角度来看，消息可以分位以下几种类型：

* `TypeUndefined` 未知消息类型
* `TypeCall`      请求消息  
* `TypeReply`     响应消息
* `TypePush`      推送消息(不需要响应)
* `TypeAuthCall`  权限认证请求
* `TypeAuthReply` 权限认证响应

## 消息的构成

消息的功能非常多，它承担了请求及响应的格式化，组装等等.

它的组成分为`消息头(Header)`和`消息体(Body)`。

### 消息头

消息头中主要承载的元素有`消息序列号(Seq)`，`消息类型(MType)`，`请求方法名(ServiceMethod)`,`消息状态(Status)`,`自定义元数据(Meta)`,

* `消息序列号(Seq)` 是int32类型，保证在一个链接中唯一且自增即可，这样方便客户端和服务端区分消息。
* `消息类型(MType)` 目前有`TypeCall`，`TypeReply`，`TypePush`，`TypeAuthCall`，`TypeAuthReply`这五种。
* `请求方法名(ServiceMethod)` 请求的服务方法名称 长度必须小于255字节。
* `消息状态(Status)` 详见`Status`章节，在传输的过程中，是以串化的形式传输。
* `自定义元数据(Meta)` 自定义的数据，数据在传输的时候是使用了序列化串，最大长度为: max len ≤ 65535

### 消息体

通常来说，消息体的结构是非常简单的，该有的信息在消息头中已经存在，消息体的作用就是传输真正的消息内容。

我们的消息体也确实比较简单，它由`编码格式(Codec)`,`消息体(Body)`组成，支持`MarshalBody(串化)`,`UnmarshalBody(反向解析)`。

### 消息


## 消息设置

对消息进行设置是常见的需求，在`drpc`中，要实现这个功能是非常简单的。

写入值到消息的元数据中
```go
sess.Call("/math/add",
    []int{1, 2, 3, 4, 5},
    &result,
    message.WithSetMeta("author", "osgochina"),
).Status()
```

我们暴露了 ```type MsgSetting func(Message)```类型，只要实现了类型，就可以作为参数传入到请求中，从而对消息进行设置。

预设了一些方法，实现需要的功能，当然你也可以自己实现该方法。
预设的方法如下：

* 设置消息的上下文对象
> func WithContext(ctx context.Context) MsgSetting

* 设置消息的服务器接口名
> func WithServiceMethod(serviceMethod string) MsgSetting 

* 设置消息的状态
> func WithStatus(stat *status.Status) MsgSetting

* 添加消息的元数据
> func WithSetMeta(key, value string) MsgSetting

* 使用数组添加元数据
> func WithSetMetas(metas map[string]interface{}) MsgSetting 

* 删除消息元数据
> func WithDelMeta(key string) MsgSetting

* 设置消息的消息体编码格式
> func WithBodyCodec(bodyCodec byte) MsgSetting

* 设置消息体的内容
> func WithBody(body interface{}) MsgSetting 

* 设置创建消息体的函数
> func WithNewBody(newBodyFunc NewBodyFunc) MsgSetting

* 设置消息的管道类型
> func WithXFerPipe(filterID ...byte) MsgSetting

以上这些方法都是对消息进行一些自定义设置。
另外，我们还可以对消息的最大长度进行设置，该方法是全局生效的。

* 设置消息最大长度
```go
//设置消息最大长度为8M
message.SetMsgSizeLimit(1024*1024*8) 
```
* 获取消息最大长度
> func MsgSizeLimit() uint32
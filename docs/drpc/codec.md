# 编码解码器(codec)

编码解码器`codec`是`encoder`和`decoder`的缩写，这里的含义"网络数据与业务消息之间相互转换"的代码。
使得消息内容能从各种语言层面的类型，转换成方便网络传输的字节流，通过网络传输后，再次转换成业务可使用的数据类型。

`codec`位于`proto`内层。

```mermaid
flowchart LR
     Message --> Proto --> TFilter --> Codec --> Socket
```

默认支持的编码解码器(codec)

| id  | name | 介绍     |
|-----|------|--------|
| f   |  form    | 表单编解码器 |
| j   |  json    | json   |
| s   |  plain    | 字符串    |
| p   |  protobuf    | protobuf    |
| x   |  xml    | xml    |



## 如何使用编码解码器

**框架默认编码解码器为`JSONCodec`。**

  可以通过`drpc.DefaultBodyCodec()`获取默认的编码解码器。

  也可以通过`err := drpc.SetDefaultBodyCodec(codec.ProtobufName)`设置默认的.

  ps: 该设置是全局设置，多个`endpoint`同时生效。

**为单个`endpoint`设置编码解码器。**

创建`endpoint`传入`drpc.EndpointConfig`配置的时候， 设置参数`DefaultBodyCodec`.

```go
    cfg :=drpc.EndpointConfig{DefaultBodyCodec:codec.XmlName}
```

**为单个消息单独设置编码解码器**

通过`drpc.WithBodyCodec(codec.PlainName)`方法，为每条消息设置编码解码器。

```go
var result int
stat := sess.Call("/math/add",
    []int{1, 2, 3, 4, 5},
    &result,
    message.WithBodyCodec(codec.PlainName),
).Status()
if !stat.OK() {
    logger.Fatalf("%v", stat)
}
```
如以上代码，为此次请求设置单独的编码解码器。

一般为单挑消息设置单独的`编码解码器`常见于自定义插件中，日常业务开发不常用。

## 实现自己的`编码解码器(codec)`

框架已经为大家实现了常用的`编码解码器(codec)`,可以直接使用。但是业务场景千变万化，需求多样复杂，为了更便捷的适应业务场景,
框架考虑到这种情况，抽象出`codec`接口。
接口定义如下：
```go
type Codec interface {
	ID() byte
	Name() string
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}
```

实现如下

```go
const (
    TestcodecName = "test"
    TestcodecId   = 't'
)
type TESTCodec struct{}

func (TESTCodec) ID() byte {
    return TestcodecId
}

func (TESTCodec) Name() string {
    return TestcodecName
}

func (TESTCodec) Marshal(v interface{}) ([]byte, error) {
    return json.Marshal(v)
}

func (TESTCodec) Unmarshal(data []byte, v interface{}) error {
    return json.Unmarshal(data, v)
}
```

首先定义`TestcodecName`,`TestcodecId`。

注意，这里的name和id必须唯一，不能更框架已定义的codec冲突，具体已定义的内容见文档开头。

实现`Marshal`,`Unmarshal`两个方法就可以。

最后把自定义的编码解码器注册到框架`drpc.Reg(new(TESTCodec))`。

使用过程参见文档`如何使用编码解码器`。




# Proto 协议

协议层是`drpc框架`的核心，框架对各种协议的支持程度，决定了框架生命力。

为了提升框架的生命力，那么就需要做到最简单的支持一种新的协议。
得益于`drpc框架`科学的架构，它支持一种新的协议是如此的清晰简单。

目前`drpc框架`已经内置支持了`http`，`json`,`jsonrpc`，`raw`，`redis`协议，后续计划支持`protobuf`，`thrift`协议。

#### 协议结构的定义

`proto`定义了一个`interface`,只要需要支持的协议实现对应的方法，就能被框架所用。
定义如下：
```go
type Proto interface {
	Version() (byte, string)
	Pack(Message) error
	Unpack(Message) error
}
```

注意`Version`方法返回了协议的id及name，不能重复，所以自定义协议的时候需要注意以后的协议id及name

已使用的协议:

协议名| id  |name 
---|-----|---
[http协议](drpc/proto_http.md) | h   | http |
[json协议](drpc/proto_json.md) | j   | json
[jsonrpc协议](drpc/proto_jsonrpc.md) | 3   | jsonrpc
[raw协议](drpc/proto_raw.md) | 6   | raw
redis协议 | r   | redis
protobuf协议 | p   | protobuf
thrift-binary协议 | b   | thrift-binary




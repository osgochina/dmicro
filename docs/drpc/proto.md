# Proto 协议

协议层是`drpc框架`的核心，框架对各种协议的支持程度，决定了框架生命力。

为了提升框架的生命力，那么就需要做到最简单的支持一种新的协议。
得益于`drpc框架`科学的架构，它支持一种新的协议是如此的清晰简单。

目前`drpc框架`已经内置支持了`http`，`json`,`jsonrpc`，`raw`，`redis`协议，后续计划支持`protobuf`，`thrift`协议。

## 协议结构的定义

`proto`定义了一个`interface`,只要需要支持的协议实现对应的方法，就能被框架所用。
定义如下：
```go
type Proto interface {
	Version() (byte, string)
	Pack(Message) error
	Unpack(Message) error
}
```

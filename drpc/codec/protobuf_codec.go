package codec

import (
	"fmt"
	"google.golang.org/protobuf/proto"
)

const (
	ProtobufName = "protobuf"
	ProtobufId   = 'p'
)

func init() {
	Reg(new(ProtoCodec))
}

type ProtoCodec struct{}

func (ProtoCodec) Name() string {
	return ProtobufName
}

func (ProtoCodec) ID() byte {
	return ProtobufId
}

// Marshal 编码
func (ProtoCodec) Marshal(v interface{}) ([]byte, error) {
	return ProtoMarshal(v)
}

// Unmarshal 解码
func (ProtoCodec) Unmarshal(data []byte, v interface{}) error {
	return ProtoUnmarshal(data, v)
}

var (
	// PbEmptyStruct 如果需要编码nil 或者struct，则不能编码，使用空的
	PbEmptyStruct = new(PbEmpty)
)

func ProtoMarshal(v interface{}) ([]byte, error) {
	switch p := v.(type) {
	case proto.Message:
		return proto.Marshal(p)
	case nil, *struct{}, struct{}:
		return proto.Marshal(PbEmptyStruct)
	}
	return nil, fmt.Errorf("protobuf 编码失败: %T 传入的对象未实现 proto.Message,有可能的原因是你传入的是一个值，而非指针对象，可以检查是否使用了&获取地址", v)
}

func ProtoUnmarshal(data []byte, v interface{}) error {
	switch p := v.(type) {
	case proto.Message:
		return proto.Unmarshal(data, p)
	case nil, *struct{}, struct{}:
		return nil
	}
	return fmt.Errorf("protobuf 编码失败: %T 传入的对象未实现 proto.Message,有可能的原因是你传入的是一个值，而非指针对象，可以检查是否使用了&获取地址", v)
}

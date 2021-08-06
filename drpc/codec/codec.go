package codec

import "fmt"

// Codec 消息内容的编解码器
type Codec interface {
	ID() byte
	Name() string
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}

//编解码器对象容器
var codecMap = struct {
	idMap   map[byte]Codec
	nameMap map[string]Codec
}{
	idMap:   make(map[byte]Codec),
	nameMap: make(map[string]Codec),
}

const (
	// NilCodecID 空的编解码器id.
	NilCodecID byte = 0
	// NilCodecName 空的编解码器名称.
	NilCodecName string = ""
)

// Get 通过编解码器的id获取编解码器对象
func Get(codecID byte) (Codec, error) {
	codec, ok := codecMap.idMap[codecID]
	if !ok {
		return nil, fmt.Errorf("unsupported codec id: %d", codecID)
	}
	return codec, nil
}

// GetByName 通过编解码器的名字获取编解码器对象
func GetByName(codecName string) (Codec, error) {
	codec, ok := codecMap.nameMap[codecName]
	if !ok {
		return nil, fmt.Errorf("unsupported codec name: %s", codecName)
	}
	return codec, nil
}

// Marshal 使用指定编解码器编码
func Marshal(codecID byte, v interface{}) ([]byte, error) {
	codec, err := Get(codecID)
	if err != nil {
		return nil, err
	}
	return codec.Marshal(v)
}

// Unmarshal 使用指定编解码器解码
func Unmarshal(codecID byte, data []byte, v interface{}) error {
	codec, err := Get(codecID)
	if err != nil {
		return err
	}
	return codec.Unmarshal(data, v)
}

// Reg 注册编解码器到容器
func Reg(codec Codec) {
	if codec.ID() == NilCodecID {
		panic(fmt.Sprintf("codec id can not be %d", NilCodecID))
	}
	if _, ok := codecMap.idMap[codec.ID()]; ok {
		panic(fmt.Sprintf("multi-register codec id: %d", codec.ID()))
	}
	if _, ok := codecMap.nameMap[codec.Name()]; ok {
		panic("multi-register codec name: " + codec.Name())
	}
	codecMap.idMap[codec.ID()] = codec
	codecMap.nameMap[codec.Name()] = codec
}

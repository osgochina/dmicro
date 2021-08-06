package codec

import "encoding/json"

var _ Codec = new(JSONCodec)

const (
	NameJson = "json"
	IdJson   = 'j'
)

type JSONCodec struct{}

func (JSONCodec) ID() byte {
	return IdJson
}

func (JSONCodec) Name() string {
	return NameJson
}

func (JSONCodec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (JSONCodec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func init() {
	Reg(new(JSONCodec))
}

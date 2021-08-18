package codec

import "encoding/xml"

const (
	NameXml = "xml"
	IdXml   = 'x'
)

func init() {
	Reg(new(XMLCodec))
}

type XMLCodec struct{}

// Name returns codec name.
func (XMLCodec) Name() string {
	return NameXml
}

// ID returns codec id.
func (XMLCodec) ID() byte {
	return IdXml
}

func (XMLCodec) Marshal(v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}

func (XMLCodec) Unmarshal(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}

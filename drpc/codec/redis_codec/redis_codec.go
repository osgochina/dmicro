package redis_codec

import (
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/errors/gerror"
	"github.com/osgochina/dmicro/drpc/codec"
)

var _ codec.Codec = new(codec.JSONCodec)

const (
	NameRedis = "redis"
	IdRedis   = 'r'
)

type REDISCodec struct{}

func (REDISCodec) ID() byte {
	return IdRedis
}

func (REDISCodec) Name() string {
	return NameRedis
}

func (REDISCodec) Marshal(v interface{}) ([]byte, error) {
	m, ok := v.(Msg)
	if !ok {
		return nil, gerror.New("转换消息失败")
	}
	return m.Bytes(), nil
}

func (REDISCodec) Unmarshal(data []byte, v interface{}) error {
	return gjson.DecodeTo(data, v)
}

type CmdLine [][]byte

func init() {
	codec.Reg(new(REDISCodec))
}

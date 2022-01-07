package websocket

import (
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/mixer/websocket/jsonSubProto"
	"github.com/osgochina/dmicro/drpc/mixer/websocket/pbSubProto"
	"github.com/osgochina/dmicro/drpc/proto"
	"path"
	"strings"
)

// Client websocket的rpc客户端
type Client struct {
	drpc.Endpoint
}

// NewClient 创建websocket协议的rpc客户端
// rootPath: url路径
// cfg: drpc框架的配置
// globalLeftPlugin: 插件
func NewClient(rootPath string, cfg drpc.EndpointConfig, globalLeftPlugin ...drpc.Plugin) *Client {
	globalLeftPlugin = append([]drpc.Plugin{NewDialPlugin(rootPath)}, globalLeftPlugin...)
	endpoint := drpc.NewEndpoint(cfg, globalLeftPlugin...)
	return &Client{
		Endpoint: endpoint,
	}
}

// DialJSON 链接到服务器，并使用json协议编解码传输
func (that *Client) DialJSON(addr string) (drpc.Session, *drpc.Status) {
	return that.Dial(addr, jsonSubProto.NewJSONSubProtoFunc())
}

// DialProtobuf 链接到服务器，并使用protobuf协议编解码传输
func (that *Client) DialProtobuf(addr string) (drpc.Session, *drpc.Status) {
	return that.Dial(addr, pbSubProto.NewPbSubProtoFunc())
}

// Dial 链接到服务器
func (that *Client) Dial(addr string, protoFunc ...proto.ProtoFunc) (drpc.Session, *drpc.Status) {
	if len(protoFunc) == 0 {
		return that.Endpoint.Dial(addr, defaultProto)
	}
	return that.Endpoint.Dial(addr, protoFunc...)
}

// 格式化root path
func fixRootPath(rootPath string) string {
	rootPath = path.Join("/", strings.TrimRight(rootPath, "/"))
	return rootPath
}

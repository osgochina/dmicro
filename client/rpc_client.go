package client

import (
	"fmt"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/message"
	"github.com/osgochina/dmicro/drpc/plugin/heartbeat"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/selector"
	"time"
)

var (
	// DefaultPoolSize 连接池默认大小
	DefaultPoolSize = 100
	// DefaultPoolTTL 链接默认存活时间
	DefaultPoolTTL = time.Minute
	// DefaultBodyCodec 默认的消息编码
	DefaultBodyCodec = "json"
	// DefaultSessionAge 默认session会话生命周期
	DefaultSessionAge = time.Duration(0)
	// DefaultContextAge 默认单次请求生命周期
	DefaultContextAge = time.Duration(0)
	// DefaultDialTimeout 作为客户端角色时，请求服务端的超时时间
	DefaultDialTimeout = time.Second * 5
	// DefaultSlowCometDuration 慢处理定义时间
	DefaultSlowCometDuration = time.Duration(0)
)

type RpcClient struct {
	serviceName string // 服务名称
	endpoint    drpc.Endpoint
	opts        Options
}

func NewRpcClient(serviceName string, opt ...Option) *RpcClient {

	opts := NewOptions(opt...)
	//如果设置了心跳包，则发送心跳
	var heartbeatPing heartbeat.Ping
	if opts.HeartbeatTime > time.Duration(0) {
		heartbeatPing = heartbeat.NewPing(int(opts.HeartbeatTime/time.Second), false)
		opts.GlobalLeftPlugin = append(opts.GlobalLeftPlugin, heartbeatPing)
	}
	endpoint := drpc.NewEndpoint(opts.EndpointConfig(), opts.GlobalLeftPlugin...)
	// 优先使用已生成的证书对象
	if opts.TLSConfig != nil {
		endpoint.SetTLSConfig(opts.TLSConfig)
	} else {
		//如果设置了证书，则必须执行证书操作
		if len(opts.TlsCertFile) > 0 && len(opts.TlsKeyFile) > 0 {
			err := endpoint.SetTLSConfigFromFile(opts.TlsCertFile, opts.TlsKeyFile)
			if err != nil {
				logger.Fatalf("%v", err)
			}
		}
	}
	rc := &RpcClient{
		serviceName: serviceName,
		opts:        opts,
		endpoint:    endpoint,
	}
	return rc
}

func (that *RpcClient) Options() Options {
	return that.opts
}

// 获取服务可用的节点列表
func (that *RpcClient) next(serviceName string) (selector.Next, *drpc.Status) {
	next, err := that.Options().Selector.Select(serviceName)
	if err != nil {
		if err == selector.ErrNotFound {
			return nil, drpc.NewStatus(drpc.CodeInternalServerError, fmt.Sprintf("dmicro.client service %s: %s", serviceName, err.Error()))
		}
		return nil, drpc.NewStatus(drpc.CodeInternalServerError, fmt.Sprintf("dmicro.client error selecting %s node: %s", serviceName, err.Error()))
	}

	return next, nil
}

func (that *RpcClient) Call(serviceMethod string, args interface{}, result interface{}, setting ...message.MsgSetting) *drpc.Status {
	next, err := that.next(that.serviceName)
	if err != nil {
		return err
	}
	node, e := next()
	if e != nil {
		if e == selector.ErrNotFound {
			return drpc.NewStatus(drpc.CodeInternalServerError, fmt.Sprintf("dmicro.client service %s: %s", that.serviceName, e.Error()))
		}
		return drpc.NewStatus(drpc.CodeInternalServerError, fmt.Sprintf("dmicro.client error selecting %s node: %s", that.serviceName, e.Error()))
	}
	sess, stat := that.endpoint.Dial(node.Address, that.opts.ProtoFunc)
	if !stat.OK() {
		return stat
	}
	callCmd := sess.Call(serviceMethod, args, result, setting...)
	return callCmd.Status()
}

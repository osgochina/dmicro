package etcd

import (
	"context"
	"encoding/json"
	"github.com/osgochina/dmicro/registry"
	"path"
	"strings"
	"time"
)

type authKey struct{}

type authCreds struct {
	Username string
	Password string
}
type logConfigKey struct{}
type leasesInterval struct{}

// Auth 生成etcd的用户名密码认证方式配置
func Auth(username, password string) registry.Option {
	return func(o *registry.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, authKey{}, &authCreds{Username: username, Password: password})
	}
}

// LeasesInterval 租约续期时间
func LeasesInterval(t time.Duration) registry.Option {
	return func(o *registry.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, leasesInterval{}, t)
	}
}

// RegisterTTL 服务key在etcd中的生存时间
func RegisterTTL(t time.Duration) registry.Option {
	return func(o *registry.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, "RegisterTTL", t)
	}
}

// 生成服务节点的路径
func nodePath(s, id string) string {
	service := strings.Replace(s, "/", "-", -1)
	node := strings.Replace(id, "/", "-", -1)
	return path.Join(prefix, service, node)
}

// 生成服务的路径，该路径会包含节点
func servicePath(s string) string {
	return path.Join(prefix, strings.Replace(s, "/", "-", -1))
}

// 对值编码
func encode(s *registry.Service) string {
	b, _ := json.Marshal(s)
	return string(b)
}

// 对值解码
func decode(ds []byte) *registry.Service {
	var s *registry.Service
	_ = json.Unmarshal(ds, &s)
	return s
}

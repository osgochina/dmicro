package graceful

import (
	"github.com/gogf/gf/container/gset"
	"github.com/osgochina/dmicro/drpc"
	"net"
)

// GracePlugin 平滑重启插件
type GracePlugin struct {
	endpointList *gset.Set
}

var (
	_ drpc.AfterNewEndpointPlugin    = new(GracePlugin)
	_ drpc.BeforeCloseEndpointPlugin = new(GracePlugin)
	_ drpc.AfterListenPlugin         = new(GracePlugin)
)

// NewGracefulPlugin 平滑重启插件
func NewGracefulPlugin() *GracePlugin {
	return &GracePlugin{
		endpointList: gset.New(true),
	}
}

func (that *GracePlugin) Name() string {
	return "graceful"
}

func (that *GracePlugin) AfterNewEndpoint(endpoint drpc.EarlyEndpoint) error {
	that.endpointList.Add(endpoint)
	return nil
}

func (that *GracePlugin) BeforeCloseEndpoint(endpoint drpc.Endpoint) error {
	that.endpointList.Remove(endpoint)
	return nil
}

func (that *GracePlugin) AfterListen(addr net.Addr) error {

	return nil
}

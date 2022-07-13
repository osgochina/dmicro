package dserver

import (
	"fmt"
	"github.com/gogf/gf/os/gfile"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/proto/pbproto"
)

func (that *DServer) endpoint() {
	unix := fmt.Sprintf("/tmp/dserver.scoket")
	if gfile.Exists(unix) {
		_ = gfile.Remove(unix)
	}
	cfg := drpc.EndpointConfig{
		Network:  "unix",
		ListenIP: unix,
	}
	that.ctrlEndpoint = drpc.NewEndpoint(cfg)
	that.ctrlEndpoint.RouteCall(new(Ctrl))
	go func() {
		_ = that.ctrlEndpoint.ListenAndServe(pbproto.NewPbProtoFunc())
	}()
}

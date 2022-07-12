package dserver

import "github.com/osgochina/dmicro/drpc"

func (that *DServer) endpoint() {
	that.ctrlEndpoint = drpc.NewEndpoint(drpc.EndpointConfig{})
}

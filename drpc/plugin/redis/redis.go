package redis

import (
	"fmt"
	"github.com/osgochina/dmicro/drpc"
	"strings"
)

func NewRedisPlugin() drpc.Plugin {
	return &redis{}
}

type redis struct{}

var (
	_ drpc.AfterNewEndpointPlugin    = new(redis)
	_ drpc.AfterReadCallHeaderPlugin = new(redis)
)

func (that *redis) Name() string {
	return "redis"
}

func (that *redis) AfterNewEndpoint(peer drpc.EarlyEndpoint) error {
	peer.SetUnknownCall(that.call)
	peer.SetUnknownPush(that.push)
	return nil
}

func (that *redis) call(ctx drpc.UnknownCallCtx) (interface{}, *drpc.Status) {
	return nil, drpc.NewStatus(404, fmt.Sprintf("ERR unknown command '%s'", ctx.ServiceMethod()))
}

func (that *redis) push(ctx drpc.UnknownPushCtx) *drpc.Status {
	return drpc.NewStatus(404, fmt.Sprintf("ERR unknown command '%s'", ctx.ServiceMethod()))
}

func (that *redis) AfterReadCallHeader(ctx drpc.ReadCtx) *drpc.Status {
	ctx.ResetServiceMethod(fmt.Sprintf("/redis/%s", strings.ToLower(ctx.ServiceMethod())))
	return nil
}

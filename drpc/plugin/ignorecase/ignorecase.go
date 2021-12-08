// Package ignorecase dynamically ignoring the case of path
package ignorecase

import (
	"github.com/osgochina/dmicro/drpc"
	"strings"
)

// NewIgnoreCase Returns a ignoreCase plugin.
func NewIgnoreCase() *ignoreCase {
	return &ignoreCase{}
}

type ignoreCase struct{}

var (
	_ drpc.AfterReadCallHeaderPlugin = new(ignoreCase)
	_ drpc.AfterReadPushHeaderPlugin = new(ignoreCase)
)

func (that *ignoreCase) Name() string {
	return "ignoreCase"
}

func (that *ignoreCase) AfterReadCallHeader(ctx drpc.ReadCtx) *drpc.Status {
	// Dynamic transformation path is lowercase
	ctx.ResetServiceMethod(strings.ToLower(ctx.ServiceMethod()))
	return nil
}

func (that *ignoreCase) AfterReadPushHeader(ctx drpc.ReadCtx) *drpc.Status {
	return that.AfterReadCallHeader(ctx)
}

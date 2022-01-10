package drpc

import (
	"context"
	"github.com/gogf/gf/container/gmap"
	"github.com/osgochina/dmicro/drpc/codec"
	"github.com/osgochina/dmicro/drpc/message"
	"time"
)

// 伪造的响应命令，用在调用出错的情况，但是为了保持api的统一性，使用它
type fakeCallCmd struct {
	output    message.Message
	result    interface{}
	stat      *Status
	inputMeta *gmap.Map
}

// NewFakeCallCmd 构建伪造的回调命令
func NewFakeCallCmd(serviceMethod string, arg, result interface{}, stat *Status) CallCmd {
	return &fakeCallCmd{
		output: message.NewMessage(
			message.WithMType(TypeCall),
			message.WithServiceMethod(serviceMethod),
			message.WithBody(arg),
		),
		result: result,
		stat:   stat,
	}
}

func (that *fakeCallCmd) TraceEndpoint() (Endpoint, bool) {
	return nil, false
}

func (that *fakeCallCmd) TraceSession() (Session, bool) {
	return nil, false
}

var closedChan = func() <-chan struct{} {
	ch := make(chan struct{})
	close(ch)
	return ch
}()

func (that *fakeCallCmd) Done() <-chan struct{} {
	return closedChan
}

func (that *fakeCallCmd) Output() message.Message {
	return that.output
}

func (that *fakeCallCmd) Context() context.Context {
	return that.output.Context()
}

func (that *fakeCallCmd) Reply() (interface{}, *Status) {
	return that.result, that.stat
}

func (that *fakeCallCmd) StatusOK() bool {
	return that.stat.OK()
}

func (that *fakeCallCmd) Status() *Status {
	return that.stat
}

func (that *fakeCallCmd) InputBodyCodec() byte {
	return codec.NilCodecID
}

func (that *fakeCallCmd) InputMeta() *gmap.Map {
	if that.inputMeta == nil {
		that.inputMeta = gmap.New(true)
	}
	return that.inputMeta
}

func (that *fakeCallCmd) CostTime() time.Duration {
	return 0
}

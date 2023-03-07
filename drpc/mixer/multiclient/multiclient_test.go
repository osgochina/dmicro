package multiclient_test

import (
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/mixer/multiclient"
	"testing"
	"time"
)

type Arg struct {
	A int
	B int `param:"<range:1:>"`
}

type P struct{ drpc.CallCtx }

func (p *P) Divide(arg *Arg) (int, *drpc.Status) {
	return arg.A / arg.B, nil
}

func TestMultiClient(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		srv := drpc.NewEndpoint(drpc.EndpointConfig{
			ListenPort: 9090,
		})
		srv.RouteCall(new(P))
		go srv.ListenAndServe()
		time.Sleep(time.Second)

		cli := multiclient.New(
			drpc.NewEndpoint(drpc.EndpointConfig{}),
			":9090",
			time.Second*5,
		)
		go func() {
			for {
				t.Logf("%d", cli.Size())
				time.Sleep(time.Millisecond * 500)
			}
		}()
		go func() {
			var result int
			for i := 0; ; i++ {
				stat := cli.Call("/p/divide", &Arg{
					A: i,
					B: 2,
				}, &result).Status()
				if !stat.OK() {
					t.Log(stat)
				} else {
					t.Logf("%d/2=%v", i, result)
				}
				time.Sleep(time.Millisecond * 500)
			}
		}()
		time.Sleep(time.Second * 6)
		cli.Close()
		time.Sleep(time.Second * 3)
	})
}

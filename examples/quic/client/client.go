package main

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/logger"
	"time"
)

func main() {

	cli := drpc.NewEndpoint(drpc.EndpointConfig{Network: "quic"})
	defer cli.Close()
	e := cli.SetTLSConfigFromFile(
		fmt.Sprintf("%s/../cert.pem", gfile.MainPkgPath()),
		fmt.Sprintf("%s/../key.pem", gfile.MainPkgPath()),
		true,
	)
	if e != nil {
		logger.Fatalf(context.TODO(), "%v", e)
	}

	cli.RoutePush(new(Push))

	sess, stat := cli.Dial(":9090")
	if !stat.OK() {
		logger.Fatalf(context.TODO(), "%v", stat)
	}

	var result int
	stat = sess.Call("/math/add",
		[]int{1, 2, 3, 4, 5},
		&result,
		drpc.WithSetMeta("author", "clownfish"),
	).Status()
	if !stat.OK() {
		logger.Fatalf(context.TODO(), "%v", stat)
	}
	logger.Printf(context.TODO(), "result: %d", result)

	logger.Printf(context.TODO(), "wait for 10s...")
	time.Sleep(time.Second * 10)
}

// Push push handler
type Push struct {
	drpc.PushCtx
}

// Status handles '/push/status' message
func (p *Push) Status(arg *string) *drpc.Status {
	logger.Printf(context.TODO(), "%s", *arg)
	return nil
}

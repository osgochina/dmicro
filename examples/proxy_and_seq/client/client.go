package main

import (
	"context"
	"fmt"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/logger"
)

//go:generate go build $GOFILE

func main() {
	cli := drpc.NewEndpoint(
		drpc.EndpointConfig{},
	)
	defer cli.Close()

	sess, stat := cli.Dial(":8080")
	if !stat.OK() {
		logger.Fatalf(context.TODO(), "%v", stat)
	}

	var result int
	stat = sess.Call("/math/add",
		[]int{1, 2, 3, 4, 5},
		&result,
	).Status()

	if !stat.OK() {
		logger.Fatalf(context.TODO(), "%v", stat)
	}
	logger.Printf(context.TODO(), "result: %d", result)

	stat = sess.Push(
		"/chat/say",
		fmt.Sprintf("I get result %d", result),
		drpc.WithSetMeta("X-ID", "client-001"),
	)
	if !stat.OK() {
		logger.Fatalf(context.TODO(), "%v", stat)
	}
}

package main

import (
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/logger"
	"time"
)

func main() {

	cli := drpc.NewEndpoint(drpc.EndpointConfig{Network: "tcp", RedialTimes: 1})
	defer cli.Close()

	sess, stat := cli.Dial("127.0.0.1:8199")
	if !stat.OK() {
		logger.Fatalf("%v", stat)
	}

	for true {
		var result string
		stat = sess.Call("/app/home",
			nil,
			&result,
		).Status()
		if !stat.OK() {
			logger.Fatalf("%v", stat)
		}
		logger.Printf("result: %s", result)
		time.Sleep(time.Second)
	}
}

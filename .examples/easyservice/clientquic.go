package main

import (
	"fmt"
	"github.com/gogf/gf/os/gfile"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/logger"
	"time"
)

func main() {

	cli := drpc.NewEndpoint(drpc.EndpointConfig{Network: "quic", RedialTimes: 1})
	defer cli.Close()
	e := cli.SetTLSConfigFromFile(
		fmt.Sprintf("%s/../quic/cert.pem", gfile.MainPkgPath()),
		fmt.Sprintf("%s/../quic/key.pem", gfile.MainPkgPath()),
		true,
	)
	if e != nil {
		logger.Fatalf("%v", e)
	}

	sess, stat := cli.Dial("127.0.0.1:8198")
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

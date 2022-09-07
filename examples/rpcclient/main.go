package main

import (
	"github.com/osgochina/dmicro/client"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/metrics"
	"github.com/osgochina/dmicro/metrics/prometheus"
	"time"
)

func main() {
	c := client.NewRpcClient("test_one",
		client.OptPrintDetail(true),
		client.OptLocalIP("127.0.0.1"),
		client.OptMetrics(prometheus.NewPromMetrics(
			metrics.OptHost("0.0.0.0"),
			metrics.OptPort(9102),
		)),
	)
	defer c.Close()
	for i := 0; i < 100; i++ {
		var result int
		stat := c.Call("/math/add",
			[]int{1, 2, 3, 4, 5},
			&result,
		).Status()
		if stat.OK() {

		}
		logger.Printf("result: %d", result)
		time.Sleep(time.Second * 1)
	}
}

package main

import (
	"fmt"
	"github.com/osgochina/dmicro/easyservice"
	"github.com/osgochina/dmicro/logger"
)

func main() {
	easyservice.Setup(func(svr *easyservice.EasyService) {
		//注册服务停止时要执行法方法
		svr.BeforeStop(func(service *easyservice.EasyService) bool {
			fmt.Println("server stop")
			return true
		})
		logger.Debug("test debug")
	})
}

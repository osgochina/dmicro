package main

import (
	"fmt"
	"github.com/osgochina/dmicro/easyservice"
	"github.com/osgochina/dmicro/logger"
	"time"
)

// 设置日志级别测试
func main() {
	easyservice.Setup(func(svr *easyservice.EasyService) {
		//注册服务停止时要执行法方法
		svr.BeforeStop(func(service *easyservice.EasyService) bool {
			fmt.Println("server stop")
			return true
		})
		for true {
			time.Sleep(3 * time.Second)
			logger.Debug("Debug")
			logger.Info("Info")
			logger.Notice("Notice")
			logger.Warning("Warning")
			logger.Error("Error")
			logger.Critical("Critical")
		}
	})
}

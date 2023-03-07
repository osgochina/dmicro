package main

import (
	"context"
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
			logger.Debug(context.TODO(), "Debug")
			logger.Info(context.TODO(), "Info")
			logger.Notice(context.TODO(), "Notice")
			logger.Warning(context.TODO(), "Warning")
			logger.Error(context.TODO(), "Error")
			logger.Critical(context.TODO(), "Critical")
		}
	})
}

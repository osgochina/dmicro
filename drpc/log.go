package drpc

import (
	"github.com/gogf/gf/os/glog"
	"github.com/osgochina/dmicro/drpc/internal"
)

// SetLogger 使用自定义的log
func SetLogger(l *glog.Logger) {
	internal.SetLogger(l)
}

// GetLogger 获取drpc组件使用的logger对象
func GetLogger() *glog.Logger {
	return internal.GetLogger()
}

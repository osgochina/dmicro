package dserver

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/logger"
)

// 日志级别的转换关系
var levelStringMap = map[string]int{
	"ALL":      glog.LEVEL_DEBU | glog.LEVEL_INFO | glog.LEVEL_NOTI | glog.LEVEL_WARN | glog.LEVEL_ERRO | glog.LEVEL_CRIT,
	"DEV":      glog.LEVEL_DEBU | glog.LEVEL_INFO | glog.LEVEL_NOTI | glog.LEVEL_WARN | glog.LEVEL_ERRO | glog.LEVEL_CRIT,
	"DEVELOP":  glog.LEVEL_DEBU | glog.LEVEL_INFO | glog.LEVEL_NOTI | glog.LEVEL_WARN | glog.LEVEL_ERRO | glog.LEVEL_CRIT,
	"PROD":     glog.LEVEL_WARN | glog.LEVEL_ERRO | glog.LEVEL_CRIT,
	"PRODUCT":  glog.LEVEL_WARN | glog.LEVEL_ERRO | glog.LEVEL_CRIT,
	"DEBU":     glog.LEVEL_DEBU | glog.LEVEL_INFO | glog.LEVEL_NOTI | glog.LEVEL_WARN | glog.LEVEL_ERRO | glog.LEVEL_CRIT,
	"DEBUG":    glog.LEVEL_DEBU | glog.LEVEL_INFO | glog.LEVEL_NOTI | glog.LEVEL_WARN | glog.LEVEL_ERRO | glog.LEVEL_CRIT,
	"INFO":     glog.LEVEL_INFO | glog.LEVEL_NOTI | glog.LEVEL_WARN | glog.LEVEL_ERRO | glog.LEVEL_CRIT,
	"NOTI":     glog.LEVEL_NOTI | glog.LEVEL_WARN | glog.LEVEL_ERRO | glog.LEVEL_CRIT,
	"NOTICE":   glog.LEVEL_NOTI | glog.LEVEL_WARN | glog.LEVEL_ERRO | glog.LEVEL_CRIT,
	"WARN":     glog.LEVEL_WARN | glog.LEVEL_ERRO | glog.LEVEL_CRIT,
	"WARNING":  glog.LEVEL_WARN | glog.LEVEL_ERRO | glog.LEVEL_CRIT,
	"ERRO":     glog.LEVEL_ERRO | glog.LEVEL_CRIT,
	"ERROR":    glog.LEVEL_ERRO | glog.LEVEL_CRIT,
	"CRIT":     glog.LEVEL_CRIT,
	"CRITICAL": glog.LEVEL_CRIT,
}

type ctrlLoggerHandler struct {
	sess   drpc.Session
	logger *glog.Logger
	Level  int
}

func newCtrlLogger(level int, sess drpc.Session) *ctrlLoggerHandler {
	c := &ctrlLoggerHandler{
		sess:  sess,
		Level: level,
	}
	c.logger = glog.NewWithWriter(c)
	return c
}

// Log 实现logger的Handler接口
func (that *ctrlLoggerHandler) Log(level int, v ...interface{}) {
	if that.checkLevel(level) {
		switch level {
		case glog.LEVEL_NOTI:
			that.logger.Notice(context.TODO(), v...)
		case glog.LEVEL_DEBU:
			that.logger.Debug(context.TODO(), v...)
		case glog.LEVEL_INFO:
			that.logger.Info(context.TODO(), v...)
		case glog.LEVEL_NONE:
			that.logger.Print(context.TODO(), v...)
		case glog.LEVEL_WARN:
			that.logger.Warning(context.TODO(), v...)
		case glog.LEVEL_ERRO:
			that.logger.Error(context.TODO(), v...)
		case glog.LEVEL_CRIT:
			that.logger.Critical(context.TODO(), v...)
		}
	}
}

// LogF 实现logger的Handler接口
func (that *ctrlLoggerHandler) LogF(level int, format string, v ...interface{}) {
	if that.checkLevel(level) {
		switch level {
		case glog.LEVEL_NOTI:
			that.logger.Noticef(context.TODO(), format, v...)
		case glog.LEVEL_DEBU:
			that.logger.Debugf(context.TODO(), format, v...)
		case glog.LEVEL_INFO:
			that.logger.Infof(context.TODO(), format, v...)
		case glog.LEVEL_NONE:
			that.logger.Printf(context.TODO(), format, v...)
		case glog.LEVEL_WARN:
			that.logger.Warningf(context.TODO(), format, v...)
		case glog.LEVEL_ERRO:
			that.logger.Errorf(context.TODO(), format, v...)
		case glog.LEVEL_CRIT:
			that.logger.Criticalf(context.TODO(), format, v...)
		}
	}
}

// glog的write接口，实现日志数据转发到ctl客户端
func (that *ctrlLoggerHandler) Write(p []byte) (n int, err error) {
	if that.sess == nil || !that.sess.Health() {
		logger.RemoveHandler("ctl_logger")
		return 0, nil
	}
	that.sess.Push("/ctl_logger_push/logger", p)
	return len(p), nil
}

// 检查是否符合需要打印的日志级别
func (that *ctrlLoggerHandler) checkLevel(level int) bool {
	return that.Level&level > 0
}

// ctl端注册的方法，接收服务端推送的日志信息
type ctrlLoggerPush struct {
	drpc.PushCtx
}

func (that *ctrlLoggerPush) Logger(arg *[]byte) *drpc.Status {
	fmt.Printf("%s", gconv.String(*arg))
	return nil
}

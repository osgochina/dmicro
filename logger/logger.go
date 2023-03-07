package logger

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"sync"
)

type Handler interface {
	Log(level int, v ...interface{})
	LogF(level int, format string, v ...interface{})
}

type logger struct {
	gLog    *glog.Logger
	mu      sync.Mutex
	handler map[string]Handler
}

var log *logger

func init() {
	log = &logger{
		gLog:    glog.DefaultLogger(),
		handler: make(map[string]Handler),
	}
}

func (that *logger) callHandler(level int, format string, v ...interface{}) {
	if len(that.handler) <= 0 {
		return
	}
	for _, h := range that.handler {
		if len(format) > 0 {
			h.LogF(level, format, v)
		} else {
			h.Log(level, v)
		}

	}
}

func DefaultLogger() *glog.Logger {
	return log.gLog
}

func SetConfigWithMap(setting g.Map) error {
	return log.gLog.SetConfigWithMap(setting)
}

func SetLevel(level int) {
	log.gLog.SetLevel(level)
}

func GetLevel() int {
	return log.gLog.GetLevel()
}

func SetLevelStr(levelStr string) error {
	return log.gLog.SetLevelStr(levelStr)
}

// Print prints <v> with newline using fmt.Sprintln.
// The parameter <v> can be multiple variables.
func Print(ctx context.Context, v ...interface{}) {
	log.gLog.Print(ctx, v...)
	log.callHandler(glog.LEVEL_NONE, "", v...)
}

// Printf prints <v> with format <format> using fmt.Sprintf.
// The parameter <v> can be multiple variables.
func Printf(ctx context.Context, format string, v ...interface{}) {
	log.gLog.Printf(ctx, format, v...)
	log.callHandler(glog.LEVEL_NONE, format, v...)
}

// Println See Print.
func Println(ctx context.Context, v ...interface{}) {
	log.gLog.Printf(ctx, "%v\n", v...)
}

// Fatal prints the logging content with [FATA] header and newline, then exit the current process.
func Fatal(ctx context.Context, v ...interface{}) {
	log.gLog.Fatal(ctx, v...)
}

// Fatalf prints the logging content with [FATA] header, custom format and newline, then exit the current process.
func Fatalf(ctx context.Context, format string, v ...interface{}) {
	log.gLog.Fatalf(ctx, format, v...)
}

// Panic prints the logging content with [PANI] header and newline, then panics.
func Panic(ctx context.Context, v ...interface{}) {
	log.gLog.Panic(ctx, v...)
}

// Panicf prints the logging content with [PANI] header, custom format and newline, then panics.
func Panicf(ctx context.Context, format string, v ...interface{}) {
	log.gLog.Panicf(ctx, format, v...)
}

// Info prints the logging content with [INFO] header and newline.
func Info(ctx context.Context, v ...interface{}) {
	log.gLog.Info(ctx, v...)
	log.callHandler(glog.LEVEL_INFO, "", v...)
}

// Infof prints the logging content with [INFO] header, custom format and newline.
func Infof(ctx context.Context, format string, v ...interface{}) {
	log.gLog.Infof(ctx, format, v...)
	log.callHandler(glog.LEVEL_INFO, format, v...)
}

// Debug prints the logging content with [DEBU] header and newline.
func Debug(ctx context.Context, v ...interface{}) {
	log.gLog.Debug(ctx, v...)
	log.callHandler(glog.LEVEL_DEBU, "", v...)
}

// Debugf prints the logging content with [DEBU] header, custom format and newline.
func Debugf(ctx context.Context, format string, v ...interface{}) {
	log.gLog.Debugf(ctx, format, v...)
	log.callHandler(glog.LEVEL_DEBU, format, v...)
}

// Notice prints the logging content with [NOTI] header and newline.
// It also prints caller stack info if stack feature is enabled.
func Notice(ctx context.Context, v ...interface{}) {
	log.gLog.Notice(ctx, v...)
	log.callHandler(glog.LEVEL_NOTI, "", v...)
}

// Noticef prints the logging content with [NOTI] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func Noticef(ctx context.Context, format string, v ...interface{}) {
	log.gLog.Noticef(ctx, format, v...)
	log.callHandler(glog.LEVEL_NOTI, format, v...)
}

// Warning prints the logging content with [WARN] header and newline.
// It also prints caller stack info if stack feature is enabled.
func Warning(ctx context.Context, v ...interface{}) {
	log.gLog.Warning(ctx, v...)
	log.callHandler(glog.LEVEL_WARN, "", v...)
}

// Warningf prints the logging content with [WARN] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func Warningf(ctx context.Context, format string, v ...interface{}) {
	log.gLog.Warningf(ctx, format, v...)
	log.callHandler(glog.LEVEL_WARN, format, v...)
}

// Error prints the logging content with [ERRO] header and newline.
// It also prints caller stack info if stack feature is enabled.
func Error(ctx context.Context, v ...interface{}) {
	log.gLog.Error(ctx, v...)
	log.callHandler(glog.LEVEL_ERRO, "", v...)
}

// Errorf prints the logging content with [ERRO] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func Errorf(ctx context.Context, format string, v ...interface{}) {
	log.gLog.Errorf(ctx, format, v...)
	log.callHandler(glog.LEVEL_ERRO, format, v...)
}

// Critical prints the logging content with [CRIT] header and newline.
// It also prints caller stack info if stack feature is enabled.
func Critical(ctx context.Context, v ...interface{}) {
	log.gLog.Critical(ctx, v...)
	log.callHandler(glog.LEVEL_CRIT, "", v...)
}

// Criticalf prints the logging content with [CRIT] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func Criticalf(ctx context.Context, format string, v ...interface{}) {
	log.gLog.Criticalf(ctx, format, v...)
	log.callHandler(glog.LEVEL_CRIT, format, v...)
}

func SetDebug(debug bool) {
	log.gLog.SetDebug(debug)
}

func SetStdoutPrint(enabled bool) {
	log.gLog.SetStdoutPrint(enabled)
}

func IsDebug() bool {
	return (log.gLog.GetLevel() & glog.LEVEL_DEBU) > 0
}

// AddHandler 添加处理接口
func AddHandler(name string, w Handler) {
	log.mu.Lock()
	defer log.mu.Unlock()
	log.handler[name] = w
}

// RemoveHandler 移除接接口
func RemoveHandler(name string) {
	log.mu.Lock()
	defer log.mu.Unlock()
	delete(log.handler, name)
}

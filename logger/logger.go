package logger

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
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
func Print(v ...interface{}) {
	log.gLog.Print(v...)
	log.callHandler(glog.LEVEL_NONE, "", v...)
}

// Printf prints <v> with format <format> using fmt.Sprintf.
// The parameter <v> can be multiple variables.
func Printf(format string, v ...interface{}) {
	log.gLog.Printf(format, v...)
	log.callHandler(glog.LEVEL_NONE, format, v...)
}

// Println See Print.
func Println(v ...interface{}) {
	log.gLog.Println(v...)
	log.callHandler(glog.LEVEL_NONE, "", v...)
}

// Fatal prints the logging content with [FATA] header and newline, then exit the current process.
func Fatal(v ...interface{}) {
	log.gLog.Fatal(v...)
}

// Fatalf prints the logging content with [FATA] header, custom format and newline, then exit the current process.
func Fatalf(format string, v ...interface{}) {
	log.gLog.Fatalf(format, v...)
}

// Panic prints the logging content with [PANI] header and newline, then panics.
func Panic(v ...interface{}) {
	log.gLog.Panic(v...)
}

// Panicf prints the logging content with [PANI] header, custom format and newline, then panics.
func Panicf(format string, v ...interface{}) {
	log.gLog.Panicf(format, v...)
}

// Info prints the logging content with [INFO] header and newline.
func Info(v ...interface{}) {
	log.gLog.Info(v...)
	log.callHandler(glog.LEVEL_INFO, "", v...)
}

// Infof prints the logging content with [INFO] header, custom format and newline.
func Infof(format string, v ...interface{}) {
	log.gLog.Infof(format, v...)
	log.callHandler(glog.LEVEL_INFO, format, v...)
}

// Debug prints the logging content with [DEBU] header and newline.
func Debug(v ...interface{}) {
	log.gLog.Debug(v...)
	log.callHandler(glog.LEVEL_DEBU, "", v...)
}

// Debugf prints the logging content with [DEBU] header, custom format and newline.
func Debugf(format string, v ...interface{}) {
	log.gLog.Debugf(format, v...)
	log.callHandler(glog.LEVEL_DEBU, format, v...)
}

// Notice prints the logging content with [NOTI] header and newline.
// It also prints caller stack info if stack feature is enabled.
func Notice(v ...interface{}) {
	log.gLog.Notice(v...)
	log.callHandler(glog.LEVEL_NOTI, "", v...)
}

// Noticef prints the logging content with [NOTI] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func Noticef(format string, v ...interface{}) {
	log.gLog.Noticef(format, v...)
	log.callHandler(glog.LEVEL_NOTI, format, v...)
}

// Warning prints the logging content with [WARN] header and newline.
// It also prints caller stack info if stack feature is enabled.
func Warning(v ...interface{}) {
	log.gLog.Warning(v...)
	log.callHandler(glog.LEVEL_WARN, "", v...)
}

// Warningf prints the logging content with [WARN] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func Warningf(format string, v ...interface{}) {
	log.gLog.Warningf(format, v...)
	log.callHandler(glog.LEVEL_WARN, format, v...)
}

// Error prints the logging content with [ERRO] header and newline.
// It also prints caller stack info if stack feature is enabled.
func Error(v ...interface{}) {
	log.gLog.Error(v...)
	log.callHandler(glog.LEVEL_ERRO, "", v...)
}

// Errorf prints the logging content with [ERRO] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func Errorf(format string, v ...interface{}) {
	log.gLog.Errorf(format, v...)
	log.callHandler(glog.LEVEL_ERRO, format, v...)
}

// Critical prints the logging content with [CRIT] header and newline.
// It also prints caller stack info if stack feature is enabled.
func Critical(v ...interface{}) {
	log.gLog.Critical(v...)
	log.callHandler(glog.LEVEL_CRIT, "", v...)
}

// Criticalf prints the logging content with [CRIT] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func Criticalf(format string, v ...interface{}) {
	log.gLog.Criticalf(format, v...)
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

package logger

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"io"
	"sync"
)

type logger struct {
	logger  *glog.Logger
	handler map[string]io.Writer
	mu      sync.Mutex
}

func (that *logger) Write(p []byte) (n int, err error) {
	//fmt.Printf("write: %s", gconv.String(p))
	// 如果有handler
	if len(that.handler) > 0 {
		for _, w := range that.handler {
			_, _ = w.Write(p)
		}
	}
	return that.logger.Write(p)
}

var log = glog.New()
var myLog = &logger{logger: glog.DefaultLogger(), handler: make(map[string]io.Writer)}

func init() {
	log.SetWriter(myLog)
}

func DefaultLogger() *glog.Logger {
	return log
}

func SetConfigWithMap(setting g.Map) error {
	return log.SetConfigWithMap(setting)
}

func SetLevel(level int) {
	log.SetLevel(level)
}

func GetLevel() int {
	return log.GetLevel()
}

func SetLevelStr(levelStr string) error {
	return log.SetLevelStr(levelStr)
}

// Print prints <v> with newline using fmt.Sprintln.
// The parameter <v> can be multiple variables.
func Print(v ...interface{}) {
	log.Print(v...)
}

// Printf prints <v> with format <format> using fmt.Sprintf.
// The parameter <v> can be multiple variables.
func Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

// See Print.
func Println(v ...interface{}) {
	log.Println(v...)
}

// Fatal prints the logging content with [FATA] header and newline, then exit the current process.
func Fatal(v ...interface{}) {
	log.Fatal(v...)
}

// Fatalf prints the logging content with [FATA] header, custom format and newline, then exit the current process.
func Fatalf(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}

// Panic prints the logging content with [PANI] header and newline, then panics.
func Panic(v ...interface{}) {
	log.Panic(v...)
}

// Panicf prints the logging content with [PANI] header, custom format and newline, then panics.
func Panicf(format string, v ...interface{}) {
	log.Panicf(format, v...)
}

// Info prints the logging content with [INFO] header and newline.
func Info(v ...interface{}) {
	log.Info(v...)
}

// Infof prints the logging content with [INFO] header, custom format and newline.
func Infof(format string, v ...interface{}) {
	log.Infof(format, v...)
}

// Debug prints the logging content with [DEBU] header and newline.
func Debug(v ...interface{}) {
	log.Debug(v...)
}

// Debugf prints the logging content with [DEBU] header, custom format and newline.
func Debugf(format string, v ...interface{}) {
	log.Debugf(format, v...)
}

// Notice prints the logging content with [NOTI] header and newline.
// It also prints caller stack info if stack feature is enabled.
func Notice(v ...interface{}) {
	log.Notice(v...)
}

// Noticef prints the logging content with [NOTI] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func Noticef(format string, v ...interface{}) {
	log.Noticef(format, v...)
}

// Warning prints the logging content with [WARN] header and newline.
// It also prints caller stack info if stack feature is enabled.
func Warning(v ...interface{}) {
	log.Warning(v...)
}

// Warningf prints the logging content with [WARN] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func Warningf(format string, v ...interface{}) {
	log.Warningf(format, v...)
}

// Error prints the logging content with [ERRO] header and newline.
// It also prints caller stack info if stack feature is enabled.
func Error(v ...interface{}) {
	log.Error(v...)
}

// Errorf prints the logging content with [ERRO] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func Errorf(format string, v ...interface{}) {
	log.Errorf(format, v...)
}

// Critical prints the logging content with [CRIT] header and newline.
// It also prints caller stack info if stack feature is enabled.
func Critical(v ...interface{}) {
	log.Critical(v...)
}

// Criticalf prints the logging content with [CRIT] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func Criticalf(format string, v ...interface{}) {
	log.Criticalf(format, v...)
}

func SetDebug(debug bool) {
	log.SetDebug(debug)
}

func SetStdoutPrint(enabled bool) {
	log.SetStdoutPrint(enabled)
}

func IsDebug() bool {
	return (log.GetLevel() & glog.LEVEL_DEBU) > 0
}

// AddHandler 添加处理接口
func AddHandler(name string, w io.Writer) {
	myLog.mu.Lock()
	defer myLog.mu.Unlock()
	myLog.handler[name] = w
}

// RemoveHandler 移除接接口
func RemoveHandler(name string) {
	myLog.mu.Lock()
	defer myLog.mu.Unlock()
	delete(myLog.handler, name)
}

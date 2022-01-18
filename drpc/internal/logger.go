package internal

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
)

var gLog *glog.Logger

func init() {
	gLog = glog.DefaultLogger()
}

func GetLogger() *glog.Logger {
	return gLog
}

func SetLogger(l *glog.Logger) {
	gLog = l
}

func SetConfigWithMap(setting g.Map) error {
	return gLog.SetConfigWithMap(setting)
}

func SetLevel(level int) {
	gLog.SetLevel(level)
}

func GetLevel() int {
	return gLog.GetLevel()
}

func SetLevelStr(levelStr string) error {
	return gLog.SetLevelStr(levelStr)
}

// Print prints <v> with newline using fmt.Sprintln.
// The parameter <v> can be multiple variables.
func Print(v ...interface{}) {
	gLog.Print(v...)
}

// Printf prints <v> with format <format> using fmt.Sprintf.
// The parameter <v> can be multiple variables.
func Printf(format string, v ...interface{}) {
	gLog.Printf(format, v...)
}

func Println(v ...interface{}) {
	gLog.Println(v...)
}

// Fatal prints the logging content with [FATA] header and newline, then exit the current process.
func Fatal(v ...interface{}) {
	gLog.Fatal(v...)
}

// Fatalf prints the logging content with [FATA] header, custom format and newline, then exit the current process.
func Fatalf(format string, v ...interface{}) {
	gLog.Fatalf(format, v...)
}

// Panic prints the logging content with [PANI] header and newline, then panics.
func Panic(v ...interface{}) {
	gLog.Panic(v...)
}

func Panicf(format string, v ...interface{}) {
	gLog.Panicf(format, v...)
}

// Info prints the logging content with [INFO] header and newline.
func Info(v ...interface{}) {
	gLog.Info(v...)
}

func Infof(format string, v ...interface{}) {
	gLog.Infof(format, v...)
}

// Debug prints the logging content with [DEBU] header and newline.
func Debug(v ...interface{}) {
	gLog.Debug(v...)
}

// Debugf prints the logging content with [DEBU] header, custom format and newline.
func Debugf(format string, v ...interface{}) {
	gLog.Debugf(format, v...)
}

// Notice prints the logging content with [NOTI] header and newline.
// It also prints caller stack info if stack feature is enabled.
func Notice(v ...interface{}) {
	gLog.Notice(v...)
}

// Noticef prints the logging content with [NOTI] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func Noticef(format string, v ...interface{}) {
	gLog.Noticef(format, v...)
}

// Warning prints the logging content with [WARN] header and newline.
// It also prints caller stack info if stack feature is enabled.
func Warning(v ...interface{}) {
	gLog.Warning(v...)
}

// Warningf prints the logging content with [WARN] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func Warningf(format string, v ...interface{}) {
	gLog.Warningf(format, v...)
}

// Error prints the logging content with [ERRO] header and newline.
// It also prints caller stack info if stack feature is enabled.
func Error(v ...interface{}) {
	gLog.Error(v...)
}

// Errorf prints the logging content with [ERRO] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func Errorf(format string, v ...interface{}) {
	gLog.Errorf(format, v...)
}

// Critical prints the logging content with [CRIT] header and newline.
// It also prints caller stack info if stack feature is enabled.
func Critical(v ...interface{}) {
	gLog.Critical(v...)
}

// Criticalf prints the logging content with [CRIT] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func Criticalf(format string, v ...interface{}) {
	gLog.Criticalf(format, v...)
}

func SetDebug(debug bool) {
	gLog.SetDebug(debug)
}

func SetStdoutPrint(enabled bool) {
	gLog.SetStdoutPrint(enabled)
}

func IsDebug() bool {
	return (gLog.GetLevel() & glog.LEVEL_DEBU) > 0
}

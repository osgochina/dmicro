package logger

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
)

type logger struct {
	gLog *glog.Logger
}

var log *logger

func init() {
	log = &logger{gLog: glog.DefaultLogger()}
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
}

// Printf prints <v> with format <format> using fmt.Sprintf.
// The parameter <v> can be multiple variables.
func Printf(format string, v ...interface{}) {
	log.gLog.Printf(format, v...)
}

// See Print.
func Println(v ...interface{}) {
	log.gLog.Println(v...)
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
}

// Infof prints the logging content with [INFO] header, custom format and newline.
func Infof(format string, v ...interface{}) {
	log.gLog.Infof(format, v...)
}

// Debug prints the logging content with [DEBU] header and newline.
func Debug(v ...interface{}) {
	log.gLog.Debug(v...)
}

// Debugf prints the logging content with [DEBU] header, custom format and newline.
func Debugf(format string, v ...interface{}) {
	log.gLog.Debugf(format, v...)
}

// Notice prints the logging content with [NOTI] header and newline.
// It also prints caller stack info if stack feature is enabled.
func Notice(v ...interface{}) {
	log.gLog.Notice(v...)
}

// Noticef prints the logging content with [NOTI] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func Noticef(format string, v ...interface{}) {
	log.gLog.Noticef(format, v...)
}

// Warning prints the logging content with [WARN] header and newline.
// It also prints caller stack info if stack feature is enabled.
func Warning(v ...interface{}) {
	log.gLog.Warning(v...)
}

// Warningf prints the logging content with [WARN] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func Warningf(format string, v ...interface{}) {
	log.gLog.Warningf(format, v...)
}

// Error prints the logging content with [ERRO] header and newline.
// It also prints caller stack info if stack feature is enabled.
func Error(v ...interface{}) {
	log.gLog.Error(v...)
}

// Errorf prints the logging content with [ERRO] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func Errorf(format string, v ...interface{}) {
	log.gLog.Errorf(format, v...)
}

// Critical prints the logging content with [CRIT] header and newline.
// It also prints caller stack info if stack feature is enabled.
func Critical(v ...interface{}) {
	log.gLog.Critical(v...)
}

// Criticalf prints the logging content with [CRIT] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func Criticalf(format string, v ...interface{}) {
	log.gLog.Criticalf(format, v...)
}

func SetDebug(debug bool) {
	log.gLog.SetDebug(debug)
}

func IsDebug() bool {
	return (log.gLog.GetLevel() & glog.LEVEL_DEBU) > 0
}

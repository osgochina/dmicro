package internal

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
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
func Print(ctx context.Context, v ...interface{}) {
	gLog.Print(ctx, v...)
}

// Printf prints <v> with format <format> using fmt.Sprintf.
// The parameter <v> can be multiple variables.
func Printf(ctx context.Context, format string, v ...interface{}) {
	gLog.Printf(ctx, format, v...)
}

func Println(ctx context.Context, v ...interface{}) {
	gLog.Printf(ctx, "%v\n", v...)
}

// Fatal prints the logging content with [FATA] header and newline, then exit the current process.
func Fatal(ctx context.Context, v ...interface{}) {
	gLog.Fatal(ctx, v...)
}

// Fatalf prints the logging content with [FATA] header, custom format and newline, then exit the current process.
func Fatalf(ctx context.Context, format string, v ...interface{}) {
	gLog.Fatalf(ctx, format, v...)
}

// Panic prints the logging content with [PANI] header and newline, then panics.
func Panic(ctx context.Context, v ...interface{}) {
	gLog.Panic(ctx, v...)
}

func Panicf(ctx context.Context, format string, v ...interface{}) {
	gLog.Panicf(ctx, format, v...)
}

// Info prints the logging content with [INFO] header and newline.
func Info(ctx context.Context, v ...interface{}) {
	gLog.Info(ctx, v...)
}

func Infof(ctx context.Context, format string, v ...interface{}) {
	gLog.Infof(ctx, format, v...)
}

// Debug prints the logging content with [DEBU] header and newline.
func Debug(ctx context.Context, v ...interface{}) {
	gLog.Debug(ctx, v...)
}

// Debugf prints the logging content with [DEBU] header, custom format and newline.
func Debugf(ctx context.Context, format string, v ...interface{}) {
	gLog.Debugf(ctx, format, v...)
}

// Notice prints the logging content with [NOTI] header and newline.
// It also prints caller stack info if stack feature is enabled.
func Notice(ctx context.Context, v ...interface{}) {
	gLog.Notice(ctx, v...)
}

// Noticef prints the logging content with [NOTI] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func Noticef(ctx context.Context, format string, v ...interface{}) {
	gLog.Noticef(ctx, format, v...)
}

// Warning prints the logging content with [WARN] header and newline.
// It also prints caller stack info if stack feature is enabled.
func Warning(ctx context.Context, v ...interface{}) {
	gLog.Warning(ctx, v...)
}

// Warningf prints the logging content with [WARN] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func Warningf(ctx context.Context, format string, v ...interface{}) {
	gLog.Warningf(ctx, format, v...)
}

// Error prints the logging content with [ERRO] header and newline.
// It also prints caller stack info if stack feature is enabled.
func Error(ctx context.Context, v ...interface{}) {
	gLog.Error(ctx, v...)
}

// Errorf prints the logging content with [ERRO] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func Errorf(ctx context.Context, format string, v ...interface{}) {
	gLog.Errorf(ctx, format, v...)
}

// Critical prints the logging content with [CRIT] header and newline.
// It also prints caller stack info if stack feature is enabled.
func Critical(ctx context.Context, v ...interface{}) {
	gLog.Critical(ctx, v...)
}

// Criticalf prints the logging content with [CRIT] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func Criticalf(ctx context.Context, format string, v ...interface{}) {
	gLog.Criticalf(ctx, format, v...)
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

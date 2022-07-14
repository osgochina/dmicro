package dserver

import (
	"fmt"
	"github.com/gogf/gf/util/gconv"
	"github.com/osgochina/dmicro/drpc"
)

type ctrlLogger struct {
	sess drpc.Session
}

func (that *ctrlLogger) Write(p []byte) (n int, err error) {
	//that.sess.Push("logger",)
	fmt.Printf("ctrlLogger: %s", gconv.String(p))
	that.sess.Push("/ctrl_logger_push/logger", p)
	return len(p), nil
}

type ctrlLoggerPush struct {
	drpc.PushCtx
}

func (that *ctrlLoggerPush) Logger(arg *[]byte) *drpc.Status {
	fmt.Printf("%s", gconv.String(*arg))
	return nil
}

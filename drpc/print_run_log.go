package drpc

import (
	"encoding/json"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
	"github.com/osgochina/dmicro/drpc/internal"
	"github.com/osgochina/dmicro/drpc/message"
	"strconv"
	"time"
)

const (
	typePushLaunch int8 = 1
	typePushHandle int8 = 2
	typeCallLaunch int8 = 3
	typeCallHandle int8 = 4
)

const (
	logFormatPushLaunch = "PUSH-> %s %s %q SEND(%s)"
	logFormatPushHandle = "PUSH<- %s %s %q RECV(%s)"
	logFormatCallLaunch = "CALL-> %s %s %q SEND(%s) RECV(%s)"
	logFormatCallHandle = "CALL<- %s %s %q RECV(%s) SEND(%s)"
)

func enablePrintRunLog() bool {
	return internal.GetLevel()&glog.LEVEL_DEBU > 0
}

//打印运行log
func (that *session) printRunLog(realIP string, costTime time.Duration, input, output message.Message, logType int8) {
	var addr = that.RemoteAddr().String()
	if realIP != "" && realIP == addr {
		realIP = "same"
	}
	if realIP == "" {
		realIP = "-"
	}
	addr += "(real:" + realIP + ")"
	var (
		costTimeStr string
		printFunc   = internal.Debugf
	)
	if that.endpoint.countTime {
		if costTime >= that.endpoint.slowCometDuration {
			costTimeStr = costTime.String() + "(slow)"
		} else {
			costTimeStr = costTime.String() + "(fast)"
		}
	} else {
		costTimeStr = "(-)"
	}

	switch logType {
	case typePushLaunch:
		printFunc(logFormatPushLaunch, addr, costTimeStr, output.ServiceMethod(), messageLogBytes(output, that.endpoint.printDetail))
	case typePushHandle:
		printFunc(logFormatPushHandle, addr, costTimeStr, input.ServiceMethod(), messageLogBytes(input, that.endpoint.printDetail))
	case typeCallLaunch:
		printFunc(logFormatCallLaunch, addr, costTimeStr, output.ServiceMethod(), messageLogBytes(output, that.endpoint.printDetail), messageLogBytes(input, that.endpoint.printDetail))
	case typeCallHandle:
		printFunc(logFormatCallHandle, addr, costTimeStr, input.ServiceMethod(), messageLogBytes(input, that.endpoint.printDetail), messageLogBytes(output, that.endpoint.printDetail))
	}
}

func messageLogBytes(message message.Message, printDetail bool) []byte {
	var b = make([]byte, 0, 128)
	b = append(b, '{')
	b = append(b, '"', 's', 'i', 'z', 'e', '"', ':')
	b = append(b, strconv.FormatUint(uint64(message.Size()), 10)...)
	if statBytes := message.Status().String(); len(statBytes) > 0 {
		b = append(b, ',', '"', 's', 't', 'a', 't', 'u', 's', '"', ':')
		b = append(b, statBytes...)
	}
	if printDetail {
		if message.Meta().Size() > 0 {
			b = append(b, ',', '"', 'm', 'e', 't', 'a', '"', ':')
			b = append(b, gconv.Bytes(message.Meta().String())...)
		}
		if bodyBytes := bodyLogBytes(message); len(bodyBytes) > 0 {
			b = append(b, ',', '"', 'b', 'o', 'd', 'y', '"', ':')
			b = append(b, bodyBytes...)
		}
	}
	b = append(b, '}')
	return b
}

func bodyLogBytes(message message.Message) []byte {
	switch v := message.Body().(type) {
	case nil:
		return nil
	case []byte:
		return gconv.Bytes(v)
	case *[]byte:
		return gconv.Bytes(v)
	}
	b, _ := json.Marshal(message.Body())
	return b
}

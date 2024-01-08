package dserver

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/genv"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/utils/signals"
	"os"
)

// 停止服务
func (that *DServer) stop(signal string) {
	pidFile := that.pidFile
	var serverPid = 0
	if gfile.IsFile(pidFile) {
		serverPid = gconv.Int(gstr.Trim(gfile.GetContents(pidFile)))
	}
	if serverPid == 0 {
		logger.Println(context.TODO(), "Server is not running.")
		os.Exit(0)
	}

	var sigNo string
	switch signal {
	case "stop":
		sigNo = "SIGTERM"
	case "reload":
		sigNo = "SIGUSR2"
	case "quit":
		sigNo = "SIGQUIT"
	default:
		logger.Printf(context.TODO(), "signal cmd `%s' not found", signal)
		os.Exit(0)
	}
	err := signals.KillPid(serverPid, signals.ToSignal(sigNo), false)
	if err != nil {
		logger.Printf(context.TODO(), "error:%v", err)
	}
	os.Exit(0)
}

// 初始化要启动的服务名
func (that *DServer) initSandboxNames(names []string) {
	// 获取要启动的sandbox名称
	for _, sandboxName := range names {
		if gstr.ContainsI(sandboxName, ",") {
			sandboxNameArray := gstr.Split(sandboxName, ",")
			if len(sandboxNameArray) > 1 {
				that.sandboxNames.Append(sandboxNameArray...)
			}
			continue
		}
		that.sandboxNames.Append(gstr.Trim(sandboxName))
	}
}

// 初始化pid文件的地址
func (that *DServer) initPidFile(pidPath string) {
	if len(pidPath) > 0 {
		that.pidFile = pidPath
		return
	}
	pidFile := fmt.Sprintf("%s.pid", that.name)
	that.pidFile = gfile.Temp(pidFile)
}

// 检查服务是否已经启动
func (that *DServer) checkStart() {
	pidFile := that.pidFile
	var serverPid = 0
	if gfile.IsFile(pidFile) {
		serverPid = gconv.Int(gstr.Trim(gfile.GetContents(pidFile)))
	}
	if serverPid == 0 {
		return
	}
	if signals.CheckPidExist(serverPid) {
		logger.Fatalf(context.TODO(), "Server [%d] is already running.", serverPid)
	}
	return
}

// 解析配置文件信息
func (that *DServer) parserConfig(config string) {
	that.initGFConf(config)
	// 设置配置文件中log的配置
	that.initLogCfg()
}

// 是否守护进程启动
func (that *DServer) parserDaemon(daemon bool) {
	if daemon {
		_ = genv.Set("DAEMON", "true")
	} else {
		_ = genv.Set("DAEMON", "false")
	}
}

// 解析运行环境
func (that *DServer) parserEnv(env string) {
	//如果命令行传入了env参数，则使用命令行参数
	if len(env) > 0 {
		_ = genv.Set("ENV_NAME", gstr.ToLower(env))
	} else {
		_ = genv.Set("ENV_NAME", "product")
	}
}

// 解析debug参数
func (that *DServer) parserDebug(debug bool) {
	if debug {
		// 如果启动命令行强制设置了debug参数，则优先级最高
		_ = genv.Set("DEBUG", "true")
	} else {
		_ = genv.Set("DEBUG", "false")
	}
}

const configNodeNameLogger = "logger"

// 把配置文件中的配置信息写入到logger配置中
func (that *DServer) initLogCfg() {
	if !g.Cfg().Available(context.TODO()) {
		return
	}
	var m map[string]interface{}
	nodeKey, _ := gutil.MapPossibleItemByKey(g.Cfg().MustGet(context.TODO(), ".").Map(), configNodeNameLogger)
	if nodeKey == "" {
		nodeKey = configNodeNameLogger
	}
	m = g.Cfg().MustGet(context.TODO(), fmt.Sprintf(`%s.%s`, nodeKey, glog.DefaultName)).Map()
	if len(m) == 0 {
		m = g.Cfg().MustGet(context.TODO(), nodeKey).Map()
	}
	if len(m) > 0 {
		if err := logger.SetConfigWithMap(m); err != nil {
			panic(err)
		}
	}
}

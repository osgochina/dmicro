package dserver

import (
	"fmt"
	"github.com/gogf/gf/os/genv"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gutil"
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
		logger.Println("Server is not running.")
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
		logger.Printf("signal cmd `%s' not found", signal)
		os.Exit(0)
	}
	err := signals.KillPid(serverPid, signals.ToSignal(sigNo), false)
	if err != nil {
		logger.Printf("error:%v", err)
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
	that.pidFile = gfile.TempDir(pidFile)
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
		logger.Fatalf("Server [%d] is already running.", serverPid)
	}
	return
}

// 解析配置文件信息
func (that *DServer) parserConfig(config string) {
	that.config = that.getGFConf(config)
	// 设置配置文件中log的配置
	that.initLogCfg()
}

// 是否守护进程启动
func (that *DServer) parserDaemon(daemon bool) {
	_ = that.config.Set("Daemon", daemon)
}

// 解析运行环境
func (that *DServer) parserEnv(env string) {
	//如果命令行传入了env参数，则使用命令行参数
	if len(env) > 0 {
		_ = genv.Set("ENV_NAME", gstr.ToLower(env))
		_ = that.config.Set("ENV_NAME", gstr.ToLower(env))
	} else if len(that.config.GetString("ENV_NAME")) <= 0 {
		//如果命令行未传入env参数,且配置文件中页不存在ENV_NAME配置，则先查找环境变量ENV_NAME，并把环境变量中的ENV_NAME赋值给配置文件
		_ = that.config.Set("ENV_NAME", gstr.ToLower(genv.Get("ENV_NAME", "product")))
	}
}

// 解析debug参数
func (that *DServer) parserDebug(debug bool) {
	if debug {
		// 如果启动命令行强制设置了debug参数，则优先级最高
		_ = genv.Set("DEBUG", "true")
		_ = that.config.Set("Debug", true)
	} else {
		// 1. 从配置文件中获取debug参数,如果获取到则使用，未获取到这进行下一步
		// 2. 先从环境变量获取debug参数
		// 3. 最终传导获取到debug值，把它设置到配置文件中
		_ = that.config.Set("Debug", that.config.GetBool("Debug", genv.GetVar("DEBUG", false).Bool()))
	}
}

const configNodeNameLogger = "logger"

// 把配置文件中的配置信息写入到logger配置中
func (that *DServer) initLogCfg() {
	if !that.config.Available() {
		return
	}
	var m map[string]interface{}
	nodeKey, _ := gutil.MapPossibleItemByKey(that.config.GetMap("."), configNodeNameLogger)
	if nodeKey == "" {
		nodeKey = configNodeNameLogger
	}
	m = that.config.GetMap(fmt.Sprintf(`%s.%s`, nodeKey, glog.DefaultName))
	if len(m) == 0 {
		m = that.config.GetMap(nodeKey)
	}
	if len(m) > 0 {
		if err := logger.SetConfigWithMap(m); err != nil {
			panic(err)
		}
	}
}

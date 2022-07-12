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

var defaultOptions = map[string]bool{
	"c,config": true,
	"env":      true,
	"pid":      true,
	"debug":    true,
	"m,model":  true,
}

// SetOptions 添加自定义的参数解析
func SetOptions(opt map[string]bool) {
	for k, v := range opt {
		defaultOptions[k] = v
	}
}

// SetOption 添加单个自定义的参数解析
func SetOption(key string, v bool) {
	defaultOptions[key] = v
}

var (
	helpContent = gstr.TrimLeft(`
USAGE
	%path% start|stop|reload|quit [OPTION] [sandboxName1|sandboxName2...] 
OPTION
	-c,--config     指定要载入的配置文件，该参数与gf.gcfg.file参数二选一，建议使用该参数
	-d,--daemon     使用守护进程模式启动
	--env           环境变量，表示当前启动所在的环境,有[dev,test,product]这三种，默认是product
	--debug         是否开启debug 默认debug=false
	--pid           设置pid文件的地址，默认是/tmp/[server].pid
	-h,--help       获取帮助信息
	-v,--version    获取编译版本信息
	-m,--model      进程模型，0表示单进程模型，1表示多进程模型
	
EXAMPLES
	%path% 
	%path% start --env=dev --debug=true --pid=/tmp/server.pid
	%path% start --gf.gcfg.file=config.product.toml
	%path% start -c=config.product.toml
	%path% start --config=config.product.toml user admin
	%path% start user
	%path% stop
	%path% quit
	%path% reload
	%path% version
	%path% help
`)
)

// SetHelpContent 自定义帮助信息
func SetHelpContent(content string) {
	helpContent = content
}

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
func (that *DServer) initPidFile(pid string) {
	if len(pid) > 0 {
		that.pidFile = pid
		return
	}
	pidFile := fmt.Sprintf("%s.pid", gfile.Basename(os.Args[0]))
	if that.sandboxNames.Len() > 0 {
		pidFile = fmt.Sprintf("%s.pid", that.sandboxNames.Join("-"))
	}
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
		// 1. 从命令行中获取debug参数,如果获取到则使用，未获取到这进行下一步
		// 2. 从配置文件中获取debug参数,如果获取到则使用，未获取到这进行下一步
		// 3. 先从环境变量获取debug参数
		// 4. 最终传导获取到debug值，把它设置到配置文件中
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

//Help 显示帮助信息
func (that *DServer) Help() {
	fmt.Print(gstr.Replace(helpContent, "%path%", gfile.Abs(os.Args[0])))
}

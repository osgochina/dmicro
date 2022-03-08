package easyservice

import (
	"fmt"
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/os/genv"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gutil"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/utils/signals"
	"os"
	"strings"
)

var defaultOptions = map[string]bool{
	"c,config": true,
	"env":      true,
	"pid":      true,
	"debug":    true,
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
	./server [start|stop|quit] [default|custom] [OPTION]
OPTION
	-c,--config     指定要载入的配置文件，该参数与gf.gcfg.file参数二选一，建议使用该参数
	-d,--daemon     使用守护进程模式启动
	--env           环境变量，表示当前启动所在的环境,有[dev,test,product]这三种，默认是product
	--debug         是否开启debug 默认debug=false
	--pid           设置pid文件的地址，默认是/tmp/[server].pid
	-h,--help       获取帮助信息
	-v,--version    获取编译版本信息
	
EXAMPLES
	/path/to/server 
	/path/to/server start --env=dev --debug=true --pid=/tmp/server.pid
	/path/to/server start --gf.gcfg.file=config.product.toml
	/path/to/server start -c=config.product.toml
	/path/to/server start user,admin --config=config.product.toml
	/path/to/server start user
	/path/to/server stop
	/path/to/server quit
	/path/to/server reload
	/path/to/server version
	/path/to/server help
`)
)

// SetHelpContent 自定义帮助信息
func SetHelpContent(content string) {
	helpContent = content
}

// 解析命令行，根据返回值判断是否继续执行
// 返回false，则结束进程，返回true继续执行
func (that *EasyService) parserArgs(parser *gcmd.Parser) bool {
	command := gcmd.GetArg(1)
	switch strings.ToLower(command) {
	case "help":
		that.help()
		return false
	case "version":
		that.version()
		return false
	case "start":
		that.initSandboxNames()
		that.initPidFile(parser)
		that.checkStart()
		return true
	case "stop":
		that.initSandboxNames()
		that.initPidFile(parser)
		that.stop("stop")
		return false
	case "reload":
		that.initSandboxNames()
		that.initPidFile(parser)
		that.stop("reload")
		return false
	case "quit":
		that.initSandboxNames()
		that.initPidFile(parser)
		that.stop("quit")
		return false
	default:
		for _, v := range gcmd.GetArgAll() {
			switch v {
			case "?", "h":
				that.help()
				return false
			case "i", "v":
				that.version()
				return false
			}
		}
		// 识别参数展示帮助信息和版本信息
		array := garray.NewStrArrayFrom(os.Args)
		if array.Search("--help") != -1 || array.Search("-h") != -1 {
			that.help()
			return false
		}
		if array.Search("--version") != -1 || array.Search("-v") != -1 {
			that.version()
			return false
		}
	}
	that.initSandboxNames()
	that.initPidFile(parser)
	that.checkStart()
	return true
}

// 停止服务
func (that *EasyService) stop(signal string) {
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
func (that *EasyService) initSandboxNames() {
	command := gcmd.GetArg(1)
	switch strings.ToLower(command) {
	case "start":
		fallthrough
	case "stop":
		fallthrough
	case "reload":
		fallthrough
	case "quit":
		// 获取要启动的服务名，并存储
		sandboxNames := gcmd.GetArg(2)
		if len(sandboxNames) > 0 {
			sandboxNames = gstr.Trim(sandboxNames)
			sandboxNameArray := gstr.Split(sandboxNames, ",")
			if len(sandboxNameArray) > 1 {
				that.sandboxNames.Append(sandboxNameArray...)
			} else {
				that.sandboxNames.Append(sandboxNames)
			}
		}
	}
	return
}

// 初始化pid文件地址
func (that *EasyService) initPidFile(parser *gcmd.Parser) {
	pidFile := "easyservice.pid"
	if that.sandboxNames.Len() > 0 {
		pidFile = fmt.Sprintf("%s.pid", that.sandboxNames.Join("-"))
	}
	that.pidFile = parser.GetOpt("pid", gfile.TempDir(pidFile))
}

// 检查服务是否已经启动
func (that *EasyService) checkStart() {
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

//解析配置文件
func (that *EasyService) parserConfig(parser *gcmd.Parser) {
	that.config = that.getGFConf(parser)
	// 设置配置文件中log的配置
	that.initLogCfg()
	array := garray.NewStrArrayFrom(os.Args)
	//判断是否需要后台运行
	index := array.Search("--daemon")
	if index != -1 {
		_ = that.config.Set("Daemon", true)
	}
	index = array.Search("-d")
	if index != -1 {
		_ = that.config.Set("Daemon", true)
	}
	//通过命令行传入环境参数
	env := parser.GetOpt("env", "")
	//如果命令行传入了env参数，则使用命令行参数
	if len(env) > 0 {
		_ = genv.Set("ENV_NAME", gstr.ToLower(env))
		_ = that.config.Set("ENV_NAME", gstr.ToLower(env))
	} else if len(that.config.GetString("ENV_NAME")) <= 0 {
		//如果命令行未传入env参数,且配置文件中页不存在ENV_NAME配置，则先查找环境变量ENV_NAME，并把环境变量中的ENV_NAME赋值给配置文件
		_ = that.config.Set("ENV_NAME", gstr.ToLower(genv.Get("ENV_NAME", "product")))
	}
	//通过启动命令判断是否开启debug
	index = array.Search("--debug")
	if index != -1 {
		_ = that.config.Set("Debug", true)
	} else {
		debug := parser.GetOptVar("debug", that.config.GetBool("Debug", false))
		_ = that.config.Set("Debug", debug.Bool())
	}
}

const configNodeNameLogger = "logger"

// 把配置文件中的配置信息写入到logger配置中
func (that *EasyService) initLogCfg() {
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

//显示帮助信息
func (that *EasyService) help() {
	fmt.Print(helpContent)
}

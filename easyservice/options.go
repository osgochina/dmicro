package easyservice

import (
	"fmt"
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/os/gcfg"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/os/genv"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"github.com/osgochina/dmicro/logger"
	"os"
	"strings"
	"syscall"
)

var defaultOptions = map[string]bool{
	"h,host": true,
	"p,port": true,
	"c,conf": true,
	"env":    true,
	"pid":    true,
	"debug":  true,
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
	./server start|stop|quit [default|custom] [OPTION]
OPTION
	-h,--host	服务监听地址，默认监听的地址为127.0.0.1
	-p,--port	服务监听端口，默认监听端口为0，表示随机监听
	-d,--daemon	debug模式开关，默认关闭debug=false
	--gf.gcfg.file  需要加载的配置文件名 如 config.dev.toml
	--env		环境变量，表示当前启动所在的环境,有[dev,test,product]这三种，默认是product
	--debug		default debug=false
	--pid		设置pid文件的地址，默认是/tmp/[server].pid
EXAMPLES	
	/path/to/server start --env=dev --debug=true --pid=/tmp/server.pid
	/path/to/server start --host=127.0.0.1 --port=8808
	/path/to/server start --host=127.0.0.1 --port=8808 --gf.gcfg.file=config.product.toml
	/path/to/server start user --host=127.0.0.1 --port=8808 
	/path/to/server start pay  --host=127.0.0.1 --port=8808
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
	//设置pid文件
	that.pidFile = parser.GetOpt("pid", gfile.TempDir(fmt.Sprintf("%s.pid", that.processName)))

	command := gcmd.GetArg(1)
	switch strings.ToLower(command) {
	case "help":
		that.help()
		return false
	case "version":
		that.version()
		return false
	case "start":
		that.checkStart()
		return true
	case "stop":
		that.stop(parser, "stop")
	case "reload":
		that.stop(parser, "reload")
	case "quit":
		that.stop(parser, "quit")
		return false
	default:
		for k := range gcmd.GetOptAll() {
			switch k {
			case "?", "h":
				fmt.Println(helpContent)
				return false
			case "i", "v":
				that.version()
				return false
			}
		}
	}
	that.help()
	return false
}

func (that *EasyService) stop(parser *gcmd.Parser, signal string) {
	pidFile := that.pidFile
	var serverPid = 0
	if gfile.IsFile(pidFile) {
		serverPid = gconv.Int(gstr.Trim(gfile.GetContents(pidFile)))
	}
	if serverPid == 0 {
		logger.Fatalf("Server is not running.")
	}

	var sigNo syscall.Signal
	switch signal {
	case "stop":
		sigNo = syscall.SIGTERM
	case "reload":
		sigNo = syscallSIGUSR2
	case "quit":
		sigNo = syscall.SIGQUIT
	default:
		logger.Fatalf("signal cmd `%s' not found", signal)
	}
	err := syscallKill(serverPid, sigNo)
	if err != nil {
		logger.Errorf("error:%v", err)
	}
	os.Exit(0)
}

func (that *EasyService) checkStart() {
	pidFile := that.pidFile
	var serverPid = 0
	if gfile.IsFile(pidFile) {
		serverPid = gconv.Int(gstr.Trim(gfile.GetContents(pidFile)))
	}
	if serverPid == 0 {
		return
	}
	if checkStart(serverPid) {
		logger.Fatalf("Server [%d] is already running.", serverPid)
	}
	return
}

//解析配置文件
func (that *EasyService) parserConfig(parser *gcmd.Parser) {
	that.config = gcfg.Instance()
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
	if len(env) > 0 {
		_ = genv.Set("ENV_NAME", gstr.ToLower(env))
		_ = that.config.Set("ENV_NAME", gstr.ToLower(env))
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

//显示帮助信息
func (that *EasyService) help() {
	fmt.Print(helpContent)
}

package dserver

import (
	"fmt"
	"github.com/desertbit/grumble"
	"github.com/osgochina/dmicro/logger"
	"os"
)

// 正常进程
func (that *DServer) initGrumble() {
	that.grumbleApp = grumble.New(&grumble.Config{
		Name: that.name,
	})
	that.grumbleApp.SetPrintASCIILogo(func(a *grumble.App) {
		that.Version()
	})
	that.grumbleApp.AddCommand(&grumble.Command{
		Name:    "version",
		Help:    "打印当前程序的版本信息",
		Aliases: []string{"v", "i", "--version", "-v"},
		Run: func(c *grumble.Context) error {
			that.Version()
			os.Exit(0)
			return nil
		},
	})

	that.grumbleApp.AddCommand(&grumble.Command{
		Name:    "help",
		Help:    "use 'help [command]' for command help",
		Aliases: []string{"?", "h", "--help"},
		Args: func(a *grumble.Args) {
			a.StringList("command", "the name of the command")
		},
		Run: func(c *grumble.Context) error {
			that.Version()
			fmt.Println()
			that.Help()
			os.Exit(0)
			return nil
		},
	})
	that.grumbleApp.AddCommand(&grumble.Command{
		Name:     "start",
		Help:     "启动服务",
		LongHelp: "启动指定的服务,传入要启动的服务名",
		Usage:    "start [OPTIONS] sandboxName [sandboxName...]",
		Args: func(a *grumble.Args) {
			a.StringList("sandboxNames", "sandbox names")
		},
		Flags: func(f *grumble.Flags) {
			f.String("c", "config", "", "指定要载入的配置文件，该参数与gf.gcfg.file参数二选一，建议使用该参数")
			f.Bool("d", "daemon", false, "使用守护进程模式启动")
			f.String("e", "env", "product", "环境变量，表示当前启动所在的环境,有[dev,test,product]这三种，默认是product")
			f.Bool("", "debug", true, "是否开启debug 默认debug=true")
			f.String("", "pid", "", " 设置pid文件的地址，默认是/tmp/[server].pid")
			f.Int("m", "model", 0, " 进程模型，0表示单进程模型，1表示多进程模型")
		},
		Run: func(c *grumble.Context) error {
			// 获取要启动的sandbox名称
			that.initSandboxNames(c.Args.StringList("sandboxNames"))
			// 初始化pid文件的路径
			that.initPidFile(c.Flags.String("pid"))
			// 判断服务进程是否已经启动
			that.checkStart()
			//解析配置文件
			that.parserConfig(c.Flags.String("config"))
			// 解析是否守护进程启动
			that.parserDaemon(c.Flags.Bool("daemon"))
			// 解析运行环境
			that.parserEnv(c.Flags.String("env"))
			// 解析debug参数
			that.parserDebug(c.Flags.Bool("debug"))
			//初始化日志配置
			if e := that.initLogSetting(that.config); e != nil {
				logger.Fatalf("error:%v", e)
			}
			// 启动
			that.run(c)
			return nil
		},
	})

	that.grumbleApp.AddCommand(&grumble.Command{
		Name: "stop",
		Help: "停止服务",
		Args: func(a *grumble.Args) {
			a.StringList("sandboxNames", "sandbox names")
		},
		Flags: func(f *grumble.Flags) {
			f.String("", "pid", "", " 设置pid文件的地址，默认是/tmp/[server].pid")
		},
		Run: func(c *grumble.Context) error {
			// 获取要启动的sandbox名称
			that.initSandboxNames(c.Args.StringList("sandboxNames"))
			// 初始化pid文件的路径
			that.initPidFile(c.Flags.String("pid"))
			that.stop("stop")
			return nil
		},
	})

	that.grumbleApp.AddCommand(&grumble.Command{
		Name: "reload",
		Help: "平滑重启服务",
		Args: func(a *grumble.Args) {
			a.StringList("sandboxNames", "sandbox names")
		},
		Flags: func(f *grumble.Flags) {
			f.String("", "pid", "", " 设置pid文件的地址，默认是/tmp/[server].pid")
		},
		Run: func(c *grumble.Context) error {
			// 获取要启动的sandbox名称
			that.initSandboxNames(c.Args.StringList("sandboxNames"))
			// 初始化pid文件的路径
			that.initPidFile(c.Flags.String("pid"))
			that.stop("reload")
			return nil
		},
	})

	that.grumbleApp.AddCommand(&grumble.Command{
		Name: "quit",
		Help: "优雅的停止服务",
		Args: func(a *grumble.Args) {
			a.StringList("sandboxNames", "sandbox names")
		},
		Flags: func(f *grumble.Flags) {
			f.String("", "pid", "", " 设置pid文件的地址，默认是/tmp/[server].pid")
		},
		Run: func(c *grumble.Context) error {
			// 获取要启动的sandbox名称
			that.initSandboxNames(c.Args.StringList("sandboxNames"))
			// 初始化pid文件的路径
			that.initPidFile(c.Flags.String("pid"))
			that.stop("quit")
			return nil
		},
	})
}

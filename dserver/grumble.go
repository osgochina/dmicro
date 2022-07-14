package dserver

import (
	"fmt"
	"github.com/desertbit/grumble"
	"github.com/fatih/color"
	"github.com/gogf/gf/os/genv"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/util/grand"
	"github.com/modood/table"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/proto/pbproto"
	"github.com/osgochina/dmicro/logger"
	"os"
	"time"
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
			f.Bool("", "debug", false, "是否开启debug 默认debug=false")
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

// 是否在打印日志
var ctrlLoging = false

// ctrl进程
func (that *DServer) initCtrlGrumble() {
	that.grumbleApp = grumble.New(&grumble.Config{
		Name:                  that.name,
		Description:           "好用的服务管理工具",
		HistoryFile:           fmt.Sprintf("%s/.%s.hist", genv.Get("HOME", "/tmp"), that.name),
		Prompt:                fmt.Sprintf("%s » ", that.name),
		PromptColor:           color.New(color.FgCyan, color.Bold),
		HelpHeadlineColor:     color.New(color.FgGreen),
		HelpHeadlineUnderline: true,
		HelpSubCommands:       true,
	})
	that.grumbleApp.SetInterruptHandler(func(a *grumble.App, count int) {
		if count >= 2 {
			_, _ = a.Println("exit success!!!")
			os.Exit(1)
		}
		if ctrlLoging {
			var result *Result
			sess, err := that.getCtrlSession()
			if err != nil {
				a.PrintError(err)
				return
			}
			stat := sess.Call("/ctrl/close_logger",
				nil,
				&result,
			).Status()
			if !stat.OK() {
				a.PrintError(err)
				return
			}
			ctrlLoging = false
			return
		}
		_, _ = a.Println("input Ctrl-c once more to exit")
	})
	that.grumbleApp.SetPrintASCIILogo(func(a *grumble.App) {
		that.Version()
	})
	that.grumbleApp.AddCommand(&grumble.Command{
		Name:    "version",
		Help:    "打印当前程序的版本信息",
		Aliases: []string{"v"},
		Run: func(c *grumble.Context) error {
			that.Version()
			return nil
		},
	})
	statusCommand := &grumble.Command{
		Name:    "info",
		Help:    "查看当前服务状态",
		Aliases: []string{"status", "ps"},
		Run: func(c *grumble.Context) error {
			var result *Infos
			sess, err := that.getCtrlSession()
			if err != nil {
				return err
			}
			stat := sess.Call("/ctrl/info",
				nil,
				&result,
			).Status()
			if !stat.OK() {
				return stat.Cause()
			}
			table.Output(result.List)
			return nil
		},
	}
	that.grumbleApp.AddCommand(statusCommand)

	startCommand := &grumble.Command{
		Name: "start",
		Help: "启动服务",
		Args: func(a *grumble.Args) {
			a.String("sandboxName", "sandbox names")
		},
		Run: func(c *grumble.Context) error {
			name := c.Args.String("sandboxName")
			var result *Result
			sess, err := that.getCtrlSession()
			if err != nil {
				return err
			}
			stat := sess.Call("/ctrl/start",
				name,
				&result,
			).Status()
			if !stat.OK() {
				return stat.Cause()
			}
			fmt.Println("启动服务成功")
			return nil
		},
	}
	that.grumbleApp.AddCommand(startCommand)

	stopCommand := &grumble.Command{
		Name: "stop",
		Help: "停止服务",
		Args: func(a *grumble.Args) {
			a.String("sandboxName", "sandbox names")
		},
		Run: func(c *grumble.Context) error {
			name := c.Args.String("sandboxName")
			var result *Result
			sess, err := that.getCtrlSession()
			if err != nil {
				return err
			}
			stat := sess.Call("/ctrl/stop",
				name,
				&result,
			).Status()
			if !stat.OK() {
				return stat.Cause()
			}
			fmt.Println("停止服务成功")
			return nil
		},
	}
	that.grumbleApp.AddCommand(stopCommand)

	reloadCommand := &grumble.Command{
		Name: "reload",
		Help: "平滑重启服务",
		Args: func(a *grumble.Args) {
			a.String("sandboxName", "sandbox name")
		},
		Run: func(c *grumble.Context) error {
			name := c.Args.String("sandboxName")
			var result *Result
			sess, err := that.getCtrlSession()
			if err != nil {
				return err
			}
			stat := sess.Call("/ctrl/reload",
				name,
				&result,
			).Status()
			if !stat.OK() {
				return stat.Cause()
			}
			fmt.Println("平滑重启成功")
			return nil
		},
	}
	that.grumbleApp.AddCommand(reloadCommand)

	debugCommand := &grumble.Command{
		Name: "debug",
		Help: "debug开关",
		Args: func(a *grumble.Args) {
			a.String("switch", "'open'，'true','1'=true; 'close','false','0'=false")
		},
		Run: func(c *grumble.Context) error {
			debug := true
			debugStr := c.Args.String("switch")
			if debugStr == "open" || debugStr == "true" || debugStr == "1" {
				debug = true
			} else if debugStr == "close" || debugStr == "false" || debugStr == "0" {
				debug = false
			}
			var result *Result
			sess, err := that.getCtrlSession()
			if err != nil {
				return err
			}
			stat := sess.Call("/ctrl/debug",
				debug,
				&result,
			).Status()
			if !stat.OK() {
				return stat.Cause()
			}
			if debug {
				fmt.Println("已开启debug模式")
			} else {
				fmt.Println("已关闭debug模式")
			}

			return nil
		},
	}
	that.grumbleApp.AddCommand(debugCommand)

	logCommand := &grumble.Command{
		Name: "log",
		Help: "打印出服务的运行日志",
		Args: func(a *grumble.Args) {
			//a.String("level", "info|debug|error")
		},
		Run: func(c *grumble.Context) error {
			level := "info"
			//level := c.Args.String("level")
			var result *Result
			sess, err := that.getCtrlSession()
			if err != nil {
				return err
			}
			stat := sess.Call("/ctrl/open_logger",
				level,
				&result,
			).Status()
			if !stat.OK() {
				return stat.Cause()
			}
			ctrlLoging = true
			return nil
		},
	}
	that.grumbleApp.AddCommand(logCommand)
}

func (that *DServer) getCtrlSession() (drpc.Session, error) {
	if that.ctrlSession != nil && that.ctrlSession.Health() {
		return that.ctrlSession, nil
	}
	localPath := gfile.TempDir(fmt.Sprintf("%s.cli.%s", that.name, grand.S(6)))
	cli := drpc.NewEndpoint(drpc.EndpointConfig{
		Network:        "unix",
		LocalIP:        localPath,
		PrintDetail:    false,
		RedialTimes:    1,
		RedialInterval: time.Second,
	},
	)
	cli.RoutePush(new(ctrlLoggerPush))
	svrPath := gfile.TempDir(fmt.Sprintf("%s.sock", that.name))
	var stat *drpc.Status
	that.ctrlSession, stat = cli.Dial(svrPath, pbproto.NewPbProtoFunc())
	if !stat.OK() {
		return nil, fmt.Errorf("链接到DServer服务[%s]失败", svrPath)
	}
	go func() {
		<-that.ctrlSession.CloseNotify()
		_, _ = that.grumbleApp.Println("DServer服务已断开链接")
		_ = cli.Close()
		that.ctrlSession = nil
		_ = gfile.Remove(localPath)
	}()

	return that.ctrlSession, nil
}

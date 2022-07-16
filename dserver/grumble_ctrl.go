package dserver

import (
	"fmt"
	"github.com/desertbit/grumble"
	"github.com/fatih/color"
	"github.com/gogf/gf/os/genv"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/grand"
	"github.com/modood/table"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/proto/pbproto"
	"os"
	"strings"
	"time"
)

var ctrlCChan = make(chan struct{})

// ctrl进程命令行
func (that *DServer) initCtrlGrumble() {
	that.grumbleApp = grumble.New(&grumble.Config{
		Name:                  that.name,
		Description:           "好用的服务管理工具",
		HistoryFile:           fmt.Sprintf("%s/.%s.hist", genv.Get("HOME", "/tmp"), that.name),
		Prompt:                fmt.Sprintf("%s » ", that.name),
		PromptColor:           color.New(color.FgGreen, color.Bold),
		HelpHeadlineColor:     color.New(color.FgGreen),
		HelpHeadlineUnderline: true,
		HelpSubCommands:       true,
	})
	that.grumbleApp.SetInterruptHandler(func(a *grumble.App, count int) {
		if count >= 2 {
			_, _ = a.Println("exit success!!!")
			os.Exit(1)
		}
		ctrlCChan <- struct{}{}
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
	that.grumbleApp.AddCommand(&grumble.Command{
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
	})

	that.grumbleApp.AddCommand(&grumble.Command{
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
	})

	that.grumbleApp.AddCommand(&grumble.Command{
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
	})

	that.grumbleApp.AddCommand(&grumble.Command{
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
	})

	that.grumbleApp.AddCommand(&grumble.Command{
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
	})

	that.grumbleApp.AddCommand(&grumble.Command{
		Name: "log",
		Help: "打印出服务的运行日志",
		Flags: func(f *grumble.Flags) {
			f.String("l", "level", "all", "日志级别")
		},
		Run: func(c *grumble.Context) error {
			level := glog.LEVEL_ALL
			levelStr := c.Flags.String("level")
			if l, ok := levelStringMap[strings.ToUpper(levelStr)]; ok {
				level = l
			}
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
			fmt.Println("开始打印服务端日志........")
			go func() {
				<-ctrlCChan

				sess, err = that.getCtrlSession()
				if err != nil {
					c.App.PrintError(err)
					return
				}
				stat = sess.Call("/ctrl/close_logger",
					nil,
					&result,
				).Status()
				if !stat.OK() {
					c.App.PrintError(err)
					return
				}
			}()

			return nil
		},
	})
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

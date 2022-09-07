package dserver

import (
	"fmt"
	"github.com/spf13/cobra"
)

func (that *DServer) initCobra() {
	that.cobraCmd = &cobra.Command{
		Use:   fmt.Sprintf("%s <command> <subcommand> [flags]", that.name),
		Short: fmt.Sprintf("%s CLI", that.name),
		Long:  fmt.Sprintf("%s is an application built on the DMicro framework", that.name),
	}
	that.cobraCmd.PersistentFlags().Bool("debug", true, "Whether to enable debug, default debug=true")
	that.cobraCmd.PersistentFlags().Bool("help", false, "Show help for command")
	that.cobraCmd.AddCommand(&cobra.Command{
		Use:     "version",
		Short:   "Print the version information of the current program",
		Aliases: []string{"v", "i", "--version", "-v"},
		Run: func(cmd *cobra.Command, args []string) {
			that.Version()
		},
	})
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "start the service.",
		Long:  "start the service.",
		Example: `
  /path/to/server start user admin --env=dev --debug=true --pid=/tmp/server.pid --config=config.product.toml
  /path/to/server start user admin --config=config.product.toml 
  /path/to/server start user admin --config=config.product.toml --mode=1
  /path/to/server start user admin --config=config.product.toml --mode=1 --daemon=true
`,
		Run: func(cmd *cobra.Command, args []string) {
			// 获取要启动的sandbox名称
			that.initSandboxNames(args)
			// 初始化pid文件的路径
			that.initPidFile(cmd.Flag("pid").Value.String())
			// 判断服务进程是否已经启动
			that.checkStart()
			//解析配置文件
			that.parserConfig(cmd.Flag("config").Value.String())
			// 解析是否守护进程启动
			that.parserDaemon(cmd.Flag("daemon").Value.String() == "true")
			// 解析运行环境
			that.parserEnv(cmd.Flag("env").Value.String())
			// 解析debug参数
			that.parserDebug(cmd.Flag("debug").Value.String() == "true")
			//初始化日志配置
			if e := that.initLogSetting(that.config); e != nil {
				cmd.PrintErrf("error:%v", e)
				return
			}
			// 启动
			that.run(cmd)
		},
	}
	startCmd.Flags().StringP("config", "c", "", "Specifies the path to the configuration file to load")
	startCmd.Flags().BoolP("daemon", "d", false, "Start in daemon mode")
	startCmd.Flags().StringP("env", "e", "product", "Environment variable, indicating the current startup environment, there are three kinds of [dev, test, product], the default is product.")
	startCmd.Flags().StringP("pid", "", "", "Set the address of the pid file, the default is /tmp/[server].pid")
	startCmd.Flags().IntP("model", "m", 0, "Process model, 0 means single process model, 1 means multi-process model")

	that.cobraCmd.AddCommand(startCmd)

	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "stop the service.",
		Long:  "stop the service.",
		Example: `
  /path/to/server stop 
  /path/to/server stop  --pid=/tmp/server.pid
  /path/to/server stop user admin --pid=/tmp/server.pid
`,
		Run: func(cmd *cobra.Command, args []string) {
			// 获取要启动的sandbox名称
			that.initSandboxNames(args)
			// 初始化pid文件的路径
			that.initPidFile(cmd.Flag("pid").Value.String())
			that.stop("stop")
		},
	}
	stopCmd.Flags().StringP("pid", "", "", "Set the address of the pid file, the default is /tmp/[server].pid")
	that.cobraCmd.AddCommand(stopCmd)

	reloadCmd := &cobra.Command{
		Use:   "reload",
		Short: "reload the service.",
		Long:  "reload the service.",
		Example: `
  /path/to/server reload 
  /path/to/server reload  --pid=/tmp/server.pid
  /path/to/server reload user admin --pid=/tmp/server.pid
`,
		Run: func(cmd *cobra.Command, args []string) {
			// 获取要启动的sandbox名称
			that.initSandboxNames(args)
			// 初始化pid文件的路径
			that.initPidFile(cmd.Flag("pid").Value.String())
			that.stop("reload")
		},
	}
	reloadCmd.Flags().StringP("pid", "", "", "Set the address of the pid file, the default is /tmp/[server].pid")
	that.cobraCmd.AddCommand(reloadCmd)

	quitCmd := &cobra.Command{
		Use:   "quit",
		Short: "gracefully stop service.",
		Long:  "gracefully stop service.",
		Example: `
  /path/to/server quit 
  /path/to/server quit  --pid=/tmp/server.pid
  /path/to/server quit user admin --pid=/tmp/server.pid
`,
		Run: func(cmd *cobra.Command, args []string) {
			// 获取要启动的sandbox名称
			that.initSandboxNames(args)
			// 初始化pid文件的路径
			that.initPidFile(cmd.Flag("pid").Value.String())
			that.stop("quit")
		},
	}
	quitCmd.Flags().StringP("pid", "", "", "Set the address of the pid file, the default is /tmp/[server].pid")
	that.cobraCmd.AddCommand(quitCmd)
}

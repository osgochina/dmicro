package dserver

import (
	"fmt"
	"github.com/gogf/gf/os/gfile"
	"github.com/spf13/cobra"
	"os"
)

var defaultServer = newDServer(fmt.Sprintf("DServer_%s", gfile.Basename(os.Args[0])))

// SetName 设置应用名
// 建议设置独特个性化的引用名，因为管理链接，日志目录等地方会用到它。
// 如果不设置，默认是"DServer_xxx",启动xxx为二进制名
func SetName(name string) {
	defaultServer.name = name
}

// Setup 启动服务
func Setup(startFunction ...StartFunc) {
	defaultServer.setup(startFunction...)
}

// CloseCtl 关闭ctl管理功能
func CloseCtl() {
	defaultServer.openCtl = false
}

// Shutdown 关闭服务
func Shutdown() {
	defaultServer.Shutdown()
}

// Cobra 注册命令
func Cobra(f func(rootCmd *cobra.Command)) {
	defaultServer.cobraRootCmdCallback = f
}

// GetSandbox 通过sandbox name 获取已注册的sandbox
func GetSandbox(sandboxName string) ISandbox {
	for _, v := range defaultServer.serviceList.Map() {
		dService := v.(*DService)
		s, found := dService.SearchSandBox(sandboxName)
		if found {
			return s.(ISandbox)
		}
	}
	return nil
}

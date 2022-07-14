package dserver

import (
	"fmt"
	"github.com/desertbit/grumble"
	"github.com/gogf/gf/os/gfile"
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
func Setup(startFunction StartFunc) {
	defaultServer.setup(startFunction)
}

// GrumbleApp 增加自定义命令
func GrumbleApp() *grumble.App {
	return defaultServer.grumbleApp
}

// CloseCtrl 关闭ctrl管理功能
func CloseCtrl() {
	defaultServer.openCtrl = false
}

// Shutdown 关闭服务
func Shutdown() {
	defaultServer.Shutdown()
}

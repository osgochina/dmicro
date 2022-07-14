package dserver

import "github.com/desertbit/grumble"

var defaultServer = newDServer()

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

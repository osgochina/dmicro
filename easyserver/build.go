package easyserver

import (
	"fmt"
	"github.com/gogf/gf/os/gfile"
)

var (
	BuildVersion     = "No Version Info"
	BuildGoVersion   = "No Version Info"
	BuildGitCommitId = "No Commit Info"
	BuildTime        = "No Time Info"
)

//显示版本信息
func (that *Server) version() {
	fmt.Printf("Server Version: %s\n", BuildVersion)
	fmt.Printf("Server Build Time: %s\n", BuildTime)
	fmt.Printf("Go version: %s\n", BuildGoVersion)
	fmt.Printf("Git commit: %s\n", BuildGitCommitId)
	fmt.Printf("Install Path: %s\n", gfile.SelfPath())
}

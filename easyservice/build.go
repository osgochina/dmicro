package easyservice

import (
	"fmt"
	"github.com/gogf/gf"
	"github.com/gogf/gf/os/gfile"
	"github.com/osgochina/dmicro"
)

var (
	BuildVersion     = "No Version Info"
	BuildGoVersion   = "No Version Info"
	BuildGitCommitId = "No Commit Info"
	BuildTime        = "No Time Info"
)

//显示版本信息
func (that *EasyService) version() {
	fmt.Printf("Server Version: %s\n", BuildVersion)
	fmt.Printf("Server Build Time: %s\n", BuildTime)
	fmt.Printf("Go version: %s\n", BuildGoVersion)
	fmt.Printf("Git commit: %s\n", BuildGitCommitId)
	fmt.Printf("DMicro Version: %s\n", dmicro.Version)
	fmt.Printf("GF Version: %s\n", gf.VERSION)
	fmt.Printf("Install Path: %s\n", gfile.SelfPath())
}

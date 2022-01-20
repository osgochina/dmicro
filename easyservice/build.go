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
	Authors          = "No Authors Info"
	Logo             = `
  ____    __  __   _                       
 |  _ \  |  \/  | (_)   ___   _ __    ___  
 | | | | | |\/| | | |  / __| | '__|  / _ \ 
 | |_| | | |  | | | | | (__  | |    | (_) |
 |____/  |_|  |_| |_|  \___| |_|     \___/
`
)

//显示版本信息
func (that *EasyService) version() {
	fmt.Print(Logo)
	fmt.Printf("Version:         %s\n", BuildVersion)
	fmt.Printf("Go Version:      %s\n", BuildGoVersion)
	fmt.Printf("DMicro Version:  %s\n", dmicro.Version)
	fmt.Printf("GF Version:      %s\n", gf.VERSION)
	fmt.Printf("Git Commit:      %s\n", BuildGitCommitId)
	fmt.Printf("Build Time:      %s\n", BuildTime)
	fmt.Printf("Authors:         %s\n", Authors)
	fmt.Printf("Install Path:    %s\n", gfile.SelfPath())
}

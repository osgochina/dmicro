package dserver

import (
	"fmt"
	"github.com/gogf/gf"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/text/gstr"
	"github.com/osgochina/dmicro"
)

var (
	BuildVersion     = "No Version Info"
	BuildGoVersion   = "No Version Info"
	BuildGitCommitId = "No Commit Info"
	BuildTime        = "No Time Info"
	Authors          = "No Authors Info"
	Logo             = `
  ____    ____                                      
 |  _ \  / ___|    ___   _ __  __   __   ___   _ __ 
 | | | | \___ \   / _ \ | '__| \ \ / /  / _ \ | '__|
 | |_| |  ___) | |  __/ | |     \ V /  |  __/ | |   
 |____/  |____/   \___| |_|      \_/    \___| |_|  
`
)

//Version 显示版本信息
func (that *DServer) Version() {
	fmt.Print(gstr.TrimLeftStr(Logo, "\n"))
	fmt.Printf("Version:         %s\n", BuildVersion)
	fmt.Printf("Go Version:      %s\n", BuildGoVersion)
	fmt.Printf("DMicro Version:  %s\n", dmicro.Version)
	fmt.Printf("GF Version:      %s\n", gf.VERSION)
	fmt.Printf("Git Commit:      %s\n", BuildGitCommitId)
	fmt.Printf("Build Time:      %s\n", BuildTime)
	fmt.Printf("Authors:         %s\n", Authors)
	fmt.Printf("Install Path:    %s\n", gfile.SelfPath())
}

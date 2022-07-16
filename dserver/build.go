package dserver

import (
	"fmt"
	"github.com/gogf/gf"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/text/gstr"
	"github.com/osgochina/dmicro"
	"os"
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

var (
	helpContent = gstr.TrimLeft(`
USAGE
  %path% start|stop|reload|quit|ctrl [OPTION] [sandboxName1|sandboxName2...] 
OPTION
  -c,--config     指定要载入的配置文件，该参数与gf.gcfg.file参数二选一，建议使用该参数
  -d,--daemon     使用守护进程模式启动
  --env           环境变量，表示当前启动所在的环境,有[dev,test,product]这三种，默认是product
  --debug         是否开启debug 默认debug=false
  --pid           设置pid文件的地址，默认是/tmp/[server].pid
  -h,--help       获取帮助信息
  -v,--version    获取编译版本信息
  -m,--model      进程模型，0表示单进程模型，1表示多进程模型	
EXAMPLES
  %path% 
  %path% start --env=dev --debug=true --pid=/tmp/server.pid
  %path% start --gf.gcfg.file=config.product.toml
  %path% start -c=config.product.toml
  %path% start --config=config.product.toml user admin
  %path% start user
  %path% stop
  %path% quit
  %path% reload
  %path% version
  %path% help
  %path% ctrl // 该命令可以链接到已启动的服务
`)
)

//Help 显示帮助信息
func (that *DServer) Help() {
	fmt.Print(gstr.Replace(helpContent, "%path%", gfile.Abs(os.Args[0])))
}

// +build !windows

package graceful

import (
	"os"
)

// AddInherited 添加需要给重启后新进程继承的文件句柄和环境变量
func (that *Graceful) AddInherited(procFiles []*os.File, envs map[string]string) {
	that.locker.Lock()
	defer that.locker.Unlock()
	for _, f := range procFiles {
		var had bool
		for _, ff := range that.inheritedProcFiles {
			if ff == f {
				had = true
				break
			}
		}
		if !had {
			that.inheritedProcFiles = append(that.inheritedProcFiles, f)
		}
	}
	for k, v := range envs {
		that.inheritedEnv[k] = v
	}
}

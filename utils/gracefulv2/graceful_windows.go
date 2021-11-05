// +build windows

package gracefulv2

import (
	"os"
)

func (that *ChangeProcessGraceful) AddInherited(procFiles []*os.File, envs map[string]string) {}

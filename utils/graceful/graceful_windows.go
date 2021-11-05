// +build windows

package graceful

import (
	"os"
)

func (that *ChangeProcessGraceful) AddInherited(procFiles []*os.File, envs map[string]string) {}

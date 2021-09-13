package procconf

import (
	"errors"
	"os/exec"
	"syscall"
)

// create command from string or []string
func createCommand(command string, args ...[]string) (*exec.Cmd, error) {
	if len(args) <= 0 {
		return nil, errors.New("empty command")
	}
	cmd := exec.Command(command)
	if len(args) > 1 {
		cmd.Args = args[0]
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{}

	return cmd, nil
}

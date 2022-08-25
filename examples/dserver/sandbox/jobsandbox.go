package sandbox

import (
	"fmt"
	"github.com/osgochina/dmicro/dserver"
)

type JobSandbox struct {
	dserver.JobSandbox
}

func (that *JobSandbox) Name() string {
	return "JobSandbox"
}

func (that *JobSandbox) Setup() error {
	fmt.Println("JobSandbox Setup")
	return nil
}

func (that *JobSandbox) Shutdown() error {
	fmt.Println("JobSandbox Shutdown")
	return nil
}

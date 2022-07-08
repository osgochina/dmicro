package sandbox

import (
	"fmt"
	"github.com/osgochina/dmicro/dserver"
)

type DefaultSandBox1 struct {
	dserver.BaseSandbox
}

func (that *DefaultSandBox1) Name() string {
	return "DefaultSandBox1"
}
func (that *DefaultSandBox1) Abc() string {
	return "DefaultSandBox1"
}
func (that *DefaultSandBox1) Setup() error {
	fmt.Println("DefaultSandBox1 Setup")
	return nil
}

func (that *DefaultSandBox1) Shutdown() error {
	fmt.Println("DefaultSandBox1 Shutdown")
	return nil
}

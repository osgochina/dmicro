package sandbox

import (
	"fmt"
	"github.com/gogf/gf/errors/gerror"
	"github.com/osgochina/dmicro/dserver"
)

type DefaultSandBox1 struct {
	dserver.BaseSandbox
	age int
}

func (that *DefaultSandBox1) Name() string {
	fmt.Println(that.age)
	return "DefaultSandBox1"
}
func (that *DefaultSandBox1) Abc() string {
	return "DefaultSandBox1"
}
func (that *DefaultSandBox1) Setup() error {
	fmt.Println("Setup")
	//fmt.Println(that.Service().Name())
	that.age = 100
	return gerror.New("Setup error")
}

func (that *DefaultSandBox1) Shutdown() error {
	fmt.Println("Shutdown")
	fmt.Println(that.age)
	return nil
}

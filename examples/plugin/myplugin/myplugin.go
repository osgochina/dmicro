package myplugin

import (
	"fmt"
	"github.com/osgochina/dmicro/drpc"
)

type myPlugin struct {
	newParams string
	myParams  string
}

// 强制检查`myPlugin`是否实现了对应的hook节点，为了防止hook点太多的时候遗漏
var (
	_ drpc.AfterNewEndpointPlugin    = new(myPlugin)
	_ drpc.BeforeCloseEndpointPlugin = new(myPlugin)
	_ drpc.AfterDialPlugin           = new(myPlugin)
	_ drpc.AfterAcceptPlugin         = new(myPlugin)
)

func NewMyPlugin(newParams string) *myPlugin {
	return &myPlugin{newParams: newParams}
}

func (that *myPlugin) Name() string {
	return "myPlugin"
}

// AfterNewEndpoint endpoint创建成功后调用
func (that *myPlugin) AfterNewEndpoint(drpc.EarlyEndpoint) error {
	that.myParams = "myParams"
	fmt.Println("AfterNewEndpoint")
	return nil
}

// BeforeCloseEndpoint endpoint关闭之前调用
func (that *myPlugin) BeforeCloseEndpoint(drpc.Endpoint) error {
	fmt.Println("BeforeCloseEndpoint")
	return nil
}

// AfterDial 客户端链接到服务端成功后调用(endpoint作为客户端时生学校)
func (that *myPlugin) AfterDial(sess drpc.EarlySession, isRedial bool) *drpc.Status {
	fmt.Println("AfterDial")
	fmt.Println(that.myParams)
	fmt.Println(sess.RemoteAddr())
	return nil
}

// AfterAccept 服务端接收客户端请求成功后调用(endpoint作为服务端时生效)
func (that *myPlugin) AfterAccept(sess drpc.EarlySession) *drpc.Status {
	fmt.Println("AfterAccept")
	fmt.Println(that.myParams)
	fmt.Println(sess.RemoteAddr())
	return nil
}

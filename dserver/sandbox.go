package dserver

import (
	"github.com/gogf/gf/errors/gerror"
	"reflect"
)

// ISandbox 服务沙盒的接口
type ISandbox interface {
	Name() string    // 沙盒名字
	Setup() error    // 启动沙盒
	Shutdown() error //关闭沙盒
}

// SandboxCtx 所有业务使用的Sandbox都需要继承SandboxCtx
type SandboxCtx interface {
	Service() *DService
	Server() *DServer
}

// 反射
type sandboxReflectValue struct {
	ctrl   reflect.Value
	ctxPtr *SandboxCtx
}

// 私有sandbox，
type privateSandbox struct {
	original   ISandbox
	cType      reflect.Type
	sandboxVal *sandboxReflectValue
	handlerCtx *handlerCtx
}

func (that *privateSandbox) Name() string {
	return that.original.Name()
}

func (that *privateSandbox) Setup() error {
	method, found := that.cType.MethodByName("Setup")
	methodFunc := method.Func
	if !found {
		return gerror.Newf("not found Setup func")
	}
	obj := that.sandboxVal
	*obj.ctxPtr = that.handlerCtx
	rets := methodFunc.Call([]reflect.Value{obj.ctrl})
	err := rets[0].Interface()
	if err != nil {
		return err.(error)
	}
	return nil
}

func (that *privateSandbox) Shutdown() error {
	method, found := that.cType.MethodByName("Shutdown")
	methodFunc := method.Func
	if !found {
		return gerror.Newf("not found Shutdown func")
	}
	obj := that.sandboxVal
	*obj.ctxPtr = that.handlerCtx
	rets := methodFunc.Call([]reflect.Value{obj.ctrl})
	err := rets[0].Interface()
	if err != nil {
		return err.(error)
	}
	return nil
}

// 实现了SandboxCtx，在调用privateSandbox的时候，使用handlerCtx替代
type handlerCtx struct {
	service *DService
	server  *DServer
}

func (that *handlerCtx) Service() *DService {
	return that.service
}

func (that *handlerCtx) Server() *DServer {
	return that.server
}

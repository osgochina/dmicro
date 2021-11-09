# Handler 处理器

在路由注册成功后，`服务名`与`处理方法`的对应关系绑定。

那么是直接绑定开发者注册的方法吗？

答案是否定的，开发者注册的方法并不能直接被使用，它们将会帮转换成`Handler`对象保存在内存中。

正常开发的时候，开发人员不太需要使用到`Handler`，在进行插件开发的时候才会接触到它。

`Handler` 对象的定义如下：
```go

type Handler struct {
	// 绑定的handler名
	name string
	//参数类型
	argElem reflect.Type
	//返回值类型 注意：只有call消息才会有
	reply reflect.Type
	//处理该消息的方法
	handleFunc func(*handlerCtx, reflect.Value)
	// 不能匹配到绑定方法时，默认的处理方法
	unknownHandleFunc func(*handlerCtx)
	// 路由类型名字
	routerTypeName string
	//是否找到绑定方法
	isUnknown bool
	// 插件容器
	pluginContainer *PluginContainer
}
```
如何把用户定义的方法转换成`Handler`对象呢？

定义方法如下：
```go
func Home(ctx drpc.CallCtx, args *string) (string, *drpc.Status) {
    return *args, nil
}
```
1. 获取参数和返回值类型
```go
    cType  := reflect.TypeOf(Home)
    cValue := reflect.ValueOf(Home)
    //参数类型，第二个参数
    argType := cType.In(1)
    argElem := argType.Elem()
    // 返回值类型，第一个返回值
    replyType := cType.Out(0)
    reply := replyType
```
2. 生成handleFunc
```go
// 获取第一个参数
var ctxTypeElem = cType.In(0).Elem()
// 判断第一个参数的名字
iType := ctxTypeElem.FieldByName("CallCtx")
// 获取参数的反射值
ctrl := reflect.New(ctxTypeElem)
// 拿到该CallCtx的对象
obj = &CallCtrlValue{
    ctrl: ctrl,
    ctxPtr: (*CallCtx)(unsafe.Pointer(uintptr(unsafe.Pointer(ctrl.Pointer())) + iType.Offset)),
}
// 构建handle执行方法，通过以下代码可以看到，是把自定义的方法通过反射调用，实现了调用绑定。
handleFunc = func(ctx *handlerCtx, argValue reflect.Value) {
			*obj.ctxPtr = ctx
			rets := cValue.Call([]reflect.Value{obj.ctrl, argValue})
			stat := (*status.Status)(unsafe.Pointer(rets[1].Pointer()))
			if !stat.OK() {
				ctx.stat = stat
				ctx.output.SetStatus(stat)
			} else {
				ctx.output.SetBody(rets[0].Interface())
			}
		}
```






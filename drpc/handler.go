package drpc

import "reflect"

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

	pluginContainer *PluginContainer

	// 路由类型名字
	routerTypeName string

	//是否找到绑定方法
	isUnknown bool
}

// RouterTypeName 获取处理器的路由方法名 pnPush/pnCall/pnUnknownPush/pnUnknownCall
func (that *Handler) RouterTypeName() string {
	return that.routerTypeName
}

// ReplyType 回复消息的类型
func (that *Handler) ReplyType() reflect.Type {
	return that.reply
}

// Name 名字
func (that *Handler) Name() string {
	return that.name
}

// NewArgValue 参数的反射值
func (that *Handler) NewArgValue() reflect.Value {
	return reflect.New(that.argElem)
}

// ArgElemType 参数的类型
func (that *Handler) ArgElemType() reflect.Type {
	return that.argElem
}

// IsCall 处理程序是否是Call
func (that *Handler) IsCall() bool {
	return that.routerTypeName == pnCall || that.routerTypeName == pnUnknownCall
}

// IsPush 处理程序是否是PUSH
func (that *Handler) IsPush() bool {
	return that.routerTypeName == pnPush || that.routerTypeName == pnUnknownPush
}

// IsUnknown 处理程序是否未找到
func (that *Handler) IsUnknown() bool {
	return that.isUnknown
}

# 路由 Router

路由`Router`是drpc框架非常核心的一个功能，它起到的作用顾名思义，绑定了`请求名`与`处理方法`之间的关系。

## ServiceMethod映射规则

映射规则目前默认支持两组，分别是http格式及rpc格式，当然你也可以自定义实现自己的规则。

### HTTPServiceMethodMapper

drpc默认使用该规则。

结构体或方法名称到服务方法名称的映射：

* `AaBb` -> `/aa_bb`
* `ABcXYz` -> `/abc_xyz`
* `Aa__Bb` -> `/aa_bb`
* `aa__bb` -> `/aa_bb`
* `ABC__XYZ` -> `/abc_xyz`
* `Aa_Bb` -> `/aa/bb`
* `aa_bb` -> `/aa/bb`
* `ABC_XYZ` -> `/abc/xyz`

### RPCServiceMethodMapper

结构体或方法名称到服务方法名称的映射：

* `AaBb` -> `AaBb`
* `ABcXYz` -> `ABcXYz`
* `Aa__Bb` -> `Aa_Bb`
* `aa__bb` -> `aa_bb`
* `ABC__XYZ` -> `ABC_XYZ`
* `Aa_Bb` -> `Aa.Bb`
* `aa_bb` -> `aa.bb`
* `ABC_XYZ` -> `ABC.XYZ`

### 其他

只需要实现`func(prefix, name string) (serviceMethod string)` 该结构的方法。

执行`droc.SetServiceMethodMapper(func(prefix, name string) (serviceMethod string))`就能设置你的自定义方法。

## 接口模板

### `CALL`类型接口的注册

以下方法的作用是把传入的数据累加，并返回累加够的值。
如果需要把它注册到路由，应该如何操作，接下来将演示各种用法。
> ps: 注意所有将要注册到路由的struct，func，都必须继承`drpc.CallCtx`或`drpc.PushCtx`.
```go

type Math struct {
	drpc.CallCtx
}

func (m *Math) Add(arg *[]int) (int, *Status) {
	var r int
	for _, a := range *arg {
		r += a
	}
	return r, nil
}

func (m *Math) Echo (arg *int) (int, *Status) {
    return *arg, nil
}
```

#### 注册`Struct`到根路由

```go
endpoint.RouteCall(new(Math))
```
以上代码将会注册两条路由。

* /math/add
* /math/echo

如果使用的是`RPCServiceMethodMapper`,则得到的路由是。

* math.add
* math.echo

#### 注册`Func`到根路由

```go
endpoint.RouteCallFunc((*Math).Add)
```
以上代码将会注册一条路由。

* /add

如果使用的是`RPCServiceMethodMapper`,则得到的路由是。

* add



#### 注册`Struct`到分组路由

```go
group := endpoint.SubRoute("github")
group.RouteCall(new(Math))
```
以上代码将会注册两条路由。

* /github/math/add
* /github/math/echo

如果使用的是`RPCServiceMethodMapper`,则得到的路由是。

* github.math.add
* github.math.echo


#### 注册`Func`到分组路由

```go
group := endpoint.SubRoute("github")
group.RouteCallFunc((*Math).Add)
```
以上代码将会注册一条路由。

* /github/add

如果使用的是`RPCServiceMethodMapper`,则得到的路由是。

* github.add

### `PUSH`类型接口的注册

```go

type MathPush struct {
    drpc.PushCtx
}

func (m *MathPush) Receive(arg *[]int) *Status {
	var r int
	for _, a := range *arg {
		r += a
	}
    fmt.Println(r)
	return nil
}

func (m *MathPush) Notify(arg *int) *Status {
    fmt.Println(*arg)
	return nil
}
```

#### 注册`Struct`到根路由

```go
endpoint.RoutePush(new(MathPush))
```
以上代码将会注册两条路由。

* /math_push/receive
* /math_push/notify

如果使用的是`RPCServiceMethodMapper`,则得到的路由是。

* MathPush.Receive
* MathPush.Notify

#### 注册`Func`到根路由

```go
endpoint.RoutePush((*MathPush).Receive)
```
以上代码将会注册一条路由。

* /receive

如果使用的是`RPCServiceMethodMapper`,则得到的路由是。

* Receive



#### 注册`Struct`到分组路由

```go
group := endpoint.SubRoute("github")
group.RoutePush(new(MathPush))
```
以上代码将会注册两条路由。

* /github/math_push/receive
* /github/math_push/notify

如果使用的是`RPCServiceMethodMapper`,则得到的路由是。

* github.MathPush.Receive
* github.MathPush.Notify


#### 注册`Func`到分组路由

```go
group := endpoint.SubRoute("github")
group.RoutePushFunc((*MathPush).Receive)
```
以上代码将会注册一条路由。

* /github/Receive

如果使用的是`RPCServiceMethodMapper`,则得到的路由是。

* github.Receive

## 未匹配如何处理

### 设置`CALL`命令的默认路由

如果请求发送到服务端，但是路由规则为匹配，需要执行什么逻辑？

可以实现 `func(UnknownCallCtx) (interface{}, *Status)` 该方法，并注册到路由中。

```go
endpoint.SetUnknownCall(func(UnknownCallCtx) (interface{}, *Status)){
    return nil,drpc.NewStatus(drpc.CodeNotFound, "Not Found", "")
})
```


### 设置`PUSH`命令的默认路由

实现 `func(SetUnknownPush) *Status` 该方法，并注册到路由中。

```go
endpoint.SetUnknownPush(func(UnknownPushCtx) *Status){
    return drpc.NewStatus(drpc.CodeNotFound, "Not Found", "")
})
```


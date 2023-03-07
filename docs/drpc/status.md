# Status 状态

在`call`、`push`、`reply`的各项操作中，需要明确操作是否成功，出错的原因，需要用到`status`对象。

## 组成
    
`status`对象的组成比较简单。

```go
type Status struct {
	code  int32
	msg   string
	cause error
	*stack
}
```

stack的作用是为了记录`error`错误的上下文。

创建简单的`status`对象只需要使用
```go
func New(code int32, msg string, cause ...interface{}) *Status
```

如通过需要携带堆栈信息，则使用：
```go
func NewWithStack(code int32, msg string, cause ...interface{}) *Status
```
## 内置类型

drpc框架已经内置了一些消息Code，当然，你也可以自定义消息Code。
```go
CodeUnknownError        int32 = -1     // 未知的错误
CodeOK                  int32 = 0      // 处理成功
CodeNoError             int32 = CodeOK // 处理成功
CodeInvalidOp           int32 = 1      // 无效的操作
CodeWrongConn           int32 = 100    // 错误的链接
CodeConnClosed          int32 = 102    // 链接已关闭
CodeWriteFailed         int32 = 104    // 写入失败
CodeDialFailed          int32 = 105    // 链接失败
CodeBadMessage          int32 = 400    // 消息格式错误
CodeUnauthorized        int32 = 401    // 认证不通过
CodeNotFound            int32 = 404    // 未找到对应的处理方法
CodeMTypeNotAllowed     int32 = 405    // 消息类型不正确
CodeHandleTimeout       int32 = 408    // 处理超时
CodeInternalServerError int32 = 500    // 内部服务器错误
CodeBadGateway          int32 = 502    // 网关错误
```

## 应用场景

`status`的应用场景分为三大类.
1. 框架内部使用
```go
//drpc/context.go文件中的409行
//传入的消息如果没有服务方法
if len(header.ServiceMethod()) == 0 {
    that.stat = statBadMessage.Copy("invalid service method for message")
    return nil
}
```
2. 处理`call`请求时的`reply`返回值
```go
   func (m *Math) Add(arg *[]int) (int, *drpc.Status) {
	// add
	var r int
	for _, a := range *arg {
		r += a
	}
	// 返回status异常的状态
	if r > 100 {
		return 0, drpc.NewStatus(998, "不允许和大于100")
	}
	// nil 表示正常，在call请求处判断status.OK()为true
	return r, nil
}
```
3. 调用`call`请求后的请求状态判断。

```go
var result int
stat = sess.Call("/math/add",
    []int{1, 2, 3, 4, 5},
    &result,
).Status()
// 只需要判断stat是否为true，就能知道该次请求是否成功
if !stat.OK() {
    logger.Fatalf(context.TODO(),"%v", stat)
}
```



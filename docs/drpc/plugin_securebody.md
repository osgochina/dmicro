### 消息体加密传输

使用该插件，能实现消息内容安全传输，使用的是`aes`算法加密消息内容。

### 使用示例

```go
package securebody_test

import (
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/plugin/securebody"
	"strconv"
	"testing"
	"time"
)

type Request struct {
	One int
	Two int
}

type Response struct {
	Three int
}

type math struct{ drpc.CallCtx }

func (m *math) Add(arg *Request) (*Response, *drpc.Status) {
	return &Response{Three: arg.One + arg.Two}, nil
}

func newSession(t *gtest.T, port uint16) drpc.Session {
	p := securebody.NewSecureBodyPlugin("cipherkey1234567")
	srv := drpc.NewEndpoint(drpc.EndpointConfig{
		ListenPort:  port,
		PrintDetail: true,
	})
	srv.RouteCall(new(math), p)
	go srv.ListenAndServe()
	time.Sleep(time.Second)

	cli := drpc.NewEndpoint(drpc.EndpointConfig{
		PrintDetail: true,
	}, p)
	sess, stat := cli.Dial(":" + strconv.Itoa(int(port)))
	if !stat.OK() {
		t.Fatal(stat)
	}
	return sess
}

func TestSecureBodyPlugin(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		sess := newSession(t, 9090)
		var result Response
		stat := sess.Call(
			"/math/add",
			&Request{One: 1, Two: 2},
			&result,
			securebody.WithSecureMeta(),
		).Status()
		t.Assert(stat.OK(), true)
		t.Assert(result.Three, 3)
		t.Logf("测试加密：1+2=%d", result.Three)
	})
}

func TestReplySecureBodyPlugin(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		sess := newSession(t, 9090)
		var result Response
		stat := sess.Call(
			"/math/add",
			&Request{One: 1, Two: 2},
			&result,
			securebody.WithReplySecureMeta(true),
		).Status()
		t.Assert(stat.OK(), true)
		t.Assert(result.Three, 3)
		t.Logf("测试加密：1+2=%d", result.Three)
	})
}

```

### 支持的方法

#### 创建`SecureBodyPlugin`插件

`NewSecureBodyPlugin(cipherKey string, statCode ...int32) drpc.Plugin {}`

参数:

* `cipherKey`  自定义加密key
* `statCode` 自定义错误码，该插件中的任何错误，都会返回该错误码


#### 强制要求加密传输

`WithSecureMeta() message.MsgSetting`

#### 强制要求服务端返回的内容加密传输

`WithReplySecureMeta(secure bool) message.MsgSetting `

#### 强制加密消息

`EnforceSecure(output message.Message) `

该方法一般用在服务端响应函数中，如：

```go

func (m *math) Add(arg *Request) (*Response, *drpc.Status) {
	//响应消息给客户端的时候，可以使用它强制加密，当然，前提是你已经加载了SecureBodyPlugin插件
    securebody.EnforceSecure(m.Output())
	return &Response{Three: arg.One + arg.Two}, nil
}
```
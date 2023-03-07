### 概述

http协议是实现HTTP风格的套接字通信协议。
它只是以HTTP方式传输数据，并没有实现完整的HTTP协议,所在在使用的时候需要注意。

### 消息内容

示例:
- 请求消息

```
POST /home/test?endpoint_id=110 HTTP/1.1
Accept-Encoding: gzip
Content-Length: 24
Content-Type: application/json;charset=utf-8
Host: localhost:9090
User-Agent: drpc-httproto/1.1
X-Mtype: 1
X-Seq: 1

{"author":"osgochina@gmail.com"}
```

- 响应消息

```
HTTP/1.1 200 OK
Content-Length: 32
Content-Type: application/json;charset=utf-8
X-Mtype: 2
X-Seq: 1

{"arg":{"author":"osgochina@gmail.com"}}
```

or

```
HTTP/1.1 299 Business Error
Content-Length: 56
Content-Type: application/json
X-Mtype: 2
X-Seq: 0

{"code":1,"msg":"test error","cause":"this is test:110"}
```

- 默认支持的 Content-Type
    - codec.ID_JSON:     application/json;charset=utf-8
    - codec.ID_FORM:     application/x-www-form-urlencoded;charset=utf-8
    - codec.ID_PLAIN:    text/plain;charset=utf-8
    - codec.ID_XML:      text/xml;charset=utf-8


-  如果要注册body的编码器，则使用`RegBodyCodec`方法

```go
func RegBodyCodec(contentType string, codecID byte)
```

### 如何使用

`import "github.com/osgochina/dmicro/drpc/proto/httproto"`

#### Test

```go
package httpproto_test

import (
  "bytes"
  "encoding/json"
  "github.com/gogf/gf/v2/test/gtest"
  "github.com/gogf/gf/v2/util/gconv"
  "github.com/osgochina/dmicro/drpc"
  "github.com/osgochina/dmicro/drpc/proto/httpproto"
  "github.com/osgochina/dmicro/logger"
  "io/ioutil"
  "net/http"
  "testing"
  "time"
)

type Home struct {
  drpc.CallCtx
}

func (h *Home) Test(arg *map[string]string) (map[string]interface{}, *drpc.Status) {
  logger.Infof(context.TODO(),"endpoint_id: %s", gconv.String(h.PeekMeta("endpoint_id")))
  return map[string]interface{}{
    "arg": *arg,
  }, nil
}

func (h *Home) TestError(arg *map[string]string) (map[string]interface{}, *drpc.Status) {
  return nil, drpc.NewStatus(1, "test error", "this is test:"+gconv.String(h.PeekMeta("endpoint_id")))
}

func TestHTTProto(t *testing.T) {
  gtest.C(t, func(t *gtest.T) {
    svr := drpc.NewEndpoint(drpc.EndpointConfig{ListenPort: 9090})
    svr.RouteCall(new(Home))
    go svr.ListenAndServe(httpproto.NewHTTProtoFunc(true))

    time.Sleep(1e9)

    cli := drpc.NewEndpoint(drpc.EndpointConfig{})
    sess, stat := cli.Dial(":9090", httpproto.NewHTTProtoFunc())

    if !stat.OK() {
      t.Fatal(stat)
    }

    var result interface{}
    var arg = map[string]string{
      "author": "liuzhiming",
    }

    testUrl := "http://localhost:9090/home/test?endpoint_id=110"
    stat = sess.Call(testUrl, arg, &result).Status()
    if !stat.OK() {
      t.Fatal(stat)
    }
    t.Logf("drpc client response: %v", result)
    b, err := json.Marshal(arg)
    if err != nil {
      return
    }
    resp, err := http.Post(testUrl, "application/json;charset=utf-8", bytes.NewReader(b))
    if err != nil {
      t.Fatal(err)
    }
    b, err = ioutil.ReadAll(resp.Body)
    _ = resp.Body.Close()
    t.Logf("http client response: %s", b)

    {
      testErrURL := "http://localhost:9090/home/test_error?endpoint_id=110"
      result = nil
      stat = sess.Call(
        testErrURL,
        arg,
        &result,
      ).Status()
      if stat.OK() {
        t.Fatal("test_error expect error")
      }
      t.Logf("erpc client response: %v, %v", stat, result)
      b, err = json.Marshal(arg)
      if err != nil {
        return
      }
      resp, err = http.Post(testUrl, "application/json;charset=utf-8", bytes.NewReader(b))
      if err != nil {
        t.Fatal(err)
      }
      b, _ = ioutil.ReadAll(resp.Body)
      resp.Body.Close()
      t.Logf("http client response: %s", b)
    }
  })
}

```

执行以下命令测试:

```sh
$ cd dmicro/drpc/proto/httpproto
$ go test -v -run=TestHTTProto
```

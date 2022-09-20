package httpproto_test

import (
	"bytes"
	"encoding/json"
	"github.com/gogf/gf/test/gtest"
	"github.com/gogf/gf/util/gconv"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/internal"
	"github.com/osgochina/dmicro/drpc/proto/httpproto"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

type Home struct {
	drpc.CallCtx
}

func (h *Home) Test(arg *map[string]string) (map[string]interface{}, *drpc.Status) {
	internal.Infof("endpoint_id: %s", gconv.String(h.PeekMeta("endpoint_id")))
	return map[string]interface{}{
		"arg": *arg,
	}, nil
}

func (h *Home) TestError(arg *map[string]string) (map[string]interface{}, *drpc.Status) {
	return nil, drpc.NewStatus(1, "test error", "this is test:"+gconv.String(h.PeekMeta("endpoint_id")))
}

func TestHTTProto(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		svr := drpc.NewEndpoint(drpc.EndpointConfig{ListenPort: 9091})
		svr.RouteCall(new(Home))
		go svr.ListenAndServe(httpproto.NewHTTProtoFunc(true))

		time.Sleep(1e9)

		cli := drpc.NewEndpoint(drpc.EndpointConfig{})
		sess, stat := cli.Dial(":9091", httpproto.NewHTTProtoFunc())

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

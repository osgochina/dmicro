// Package jsonSubProto 实现JSON套接字通信协议的。
package jsonSubProto

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/osgochina/dmicro/drpc/proto"
	"io/ioutil"
	"sync"
)

// NewJSONSubProtoFunc JSON socket协议的创建函数。
func NewJSONSubProtoFunc() proto.ProtoFunc {
	return func(rw proto.IOWithReadBuffer) proto.Proto {
		return &jsonSubProto{
			id:   'j',
			name: "json",
			rw:   rw,
		}
	}
}

type jsonSubProto struct {
	id   byte
	name string
	rw   proto.IOWithReadBuffer
	rMu  sync.Mutex
}

// Version 返回协议的id和名字
func (that *jsonSubProto) Version() (byte, string) {
	return that.id, that.name
}

//协议的编码格式
const format = `{"seq":%d,"mtype":%d,"serviceMethod":%q,"meta":%q,"bodyCodec":%d,"body":"%s","ptf":%s}`

// Pack 打包消息，并且把消息写入
func (that *jsonSubProto) Pack(m proto.Message) error {
	// marshal body
	bodyBytes, err := m.MarshalBody()
	if err != nil {
		return err
	}
	// do transfer pipe
	bodyBytes, err = m.PipeTFilter().OnPack(bodyBytes)
	if err != nil {
		return err
	}
	// marshal transfer pipe ids
	var pipeTFilterIDs = make([]int, m.PipeTFilter().Len())
	for i, id := range m.PipeTFilter().IDs() {
		pipeTFilterIDs[i] = int(id)
	}
	pipeTFilterIDsBytes, err := json.Marshal(pipeTFilterIDs)
	if err != nil {
		return err
	}

	// join json format
	s := fmt.Sprintf(format,
		m.Seq(),
		m.MType(),
		m.ServiceMethod(),
		m.Meta().String(),
		m.BodyCodec(),
		bytes.Replace(bodyBytes, []byte{'"'}, []byte{'\\', '"'}, -1),
		pipeTFilterIDsBytes,
	)

	b := gconv.Bytes(s)

	_ = m.SetSize(uint32(len(b)))

	_, err = that.rw.Write(b)
	return err
}

// Unpack 从链接中读取消息，并且解包
// NOTE: Concurrent unsafe!
func (that *jsonSubProto) Unpack(m proto.Message) error {
	that.rMu.Lock()
	defer that.rMu.Unlock()
	b, err := ioutil.ReadAll(that.rw)
	if err != nil {
		return err
	}

	_ = m.SetSize(uint32(len(b)))

	j := gjson.New(b)

	// read transfer pipe
	pipeTFilter := j.Get("ptf").Array()
	for _, r := range pipeTFilter {
		_ = m.PipeTFilter().Append(byte(gconv.Int(r)))
	}

	// read body
	m.SetBodyCodec(byte(j.Get("bodyCodec").Int()))
	bodyBytes, err := m.PipeTFilter().OnUnpack(j.Get("body").Bytes())
	if err != nil {
		return err
	}

	// read other
	m.SetSeq(j.Get("seq").Int32())
	m.SetMType(byte(j.Get("mtype").Int8()))
	m.SetServiceMethod(j.Get("serviceMethod").String())
	meta := j.Get("meta").Map()
	if len(meta) > 0 {
		for k, v := range meta {
			m.Meta().Set(k, v)
		}
	}
	// unmarshal new body
	err = m.UnmarshalBody(bodyBytes)
	return err
}

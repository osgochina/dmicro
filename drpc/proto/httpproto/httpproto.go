package httpproto

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/codec"
	"github.com/osgochina/dmicro/drpc/message"
	"github.com/osgochina/dmicro/drpc/proto"
	"github.com/osgochina/dmicro/drpc/tfilter"
	"github.com/osgochina/dmicro/drpc/tfilter/gzip"
	"github.com/osgochina/dmicro/utils/dbuffer"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

var (
	bodyCodecMapping = map[string]byte{
		//"application/x-protobuf":            codec.ID_PROTOBUF,
		"application/json": codec.IdJson,
		//"application/x-www-form-urlencoded": codec.ID_FORM,
		//"text/plain":                        codec.ID_PLAIN,
		//"text/xml":                          codec.ID_XML,
	}
	contentTypeMapping = map[byte]string{
		//codec.ID_PROTOBUF: "application/x-protobuf;charset=utf-8",
		codec.IdJson: "application/json;charset=utf-8",
		//codec.ID_FORM:     "application/x-www-form-urlencoded;charset=utf-8",
		//codec.ID_PLAIN:    "text/plain;charset=utf-8",
		//codec.ID_XML:      "text/xml;charset=utf-8",
	}
)

// RegBodyCodec 注册新的编解码器
func RegBodyCodec(contentType string, codecID byte) {
	bodyCodecMapping[contentType] = codecID
	contentTypeMapping[codecID] = contentType
}

// GetBodyCodec 根据 ContentType 获取对应的编解码器
func GetBodyCodec(contentType string, defCodecID byte) byte {
	idx := strings.Index(contentType, ";")
	if idx != -1 {
		contentType = contentType[:idx]
	}
	codecID, ok := bodyCodecMapping[contentType]
	if !ok {
		return defCodecID
	}
	return codecID
}

// GetContentType 通过编解码器获取对应的ContentType
func GetContentType(codecID byte, defContentType string) string {
	contentType, ok := contentTypeMapping[codecID]
	if !ok {
		return defContentType
	}
	return contentType
}

// NewHTTProtoFunc 创建http协议支持
//  Only support t filter : gzip
//  Must use HTTP service method mapper
func NewHTTProtoFunc(printMessage ...bool) proto.ProtoFunc {
	drpc.SetServiceMethodMapper(drpc.HTTPServiceMethodMapper)
	var printable bool
	if len(printMessage) > 0 {
		printable = printMessage[0]
	}

	return func(rw proto.IOWithReadBuffer) proto.Proto {
		return &httpProto{
			id:           'h',
			name:         "http",
			rw:           rw,
			printMessage: printable,
		}
	}
}

type httpProto struct {
	rw           proto.IOWithReadBuffer
	rMu          sync.Mutex
	name         string
	id           byte
	printMessage bool
}

var (
	methodBytes  = []byte("POST")
	versionBytes = []byte("HTTP/1.1")
	crlfBytes    = []byte("\r\n")
)

var (
	okBytes     = []byte("200 OK")
	bizErrBytes = []byte("299 Business Error")
)

// Version 协议版本
func (that *httpProto) Version() (byte, string) {
	return that.id, that.name
}

// Pack 对数据进行打包
func (that *httpProto) Pack(msg proto.Message) error {
	bodyBytes, err := msg.MarshalBody()
	if err != nil {
		return err
	}
	var header = make(http.Header, msg.Meta().Size())
	msg.PipeTFilter().Iterator(func(idx int, filter tfilter.TransferFilter) bool {
		//是否支持gzip
		if !gzip.Is(filter.ID()) {
			err = fmt.Errorf("unsupport tfilter : %s", filter.Name())
			return false
		}
		//如果支持giz，则使用giz压缩body
		bodyBytes, err = filter.OnPack(bodyBytes)
		if err != nil {
			return false
		}
		header.Set("Content-Encoding", "gzip")
		header.Set("X-Content-Encoding", filter.Name())
		return true
	})
	if err != nil {
		return err
	}
	//把请求序列号在http头中加入
	header.Set("X-Seq", strconv.FormatInt(int64(msg.Seq()), 10))
	//把消息类型加入
	header.Set("X-MType", strconv.Itoa(int(msg.MType())))
	//把元数据放入http头中
	msg.Meta().Iterator(func(k interface{}, v interface{}) bool {
		header.Add(gconv.String(k), gconv.String(v))
		return true
	})

	bb := dbuffer.GetByteBuffer()
	defer dbuffer.ReleaseByteBuffer(bb)

	switch msg.MType() {
	case message.TypeCall, message.TypeAuthCall:
		err = that.packRequest(msg, header, bb, bodyBytes)
	case message.TypeReply, message.TypeAuthReply:
		err = that.packResponse(msg, header, bb, bodyBytes)
	default:
		return fmt.Errorf("unsupport message type: %d(%s)", msg.MType(), message.TypeText(msg.MType()))
	}
	if err != nil {
		return err
	}
	_ = msg.SetSize(uint32(bb.Len()))

	if that.printMessage {
		glog.Printf("Send HTTP Message:\n%s", gconv.String(bb.B))
	}
	_, err = that.rw.Write(bb.B)
	return err
}

//打包请求头
func (that *httpProto) packRequest(msg message.Message, header http.Header, bb *dbuffer.ByteBuffer, bodyBytes []byte) error {
	u, err := url.Parse(msg.ServiceMethod())
	if err != nil {
		return err
	}
	if u.Host != "" {
		header.Set("Host", u.Host)
	}
	header.Set("User-Agent", "drpc-httpproto/1.1")
	_, _ = bb.Write(methodBytes)
	_ = bb.WriteByte(' ')
	if u.RawQuery == "" {
		_, _ = bb.WriteString(u.Path)
	} else {
		_, _ = bb.WriteString(u.Path + "?" + u.RawQuery)
	}
	_ = bb.WriteByte(' ')
	_, _ = bb.Write(versionBytes)
	_, _ = bb.Write(crlfBytes)
	header.Set("Content-Type", GetContentType(msg.BodyCodec(), "text/plain;charset=utf-8"))
	header.Set("Content-Length", strconv.Itoa(len(bodyBytes)))
	header.Set("Accept-Encoding", "gzip")
	_ = header.Write(bb)
	_, _ = bb.Write(crlfBytes)
	_, _ = bb.Write(bodyBytes)
	return nil
}

// 打包响应
func (that *httpProto) packResponse(msg message.Message, header http.Header, bb *dbuffer.ByteBuffer, bodyBytes []byte) error {
	_, _ = bb.Write(versionBytes)
	_ = bb.WriteByte(' ')
	if stat := msg.Status(); !stat.OK() {
		statBytes, _ := stat.MarshalJSON()
		_, _ = bb.Write(bizErrBytes)
		_, _ = bb.Write(crlfBytes)
		if gzipName := header.Get("X-Content-Encoding"); gzipName != "" {
			gz, _ := tfilter.GetByName(gzipName)
			statBytes, _ = gz.OnPack(statBytes)
		}
		header.Set("Content-Type", "application/json")
		header.Set("Content-Length", strconv.Itoa(len(statBytes)))
		_ = header.Write(bb)
		_, _ = bb.Write(crlfBytes)
		_, _ = bb.Write(statBytes)
		return nil
	}
	_, _ = bb.Write(okBytes)
	_, _ = bb.Write(crlfBytes)
	header.Set("Content-Type", GetContentType(msg.BodyCodec(), "text/plain"))
	header.Set("Content-Length", strconv.Itoa(len(bodyBytes)))
	_ = header.Write(bb)
	_, _ = bb.Write(crlfBytes)
	_, _ = bb.Write(bodyBytes)
	return nil
}

var respPrefix = []byte("HTTP/")

func (that *httpProto) Unpack(m proto.Message) error {
	that.rMu.Lock()
	defer that.rMu.Unlock()

	bb := dbuffer.GetByteBuffer()
	defer dbuffer.ReleaseByteBuffer(bb)

	var size = 5
	bb.ChangeLen(size)
	_, err := io.ReadFull(that.rw, bb.B)
	if err != nil {
		return err
	}
	prefixBytes := make([]byte, 5, 128)
	copy(prefixBytes, bb.B)
	err = that.readLine(bb)
	if err != nil {
		return err
	}
	size += bb.Len()
	firstLine := append(prefixBytes, bb.B...)

	var msg []byte

	if bytes.Equal(prefixBytes, respPrefix) {
		m.SetMType(message.TypeReply)
		// status line
		var ok bool
		a := bytes.SplitN(firstLine, spaceBytes, 2)
		if len(a) != 2 {
			return errBadHTTPMsg
		}
		if bytes.Equal(a[1], okBytes) {
			ok = true
		} else if !bytes.Equal(a[1], bizErrBytes) {
			return errUnSupportHTTPCode
		}
		size, msg, err = that.unpack(m, bb)
		if err != nil {
			return err
		}
		if that.printMessage {
			glog.Printf("Recv HTTP Message:\n%s\r\n%s",
				gconv.String(firstLine), gconv.String(msg))
		}
		size += len(firstLine)
		_ = m.SetSize(uint32(size))
		if ok {
			return m.UnmarshalBody(bb.B)
		}
		_ = m.UnmarshalBody(nil)
		return m.Status(true).UnmarshalJSON(bb.B)
	}
	// request
	m.SetMType(message.TypeCall)
	a := bytes.SplitN(firstLine, spaceBytes, 3)
	if len(a) != 3 {
		return errBadHTTPMsg
	}
	u, err := url.Parse(gconv.String(a[1]))
	if err != nil {
		return err
	}
	m.SetServiceMethod(u.Path)
	if u.RawQuery != "" {
		for _, val := range gstr.SplitAndTrim(u.RawQuery, "&") {
			v := gstr.SplitAndTrim(val, "=")
			if len(v) == 2 {
				m.Meta().Set(v[0], v[1])
			}
		}
	}
	size, msg, err = that.unpack(m, bb)
	if err != nil {
		return err
	}
	if that.printMessage {
		glog.Printf("Recv HTTP Message:\n%s\r\n%s",
			gconv.String(firstLine), gconv.String(msg))
	}
	size += len(firstLine)
	_ = m.SetSize(uint32(size))
	return m.UnmarshalBody(bb.B)
}

func (that *httpProto) unpack(m message.Message, bb *dbuffer.ByteBuffer) (size int, msg []byte, err error) {
	var bodySize int
	var a [][]byte
	for i := 0; true; i++ {
		err = that.readLine(bb)
		if err != nil {
			return 0, nil, err
		}
		if that.printMessage {
			msg = append(msg, bb.B...)
			msg = append(msg, '\r', '\n')
		}
		size += bb.Len()
		// blank line, to read body
		if bb.Len() == 0 {
			break
		}
		// header
		a = bytes.SplitN(bb.B, colonBytes, 2)
		if len(a) != 2 {
			return 0, nil, errBadHTTPMsg
		}
		a[1] = bytes.TrimSpace(a[1])
		if bytes.Equal(contentTypeBytes, a[0]) {
			m.SetBodyCodec(GetBodyCodec(gconv.String(a[1]), codec.NilCodecID))
			continue
		}
		if bytes.Equal(contentLengthBytes, a[0]) {
			bodySize, err = strconv.Atoi(gconv.String(a[1]))
			if err != nil {
				return 0, nil, errBadHTTPMsg
			}
			size += bodySize
			continue
		}
		if bytes.Equal(xContentEncodingBytes, a[0]) {
			zg, err := tfilter.GetByName(gconv.String(a[1]))
			if err != nil {
				return 0, nil, err
			}
			_ = m.PipeTFilter().Append(zg.ID())
			continue
		}
		if bytes.Equal(xSeqBytes, a[0]) {
			var seq int
			seq, err = strconv.Atoi(gconv.String(a[1]))
			if err != nil {
				return 0, nil, errBadHTTPMsg
			}
			m.SetSeq(int32(seq))
			continue
		}
		if bytes.Equal(xMTypeBytes, a[0]) {
			var mtype int
			mtype, err = strconv.Atoi(gconv.String(a[1]))
			if err != nil {
				return 0, nil, errBadHTTPMsg
			}
			m.SetMType(byte(mtype))
			continue
		}
		m.Meta().Set(gconv.String(a[0]), gconv.String(a[1]))
	}
	if bodySize == 0 {
		return size, msg, nil
	}
	bb.ChangeLen(bodySize)
	_, err = io.ReadFull(that.rw, bb.B)
	if err != nil {
		return 0, nil, err
	}
	if that.printMessage {
		msg = append(msg, bb.B...)
		msg = append(msg, '\r', '\n')
	}
	bb.B, err = m.PipeTFilter().OnUnpack(bb.B)
	return size, msg, err
}

func (that *httpProto) readLine(bb *dbuffer.ByteBuffer) error {
	bb.Reset()
	oneByte := make([]byte, 1)
	var err error
	for {
		_, err = io.ReadFull(that.rw, oneByte)
		if err != nil {
			return err
		}
		if oneByte[0] == '\n' {
			n := bb.Len()
			if n > 0 && bb.B[n-1] == '\r' {
				bb.B = bb.B[:n-1]
			}
			return nil
		}
		_, _ = bb.Write(oneByte)
	}
}

var (
	spaceBytes            = []byte(" ")
	colonBytes            = []byte(":")
	contentTypeBytes      = []byte("Content-Type")
	contentLengthBytes    = []byte("Content-Length")
	xContentEncodingBytes = []byte("X-Content-Encoding")
	xSeqBytes             = []byte("X-Seq")
	xMTypeBytes           = []byte("X-MType")
	errBadHTTPMsg         = errors.New("bad HTTP message")
	errUnSupportHTTPCode  = errors.New("unSupport HTTP status code")
)

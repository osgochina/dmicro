package rawproto

import (
	"encoding/binary"
	"errors"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/osgochina/dmicro/drpc/proto"
	"github.com/osgochina/dmicro/utils/dbuffer"
	"io"
	"math"
	"strconv"
	"sync"
)

/**
rpc协议的原始格式 使用网络字节序，大端
{4 bytes 表示整个消息的长度}
{1 byte  表示协议的版本}
{1 byte  表示传输管道过滤器id的长度}
{ 传输管道过滤器id序列化后的内容}
# 以下的内容都是经过传输管道过滤器处理过的数据
{1 bytes 表示消息序列号长度}
{序列号 32进制的int32数字表示}
{1 byte 表示消息类型} # CALL:1;REPLY:2;PUSH:3
{1 byte 表示请求服务的方法长度 service method length}
{service method}
{2 bytes status length}
{status(json)}
{2 bytes metadata length}
{metadata(json)}
{1 byte bode codec id}
{body}
**/

func NewRawProtoFunc() proto.ProtoFunc {
	return RawProtoFunc
}

var RawProtoFunc = func(rw proto.IOWithReadBuffer) proto.Proto {
	return &rawProto{
		id:   6,
		name: "raw",
		r:    rw,
		w:    rw,
	}
}

var _ proto.Proto = new(rawProto)

var version = 1

type rawProto struct {
	r    io.Reader
	w    io.Writer
	rMu  sync.Mutex
	name string
	id   byte
}

func (that *rawProto) Version() (byte, string) {
	return that.id, that.name
}

// Pack 打包
func (that *rawProto) Pack(m proto.Message) error {
	bb := dbuffer.GetByteBuffer()
	defer dbuffer.ReleaseByteBuffer(bb)

	//占一个uint32类型长度的位置
	err := binary.Write(bb, binary.BigEndian, uint32(0))
	if err != nil {
		return err
	}
	//写入协议版本
	_ = bb.WriteByte(byte(version))
	//写入管道id长度
	_ = bb.WriteByte(byte(m.PipeTFilter().Len()))
	//写入管道id
	_, _ = bb.Write(m.PipeTFilter().IDs())

	prefixLen := bb.Len()

	// 写入消息头
	err = that.writeHeader(bb, m)
	if err != nil {
		return err
	}

	// 写入消息体
	err = that.writeBody(bb, m)
	if err != nil {
		return err
	}

	//通过管道处理打包，并把原始数据替换成处理后的数据
	payload, err := m.PipeTFilter().OnPack(bb.B[prefixLen:])
	if err != nil {
		return err
	}
	bb.B = append(bb.B[:prefixLen], payload...)

	// 设置消息长度
	err = m.SetSize(uint32(bb.Len()))
	if err != nil {
		return err
	}
	// 重新设置真是的长度
	binary.BigEndian.PutUint32(bb.B, m.Size())

	//写入数据到链接
	_, err = that.w.Write(bb.B)
	if err != nil {
		return err
	}

	return nil
}

//写入消息头
func (that *rawProto) writeHeader(bb *dbuffer.ByteBuffer, m proto.Message) error {

	seqStr := strconv.FormatInt(int64(m.Seq()), 36)
	//写入序列号长度
	_ = bb.WriteByte(byte(len(seqStr)))
	//写入序列号
	_, _ = bb.Write(gconv.Bytes(seqStr))
	//写入消息类型
	_ = bb.WriteByte(m.MType())

	serviceMethod := gconv.Bytes(m.ServiceMethod())
	serviceMethodLength := len(serviceMethod)
	if serviceMethodLength > math.MaxUint8 {
		return errors.New("raw proto: not support service method longer than 255")
	}
	//写入服务名长度
	_ = bb.WriteByte(byte(serviceMethodLength))
	//写入服务名
	_, _ = bb.Write(serviceMethod)

	statusBytes := gconv.Bytes(m.Status(true).String())
	//写入状态字符串的长度
	_ = binary.Write(bb, binary.BigEndian, uint16(len(statusBytes)))
	//写入状态字符串
	_, _ = bb.Write(statusBytes)

	metaBytes := gconv.Bytes(m.Meta().String())
	//写入元数据长度
	_ = binary.Write(bb, binary.BigEndian, uint16(len(metaBytes)))
	//写入元数据
	_, _ = bb.Write(metaBytes)
	return nil
}

//写入消息体
func (that *rawProto) writeBody(bb *dbuffer.ByteBuffer, m proto.Message) error {
	_ = bb.WriteByte(m.BodyCodec())
	bodyBytes, err := m.MarshalBody()
	if err != nil {
		return err
	}
	_, _ = bb.Write(bodyBytes)
	return nil
}

// Unpack 解包
func (that *rawProto) Unpack(m proto.Message) error {
	bb := dbuffer.GetByteBuffer()
	defer dbuffer.ReleaseByteBuffer(bb)
	// 读取消息
	err := that.readMessage(bb, m)
	if err != nil {
		return err
	}
	// 使用管道过滤器处理内容
	data, err := m.PipeTFilter().OnUnpack(bb.B)
	if err != nil {
		return err
	}

	// 读取header
	data, err = that.readHeader(data, m)
	if err != nil {
		return err
	}
	// 读取消息内容
	return that.readBody(data, m)
}

//读取消息内容
func (that *rawProto) readMessage(bb *dbuffer.ByteBuffer, m proto.Message) error {
	that.rMu.Lock()
	defer that.rMu.Unlock()

	bb.ChangeLen(4)
	_, err := io.ReadFull(that.r, bb.B)
	if err != nil {
		return err
	}
	lastSize := binary.BigEndian.Uint32(bb.B)
	if err = m.SetSize(lastSize); err != nil {
		return err
	}
	lastSize -= 4
	bb.ChangeLen(int(lastSize))

	// version和pipe len
	_, err = io.ReadFull(that.r, bb.B[:2])
	//var version = bb.B[0]
	var pipeLen = bb.B[1]
	if pipeLen > 0 {
		_, err = io.ReadFull(that.r, bb.B[:pipeLen])
		if err != nil {
			return err
		}
		err = m.PipeTFilter().Append(bb.B[:pipeLen]...)
		if err != nil {
			return err
		}
	}
	lastSize -= 2 + uint32(pipeLen)
	bb.ChangeLen(int(lastSize))

	_, err = io.ReadFull(that.r, bb.B)
	return err
}

//读取包头
func (that *rawProto) readHeader(data []byte, m proto.Message) ([]byte, error) {
	// 读取序列号
	seqLen := data[0]
	data = data[1:]
	seq, err := strconv.ParseInt(gconv.String(data[:seqLen]), 36, 32)
	if err != nil {
		return nil, err
	}
	m.SetSeq(int32(seq))
	data = data[seqLen:]

	// 设置消息类型
	m.SetMType(data[0])
	data = data[1:]

	// 服务名
	serviceMethodLen := data[0]
	data = data[1:]
	m.SetServiceMethod(string(data[:serviceMethodLen]))
	data = data[serviceMethodLen:]

	// 状态
	statusLen := binary.BigEndian.Uint16(data)
	data = data[2:]
	err = m.Status(true).UnmarshalJSON(data[:statusLen])
	if err != nil {
		return nil, err
	}
	data = data[statusLen:]

	// meta
	metaLen := binary.BigEndian.Uint16(data)
	data = data[2:]
	err = m.Meta().UnmarshalJSON(data[:metaLen])
	if err != nil {
		return nil, err
	}
	data = data[metaLen:]

	return data, nil
}

//读取消息内容
func (that *rawProto) readBody(data []byte, m proto.Message) error {
	m.SetBodyCodec(data[0])
	return m.UnmarshalBody(data[1:])
}

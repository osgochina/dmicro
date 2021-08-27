package redis_codec

import (
	"github.com/osgochina/dmicro/utils/dbuffer"
	"strconv"
)

const CRLF = "\r\n"

type Msg interface {
	Bytes() []byte
}

// SuccessMsg 状态消息
type SuccessMsg struct {
	Status string
}

func (that *SuccessMsg) Bytes() []byte {
	return []byte("+" + that.Status + CRLF)
}

// MakeSuccessMsg 创建正确状态消息
func MakeSuccessMsg(status string) *SuccessMsg {
	return &SuccessMsg{
		Status: status,
	}
}

// ErrorMsg 错误的消息
type ErrorMsg struct {
	Status string
}

func (that *ErrorMsg) Bytes() []byte {
	return []byte("-" + that.Status + CRLF)
}

// MakeErrorMsg 创建错误的消息
func MakeErrorMsg(status string) *ErrorMsg {
	return &ErrorMsg{
		Status: status,
	}
}

// NumberMsg 整数消息
type NumberMsg struct {
	Number int64
}

func (that *NumberMsg) Bytes() []byte {
	return []byte(":" + strconv.FormatInt(that.Number, 10) + CRLF)
}

// MakeNumberMsg 创建整数消息
func MakeNumberMsg(number int64) *NumberMsg {
	return &NumberMsg{
		Number: number,
	}
}

type BulkMsg struct {
	Arg []byte
}

// MakeBulkMsg creates  BulkMsg
func MakeBulkMsg(arg []byte) *BulkMsg {
	return &BulkMsg{
		Arg: arg,
	}
}

var nullBulkReplyBytes = []byte("$-1")

// Bytes 串化
func (that *BulkMsg) Bytes() []byte {
	if len(that.Arg) == 0 {
		return nullBulkReplyBytes
	}
	return []byte("$" + strconv.Itoa(len(that.Arg)) + CRLF + string(that.Arg) + CRLF)
}

// MultiBulkMsg 多行消息
type MultiBulkMsg struct {
	Args [][]byte
}

func MakeMultiBulkMsg(args [][]byte) *MultiBulkMsg {
	return &MultiBulkMsg{
		Args: args,
	}
}

// Bytes 把多行消息转换成bytes
func (that *MultiBulkMsg) Bytes() []byte {
	bb := dbuffer.GetByteBuffer()
	defer dbuffer.ReleaseByteBuffer(bb)
	//写入消息块
	_, _ = bb.WriteString("*" + strconv.Itoa(len(that.Args)) + CRLF)
	for _, arg := range that.Args {
		if arg == nil {
			_, _ = bb.WriteString("$-1" + CRLF)
		} else {
			_, _ = bb.WriteString("$" + strconv.Itoa(len(arg)) + CRLF + string(arg) + CRLF)
		}
	}
	return bb.Bytes()
}

var emptyMultiBulkBytes = []byte("*0\r\n")

// EmptyMultiBulkMsg 空的消息
type EmptyMultiBulkMsg struct{}

// Bytes 序列化
func (r *EmptyMultiBulkMsg) Bytes() []byte {
	return emptyMultiBulkBytes
}

var nullBulkBytes = []byte("$-1\r\n")

// NullBulkMsg 空字符串
type NullBulkMsg struct{}

// Bytes 序列化
func (r *NullBulkMsg) Bytes() []byte {
	return nullBulkBytes
}

// MakeNullBulkMsg creates a new NullBulkMsg
func MakeNullBulkMsg() *NullBulkMsg {
	return &NullBulkMsg{}
}

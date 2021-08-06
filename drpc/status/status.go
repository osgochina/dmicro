package status

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gogf/gf/util/gconv"
	"io"
	"strconv"
)

const (
	// OK 成功状态
	OK int32 = 0

	// UnknownError 未知错误
	UnknownError int32 = -1
)

type Status struct {
	code  int32
	msg   string
	cause error
	*stack
}

// New 传入错误码，错误消息，错误原因，创建一个状态结构对象
func New(code int32, msg string, cause ...interface{}) *Status {
	s := &Status{
		code: code,
		msg:  msg,
	}
	if len(cause) > 0 {
		s.cause = toError(cause[0])
	}
	return s
}

// NewWithStack 传入错误码，错误消息，错误栈堆，创建一个状态结构对象
func NewWithStack(code int32, msg string, cause ...interface{}) *Status {
	return New(code, msg, cause...).TagStack(1)
}

// SetCode 设置错误码
func (that *Status) SetCode(code int32) *Status {
	if that != nil {
		that.code = code
	}
	return that
}

// SetMsg 设置错误消息
func (that *Status) SetMsg(msg string) *Status {
	if that != nil {
		that.msg = msg
	}
	return that
}

// SetCases 设置错误原因
func (that *Status) SetCases(cause interface{}) *Status {
	if that != nil {
		that.cause = toError(cause)
	}
	return that
}

// Code 获取错误码
func (that *Status) Code() int32 {
	if that == nil {
		return OK
	}
	return that.code
}

// Msg 获取错误消息
func (that *Status) Msg() string {
	if that == nil {
		return ""
	}
	if that.msg == "" && that.cause != nil {
		return that.cause.Error()
	}
	return that.msg
}

// Cause 获取错误原因
func (that *Status) Cause() error {
	if that == nil {
		return nil
	}
	if that.cause == nil && that.code != OK {
		return errors.New(that.msg)
	}
	return that.cause
}

// Clear 清除状态结构
func (that *Status) Clear() {
	*that = Status{}
}

// OK 判断状态是否成功
func (that *Status) OK() bool {
	return that.Code() == OK
}

// UnknownError 是否为未知错误
func (that *Status) UnknownError() bool {
	return that.Code() == UnknownError
}

// TagStack 栈堆结构
func (that *Status) TagStack(skip ...int) *Status {
	depth := 3
	if len(skip) > 0 {
		depth += skip[0]
	}
	that.stack = callers(depth)
	return that
}

// StackTrace 获取错误的栈堆
func (that *Status) StackTrace() StackTrace {
	if that == nil || that.stack == nil {
		return nil
	}
	return that.stack.StackTrace()
}

//把状态对象变成字符串输出
func (that *Status) String() string {
	if that == nil {
		return "<nil>"
	}
	b, _ := that.MarshalJSON()
	return string(b)
}

// Format 格式化输出
func (that *Status) Format(state fmt.State, verb rune) {
	switch verb {
	case 'v':
		if state.Flag('+') {
			_, _ = fmt.Fprintf(state, "%+v", that.String())
			if that.stack != nil {
				that.stack.Format(state, verb)
			}
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(state, that.String())
	case 'q':
		_, _ = fmt.Fprintf(state, "%q", that.String())
	}
}

type exportStatus struct {
	Code  int32  `json:"code"`
	Msg   string `json:"msg"`
	Cause string `json:"cause"`
}

var (
	reA  = []byte(`{"code":`)
	reB  = []byte(`,"msg":`)
	reC  = []byte(`,"cause":`)
	null = []byte("null")
)

//接口断言  status必须实现 json.Marshaler,json.Unmarshaler 这两个方法
var (
	_ json.Marshaler   = new(Status)
	_ json.Unmarshaler = new(Status)
)

// MarshalJSON json编码
func (that *Status) MarshalJSON() ([]byte, error) {
	if that == nil {
		return null, nil
	}
	b := append(reA, strconv.FormatInt(int64(that.code), 10)...)

	b = append(b, reB...)
	msg := strconv.Quote(gconv.String(that.msg))
	b = append(b, []byte(msg)...)

	var cause string
	if that.cause != nil {
		cause = that.cause.Error()
	}
	b = append(b, reC...)
	c := strconv.Quote(gconv.String(cause))
	b = append(b, []byte(c)...)
	b = append(b, '}')
	return b, nil

}

// UnmarshalJSON json解码
func (that *Status) UnmarshalJSON(b []byte) error {
	if that == nil {
		return nil
	}
	if len(b) == 0 {
		that.Clear()
		return nil
	}
	var v exportStatus
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	that.code = v.Code
	that.msg = v.Msg
	if v.Cause != "" {
		that.cause = errors.New(v.Cause)
	} else {
		that.cause = nil
	}
	return nil
}

// Copy copy 一个新的状态对象
func (that *Status) Copy(newCause interface{}, newStackSkip ...int) *Status {
	if that == nil {
		return nil
	}
	if newCause == nil {
		newCause = that.cause
	}
	cp := New(that.code, that.msg, newCause)
	if len(newStackSkip) != 0 {
		cp.stack = callers(3 + newStackSkip[0])
	}
	return cp
}

func toError(cause interface{}) error {
	switch v := cause.(type) {
	case nil:
		return nil
	case error:
		return v
	case string:
		return errors.New(v)
	case *Status:
		return v.cause
	case Status:
		return v.cause
	default:
		return fmt.Errorf("%v", v)
	}
}

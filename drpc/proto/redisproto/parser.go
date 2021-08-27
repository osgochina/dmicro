package redisproto

import (
	"bufio"
	"github.com/gogf/gf/errors/gerror"
	"github.com/osgochina/dmicro/drpc/codec/redis_codec"
	"io"
	"strconv"
	"strings"
)

type readState struct {
	readingMultiLine  bool     // 是否是多次读取
	expectedArgsCount int      //预期消息块的数量
	msgType           byte     //消息类型
	args              [][]byte // 消息参数
	bulkLen           int64    //
}

//判断当前消息是否读取完毕
func (that *readState) finished() bool {
	return that.expectedArgsCount > 0 && len(that.args) == that.expectedArgsCount
}

//从链接中解析消息
func parseMsgByIO(bufReader *bufio.Reader) (redis_codec.Msg, error) {
	var (
		state readState
		err   error
		msg   []byte
	)
	for {
		msg, err = readLine(bufReader, &state)
		if err != nil {
			return nil, err
		}
		//如果不是多行消息,则表示这是一个新的请求
		if !state.readingMultiLine {
			//如果该消息有多个消息块
			if msg[0] == '*' {
				err = parseMultiBulkHeader(msg, &state)
				if err != nil {
					return nil, gerror.Newf("protocol error: %s,%v", msg, err)
				}
				if state.expectedArgsCount == 0 {
					return &redis_codec.EmptyMultiBulkMsg{}, nil
				}
				continue
			}
			//如果该消息是一个块数据
			if msg[0] == '$' {
				err = parseBulkHeader(msg, &state)
				if err != nil {
					if err != nil {
						return nil, gerror.Newf("protocol error: %s,%v", msg, err)
					}
				}
				if state.bulkLen == -1 { // null bulk reply
					return &redis_codec.NullBulkMsg{}, nil
				}
				continue
			}
			//单行消息
			return parseSingleLine(msg)
		}
		err = readBody(msg, &state)
		if err != nil {
			if err != nil {
				return nil, gerror.Newf("protocol error: %s,%v", msg, err)
			}
		}
		if state.finished() {
			var result redis_codec.Msg
			if state.msgType == '*' {
				result = redis_codec.MakeMultiBulkMsg(state.args)
			} else if state.msgType == '$' {
				result = redis_codec.MakeBulkMsg(state.args[0])
			}
			return result, nil
		}
	}
}

// 从链接中读取一行数据
func readLine(bufReader *bufio.Reader, state *readState) (line []byte, err error) {
	var (
		lineSize = 0
	)
	//如果是正常的单行消息，则只读一行就结束
	if state.bulkLen == 0 {
		line, err = bufReader.ReadBytes('\n')
		if err != nil {
			return nil, err
		}
		lineSize = len(line)
		if lineSize == 0 || line[lineSize-2] != '\r' {
			return nil, gerror.Newf("protocol error: %s", line)
		}
	} else {
		//如果是二进制安全的消息，则使用消息长度+2(\r\n)，把该消息一次性读出
		line = make([]byte, state.bulkLen+2)
		_, err = io.ReadFull(bufReader, line)
		if err != nil {
			return nil, err
		}
		lineSize = len(line)
		//判断消息格式是否正确
		if lineSize == 0 || line[lineSize-2] != '\r' || line[lineSize-1] != '\n' {
			return nil, gerror.Newf("protocol error: %s", line)
		}
		//该行消息已经全部读出来了，把标示位值置0
		state.bulkLen = 0
	}

	return line, nil
}

//解析当行消息
func parseSingleLine(line []byte) (msg redis_codec.Msg, err error) {
	//去除结尾的"\r\n"
	str := strings.TrimSuffix(string(line), "\r\n")
	switch line[0] {
	case '+': //正确的状态
		msg = redis_codec.MakeSuccessMsg(str[1:])
	case '-': //错误的状态
		msg = redis_codec.MakeErrorMsg(str[1:])
	case ':': // 整数
		val, err := strconv.ParseInt(str[1:], 10, 64)
		if err != nil {
			return nil, gerror.Newf("protocol error: %s", line)
		}
		msg = redis_codec.MakeNumberMsg(val)
	default: //空格分割的参数
		strs := strings.Split(str, " ")
		args := make([][]byte, len(strs))
		for i, s := range strs {
			args[i] = []byte(s)
		}
		msg = redis_codec.MakeMultiBulkMsg(args)
	}
	return msg, nil
}

// 解析多消息块的消息头
func parseMultiBulkHeader(msg []byte, state *readState) error {
	var err error
	var expectedLine uint64

	msgSize := len(msg)
	// 预期有多少个消息块
	expectedLine, err = strconv.ParseUint(string(msg[1:msgSize-2]), 10, 32)
	if err != nil {
		return gerror.Newf("protocol error: %s", msg)
	}
	if expectedLine == 0 {
		state.expectedArgsCount = 0
		return nil
	} else if expectedLine > 0 {
		// 第一个字符表示消息类型
		state.msgType = msg[0]
		state.readingMultiLine = true
		state.expectedArgsCount = int(expectedLine)
		state.args = make([][]byte, 0, expectedLine)
		return nil
	}
	return gerror.Newf("protocol error: %s", msg)
}

func parseBulkHeader(msg []byte, state *readState) (err error) {
	msgSize := len(msg)
	// 获取块数据的长度
	state.bulkLen, err = strconv.ParseInt(string(msg[1:msgSize-2]), 10, 64)
	if err != nil {
		return gerror.Newf("protocol error: %s", msg)
	}
	if state.bulkLen == -1 { // 空的数据块
		return nil
	} else if state.bulkLen > 0 {

		state.msgType = msg[0] // 消息类型
		state.readingMultiLine = true
		state.expectedArgsCount = 1
		state.args = make([][]byte, 0, 1)
		return nil
	}
	return gerror.Newf("protocol error: %s", msg)
}

// 读取批量消息的非第一行
func readBody(msg []byte, state *readState) (err error) {
	line := msg[0 : len(msg)-2]

	if line[0] == '$' {
		state.bulkLen, err = strconv.ParseInt(string(line[1:]), 10, 64)
		if err != nil {
			return gerror.Newf("protocol error: %s", msg)
		}
		if state.bulkLen <= 0 { // null bulk in multi bulks
			state.args = append(state.args, []byte{})
			state.bulkLen = 0
		}
	} else {
		state.args = append(state.args, line)
	}
	return nil
}

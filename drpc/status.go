package drpc

import "github.com/osgochina/dmicro/drpc/status"

type Status = status.Status

var (
	NewStatus = status.New
)

//框架保留 Status code
//建议自定义代码大于1000
// 未知错误：-1
// 发送错误: 100-199
// 消息处理错误: 400-499
// 接收错误: 500-599
const (
	CodeUnknownError        int32 = -1
	CodeOK                  int32 = 0      // nil error (ok)
	CodeNoError             int32 = CodeOK // nil error (ok)
	CodeInvalidOp           int32 = 1
	CodeWrongConn           int32 = 100
	CodeConnClosed          int32 = 102
	CodeWriteFailed         int32 = 104
	CodeDialFailed          int32 = 105
	CodeBadMessage          int32 = 400
	CodeUnauthorized        int32 = 401
	CodeNotFound            int32 = 404
	CodeMTypeNotAllowed     int32 = 405
	CodeHandleTimeout       int32 = 408
	CodeInternalServerError int32 = 500
	CodeBadGateway          int32 = 502

	CodeConflict int32 = 409
	// CodeUnsupportedTx                 int32 = 410
	// CodeUnsupportedCodecType          int32 = 415
	// CodeServiceUnavailable            int32 = 503
	// CodeGatewayTimeout                int32 = 504
	// CodeVariantAlsoNegotiates         int32 = 506
	// CodeInsufficientStorage           int32 = 507
	// CodeLoopDetected                  int32 = 508
	// CodeNotExtended                   int32 = 510
	// CodeNetworkAuthenticationRequired int32 = 511
)

// NewStatusByCodeText 通过错误码生成错误信息
func NewStatusByCodeText(code int32, cause interface{}, tagStack bool) *Status {
	stat := NewStatus(code, CodeText(code), cause)
	if tagStack {
		stat.TagStack(1)
	}
	return stat
}

func CodeText(statCode int32) string {
	switch statCode {
	case CodeNoError:
		return ""
	case CodeInvalidOp:
		return "Invalid Operation"
	case CodeBadMessage:
		return "Bad Message"
	case CodeUnauthorized:
		return "Unauthorized"
	case CodeDialFailed:
		return "Dial Failed"
	case CodeWrongConn:
		return "Wrong Connection"
	case CodeConnClosed:
		return "Connection Closed"
	case CodeWriteFailed:
		return "Write Failed"
	case CodeNotFound:
		return "Not Found"
	case CodeHandleTimeout:
		return "Handle Timeout"
	case CodeMTypeNotAllowed:
		return "Message Type Not Allowed"
	case CodeInternalServerError:
		return "Internal Server Error"
	case CodeBadGateway:
		return "Bad Gateway"
	case CodeUnknownError:
		fallthrough
	default:
		return "Unknown Error"
	}
}

var (
	statInvalidOpError      = NewStatus(CodeInvalidOp, CodeText(CodeInvalidOp), "")
	statUnknownError        = NewStatus(CodeUnknownError, CodeText(CodeUnknownError), "")
	statDialFailed          = NewStatus(CodeDialFailed, CodeText(CodeDialFailed), "")
	statConnClosed          = NewStatus(CodeConnClosed, CodeText(CodeConnClosed), "")
	statWriteFailed         = NewStatus(CodeWriteFailed, CodeText(CodeWriteFailed), "")
	statBadMessage          = NewStatus(CodeBadMessage, CodeText(CodeBadMessage), "")
	statNotFound            = NewStatus(CodeNotFound, CodeText(CodeNotFound), "")
	statCodeMTypeNotAllowed = NewStatus(CodeMTypeNotAllowed, CodeText(CodeMTypeNotAllowed), "")
	statHandleTimeout       = NewStatus(CodeHandleTimeout, CodeText(CodeHandleTimeout), "")
	statInternalServerError = NewStatus(CodeInternalServerError, CodeText(CodeInternalServerError), "")
	// 必须要在 post dial和post accept阶段调用，不然就报错
	statUnpreparedError = statInvalidOpError.Copy("Cannot be called during the Non-PostDial and Non-PostAccept phase")
)

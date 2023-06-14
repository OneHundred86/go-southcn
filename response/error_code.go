package response

import (
	"fmt"
	"sync"
)

// 预定义标准正确和错误的错误码
const (
	CodeError = -1
	CodeOK    = 0
)

var preDefinedErrMsg = map[int]string{
	CodeError: "error",
	CodeOK:    "ok",
}

var errMap *sync.Map

func init() {
	errMap = &sync.Map{}

	for code, msg := range preDefinedErrMsg {
		RegisterErrCode(code, msg)
	}
}

// RegisterErrCode 注册错误码
func RegisterErrCode(code int, msg string) {
	errMap.Store(code, msg)
}

// ParseErrCode 解释错误码为错误信息
func ParseErrCode(code int) string {
	if msg, ok := errMap.Load(code); ok {
		return msg.(string)
	}

	return fmt.Sprintf("undefined: %d", code)
}

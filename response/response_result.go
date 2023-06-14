package response

type Result struct {
	Code int         `json:"errcode"`
	Msg  string      `json:"errmessage"`
	Data interface{} `json:"data"`
}

func Ok(data interface{}) Result {
	return Result{
		Code: CodeOK,
		Msg:  ParseErrCode(CodeOK),
		Data: data,
	}
}

func Error(code int, msg string, data interface{}) Result {
	if msg == "" {
		msg = ParseErrCode(code)
	}

	return Result{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

func ErrCode(code int) Result {
	return Result{
		Code: code,
		Msg:  ParseErrCode(code),
	}
}

func ErrMsg(msg string) Result {
	return Result{
		Code: CodeError,
		Msg:  msg,
		Data: nil,
	}
}

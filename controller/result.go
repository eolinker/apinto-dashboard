package controller

const (
	ordinaryCode            = -1 //常规错误(参数错误、sql错误、业务逻辑错误等)
	accessCode              = -2 //没有权限
	CodeLoginInvalid        = -3 //无效token(需要重新登录)
	CodeLoginUserNoExistent = -4 //找不到用户
	CodeLoginPwdErr         = -5 //密码错误
	CodeLoginCodeErr        = -6 //验证校验失败
	codeCertExceedErr       = -7 //证书过期
)

type Result struct {
	Code int         `json:"code"`
	Data interface{} `json:"data,omitempty"`
	Msg  string      `json:"msg"`
}

func NewSuccessResult(data interface{}) *Result {
	return &Result{
		Data: data,
		Msg:  "success",
	}
}

func NewErrorResult(msg string) *Result {
	return &Result{
		Code: ordinaryCode,
		Msg:  msg,
	}
}

func NewNoAccessError(msg string) *Result {
	return &Result{
		Code: accessCode,
		Msg:  msg,
	}
}

func NewCertExceedError(msg string) *Result {
	return &Result{
		Code: codeCertExceedErr,
		Msg:  msg,
	}
}

func NewLoginInvalidError(code int, msg string) *Result {
	return &Result{
		Code: code,
		Msg:  msg,
	}
}

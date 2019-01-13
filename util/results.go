package util

//Result represents HTTP response body.
type Result struct {
	Msg  string      `json:"msg"`
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

const (
	Success = iota
	InvalidRequstParamter
	InternalError
)

// NewResult creates a new result with Code = 0, Msg = "" and Data = nil.
func NewResult() *Result {
	return &Result{
		Code: Success,
		Msg:  "",
		Data: nil,
	}
}

func NewInternalErrorResult() *Result {
	return &Result{
		Code: InternalError,
		Msg:  "internal error.",
		Data: nil,
	}
}

func NewInvalidRequestParameter(desp string) *Result {
	return &Result{
		Code: InvalidRequstParamter,
		Msg:  desp,
		Data: nil,
	}
}

package errno

type BusinessErrorInfo struct {
	Code int
	Msg  string
}

func (e *BusinessErrorInfo) Error() string {
	return e.Msg
}

func (e *BusinessErrorInfo) ReplaceErrorMsgWith(newMsg string) *BusinessErrorInfo {
	return &BusinessErrorInfo{e.Code, newMsg}
}

func (e *BusinessErrorInfo) AppendErrorMsg(appendMsg string) *BusinessErrorInfo {
	return &BusinessErrorInfo{e.Code, e.Msg + " [" + appendMsg + "]"}
}

func DecodeErr(err error) (int, string) {
	if err == nil {
		return Success.Code, Success.Msg
	}
	switch typed := err.(type) {
	case *BusinessErrorInfo:
		return typed.Code, typed.Msg
	default:
		return InternalServerError.Code, err.Error()
	}
}

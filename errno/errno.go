package errno

import "strings"

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

func (e *BusinessErrorInfo) ReplaceErrorByNewError(er error) *BusinessErrorInfo {
	return &BusinessErrorInfo{e.Code, er.Error()}
}

func (e *BusinessErrorInfo) ReplaceErrorByErrors(errs []error) *BusinessErrorInfo {
	strs := make([]string, len(errs))
	for i, err := range errs {
		strs[i] = err.Error()
	}
	return &BusinessErrorInfo{e.Code, strings.Join(strs, ", ")}
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

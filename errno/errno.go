package errno

import "fmt"

type BusinessErrorInfo struct {
	Code int
	Msg  string
}

func (e *BusinessErrorInfo) Error() string {
	return e.Msg
}

func (e *BusinessErrorInfo) ReplaceErrnoMsgWith(newMsg string) *BusinessErrorInfo {
	return &BusinessErrorInfo{e.Code, newMsg}
}

func (e *BusinessErrorInfo) AppendErrorMsg(appendMsg string) *BusinessErrorInfo {
	return &BusinessErrorInfo{e.Code, e.Msg + " [" + appendMsg + "]"}
}

//// next.

type Err struct {
	Code int
	Msg  string
	Err  error
}

func (err *Err) Error() string {
	return fmt.Sprintf("Err => code: %d, Msg: %s, error: %s", err.Code, err.Msg, err.Err)
}

func New(errno *BusinessErrorInfo, err error) *Err {
	return &Err{Code: errno.Code, Msg: errno.Msg, Err: err}
}

func (err *Err) Add(message string) error {
	err.Msg += " " + message
	return err
}

func (err *Err) Addf(format string, arg ...interface{}) error {
	err.Msg += " " + fmt.Sprintf(format, arg...)
	return err
}

func DecodeErr(err error) (int, string) {
	if err == nil {
		return Success.Code, Success.Msg
	}
	switch typed := err.(type) {
	case *Err:
		return typed.Code, typed.Msg
	case *BusinessErrorInfo:
		return typed.Code, typed.Msg
	default:

	}
	return InternalServerError.Code, err.Error()
}

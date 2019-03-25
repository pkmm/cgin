package errno

import "fmt"

type Errno struct {
	Code int
	Msg  string
}

func (e *Errno) Error() string {
	return e.Msg
}

func (e *Errno) ReplaceErrnoMsgWith(newMsg string) *Errno {
	return &Errno{e.Code, newMsg}
}

func (e *Errno) AppendErrorMsg(appendMsg string) *Errno {
	return &Errno{e.Code, e.Msg + " [" + appendMsg + "]"}
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

func New(errno *Errno, err error) *Err {
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
	case *Errno:
		return typed.Code, typed.Msg
	default:

	}
	return InternalServerError.Code, err.Error()
}

package errno

var (
	// common errors
	Success             = &Errno{Code: 0, Msg: "success"}
	InternalServerError = &Errno{Code: 10001, Msg: "Internal server error."}
	ErrBind             = &Errno{Code: 10002, Msg: "Error occurred while binding the request body to struct."}

	// user errors
	ErrUserNotFound         = &Errno{Code: 20001, Msg: "The user was not found."}
	ErrCheckZfAccountFailed = &Errno{Code: 20002, Msg: "check account failed of zcmu."}
)

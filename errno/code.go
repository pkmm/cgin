package errno

var (
	// common errors
	Success                = &Errno{Code: 0, Msg: "success"}
	InternalServerError    = &Errno{Code: 10001, Msg: "Internal server error"}
	ErrBind                = &Errno{Code: 10002, Msg: "Error occurred while binding the request body to struct"}
	InvalidParameters      = &Errno{10003, "Invalid Parameters"}
	TokenNotValid          = &Errno{10004, "Token is not valid"}
	GenerateJwtTokenFailed = &Errno{10004, "Generate token was failed"}
	UserNotAuth            = &Errno{10004, "You must login"}

	// user errors
	ErrUserNotFound         = &Errno{Code: 20001, Msg: "The user was not found"}
	ErrCheckZfAccountFailed = &Errno{Code: 20002, Msg: "Check account failed of zcmu"}
)

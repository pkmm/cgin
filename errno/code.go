package errno

var (
	// common errors
	// error code begin with 1... for system error
	Success                = &Errno{Code: 0, Msg: "success"}
	InternalServerError    = &Errno{Code: 10001, Msg: "Internal server error"}
	InvalidParameters      = &Errno{10003, "Invalid Parameters"}
	TokenNotValid          = &Errno{10004, "Token is not valid"}
	GenerateJwtTokenFailed = &Errno{10004, "Generate token was failed"}
	UserNotAuth            = &Errno{10004, "You must login"}

	// user errors
	// error code begin with 2... for business error
	UserNotFoundException         = &Errno{Code: 20001, Msg: "The user was not found"}
	CheckZfAccountFailedException = &Errno{Code: 20002, Msg: "Check account failed of zcmu"}
)

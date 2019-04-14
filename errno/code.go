package errno

var (
	// common errors
	// error code begin with 1... for system error
	Success                = &Errno{00000, "success"}
	InternalServerError    = &Errno{10001, "Internal server error"}
	InvalidParameters      = &Errno{10002, "Invalid Parameters"}
	TokenNotValid          = &Errno{10003, "Token is not valid"}
	GenerateJwtTokenFailed = &Errno{10004, "Generate token was failed"}
	UserNotAuth            = &Errno{10005, "You must login"}
	Welcome                = &Errno{10006, "Welcome to Cgin."}
	NotSuchRouteException  = &Errno{10007, "Not found resource."}

	// user errors
	// error code begin with 2... for business error
	UserNotFoundException         = &Errno{20001, "The user was not found"}
	CheckZfAccountFailedException = &Errno{20002, "Check account failed of zcmu"}
)

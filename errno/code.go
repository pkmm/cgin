package errno

var (
	// common errors
	// error code begin with 1... for system error
	Success                = &BusinessErrorInfo{00000, "success."}
	InternalServerError    = &BusinessErrorInfo{10001, "Internal server error."}
	InvalidParameters      = &BusinessErrorInfo{10002, "Invalid Parameters."}
	TokenNotValid          = &BusinessErrorInfo{10003, "Token is not valid."}
	GenerateJwtTokenFailed = &BusinessErrorInfo{10004, "Generate token was failed."}
	UserNotAuth            = &BusinessErrorInfo{10005, "You must login."}
	Welcome                = &BusinessErrorInfo{10006, "Hello, I am still alive."}
	NotSuchRouteException  = &BusinessErrorInfo{10007, "Not found resource."}
	NormalException        = &BusinessErrorInfo{10008, "exception:"}

	// user errors
	// error code begin with 2... for business error
	UserNotFoundException         = &BusinessErrorInfo{20001, "The user was not found."}
	CheckZfAccountFailedException = &BusinessErrorInfo{20002, "Check account failed of education system."}
)

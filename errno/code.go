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
	NormalException        = &BusinessErrorInfo{10008, "Exception:"}

	// user errors
	// error code begin with 2... for business error
	UserNotFoundException         = &BusinessErrorInfo{20001, "The user was not found."}
	CheckZfAccountFailedException = &BusinessErrorInfo{20002, "Check account failed of education system."}

	// 背单词的错误信息
	RememberTaskNotBegin   = &BusinessErrorInfo{3001, "任务还未开始."}
	RememberTaskHasDone    = &BusinessErrorInfo{3002, "任务已经结束了."}
	RememberRecordNotFound = &BusinessErrorInfo{3003, "未找到相关的任务计划"}
)

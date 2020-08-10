package errno

var (
	// common errors
	// error code begin with 1... for system error
	Success                = &BusinessErrorInfo{10000, "OK"}
	InternalServerError    = &BusinessErrorInfo{10001, "internal server error"}
	InvalidParameters      = &BusinessErrorInfo{10002, "invalid parameters"}
	ErrorTokenNotValid     = &BusinessErrorInfo{10003, "token is not valid"}
	GenerateJwtTokenFailed = &BusinessErrorInfo{10004, "generate token was failed"}
	UserNotAuth            = &BusinessErrorInfo{10005, "you must login"}
	Welcome                = &BusinessErrorInfo{10006, "services developed by golang"}
	NotSuchRouteException  = &BusinessErrorInfo{10007, "not found resource"}
	NormalException        = &BusinessErrorInfo{10008, "exception:"}
	PermissionDenied       = &BusinessErrorInfo{10009, "permission denied"}
	ErrorTokenExpired      = &BusinessErrorInfo{10010, "token expired"}

	// user errors
	// error code begin with 2... for business error
	UserNotFoundException         = &BusinessErrorInfo{20001, "the user was not found"}
	CheckZfAccountFailedException = &BusinessErrorInfo{20002, "check account failed of education system"}

	// 背单词的错误信息
	RememberTaskNotBegin   = &BusinessErrorInfo{3001, "任务还未开始"}
	RememberTaskHasDone    = &BusinessErrorInfo{3002, "任务已经结束了"}
	RememberRecordNotFound = &BusinessErrorInfo{3003, "未找到相关的任务计划"}
)

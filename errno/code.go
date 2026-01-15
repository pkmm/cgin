package errno

var (
	// common errors
	// error code begin with 1... for system error
	Success             = &BusinessErrorInfo{10000, "OK"}
	InternalServerError = &BusinessErrorInfo{10001, "internal server error"}
)

package conf

import (
	"time"
)

const (
	// 运行环境的参数设置
	appEnvironment = "appEnv"
	appEnvProd     = "prod"
	appEnvDev      = "dev"

	// mysql
	mysqlHost     = "mysql.host"
	mysqlPort     = "mysql.port"
	mysqlUser     = "mysql.username"
	mysqlPassword = "mysql.password"
	mysqlDatabase = "mysql.database"
	mysqlTimezone = "mysql.timezone"

	// jwt config
	jwtExpireAt = "jwt.day"
	JwtSignKey  = "jwt.secret"
)

// 封装一些函数

// 返回app环境 默认是prod
func AppEnvironment() string {
	return AppConfig.DefaultString(appEnvironment, appEnvProd)
}

func IsProd() bool {
	return AppConfig.String(appEnvironment) == appEnvProd
}

func IsDev() bool {
	return AppConfig.String(appEnvironment) == appEnvDev
}

// 天数对应的小时 1天
// 返回值单位小时
func GetJwtExpiresAt() time.Duration {
	day := AppConfig.DefaultInt64(jwtExpireAt, 1)
	return time.Duration(day) * 24 * time.Hour
}

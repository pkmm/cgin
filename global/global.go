package global

import (
	"cgin/config"
	"cgin/service/workerpool"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	DB         *gorm.DB
	Redis      *redis.Client
	Config     config.Server
	G_VP       *viper.Viper
	GLog       *zap.Logger
	WorkerPool *workerpool.Pool // 协程池子
)

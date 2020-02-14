package context_helper

import (
	"cgin/errno"
	"cgin/middleware"
	"cgin/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type contextHelper struct {
	UserId uint64
	Params map[string]interface{}
	ctx    *gin.Context
}

func New(c *gin.Context) *contextHelper {
	return &contextHelper{UserId: 0, Params: nil, ctx:c}
}

func (b *contextHelper) checkContext() {
	if b.ctx == nil {
		panic("must set gin.context first.")
	}
}

func (b *contextHelper) get(key string) interface{} {
	b.checkContext()
	switch b.ctx.Request.Method {
	case http.MethodGet:
		v := b.ctx.Query(key)
		if v == "" {
			panic(errno.NormalException.AppendErrorMsg("解析参数不存在:" + key))
		}
		return v
	case http.MethodPost:
		if b.Params == nil {
			b.Params = map[string]interface{}{}
			if err := b.ctx.ShouldBindJSON(&b.Params); err != nil {
				panic(errno.NormalException.AppendErrorMsg("解析参数错误:" + err.Error()))
			}
		}
		return b.Params[key]
	}
	panic(errno.NormalException.AppendErrorMsg("解析参数不存在:" + key))
}

func (b *contextHelper) GetInt(key string) int {
	return int(b.GetInt64(key))
}

func (b *contextHelper) GetInt64(key string) int64 {
	v := b.get(key)
	switch v.(type) {
	case string:
		v64, _ := strconv.ParseInt(v.(string), 10, 64)
		return v64
	case float64:
		return int64(v.(float64))
	default:
		panic("无法解析的类型：" + key)
	}
}

func (b *contextHelper) GetString(key string) string {
	return b.get(key).(string)
}

func (b *contextHelper) GetAuthUserId() uint64 {
	if b.UserId != 0 {
		return b.UserId
	}
	val, ok := b.ctx.Get(middleware.UID)
	if !ok {
		panic(errno.UserNotAuth)
	}

	userId, ok := val.(uint64)
	if !ok || userId == 0 {
		panic(errno.UserNotAuth)
	}
	b.UserId = userId
	return userId
}

// 处理成功的请求
func (b *contextHelper) Response(responseData interface{}) {
	service.SendResponseSuccess(b.ctx, responseData)
}

package contextHelper

import (
	"cgin/errno"
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
	return &contextHelper{UserId: 0, Params: nil, ctx: c}
}

func (b *contextHelper) checkContext() {
	if b.ctx == nil {
		panic("must set gin.context first.")
	}
}

// 获取query或者body中的参数
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

// 获取body 或者是 query中参数方法结束

//--- 获取当前认证的用户id ----
func (b *contextHelper) GetAuthUserId() uint64 {
	if b.UserId != 0 {
		return b.UserId
	}
	val := b.ctx.MustGet("claims")
	awaitUse, ok := val.(*service.AuthClaims)
	if !ok || awaitUse.Uid == 0 {
		panic(errno.UserNotAuth)
	}
	b.UserId = awaitUse.Uid
	return b.UserId
}

// 处理成功的请求
func (b *contextHelper) Response(responseData interface{}) {
	service.SendResponseSuccess(b.ctx, responseData)
}

// 需要认证的 check
func (b *contextHelper) NeedAuthOrPanic() {
	b.GetAuthUserId()
}

// 获取 path params
func (b *contextHelper) GetPathInt(key string) int {
	i, err := strconv.Atoi(b.ctx.Param(key))
	if err != nil {
		panic(errno.InvalidParameters.ReplaceErrorMsgWith("参数不合法"))
	}
	return i
}

func (b *contextHelper) GetPathInt64(key string) int64 {
	if i64, err := strconv.ParseInt(b.ctx.Param(key), 10, 64); err == nil {
		return i64
	} else {
		panic(errno.InvalidParameters.ReplaceErrorMsgWith("参数不合法" + err.Error()))
	}
}

func (b *contextHelper) GetPathUint64(key string) uint64 {
	return uint64(b.GetPathInt64(key))
}

// 结束获取params

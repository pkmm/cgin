package middleware

import (
	"cgin/errno"
	"cgin/model"
	"cgin/service"
	"github.com/gin-gonic/gin"
	"regexp"
)

// RESTFul api 权限控制

func PermissionCheck() gin.HandlerFunc {
	return func(context *gin.Context) {
		claims, ok := context.MustGet("claims").(*service.AuthClaims)
		if !ok {
			panic(errno.TokenNotValid)
		}
		if claims.RoleId == model.RoleAdmin {
			context.Next()
			return
		}
		// not admin user
		// 获取请求的资源
		urlPath := context.Request.URL.Path
		// add slash
		if urlPath[len(urlPath)-1] != '/' {
			urlPath = urlPath + "/"
		}
		pattern := regexp.MustCompile(`api/v.*?/(.*?)/`)
		resourceIdentity := pattern.FindStringSubmatch(urlPath)
		if nil == resourceIdentity {
			panic(errno.NotSuchRouteException)
		}
		err, permission := model.GetPermissionByResourceIdentityAndMethod(resourceIdentity[1], context.Request.Method)
		if err != nil {
			panic(errno.PermissionDenied)
		}
		if permission.Effect == model.EffectAllow {
			context.Next()
		} else if permission.Effect == model.EffectOwner {
			ok = model.HasPermission(permission.ID, claims.RoleId)
			if ok {
				context.Next()
			} else {
				panic(errno.PermissionDenied)
			}
		} else {
			panic(errno.PermissionDenied)
		}
	}
}

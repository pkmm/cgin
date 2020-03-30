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
	// TODO: admin can operate all resources.
	return func(context *gin.Context) {
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
			claims, ok := context.MustGet("claims").(*service.AuthClaims)
			if !ok {
				panic(errno.PermissionDenied)
			}
			roleIds := claims.RoleIds
			flag := true
			for _, roleId := range roleIds {
				ok := model.HasPermission(permission.ID, roleId)
				if ok {
					flag = false
					context.Next()
					break
				}
			}
			if flag {
				panic(errno.PermissionDenied)
			}
		}
	}
}

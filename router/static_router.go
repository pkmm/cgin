package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func mapStaticRouter(router *gin.Engine) {
	router.StaticFS("/static", http.Dir("static"))
	router.StaticFile("/favicon.ico", "./static/favicon.ico")
	router.StaticFile("/admin", "static/web_front/index.html")
}

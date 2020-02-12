package main

import (
	"cgin/conf"
	"cgin/router"
	_ "cgin/task"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

var env = conf.AppEnvDev
var port = "8654"

func init() {
	env = conf.AppConfig.DefaultString(conf.AppEnvironment, conf.AppEnvProd)
	if env == conf.AppEnvProd {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
}

// @title My Server cgin
// @version 1.0
// @description this is a custom server of my interesting.

// @host localhost:8654
// @BasePath /api
func main() {

	handlers := router.MapRoute()
	server := &http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: handlers,
	}

	conf.Logger.Info("cgin is running at [%s]", "http://localhost:"+port)
	fmt.Printf("cgin is running at [%s]", "http://localhost:"+port)
	server.ListenAndServe()
}

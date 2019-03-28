package main

import (
	"cgin/conf"
	"cgin/controller"
	_ "cgin/task"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

var env = "prod"
var port = "8654"

func init() {
	env = conf.AppConfig.DefaultString("appEnv", "prod")
	if env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
}

func main() {

	router := controller.MapRoute()
	server := &http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: router,
	}

	conf.AppLogger.Info( "cgin is running at [%s]", "http://localhost:"+port)
	fmt.Printf("cgin is running at [%s]", "http://localhost:"+port)
	server.ListenAndServe()
}

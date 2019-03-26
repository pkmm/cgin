package main

import (
	"cgin/conf"
	"cgin/controller"
	_ "cgin/task"
	"github.com/gin-gonic/gin"
	"net/http"
)

func init() {
	env := conf.AppConfig.String("appEnv")
	if env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
}

func main() {

	var myPort = "8654"
	//env := conf.AppConfig.DefaultString("appEnv", "prod")
	//if env == "prod" {
	//	myPort = "8189"
	//}

	router := controller.MapRoute()
	server := &http.Server{
		Addr:    "0.0.0.0:" + myPort,
		Handler: router,
	}

	conf.AppLogger.Info("pkmm gin is running [%s]", "http://localhost:"+myPort)
	server.ListenAndServe()
}

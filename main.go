package main

import (
	"cgin/conf"
	"cgin/router"
	"cgin/task"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

var port = "8654"

// @title My Server cgin
// @version 1.0
// @description this is a custom server of my interesting.

// @host localhost:8654
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	defer func() {
		// release some resource.
		fmt.Println("do some clean work.")
		task.CleanPool()
		_ = conf.DB.Close()
	}()

	if conf.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	handlers := router.InitRouter()
	server := &http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: handlers,
	}

	conf.Logger.Info("cgin is running at [%s]", "http://localhost:"+port)
	fmt.Printf("cgin is running at [%s]", "http://localhost:"+port)
	server.ListenAndServe()
}

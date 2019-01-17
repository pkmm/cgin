package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"net/http"
	"pkmm_gin/controller"
	_ "pkmm_gin/task"
)



func init() {
	//log.SetFlags(log.Ldate | log.Ltime | log.LUTC | log.Lshortfile)

	// 设置beego logs
	logs.SetLogger(logs.AdapterFile, `{"filename":"logs/pkmm_gin.log","level":7,"daily":true,"maxdays":2}`)
	logs.EnableFuncCallDepth(true)
	logs.Async(1e3)

	gin.SetMode(gin.ReleaseMode)

	// Usage like this. (beego myConfig model)
	//appName := AppConfig.String("appName")
	//fmt.Println(appName)
}

func main() {

	router := controller.MapRoute()
	server := &http.Server{
		Addr:    "0.0.0.0:" + "8654",
		Handler: router,
	}

	logs.Info("pkmm gin is running [%s]", "http://localhost:8654")
	server.ListenAndServe()
}

package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"image/gif"
	"net/http"
	"pkmm_gin/utility"
	"time"

	"os"
	"path/filepath"
	"pkmm_gin/controller"
	"pkmm_gin/middleware"
	"pkmm_gin/model"
	_ "pkmm_gin/task"
)

func main() {

	//gin.DisableConsoleColor()
	//f, _ := os.Create("gin.log")
	//gin.DefaultWriter = io.MultiWriter(f)

	r := gin.Default()

	// gin.H is a shortcut for map[string]interface{}
	r.GET("/someJSON", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	})

	r.GET("/moreJSON", func(c *gin.Context) {
		// You also can use a struct
		var msg struct {
			Name    string `json:"user"`
			Message string
			Number  int
			Date    time.Time `json:"date" time_format:"2006-01-02 15:04:05"`
		}
		msg.Name = "Lena"
		msg.Message = "hey"
		msg.Number = 123562
		msg.Date = time.Now()
		// Note that msg.Name becomes "user" in the JSON
		// Will output  :   {"user": "Lena", "Message": "hey", "Number": 123}

		// DB code
		var user [] model.Student
		db := model.GetDB()
		db.Find(&user)
		c.JSON(http.StatusOK, user)
	})

	r.GET("/someXML", func(c *gin.Context) {
		c.XML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	})

	r.GET("/someYAML", func(c *gin.Context) {
		c.YAML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	})

	r.GET("/long_async", func(context *gin.Context) {
		contextCopy := context.Copy()
		go func() {
			time.Sleep(5 * time.Second)
			fmt.Println("Done! in path" + contextCopy.Request.URL.Path)
		}()
		context.JSON(http.StatusOK, gin.H{"status": "开始服务了"})
	})

	r.GET("/xx", func(context *gin.Context) {
		appPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		workPath, _ := os.Getwd()
		context.JSON(http.StatusOK, gin.H{
			"app_path":  appPath,
			"work_path": workPath,
		})
	})

	// 自己定义的restFul路由设计
	v1 := r.Group("/v1", middleware.ApiAuth())
	{
		v1.POST("/student/login", controller.Login)
		v1.POST("/student", controller.AddStudent)
		//v1.GET("/student/:id", controller.GetStudent)
		v1.GET("/student/:id/scores", controller.GetStudentScores)
	}

	r.GET("/zxx", func(context *gin.Context) {
		rep, err := http.Get("http://zfxk.zjtcm.net/CheckCode.aspx?")
		if err != nil {
			fmt.Println(err)
		}
		defer rep.Body.Close()
		im, err := gif.Decode(rep.Body)
		fmt.Println(err)
		r, err := utility.Predict(im, false)
		fmt.Println(r, err)
		context.String(200, r)
	})

	r.Static("/static", "./static")

	// Listen and serve on 0.0.0.0:8080
	r.Run(":80")
}

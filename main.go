package main

import (
	"github.com/astaxie/beego/logs"
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

	//gin.SetMode(gin.ReleaseMode)

	// Usage like this. (beego myConfig model)
	//appName := AppConfig.String("appName")
	//fmt.Println(appName)
}

func main() {

	router := controller.MapRoute()
	server := &http.Server{
		Addr:    "0.0.0.0:" + "80",
		Handler: router,
	}

	logs.Info("pkmm gin is running [%s]", "http://localhost:80")
	server.ListenAndServe()

	//	//gin.DisableConsoleColor()
	//	//f, _ := os.Create("gin.log")
	//	//gin.DefaultWriter = io.MultiWriter(f)
	//
	//	r := gin.Default()
	//
	//	// gin.H is a shortcut for map[string]interface{}
	//	r.GET("/someJSON", func(c *gin.Context) {
	//		c.JSON(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	//	})
	//
	//	r.GET("/moreJSON", func(c *gin.Context) {
	//		// You also can use a struct
	//		var msg struct {
	//			name    string `json:"user"`
	//			Message string
	//			Number  int
	//			Date    time.Time `json:"date" time_format:"2006-01-02 15:04:05"`
	//		}
	//		msg.name = "Lena"
	//		msg.Message = "hey"
	//		msg.Number = 123562
	//		msg.Date = time.Now()
	//		// Note that msg.name becomes "user" in the JSON
	//		// Will output  :   {"user": "Lena", "Message": "hey", "Number": 123}
	//
	//		// DB code
	//		var user []model.Student
	//		db := model.GetDB()
	//		db.Find(&user)
	//		c.JSON(http.StatusOK, user)
	//	})
	//
	//	r.GET("/someXML", func(c *gin.Context) {
	//		c.XML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	//	})
	//
	//	r.GET("/someYAML", func(c *gin.Context) {
	//		c.YAML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	//	})
	//
	//	r.GET("/long_async", func(context *gin.Context) {
	//		contextCopy := context.Copy()
	//		go func() {
	//			time.Sleep(5 * time.Second)
	//			fmt.Println("Done! in path" + contextCopy.Request.URL.Path)
	//		}()
	//		context.JSON(http.StatusOK, gin.H{"status": "开始服务了"})
	//	})
	//
	//	r.GET("/xx", func(context *gin.Context) {
	//		appPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	//		workPath, _ := os.Getwd()
	//		context.JSON(http.StatusOK, gin.H{
	//			"app_path":  appPath,
	//			"work_path": workPath,
	//		})
	//	})
	//
	//	// 自己定义的restFul路由设计
	//	v1 := r.Group("/v1", middleware.ApiAuth())
	//	{
	//		v1.POST("/wx/wx_login", controller.WxLogin)
	//		v1.GET("/student/:id/scores", controller.GetStudentScores)
	//		v1.GET("/get_scores", controller.GetScores)
	//		v1.POST("/student/update_zcmu_account", controller.UpdateZcmuAccount)
	//	}
	//	//r.Group("/v1")
	//	{
	//		// not need middleware
	//		//r.POST("/v1/student/update_zcmu_account", controller.UpdateZcmuAccount)
	//	}
	//
	//	r.GET("/zxx", func(context *gin.Context) {
	//		c := zf.NewCrawl("201412200903014", "ZYPzhouyaping123")
	//		t, err := c.touchScorePageForGetViewState()
	//		fmt.Println(err)
	//		//t := model.GetScore(11)
	//		context.JSON(200, gin.H{
	//			"scores": t,
	//		})
	//		//rep, err := http.Get("http://zfxk.zjtcm.net/CheckCode.aspx?")
	//		//if err != nil {
	//		//	fmt.Println(err)
	//		//}
	//		//defer rep.Body.Close()
	//		//im, err := gif.Decode(rep.Body)
	//		//fmt.Println(err)
	//		//r, err := util.Predict(im, false)
	//		//fmt.Println(r, err)
	//		//context.String(200, r)
	//	})
	//
	//
	//	r.GET("/test", func(c *gin.Context) {
	//		//var student model.Student
	//		//student.Openid = "23"
	//		//student.City = "233333"
	//		//student.AvatarUrl = "https://www.baidu.com/we"
	//		//r := model.UpdateUserWeChatInfo(student, "23")
	//		//c.JSON(200, gin.H{
	//		//	"w": r,
	//		//})
	//		stu := model.Student{Id: 21}
	//		//var scores []model.Score
	//		model.GetDB().Debug().Model(&stu).Preload("Scores").Where("id = ?", 21).First(&stu)
	//		c.JSON(200, gin.H{
	//			"w": stu,
	//		})
	//	})
	//
	//
	//	// Listen and serve on 0.0.0.0:8080
	//	r.Run(":8056")
	//}
}

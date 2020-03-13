package router

import (
	"cgin/controller"
	"cgin/middleware"
	"github.com/gin-gonic/gin"
)

func mapStudentRouter(router *gin.Engine) {
	apiStudent := router.Group(Student).Use(middleware.Auth)
	{
		// 尝试使用路径参数restful api
		// TODO: 权限控制
		// 普通的是比较简单的
		apiStudent.GET("/:studentId", controller.Student.GetStudent)
		//apiStudent.GET("/scores", controller.Student.GetScores)
		apiStudent.POST("/update_edu_account", controller.Student.UpdateEduAccount)
	}
}

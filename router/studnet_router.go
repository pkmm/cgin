package router

import (
	"cgin/controller/api/v1"
	"cgin/middleware"
	"github.com/gin-gonic/gin"
)

func initStudentRouter(studentRouter *gin.Engine) {
	apiStudent := studentRouter.Group(Student).Use(middleware.Auth())
	{
		// 尝试使用路径参数restful api
		// TODO: 权限控制
		// 普通的是比较简单的
		apiStudent.GET("/:studentId", v1.Student.GetStudent)
		apiStudent.POST("/update_edu_account", v1.Student.UpdateOrCreateEduAccount)
	}
}

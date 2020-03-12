package router

import (
	"cgin/controller"
	"cgin/middleware"
	"github.com/gin-gonic/gin"
)

func mapStudentRouter(router *gin.Engine) {
	apiStudent := router.Group(Student).Use(middleware.Auth)
	{
		apiStudent.GET("/:studentId", controller.Student.GetStudent)
		//apiStudent.GET("/scores", controller.Student.GetScores)
		apiStudent.POST("/update_edu_account", controller.Student.UpdateEduAccount)
	}
}

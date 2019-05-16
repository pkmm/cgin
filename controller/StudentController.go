package controller

import (
	"cgin/service"
	"github.com/gin-gonic/gin"
)

type studentController struct {
	BaseController
}

var Student = &studentController{}


func (s *studentController) GetStudent(c *gin.Context) {
	s.BaseController.GetAuthUserId(c)
	student := service.User.GetStudent(s.UserId)
	s.respData["student"] =  student
	service.SendResponseSuccess(c, s.respData)
}

func (s *studentController) GetScores(c *gin.Context) {
	s.BaseController.GetAuthUserId(c)
	scores := service.ScoreService.GetOwnScores(s.UserId)
	s.respData["scores"] = scores
	service.SendResponseSuccess(c, s.respData)
}
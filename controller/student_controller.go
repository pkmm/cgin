package controller

import (
	"cgin/controller/context_helper"
	"cgin/errno"
	"cgin/service"
	"cgin/zcmu"
	"github.com/gin-gonic/gin"
)

type studentController struct {}

var Student = &studentController{}

func (s *studentController) GetStudent(c *gin.Context) {
	helper := context_helper.New(c)
	student := service.User.GetStudentByUserId(helper.GetAuthUserId())
	data := gin.H{
		"student": student,
	}
	helper.Response(data)
}

func (s *studentController) GetScores(c *gin.Context) {
	helper := context_helper.New(c)
	scores := service.ScoreService.GetOwnScores(helper.GetAuthUserId())
	if len(scores) == 0 {
		student := service.User.GetStudentByUserId(helper.GetAuthUserId())
		if student == nil {
			panic(errno.NormalException.AppendErrorMsg("用户没有学生信息"))
		}
		worker, err := zcmu.NewCrawl(student.Number, student.Password)
		if err != nil {
			panic(errno.NormalException.AppendErrorMsg(err.Error()))
		}
		if scores, err := worker.GetScores(); err == nil {
			modelScores := service.ScoreService.SaveStudentScoresFromCrawl(scores, student.Id)
			service.SendResponseSuccess(c, gin.H{
				"scores": modelScores,
			})
			return
		} else {
			panic(errno.NormalException.AppendErrorMsg(err.Error()))
		}
	}
	data := gin.H{
		"scores": scores,
	}
	helper.Response(data)
}

func (s *studentController) UpdateEduAccount(c *gin.Context) {
	helper := context_helper.New(c)
	var (
		studentNumber, password string
		err                     error
	)
	studentNumber = helper.GetString("student_number")
	password = helper.GetString("password")
	// 调用zcmu接口检测账号密码是否正确
	checker, err := zcmu.NewCrawl(studentNumber, password)
	if err != nil {
		panic(errno.NormalException.AppendErrorMsg(err.Error()))
	}
	if msg := checker.CheckAccount(); msg != "" {
		panic(errno.NormalException.ReplaceErrorMsgWith(msg))
	}
	student := service.User.UpdateStudentInfoByUserId(studentNumber, password, helper.GetAuthUserId())
	if student == nil {
		panic(errno.NormalException.AppendErrorMsg("创建学生失败"))
	}
	data := gin.H{
		"student": student,
	}
	helper.Response(data)
}

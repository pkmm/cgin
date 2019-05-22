package controller

import (
	"cgin/errno"
	"cgin/service"
	"cgin/zcmu"
	"github.com/gin-gonic/gin"
)

type studentController struct {
	BaseController
}

var Student = &studentController{}


func (s *studentController) GetStudent(c *gin.Context) {
	s.GetAuthUserId(c)
	student := service.User.GetStudent(s.UserId)
	data := gin.H{
		"student": student,
	}
	service.SendResponseSuccess(c, data)
}

func (s *studentController) GetScores(c *gin.Context) {
	s.GetAuthUserId(c)
	scores := service.ScoreService.GetOwnScores(s.UserId)
	if len(scores) == 0 {
		student := service.User.GetStudent(s.UserId)
		if student == nil {
			panic(errno.NormalException.AppendErrorMsg("用户没有学生信息"))
		}
		worker, err := zcmu.NewCrawl(student.Number, student.Password)
		if err != nil {
			panic(errno.NormalException.AppendErrorMsg(err.Error()))
		}
		if scores, err := worker.GetScores(); err == nil {
			service.SendResponseSuccess(c, gin.H{
				"scores": scores,
			})
			return
		} else {
			panic(errno.NormalException.AppendErrorMsg(err.Error()))
		}
	}
	data := gin.H{
		"scores": scores,
	}
	service.SendResponseSuccess(c, data)
}

func (s *studentController) UpdateEduAccount(c *gin.Context) {
	s.GetAuthUserId(c)
	s.ProcessParams(c)
	var (
		studentNumber, password string
		ok bool
		err error
	)
	if studentNumber, ok = s.Params["student_number"].(string); !ok {
		panic(errno.InvalidParameters.AppendErrorMsg("参数student number错误"))
	}
	if password, ok = s.Params["password"].(string); !ok {
		panic(errno.InvalidParameters.AppendErrorMsg("参数password错误"))
	}
	student := service.User.UpdateStudentInfo(studentNumber, password, s.UserId)
	if student == nil {
		panic(errno.NormalException.AppendErrorMsg("创建学生失败"))
	}
	// 调用zcmu接口检测账号密码是否正确
	checker, err := zcmu.NewCrawl(studentNumber, password)
	if err != nil {
		panic(errno.NormalException.AppendErrorMsg(err.Error()))
	}
	if msg := checker.CheckAccount(); msg != "" {
		panic(errno.NormalException.ReplaceErrorMsgWith(msg))
	}
	data := gin.H{
		"student": student,
	}
	service.SendResponseSuccess(c, data)
}
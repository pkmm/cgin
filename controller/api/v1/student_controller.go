package v1

import (
	"cgin/controller/context_helper"
	"cgin/errno"
	"cgin/model"
	"cgin/service"
	"cgin/zcmu"
	"github.com/gin-gonic/gin"
)

type studentController struct{}

var Student = &studentController{}

// @Summary 获取当前用户的学生信息
// @Produce json
// @Router /students/{studentId} [get]
// @Param studentId path uint64 true "student id"
// @Success 200 {object} service.Response
// @Security ApiKeyAuth
func (s *studentController) GetStudent(c *gin.Context) {
	helper := context_helper.New(c)
	err, student := model.GetStudentByUserId(helper.GetAuthUserId())
	if err != nil {
		panic(errno.NormalException.ReplaceErrorMsgWith(err.Error()))
	}
	data := gin.H{
		"student": student,
	}
	helper.Response(data)
}

// @Summary 获取学生的成绩
// @Produce json
// @Router /scores [get]
// @Success 200 {object} service.Response
// @Security ApiKeyAuth
func (s *studentController) GetScores(c *gin.Context) {
	helper := context_helper.New(c)
	err, scores := service.ScoreService.GetOwnScores(helper.GetAuthUserId())
	if err != nil {
		panic(errno.NormalException.ReplaceErrorMsgWith(err.Error()))
	}
	if len(*scores) == 0 {
		err, student := model.GetStudentByUserId(helper.GetAuthUserId())
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
	helper.Response(gin.H{
		"scores": scores,
	})
}

// @Summary 更新学生的信息
// @Produce json
// @Security ApiKeyAuth
// @Router /students/update_edu_account [post]
// @Param auth body co.EduAccount true "update edu account info"
// @Success 200 {object} service.Response
func (s *studentController) UpdateOrCreateEduAccount(c *gin.Context) {
	helper := context_helper.New(c)
	var (
		number, password string
	)
	number = helper.GetString("student_number")
	password = helper.GetString("password")
	// 检测账号密码是否正确
	checker, err := zcmu.NewCrawl(number, password)
	if err != nil {
		panic(errno.NormalException.AppendErrorMsg(err.Error()))
	}
	if msg := checker.CheckAccount(); msg != "" {
		panic(errno.NormalException.ReplaceErrorMsgWith(msg))
	}

	err, student := (&model.Student{
		UserId:   helper.GetAuthUserId(),
		Password: password,
		Number:   number,
	}).UpdateOrCreate()

	if err != nil {
		panic(errno.NormalException.AppendErrorMsg("创建学生失败: " + err.Error()))
	}
	data := gin.H{
		"student": student,
	}
	helper.Response(data)
}

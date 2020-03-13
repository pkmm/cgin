package v1

import (
	"cgin/controller/context_helper"
	"cgin/errno"
	"cgin/service"
	"cgin/zcmu"
	"github.com/gin-gonic/gin"
)

type studentController struct {}

var Student = &studentController{}

// @Summary 获取当前用户的学生信息
// @Produce json
// @Router /students/{studentId} [get]
// @Param studentId path uint64 true "student id"
// @Success 200 {object} service.Response
// @Security ApiKeyAuth
func (s *studentController) GetStudent(c *gin.Context) {
	helper := context_helper.New(c)
	studentId := helper.GetPathUint64("studentId")
	student := service.User.GetStudentByUserId(studentId)
	// TODO: 如果student id 是0 且是该查询者是admin显示分页的学生信息
	if studentId != helper.GetAuthUserId() {
		panic(errno.NormalException.ReplaceErrorMsgWith("没有权限"))
	}
	data := gin.H{
		"student": student,
	}
	helper.Response(data)
}

// @Summary 获取学生的成绩
// @Produce json
// @Router /scores [get]
// @Param studentId query uint64 true "id of student"
// @Success 200 {object} service.Response
// @Security ApiKeyAuth
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

// @Summary 更新学生的信息
// @Produce json
// @Router /students/update_edu_account [post]
// @Param auth body co.EduAccount true "update edu account info"
// @Success 200 {object} service.Response
func (s *studentController) UpdateEduAccount(c *gin.Context) {
	helper := context_helper.New(c)
	var (
		studentNumber, password string
		//err                     error
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
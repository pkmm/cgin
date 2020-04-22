package v1

import (
	"cgin/controller/contextHelper"
	"cgin/errno"
	"cgin/model"
	"cgin/service"
	"cgin/zcmu"
	"github.com/gin-gonic/gin"
)

type studentController struct{}

var Student = &studentController{}

// @Summary 获取当前用户的学生信息
// @Tags Student
// @Produce json
// @Router /students/{studentId} [get]
// @Param studentId path uint64 true "student id"
// @Success 200 object service.Response
// @Security ApiKeyAuth
func (s *studentController) GetStudent(c *gin.Context) {
	helper := contextHelper.New(c)
	studentId := helper.GetPathUint64("studentId")
	err, student := model.GetStudentByUserId(helper.GetAuthUserId())
	if student.Id != studentId && !model.IsAdmin(helper.GetAuthUserId()) {
		panic(errno.PermissionDenied)
	}

	if err != nil {
		panic(errno.NormalException.ReplaceErrorMsgWith(err.Error()))
	}
	helper.Response(gin.H{
		"student": student,
	})
}

// @Summary 获取学生的成绩
// @Produce json
// @Tags Student
// @Router /scores [get]
// @Success 200 object service.Response
// @Security ApiKeyAuth
func (s *studentController) GetScores(c *gin.Context) {
	helper := contextHelper.New(c)
	err, student := model.GetStudentByUserId(helper.GetAuthUserId())
	if student == nil {
		panic(errno.NormalException.AppendErrorMsg("用户没有学生信息"))
	}
	robot := zcmu.New(student.Number, student.Password)
	err = robot.Login()
	if err != nil {
		panic(errno.CheckZfAccountFailedException.AppendErrorMsg(err.Error()))
	}
	if ret, err := robot.GetKcs(); err == nil {
		service.SendResponseSuccess(c, gin.H{
			"scores": ret.Items,
		})
	} else {
		panic(errno.NormalException.AppendErrorMsg(err.Error()))
	}
}

// @Summary 更新学生的信息
// @Produce json
// @Tags Student
// @Security ApiKeyAuth
// @Router /students/update_edu_account [post]
// @Param auth body co.EduAccount true "update edu account info"
// @Success 200 object service.Response
func (s *studentController) UpdateOrCreateEduAccount(c *gin.Context) {
	helper := contextHelper.New(c)
	var (
		number, password string
	)
	number = helper.GetString("student_number")
	password = helper.GetString("password")
	// 检测账号密码是否正确
	robotgg := zcmu.New(number, password)
	if err := robotgg.Login(); err != nil {
		panic(errno.CheckZfAccountFailedException.AppendErrorMsg(err.Error()))
	}
	err, student := (&model.Student{
		UserId:   helper.GetAuthUserId(),
		Password: password,
		Number:   number,
	}).UpdateOrCreate()

	if err != nil {
		panic(errno.NormalException.AppendErrorMsg("创建学生失败: " + err.Error()))
	}
	helper.Response(gin.H{
		"student": student,
	})
}

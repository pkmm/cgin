package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pkmm_gin/model"
	"pkmm_gin/utility/zf"
	"strconv"
)

func GetStudentScores(context *gin.Context) {
	stuId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"message": "Invalid id of student.",
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"scores": model.GetScore(stuId),
		})
	}
}

func GetStudent(c *gin.Context) {
	if stuId, err := strconv.ParseInt(c.Param("id"), 10, 64); err == nil {
		c.JSON(http.StatusOK, gin.H{
			"student": model.GetStudent(stuId),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Invalid id of Student.",
		})
	}
}

func UpdateStudent(c *gin.Context) {
	var stu model.Student
	if err := c.BindJSON(&stu); err == nil {
		model.AddStudent(stu)
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Field [num, pwd] must supply.",
		})
	}
}

func GetScores(c *gin.Context) {
	var stu model.Student
	u := c.DefaultQuery("uid", "")
	uid, err := strconv.ParseInt(u, 10, 64)
	stu = model.GetStudentByUserId(uid)
	if err != nil || stu.Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "not found user.",
		})
		return
	}
	scores := model.GetScore(stu.Id)

	c.JSON(200, gin.H{
		"scores": scores,
	})
}

func UpdateZcmuAccount(c *gin.Context) {
	type tmp struct {
		Id         int64
		Num        string
		Pwd        string
		IsAutoSync bool
	}
	t := tmp{}
	c.BindJSON(&t)
	if t.Id <= 0 {
		c.JSON(200, gin.H{
			"message": "必须登录",
		})
		return
	}
	spider := zf.NewCrawl(t.Num, t.Pwd)
	_, err := spider.LoginScorePage()
	for try := 3; try > 0; try-- {
		if err == nil || err.Error() == zf.LOGIN_ERROR_MSG_WRONG_PASSWORD ||
			err.Error() == zf.LOGIN_ERROR_MSG_NOT_VALID_USER {
			break
		}
	}
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
		})
		return
	}

	stu := model.UpdateStudentZcmuAccount(t.Num, t.Pwd, t.Id)
	c.JSON(200, gin.H{
		"message": "success",
		"student": stu,
	})
}

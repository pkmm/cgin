package controller

import (
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"net/http"
	"pkmm_gin/model"
	"pkmm_gin/utility"
	"strconv"
	"time"
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

func AddStudent(c *gin.Context) {
	var stu model.Student
	if err := c.Bind(&stu); err == nil {
		model.AddStudent(stu)
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Field [num, pwd] must supply. Num must be integer.",
		})
	}
}

func UpdateStudent(c *gin.Context) {
	var stu model.Student
	if err := c.Bind(&stu); err == nil {
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

func Login(c *gin.Context) {
	var stu model.Student
	err := c.Bind(&stu)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
		})
	} else if stu.Num == 0 || stu.Pwd == "" || ! model.ExistStudent(stu.Num, stu.Pwd) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Student not exist or Pwd is error.",
		})
	} else {
		securityKey, id := utility.GenerateSignatureAndId(map[string]string{
			"num": strconv.FormatInt(stu.Num, 10),
			"pwd": stu.Pwd,
		})
		myRedis := model.GetRedis()
		defer myRedis.Close()

		// 删除历史的数据
		lastId, _ := redis.String(myRedis.Do("GET", stu.Num))
		myRedis.Do("DEL", stu.Num)
		myRedis.Do("DEL", lastId)

		// 最新的数据
		expireAt := time.Now().Add(time.Hour * 24).Unix()
		myRedis.Do("SET", id, securityKey)
		myRedis.Do("SET", stu.Num, id)
		myRedis.Do("EXPIRE", stu.Num, expireAt)
		myRedis.Do("EXPIRE", id, expireAt)

		c.JSON(http.StatusOK, gin.H{
			"security_key": securityKey,
			"code":         id,
		})
	}

}

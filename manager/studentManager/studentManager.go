package studentManager

import (
	"pkmm_gin/model"
)

func GetStudentByOpenId(openId string) (student model.Student) {
	var stu model.Student
	model.GetDB().Where("open_id = ?", openId).First(&stu)
	return stu
}



func UpdateStudentZcmuAccount(num, pwd string, id int64) model.Student {
	model.GetDB().Table("stu").Where("id = ?", id).Updates(map[string]interface{}{
		"num": num,
		"pwd": pwd,
	})
	var stu model.Student
	model.GetDB().Where("id = ?", id).Find(&stu)
	return stu
}

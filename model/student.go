package model

type Student struct {
	Id        int64  `json:"id"`
	Name      string `json:"name" form:"name"`
	Num       int64  `json:"num,string" form:"num" binding:"required"` // tags "json:,string" 表示别的类型也能解析到
	Pwd       string `json:"pwd" form:"pwd" binding:"required"`
	WxID      string `json:"wx_id" form:"wx_id"`
	CreatedAt string `json:"created_at" gorm:"default:null"`
	UpdatedAt string `json:"updated_at" gorm:"default:null"`
}

func (Student) TableName() string {
	return "stu"
}

func GetStudent(id int64) Student {
	var stu Student
	db.Find(&stu, id)
	return stu
}

func AddStudent(stu Student) {
	db.Create(&stu)
}

func UpdateStudent(stu Student) {
	db.Where("id = ?", stu.Id).Updates(stu)
}

func ExistStudent(num int64, pwd string) bool {
	var stu Student
	db.Where("num = ? AND pwd = ?", num, pwd).Select("id").First(&stu)
	return stu.Id != 0
}

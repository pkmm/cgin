package model

import "time"

// 数据库查询使用
type Student struct {
	Id        int64     `json:"-" gorm:"default:null"`
	UserId    int64     `json:"user_id"`
	Num       string    `json:"num" gorm:"default:null"` // tags "json:,string" 表示别的类型也能解析到
	Pwd       string    `json:"pwd" grom:"default:null"` // 教务系统的密码
	Name      string    `json:"name"`
	CanSync   int       `json:"sync_status"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`

	// 关联表
	Scores []Score `gorm:"ForeignKey:StudentId;AssociationForeignKey:Id" json:"scores"`
}

func (Student) TableName() string {
	return "students"
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

func ExistStudent(num, pwd string) bool {
	var stu Student
	db.Where("num = ? AND pwd = ?", num, pwd).Select("id").First(&stu)
	return stu.Id != 0
}

//// 更新或者创建一个学生
//func UpdateUserWeChatInfo(student Student, openId string) Student {
//	db.Where(Student{OpenId: openId}).Assign(student).FirstOrCreate(&student)
//	return student
//}

// 设置学生当天不再同步
func UpdateUserCanSync(id int64, canSync int) error {
	err := db.Model(&Student{}).Where(Student{Id: id}).UpdateColumn("can_sync", canSync).Error
	return err
}

// 设置所有can_sync
func SetCanSync(canSync int) int64 {
	return db.Model(&Student{}).Updates(Student{CanSync: canSync}).RowsAffected
}

// 可以同步的学生的数量
func CanSyncCount() (int64, error) {
	var ans int64
	err := db.Model(&Student{}).Where(Student{CanSync: 1}).Count(&ans).Error
	return ans, err
}

func GetStudentByStudentNumber(num string) Student {
	var stu Student
	db.Where(Student{Num: num}).Find(&stu)
	return stu
}

func GetStudentByOpenId(openId string) Student {
	var stu Student
	db.Table("students").Joins("JOIN wechat_users AS WU ON WU.user_id = students.id").
		Where("WU.open_id = ?", openId).First(&stu)

	return stu
}

func GetStudentByUserId(userId int64) Student {
	var stu Student
	db.Model(&Student{}).Where(Student{UserId: userId}).First(&stu)
	return stu
}

// update or create.
func UpdateStudentZcmuAccount(num, pwd string, userId int64) Student {
	stu := Student{
		Num:    num,
		Pwd:    pwd,
		UserId: userId,
		CanSync: 1,
	}
	db.Where(Student{UserId: userId}).Assign(stu).FirstOrCreate(&stu)
	return stu
}

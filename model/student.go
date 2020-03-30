package model

import (
	"cgin/conf"
	"cgin/model/modelInterface"
)

type Student struct {
	Id       uint64 `json:"id" gorm:"primary_key;auto_increment;"`
	UserId   uint64 `json:"user_id" gorm:"index:idx_user_id;not null;"`
	Name     string `json:"name" gorm:"type:varchar(64);default:null;"`
	Number   string `json:"number" gorm:"default:null;type:varchar(64);"`
	Password string `json:"password" gorm:"default:null;type:varchar(64);"`
	IsSync   bool   `json:"is_sync" gorm:"type:bool"`
	Model

	Scores     []Score    `gorm:"ForeignKey:StudentId;AssociationForeignKey:Id" json:"scores"`
	SyncDetail *SyncDetail `json:"sync_detail" gorm:"ForeignKey:StudentId;AssociationForeignKey:Id"`
}

func (s *Student) UpdateOrCreate() (err error, _s *Student) {
	err = conf.DB.Where(Student{UserId: s.UserId}).Assign(*s).FirstOrCreate(s).Error
	return err, s
}

func (s *Student) GetList(info modelInterface.PageSizeInfo) (error, interface{}, int) {
	err, query, total := basicPagination(info, s)
	if err != nil {
		return err, nil, 0
	} else {
		var list []*Student
		err = query.Find(&list).Error
		return err, list, total
	}
}

func (s *Student) GetStudentsNeedSyncScore(page, size int) (err error, _s *[]*Student, t int) {
	err, query, total := basicPagination(modelInterface.PageSizeInfo{Page: page, PageSize: size}, s)
	if err != nil {
		return err, nil, 0
	} else {
		var result []*Student
		err = query.Where("is_sync = ?", true).Find(&result).Error
		return err, &result, total
	}
}
func ResetStudentSyncScoreStatus() error {
	err := conf.DB.Model(&Student{}).Update("is_sync", true).Error
	return err
}
func GetStudentByUserId(userId uint64) (err error, _s *Student) {
	var student Student
	err = conf.DB.Where("user_id = ?", userId).
		Preload("Scores").
		Preload("SyncDetail").
		First(&student).Error
	return err, &student
}

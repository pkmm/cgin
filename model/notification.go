package model

import (
	"cgin/conf"
	"cgin/model/modelInterface"
	"cgin/util"
)

type Notification struct {
	Id       uint64        `json:"id" gorm:"primary_key;auto_increment;"`
	Content  string        `json:"content" gorm:"type:varchar(255);default:null"`
	Disabled bool          `json:"disabled" gorm:"default:false;"`
	StartAt  util.JSONTime `json:"start_at" gorm:"type:timestamp;"`
	EndAt    util.JSONTime `json:"end_at" gorm:"type:timestamp;"`
	Model
}

func (n *Notification) CreateNotification() (err error, notification *Notification) {
	err = conf.DB.Create(n).Error
	return err, n
}

func (n *Notification) GetLatest() (error, *Notification) {
	var out Notification
	err := conf.DB.Last(&out).Error
	return err, &out
}

func (n *Notification) UpdateNotification(nid uint64) (err error, _n *Notification) {
	var result Notification
	err = conf.DB.Where(Notification{Id: nid}).
		Assign(*n).FirstOrCreate(&result).Error
	return err, &result
}

func (n *Notification) GetList(info modelInterface.PageSizeInfo) (error, interface{}, int) {
	err, gq, total := basicPagination(info, n)
	if err != nil {
		return err, nil, 0
	} else {
		var data []*Notification
		err = gq.Find(&data).Error
		return err, data, total
	}
}

func GetNotificationById(nId uint64) (err error, _n *Notification) {
	var result Notification
	err = conf.DB.Find(&result, Notification{Id: nId}).Error
	return err, &result
}
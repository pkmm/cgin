package model

import (
	"cgin/conf"
	"github.com/jinzhu/gorm"
)

// all app models.
// register model here.
var Models = []interface{}{
	&User{},
	&Score{},
	&SyncDetail{},
	&Student{},
	&Menu{},
	&IndexConfig{},
	&Notification{},
	&HermannMemorial{},
	&Sponsor{},
	&Tieba{},
	&Thinking{},
}

func init() {
	err := registerTables(conf.DB, Models)
	if err != nil {
		panic("[error]: 注册mysql，table失败: " + err.Error())
	}
	conf.Logger.Info("注册数据库表成功")
}
func registerTables(db *gorm.DB, models []interface{}) (err error) {
	err = db.
		Set("gorm:table_options", "ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=1;").
		AutoMigrate(models...).Error
	return
}

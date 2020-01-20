package task

import (
	"cgin/conf"
	"cgin/model"
	"cgin/service"
	"fmt"
	"github.com/pkmm/gb/baidu"
)

// 签到百度贴吧
func SignBaiduForums() {
	budsss := service.TiebaService.GetAll()
	if len(budsss) == 0 {
		budsss = append(budsss, &model.Tieba{Bduss: conf.AppConfig.String("bduss")})
	}
	for _, record := range budsss {
		if record.Bduss == "" {
			continue
		}
		worker, err := baidu.NewBaiduTiebaCrawl(record.Bduss)
		if err != nil {
			// TODO 邮件通知 或者 记录到DB中
			conf.AppLogger.Error("初始化百度签到worker失败: ", err.Error())
			return
		}
		tiebas, err := worker.RetrieveTiebas()
		if err != nil {
			conf.AppLogger.Error("获取贴吧失败： ", err.Error())
			return
		}
		// TODO: 处理签到的结果
		ret := worker.SignAll(tiebas)
		for k, v := range *ret {
			fmt.Printf("%#v, %#v", k, v)
		}
		//fmt.Printf("%#v", ret)
	}

}

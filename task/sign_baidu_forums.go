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
			conf.Logger.Error("初始化百度签到worker失败: ", err.Error())
			return
		}
		tiebas, err := worker.RetrieveTiebas()
		if err != nil {
			conf.Logger.Error("获取贴吧失败： ", err.Error())
			return
		}

		//  ====  使用pool的版本 =====
		// ==========================
		ts := make([]*Task, len(tiebas))
		for i, t := range tiebas {
			y := t
			ts[i] = NewTask(func() error {
				resp := worker.SignOne(y)
				fmt.Printf("Sign [%s] result is %s\n", y, resp)
				return nil
			})
		}
		pool.AddTasks(ts)
		//  ====== 使用pool  =====
		//  =====================

		// 之前的版本 ======================
		//// TODO: 处理签到的结果
		//ret := worker.SignAll(tiebas)
		//for k, v := range *ret {
		//	fmt.Printf("%#v, %#v", k, v)
		//}
		////fmt.Printf("%#v", ret)
		// ================================
	}
}

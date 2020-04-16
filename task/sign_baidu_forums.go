package task

import (
	"cgin/conf"
	"cgin/model"
	"cgin/service/workerpool"
	"fmt"
	"github.com/pkmm/gb/baidu"
)

// 签到百度贴吧
func SignBaiduForums() {
	_, data, _ := new(model.Tieba).GetList(1, 500)
	tiebaUsers := data.([]*model.Tieba)
	if len(tiebaUsers) == 0 {
		tiebaUsers = append(tiebaUsers, &model.Tieba{Bduss: conf.AppConfig.String("bduss")})
	}
	for _, record := range tiebaUsers {
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
		ts := make([]*workerpool.Task, len(tiebas))
		for i, t := range tiebas {
			y := t
			ts[i] = workerpool.NewTask(func() {
				resp := worker.SignOne(y)
				fmt.Printf("Sign [%s] result is %s\n", y, resp)
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

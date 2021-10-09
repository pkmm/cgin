package schedule

import (
	"cgin/global"
	"cgin/service/system"
	"fmt"
	"github.com/pkmm/gb/baidu"
)

// SignBaiduForums 签到百度贴吧
func SignBaiduForums() {
	err, users := system.BaiduServiceApp.GetUsers()
	if err != nil {
		fmt.Printf("加载baidu用户失败 %v", err)
		return
	}
	// TODO：处理所有的用户，当前只有一个人
	bs := users[0].Bduss
	worker, err := baidu.NewBaiduTiebaCrawl(bs)
	if err != nil {
		return
	}
	ties, err := worker.RetrieveTiebas()
	if err != nil {
		return
	}
	for i, t := range ties {
		y := t
		err = global.WorkerPool.Submit(func() {
			resp := worker.SignOne(y)
			fmt.Printf("Sign[%d], name is [%s], result is %v\n", i, y, resp)
		})
		if err != nil {
			fmt.Printf("Submit task to pool failed: %v", err)
		}
	}
}

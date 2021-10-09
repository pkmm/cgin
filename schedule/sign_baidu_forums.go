package schedule

import (
	"cgin/global"
	"fmt"
	"github.com/pkmm/gb/baidu"
)

// SignBaiduForums 签到百度贴吧
func SignBaiduForums() {
	worker, err := baidu.NewBaiduTiebaCrawl("")
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

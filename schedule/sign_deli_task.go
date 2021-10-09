package schedule

import (
	"cgin/global"
	"cgin/service/system"
	"math/rand"
	"time"
)

func SignDeli() {

	// 关闭了这个功能
	if global.Config.Deli.Stop {
		return
	}

	if err, users := system.DeliAutoSignApp.GetAllUsers(); err != nil {
		for _, user := range users {
			_user := user
			_ = global.WorkerPool.Submit(func() {
				rnd := rand.Intn(10) + 10 // [10, 20) 分钟
				time.Sleep(time.Minute * time.Duration(rnd))
				if err, _ := system.DeliAutoSignApp.SignOne(&_user); err != nil {
					// notify user of sign result.
				} else {
					// send html as result.
				}
			})
		}
	}
}

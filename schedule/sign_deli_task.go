package schedule

import (
	"cgin/global"
	"cgin/service/system"
	"cgin/service/wechatpush"
	"go.uber.org/zap"
)

func SignDeli() {

	// 关闭了这个功能
	if global.Config.Deli.Stop {
		return
	}

	if err, users := system.DeliAutoSignApp.GetAllUsers(); err == nil {
		for _, user := range users {
			global.GLog.Info("deli user", zap.Any("user", user))
			_user := user
			_ = global.WorkerPool.Submit(func() {
				//rnd := rand.Intn(10) + 10 // [10, 20) 分钟
				//time.Sleep(time.Minute * time.Duration(rnd))
				if err, _ := system.DeliAutoSignApp.SignOne(&_user); err != nil {
					// notify user of sign result.
					global.GLog.Error("签到失败！", zap.Any("error", err))
				} else {
					// send html as result.
					global.GLog.Info("签到成功！", zap.Any("username", _user.Username))
					bear := wechatpush.NewPushBear([]int{166}, wechatpush.Html)
					_, err2 := bear.Send("签到成功", _user.Username)
					if err2 != nil {
						global.GLog.Error("wxpusher 发送签到信息失败", zap.Any("error", err2))
					}
				}
			})
		}
	} else {
		global.GLog.Error("获取deli用户失败！", zap.Any("error", err))
	}
}

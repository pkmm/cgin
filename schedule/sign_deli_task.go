package schedule

import (
	"cgin/global"
	"cgin/service/system"
	"github.com/wxpusher/wxpusher-sdk-go"
	"github.com/wxpusher/wxpusher-sdk-go/model"
	"go.uber.org/zap"
	"math/rand"
	"time"
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
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				rnd := r.Intn(10) + 10 // [10, 20) 分钟
				time.Sleep(time.Minute * time.Duration(rnd))
				if err, html := system.DeliAutoSignApp.SignOne(&_user); err != nil {
					// notify user of sign result.
					global.GLog.Error("签到失败！", zap.Any("error", err))
				} else {
					// send html as result.
					global.GLog.Info("签到成功！", zap.Any("username", _user.Username))
					msg := model.NewMessage(global.Config.Wxpusher.AppToken).
						SetSummary("任务状态: " + _user.Username).
						SetContentType(2).SetContent(html).AddUId(_user.Uid)
					_, err2 := wxpusher.SendMessage(msg)
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

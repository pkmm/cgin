package schedule

import (
	"cgin/global"
	"cgin/service/system"
	"cgin/util"
	"github.com/wxpusher/wxpusher-sdk-go"
	"github.com/wxpusher/wxpusher-sdk-go/model"
	"go.uber.org/zap"
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

				x := util.GetInt64()
				x %= 25
				if x < 10 {
					x += 10
				}

				if x < 10 || x > 25 {
					x = 18
				}

				global.GLog.Info("用户签到deli休眠的时间是：", zap.Any("username", _user.Username), zap.Any("time sleep", x))

				// 时间的范围在[10, 25)之间
				time.Sleep(time.Minute * time.Duration(x))

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

package schedule

// 后台运行的任务
import (
	"cgin/global"
	"cgin/util"
	"fmt"
	"github.com/robfig/cron"
)

const (
	FlagFirst = iota << 1
	FlagBaiduTBSign
	FlagDeli
)

var SC *Schedule

// Schedule spec 秒 分 时 日 月 星期
type Schedule struct {
	slots *util.SafeMap // 分配每一个goroutine 一个id, 避免任务的重叠运行
	sc    *cron.Cron
}

func NewSchedule() *Schedule {
	return &Schedule{
		slots: util.NewSafeMap(),
		sc:    cron.New(),
	}
}

func (s *Schedule) taskWrapper(cmd func(), flag int32) func() {
	return func() {
		if _, ok := s.slots.ReadSafeMap(flag); ok {
			return
		}
		s.slots.WriteSafeMap(flag, 1)
		cmd()
		s.slots.DeleteKey(flag)
	}
}

func (s *Schedule) AddFunc(spec string, fn func(), flag int32) error {
	return s.sc.AddFunc(spec, s.taskWrapper(fn, flag))
}

func (s *Schedule) Stop() {
	s.sc.Stop()
}

// StartJobs 启动配置
// 需要作为定位的任务，只需要添加在这个文件中即可，
// TODO：添加新的任务就需要重新编译启动服务，能否做成动态的
func (s *Schedule) StartJobs() {

	// 百度贴吧签到
	s.AddFunc("0 0 0 * * *", SignBaiduForums, FlagBaiduTBSign)

	// am
	s.AddFunc("0 2 8,12 * * 1-6", SignDeli, FlagDeli)

	// pm
	if global.Config.Deli.Season == "winter" {
		// winter
		fmt.Println("winter")
		s.AddFunc("0 32 13,17 * * 1-6", SignDeli, FlagDeli)
	} else {
		// summer
		fmt.Println("summer")
		s.AddFunc("0 2 14,18 * * 1-6", SignDeli, FlagDeli)
	}

	/// 在此函数上面进行任务的配置
	s.sc.Start()
}

func (s *Schedule) Reload() {
	fmt.Println("schedule已经被重新加载那")
	s.Stop()
	s.StartJobs()
}

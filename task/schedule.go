package task

// 后台运行的任务
import (
	"fmt"
	"github.com/robfig/cron"
	"pkmm_gin/model"
	"time"
)

const (
	TEST = "test"
)

func init() {
	// 分配每一个goroutine 一个id, 避免任务的重叠运行
	runningTask := make(map[string]struct{})

	c := cron.New()
	// 定义任务列表

	// 移除过期的 nonce
	c.AddFunc("*/60 * * * * *", func() {
		ts := time.Now().Add(-1 * time.Minute).Unix()
		myRedis := model.GetRedis()
		defer myRedis.Close()
		myRedis.Do("ZREMRANGEBYRANK", "api_collection_ttl", 0, ts)
		fmt.Print("clear redis item.")
	})

	var index = 1
	c.AddFunc("@every 2s", func() {
		// avoid overlapping execute job.
		if _, ok := runningTask[TEST]; ok {
			return
		}
		runningTask[TEST] = struct{}{}
		//fmt.Println(index)
		index++
		// todo task
		time.Sleep(5 * time.Second)
		index--
		delete(runningTask, TEST)
	})

	//  开始任务
	c.Start()
}

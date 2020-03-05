package task

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
	"testing"
)

func toSleep(i int) error {
	fmt.Printf("sleep %d\n", i)
	gorequest.New().Get("http://www.baidu.com").End()
	return nil
}

func TestSimplePool_RunPool(t *testing.T) {
	tasks := make([]*Task, 10)
	for i := 0; i < 10; i++ {
		tasks[i] = NewTask(func() error {
			return toSleep(2)
		})
	}

	pool := NewSimplePool(2)
	pool.RunPool()
	pool.AddTasks(tasks)
	//pool.Stop()

	select {

	}

}
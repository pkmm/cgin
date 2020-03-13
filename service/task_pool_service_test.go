package service

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var wg sync.WaitGroup
func toSleep(i int) error {
	time.Sleep(time.Duration(i) * time.Second)
	wg.Done()
	return nil
}

func TestSimplePool_RunPool(t *testing.T) {
	var count int = 5
	tasks := make([]*Task, count)
	for i := 0; i < count; i++ {
		wg.Add(1)
		tasks[i] = NewTask(func() error {
			return toSleep(2)
		})
	}

	at := time.Now()
	pool := NewSimplePool(2)
	pool.RunPool()
	pool.AddTasks(tasks)
	wg.Wait()
	fmt.Println("pool run finished")
	pool.Stop()
	pool.Stop()
	fmt.Println("time use: ", time.Since(at).String())
}
package workerpool

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var wg sync.WaitGroup

func toSleep(i int) {
	time.Sleep(time.Duration(i) * time.Second)
	wg.Done()
}

func TestSimplePool_RunPool(t *testing.T) {
	var count int = 5
	tasks := make([]*Task, count)
	wg.Add(count)
	for i := 0; i < count; i++ {
		tasks[i] = NewTask(func() {
			toSleep(2)
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
func TestWorkerPool_Start(t *testing.T) {
	var tt sync.WaitGroup
	count := 200
	wp := &workerPool{MaxWorkersCount: 50, lock: &sync.Mutex{}}
	wp.cond = sync.NewCond(wp.lock)
	wp.Start()
	tt.Add(count)
	for i := 0; i < count; i++ {
		y := i
		wp.Serve(func() {
			time.Sleep(1 * time.Second)
			fmt.Printf("task %d\n", y)
			tt.Done()
		})
	}
	//time.Sleep(2 * time.Second)
	//tt.Wait()
	//wp.Stop()
	for {
		select {
		case <-time.Tick(time.Second * 2):
			fmt.Printf("worker chan count %d\n", wp.workersCount)
		}
	}
	//fmt.Printf("stopped\n")
}

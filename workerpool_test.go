package main

import (
	"cgin/service/workerpool"
	"fmt"
	"math/rand"
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
	t.Parallel()

	var count int = 5
	tasks := make([]*workerpool.Task, count)
	wg.Add(count)
	for i := 0; i < count; i++ {
		tasks[i] = workerpool.NewTask(func() {
			toSleep(2)
		})
	}

	at := time.Now()
	pool := workerpool.NewSimplePool(2)
	pool.RunPool()
	pool.AddTasks(tasks)
	wg.Wait()
	fmt.Println("pool run finished")
	pool.Stop()
	pool.Stop()
	fmt.Println("time use: ", time.Since(at).String())
}
func TestWorkerPool_Start(t *testing.T) {
	t.Parallel()
	var tt sync.WaitGroup
	jobCount := 86
	wp := workerpool.NewWorkerPool(17)

	for i := 0; i < 10; i++ {
		wp.Start()
		wp.Stop()
	}

	wp.Start()

	go func() {
		for {
			select {
			case <-time.Tick(time.Second * 1):
				fmt.Printf("=====>!!! worker chan jobCount %d\n", wp.GetActiveWorkerCount())
			}
		}
	}()

	tt.Add(jobCount)
	fmt.Printf("begin\n")
	go func() {
		for i := 0; i < jobCount; i++ {
			y := i
			wp.Execute(func() {
				rd := rand.Intn(8) + 1
				if y > 40 {
					rd = 1
				}
				time.Sleep(time.Duration(rd) * time.Second)
				if rd == 6 {
					tt.Done()
					panic("panic in user func.")
				}
				fmt.Printf("task_id %d, 休眠的的时间: %d\n", y+1, rd)
				tt.Done()
			})
		}
	}()
	//time.Sleep(10 * time.Second)

	//time.Sleep(2 * time.Second)
	tt.Wait()
	wp.Stop()
	fmt.Printf("over\n")
	//select{
	//case <-time.After(60 * time.Second):
	//	break
	//
	//}
	//fmt.Printf("stopped\n")
}

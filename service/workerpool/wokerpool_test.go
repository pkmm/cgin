package workerpool

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestWorkerPool(t *testing.T) {
	wg := sync.WaitGroup{}
	pool, _ := NewPool(20, time.Second*2)
	defer pool.Close()
	for i := 0; i < 100; i++ {
		wg.Add(1)
		pool.Submit(sendMail(i, &wg))
	}
	wg.Wait()
}

func sendMail(i int, wg *sync.WaitGroup) func() {
	return func() {
		time.Sleep(time.Second * 5)
		fmt.Println("send mail to ", i)
		wg.Done()
	}
}

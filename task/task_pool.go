package task

import (
	"cgin/conf"
)

// 可靠性未知

// 用于线程池执行的任务task
type Task struct {
	f func() error
}

func NewTask(fn func() error) *Task {
	return &Task{f: fn}
}

func (t *Task) Execute() error {
	return t.f()
}

// 线程池
type SimplePool struct {
	// 协程的数量
	Count int
	// 任务队列
	JobQueue chan *Task
	Stopped  chan interface{}
}

func NewSimplePool(cap int) *SimplePool {
	return &SimplePool{Count: cap, JobQueue: make(chan *Task, cap), Stopped: make(chan interface{})}
}

func (s *SimplePool) AddTasks(ts []*Task) {
	go func() {
		for _, t := range ts {
			select {
			case <-s.Stopped:
				return
			default:
				s.JobQueue <- t
			}
		}
	}()
}

func (s *SimplePool) worker(workerId int) {
	for t := range s.JobQueue {
		if err := t.Execute(); err != nil {
			conf.Logger.Debug("task execute err: %s", err.Error())
		} else {
			conf.Logger.Debug("worker %d finished", workerId)
		}
	}
	conf.Logger.Debug("worker %d stopped", workerId)
}

func (s *SimplePool) RunPool() {
	for i := 0; i < s.Count; i++ {
		go s.worker(i)
	}
}

// 调用stop之后不能再调用addTask..否则panic
func (s *SimplePool) Stop() {
	s.Stopped <- ""
	close(s.Stopped)
	close(s.JobQueue)
}

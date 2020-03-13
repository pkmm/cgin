package service

import (
	"cgin/conf"
	"sync"
)

var once sync.Once

var TaskPool *SimplePool

func init() {
	TaskPool = NewSimplePool(30)
}

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
	count int
	// 任务队列
	jobQueue chan *Task
	stopped  chan interface{}
}

func NewSimplePool(cap int) *SimplePool {
	return &SimplePool{
		count: cap,
		jobQueue: make(chan *Task, cap<<1), // 2倍的任务队列
		stopped: make(chan interface{}, 1),
	}
}

// 添加一个任务 阻塞
func (s *SimplePool) AddTask(t *Task) {
	select {
	case <-s.stopped:
		return
	default:
		s.jobQueue <- t
	}
}

// 启动一个goroutine非阻塞
func (s *SimplePool) AddTasks(ts []*Task) {
	go func() {
		for _, t := range ts {
			select {
			case <-s.stopped:
				return
			default:
				s.jobQueue <- t
			}
		}
	}()
}

func (s *SimplePool) worker(workerId int) {
	defer func() {
		if e := recover(); e != nil {
			conf.Logger.Error("worker [%d] of pool panic, msg[%#v]", workerId, e)
		}
	}()
	for t := range s.jobQueue {
		if err := t.Execute(); err != nil {
			conf.Logger.Debug("task execute err: %s", err.Error())
		}
	}
	conf.Logger.Debug("worker %d stopped", workerId)
}

func (s *SimplePool) RunPool() {
	for i := 0; i < s.count; i++ {
		go s.worker(i)
	}
}

// 调用stop之后不能再调用addTask..否则panic
func (s *SimplePool) Stop() {
	once.Do(func() {
		s.stopped <- ""
		close(s.stopped)
		close(s.jobQueue)
	})
}

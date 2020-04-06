package workerpool

import (
	"cgin/conf"
	"fmt"
	"runtime"
	"sync"
	"time"
)

var once sync.Once

var TaskPool *SimplePool

func init() {
	TaskPool = NewSimplePool(30)
}

// 可靠性未知

// 用于线程池执行的任务task
type Task struct {
	f func()
}

func NewTask(fn func()) *Task {
	return &Task{f: fn}
}

func (t *Task) Execute() {
	t.f()
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
		count:    cap,
		jobQueue: make(chan *Task, cap<<1), // 2倍的任务队列
		stopped:  make(chan interface{}, 1),
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
		t.Execute()
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

//////////////////////////////////////////////
//// new worker pool inspired from fast http//
//////////////////////////////////////////////
type workerPool struct {
	MaxWorkersCount        int
	MaxIdleWorkersCount    int
	MaxIdleWorkersDuration time.Duration
	workersCount           int
	lock                   sync.Locker
	mustStop               bool

	ready []*workerChan

	stopChan chan struct{}

	workerChanPool sync.Pool

	cond *sync.Cond
}

type workerChan struct {
	lastUseTime time.Time
	ch          chan func()
}

func (wp *workerPool) Start() {
	if wp.stopChan != nil {
		panic("BUG: worker pool already started")
	}
	wp.stopChan = make(chan struct{})
	stopChan := wp.stopChan
	wp.workerChanPool.New = func() interface{} {
		return &workerChan{ch: make(chan func(), workerChanCap)}
	}
	// 单独一个协成清理过期的
	go func() {
		var scratch []*workerChan
		for {
			wp.clean(&scratch)
			select {
			case <-stopChan:
				return
			default:
				time.Sleep(wp.getMaxIdleWorkersDuration())
			}
		}
	}()
}

func (wp *workerPool) Stop() {
	if wp.stopChan == nil {
		panic("BUG: worker pool already stopped")
	}
	close(wp.stopChan)
	wp.stopChan = nil
	wp.lock.Lock()
	ready := wp.ready
	for i := range ready {
		ready[i].ch <- nil
		ready[i] = nil
	}
	wp.ready = ready[:0]
	wp.mustStop = true
	wp.cond.Broadcast()
	wp.lock.Unlock()
}

func (wp *workerPool) Serve(f func()) bool {
	ch := wp.getCh()
	if ch == nil {
		return false
	}
	ch.ch <- f
	return true
}

func (wp *workerPool) getCh() *workerChan {
	var ch *workerChan
	createWorker := false

	wp.lock.Lock()
Reentry:
	ready := wp.ready
	n := len(ready) - 1
	if n < 0 {
		if wp.workersCount < wp.MaxWorkersCount {
			createWorker = true
			wp.workersCount++
		} else {
			// 阻塞等待
			wp.cond.Wait()
			goto Reentry
		}
	} else {
		ch = ready[n]
		ready[n] = nil
		wp.ready = ready[:n]
	}
	wp.lock.Unlock()

	if ch == nil {
		if !createWorker {
			fmt.Println("已经超过负载,等待可用协程")
			return nil // todo 满负荷应该怎么做
		}
		vch := wp.workerChanPool.Get()
		ch = vch.(*workerChan)
		go func() {
			wp.workerFunc(ch)
			// run to wait task
			// do something

			// revert
			wp.workerChanPool.Put(ch)
		}()
	}

	return ch
}

func (wp *workerPool) workerFunc(ch *workerChan) {
	// todo add defer handle panic
	var c func()
	for c = range ch.ch {
		if c == nil {
			break
		}
		// todo
		c() // execute func
		if ok := wp.release(ch); !ok {
			break
		}
	}
	wp.lock.Lock()
	wp.workersCount--
	wp.lock.Unlock()
}

func (wp *workerPool) release(ch *workerChan) bool {
	fmt.Printf("call release\n")
	ch.lastUseTime = time.Now()
	wp.lock.Lock()
	if wp.mustStop {
		wp.lock.Unlock()
		return false
	}
	wp.ready = append(wp.ready, ch)
	wp.cond.Signal()
	wp.lock.Unlock()
	return true
}

func (wp *workerPool) clean(scratch *[]*workerChan) {
	maxIdleWorkersDuration := wp.getMaxIdleWorkersDuration()
	criticalTime := time.Now().Add(-maxIdleWorkersDuration)
	wp.lock.Lock()

	ready := wp.ready
	n := len(ready)
	l, r, mid := 0, n-1, 0
	for l <= r {
		mid = l + (r-l)/2
		if criticalTime.After(ready[mid].lastUseTime) {
			l = mid + 1
		} else {
			r = mid - 1
		}
	}
	i := r
	if i == -1 { // 是空的
		wp.lock.Unlock()
		return
	}
	*scratch = append((*scratch)[:0], ready[:i+1]...) // 需要删除的
	m := copy(ready, ready[i+1:])
	for i := m; i < n; i++ {
		ready[i] = nil
	}
	wp.ready = ready[:m]
	wp.lock.Unlock()

	tmp := *scratch // 关闭空闲超时的
	for i := range tmp {
		tmp[i].ch <- nil
		tmp[i] = nil
	}
}

func (wp *workerPool) getMaxIdleWorkersDuration() time.Duration {
	if wp.MaxIdleWorkersDuration <= 0 {
		return 10 * time.Second
	}
	return wp.MaxIdleWorkersDuration
}

var workerChanCap = func() int {
	if runtime.GOMAXPROCS(0) == 1 {
		return 0
	} else {
		return 1
	}
}()

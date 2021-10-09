package workerpool

import (
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type Pool struct {
	capacity       int32             // 协程池的大小
	running        int32             // 正在运行的协程的数量
	expireDuration time.Duration     // 空闲协程的最大时间
	workers        []*Worker         // 空闲的协程
	stopped        int32             // 是不是被关闭了
	lock           sync.Locker       // 同步锁保护
	cond           *sync.Cond        // 条件变量
	once           sync.Once         // 保护pool只会被关闭一次
	workerCache    sync.Pool         // 缓存worker，可能是空
	PanicHandler   func(interface{}) // hook处理pool异常
}

var (
	workerChanCap = func() int {
		// Use blocking channel if GOMAXPROCS=1.
		// This switches context from sender to receiver immediately,
		// which results in higher performance (under go1.5 at least).
		if runtime.GOMAXPROCS(0) == 1 {
			return 0
		}

		// Use non-blocking workerChan if GOMAXPROCS>1,
		// since otherwise the sender might be dragged down if the receiver is CPU-bound.
		return 1
	}()
)

func NewPool(size int32, expireDuration time.Duration) (p *Pool, err error) {
	p = &Pool{
		capacity:       size,
		running:        0,
		expireDuration: expireDuration,
		stopped:        0,
		workers:        []*Worker{},
		lock:           &sync.Mutex{},
	}
	p.workerCache.New = func() interface{} {
		return &Worker{
			pool: p,
			task: make(chan func(), workerChanCap),
		}
	}
	p.cond = sync.NewCond(p.lock)
	go p.purge() // 启动清理协程
	return p, nil
}

func (p *Pool) incRunning() {
	atomic.AddInt32(&p.running, 1)
}

func (p *Pool) decRunning() {
	atomic.AddInt32(&p.running, -1)
}

func (p *Pool) Running() int {
	return int(atomic.LoadInt32(&p.running))
}

func (p *Pool) Cap() int {
	return int(atomic.LoadInt32(&p.capacity))
}
func (p *Pool) IsClosed() bool {
	return int(atomic.LoadInt32(&p.stopped)) == 1
}

// 从pool中提取一个worker
func (p *Pool) retrieveWorker() (w *Worker) {
	spawnWorker := func() {
		w = p.workerCache.Get().(*Worker)
		w.run()
	}
	p.lock.Lock()
	freeWorkers := p.workers
	n := len(freeWorkers) - 1
	if n >= 0 {
		w = freeWorkers[n]
		freeWorkers[n] = nil
		p.workers = freeWorkers[:n]
		p.lock.Unlock()
	} else if p.Cap() > p.Running() { // 目前没有可用的worker，同时还没有到达上限
		p.lock.Unlock()
		spawnWorker()
	} else { // 当前没有可用的worker，而且worker数量已经到达上限，那么就需要进行阻塞等待【默认】。
		for {
			p.cond.Wait()
			var nw int
			// 没有在运行的协程,可能是后台purge线程清理完了所有的线程，这个时候，如果没有关闭pool需要新建协程
			if nw = p.Running(); nw == 0 {
				p.lock.Unlock()
				if !p.IsClosed() {
					spawnWorker()
				}
				return
			}
			// 还存在运行中的协程

			l := len(p.workers) - 1
			if l < 0 {
				// 有在运行中的协程，但是空闲的协是没有的
				if nw < p.Cap() {
					p.lock.Unlock()
					spawnWorker()
					return
				}
				continue
			}
			// 存在空闲的协程，直接使用
			w = p.workers[l]
			p.workers[l] = nil
			p.workers = p.workers[:l]
			break
		}
		p.lock.Unlock()
	}
	return w
}

// 放回去一个worker
func (p *Pool) revertWorker(w *Worker) bool {
	// 满的，或者是已经关闭了
	if c := p.Cap(); (c > 0 && p.Running() > c) || p.IsClosed() {
		return false
	}

	w.recycleTime = time.Now()
	p.lock.Lock()

	if p.IsClosed() {
		p.lock.Unlock()
		return false
	}

	p.workers = append(p.workers, w)
	p.cond.Signal()
	p.lock.Unlock()
	return true
}

// 后台清理空闲的worker线程
func (p *Pool) purge() {
	heartbeat := time.NewTicker(p.expireDuration)
	defer heartbeat.Stop()
	for range heartbeat.C {
		//fmt.Printf("current running %d, idle worker: %d\n", p.running, len(p.workers))
		if p.IsClosed() {
			break
		}

		p.lock.Lock()
		currentTime := time.Now()
		idleWorkers := p.workers
		if len(idleWorkers) == 0 && p.Running() == 0 && p.IsClosed() {
			p.lock.Unlock()
			return
		}
		n := -1
		for i, w := range idleWorkers {
			if currentTime.Sub(w.recycleTime) <= p.expireDuration {
				break
			}
			n = i
			w.task <- nil
			idleWorkers[i] = nil
		}
		// [0, n] 是过期的
		if n > -1 {
			if n >= len(idleWorkers)-1 {
				p.workers = idleWorkers[:0]
			} else {
				p.workers = idleWorkers[n+1:]
			}
		}
		p.lock.Unlock()

		// 没有在运行中的worker，但是还有worker卡在p.cond.wait处
		if p.Running() == 0 {
			p.cond.Broadcast()
		}
	}
}

func (p *Pool) Close() {
	atomic.StoreInt32(&p.stopped, 1)
	p.lock.Lock()
	p.workers = nil
	p.lock.Unlock()
	// 可能还有一些协程阻塞等待在retrieveWorker函数中，因此在这里进行唤醒
	p.cond.Broadcast()
}

func (p *Pool) Submit(task func()) error {
	if p.IsClosed() {
		return errors.New("Pool is closed")
	}
	p.retrieveWorker().task <- task
	return nil
}

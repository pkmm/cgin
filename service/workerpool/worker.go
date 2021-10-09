package workerpool

import (
	"log"
	"time"
)

type Worker struct {
	pool        *Pool       // 隶属于的协程池
	task        chan func() // 分配给自己的任务
	recycleTime time.Time   // 被回收的时间，（最后一次运行结束的时间）
}

func (w *Worker) run() {
	w.pool.incRunning()
	go func() {
		defer func() {
			w.pool.decRunning()
			w.pool.workerCache.Put(w)
			if p := recover(); p != nil {
				if w.pool.PanicHandler != nil {
					w.pool.PanicHandler(p)
				} else {
					log.Printf("worker exits from a panic: %v", p)
				}
			}
			w.pool.cond.Signal()
		}()
		for fn := range w.task {
			if fn == nil {
				return
			}
			fn()
			if ok := w.pool.revertWorker(w); !ok {
				return
			}
		}
	}()
}

package workerpool

func Submit(ts []*Task) {
	TaskPool.AddTasks(ts)
}

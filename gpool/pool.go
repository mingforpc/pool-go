package gpool

// 池, 需要支持新建时只启动部分Coroutine
type Pool struct {
	size int

	taskChan chan Task

	closeChan chan int
}

func (pool *Pool) Serv() {

	pool.taskChan = make(chan Task, pool.size)

	go pool.distribute()

}

func (pool *Pool) distribute() {

	for {

		select {
		case task, ok := <-pool.taskChan:

			if !ok {
				return
			}

			cor := NewCoroutine(task)

			go cor.Start()

		case <-pool.closeChan:
			return
		}

	}

}

func (pool *Pool) Run(runnable Runnable) Task {

	task := NewTask(runnable)

	pool.taskChan <- task

	return task
}

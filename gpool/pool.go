package gpool

// 池, 需要支持新建时只启动部分Coroutine
type Pool struct {
	size int

	taskChan chan *Task

	closed bool

	closeChan chan int
}

func NewPool(size int) Pool {
	return Pool{size: size}
}

// 初始化需要的channel，和开启一个分发Task的goroutine
func (pool *Pool) Serv() {

	pool.taskChan = make(chan *Task, pool.size)
	pool.closeChan = make(chan int)

	go pool.distribute()

}

func (pool *Pool) distribute() {

	for !pool.closed {

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

func (pool *Pool) Run(runnable Runnable) (*Task, error) {

	if pool.closed {
		return nil, ERR_POOL_COSED
	}
	task := NewTask(runnable)

	pool.taskChan <- task

	return task, nil

}

func (pool *Pool) Close() {
	pool.closed = true
	close(pool.closeChan)
	close(pool.taskChan)

}

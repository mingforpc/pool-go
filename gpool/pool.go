package gpool

import "context"

// Pool 池, 需要支持新建时只启动部分Coroutine
type Pool struct {
	size int

	taskChan chan *Task

	closed bool

	closeChan chan int
}

// NewPool 新建一个Pool
func NewPool(size int) Pool {
	return Pool{size: size}
}

// Serv 初始化需要的channel，和开启一个分发Task的goroutine
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

			cor := NewGoroutine(task)

			go cor.Start()

		case <-pool.closeChan:
			return
		}

	}

}

// Run run a runnable
func (pool *Pool) Run(runnable Runnable, args ...interface{}) (*Task, error) {

	if pool.closed {
		return nil, ERR_POOL_COSED
	}
	task := NewTask(runnable, args)

	pool.taskChan <- task

	return task, nil

}

// RunFunc run function
func (pool *Pool) RunFunc(runFun func(ctx context.Context, args ...interface{}) (interface{}, error), args ...interface{}) (*Task, error) {

	if pool.closed {
		return nil, ERR_POOL_COSED
	}

	runnable := &funcRunnable{runFun: runFun}
	task := NewTask(runnable, args)

	pool.taskChan <- task

	return task, nil

}

// Close 关闭所有
func (pool *Pool) Close() {
	pool.closed = true
	close(pool.closeChan)
	close(pool.taskChan)

}

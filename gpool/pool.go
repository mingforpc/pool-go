package gpool

// 池
type Pool struct {
	maxSize int
}

func NewPool(maxSize int) *Pool {

	pool := &Pool{maxSize: maxSize}

	return pool
}

// 任务
type Task interface {
	Run()
}

// 协程
type Coroutine struct {
	id int

	// 分配得到的任务
	task *Task
}

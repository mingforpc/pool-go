package gpool

import (
	"sync"
)

const working = 1
const waiting = 0

// 池, 需要支持新建时只启动部分Coroutine
type Pool struct {
	initSize int // 初始化时的Coroutine数量
	maxSize  int // 最大Coroutine数量

	taskBenc chan Task // 任务分配，无缓存channel

	workingList chan int // 用来统计正在执行task的Coroutine数量
	waitingList chan int // 用来统计正在等待task的Coroutine数量

	coros []*Coroutine

	exitAllChan chan int // 通知所有oroutine退出的channel

	lock sync.Mutex

	once       sync.Once        // 用来执行关闭操作的
	closedSign chan interface{} // 用来判断池是否已经关闭
}

// 新建池
func NewPool(initSize, maxSize int) *Pool {

	pool := &Pool{initSize: initSize, maxSize: maxSize}

	return pool
}

// 初始化Coroutine等操作
func (pool *Pool) Serv() {

	pool.taskBenc = make(chan Task)
	pool.workingList = make(chan int, pool.maxSize)
	pool.waitingList = make(chan int, pool.maxSize)
	pool.exitAllChan = make(chan int)

	pool.closedSign = make(chan interface{}, 1)

	for i := 0; i < pool.initSize; i++ {
		pool.addCoro()
	}

}

// 往池中添加Coroutine
func (pool *Pool) addCoro() error {

	pool.lock.Lock()
	defer pool.lock.Unlock()

	if pool.IsClosed() {
		return CLOSE_ERR
	}

	coro := newCoroutine(pool)

	if pool.GetCorosCount() < pool.maxSize {

		pool.coros = append(pool.coros, coro)
		pool.waitingList <- waiting
		coro.start()
	}

	return nil
}

// 将任务分发给Coroutine
func (pool *Pool) Run(task Task) error {

	if pool.IsClosed() {
		return CLOSE_ERR
	}

	if len(pool.waitingList) == 0 && pool.GetCorosCount() < pool.maxSize {
		// 如果已经没有空闲的Coroutine，而且池中创建的Coroutine数量没有超过最大限制
		// 则创建一个新的Coroutine加入到池
		pool.addCoro()
	}

	pool.taskBenc <- task

	return nil
}

// 关闭池，由于goroutine的原因，无法强制正在执行的Task，只有等task完成，goroutine才会退出
func (pool *Pool) Close() {
	pool.once.Do(func() {
		pool.lock.Lock()
		defer pool.lock.Unlock()

		pool.closedSign <- nil

		close(pool.exitAllChan)

		idleCoros := pool.coros

		for i, coro := range idleCoros {
			coro.setRunning(false)
			idleCoros[i] = nil
		}
		pool.coros = nil

		close(pool.waitingList)
		close(pool.workingList)
	})
}

// 判断池是否已经关闭
func (pool *Pool) IsClosed() bool {

	if len(pool.closedSign) > 0 {
		return true
	} else {
		return false
	}

}

// 获取当前已经启动了的Coroutine数量(正在执行任务 + 等待任务)
func (pool *Pool) GetCorosCount() int {
	count := len(pool.waitingList) + len(pool.workingList)
	return count
}

// waitinglist 加1
func (pool *Pool) waitingIncr() {

	if len(pool.waitingList) < pool.maxSize {
		pool.waitingList <- waiting
	}

}

// waitinglist 减1
func (pool *Pool) waitingDecr() {

	if len(pool.waitingList) > 0 {
		<-pool.waitingList
	}

}

// workinglist 加1
func (pool *Pool) workingIncr() {
	if len(pool.workingList) < pool.maxSize {
		pool.workingList <- working
	}
}

// workinglist 减1
func (pool *Pool) workingDecr() {
	if len(pool.workingList) > 0 {
		<-pool.workingList
	}

}

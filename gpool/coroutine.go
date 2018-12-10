package gpool

import "sync"

// 协程
type Coroutine struct {

	// 分配得到的任务
	task *Task

	// 所属的池
	pool *Pool

	// Coroutine是否已经关闭
	running bool

	// Coroutine是否正在执行task, true: 正在执行, fasle: 等待分配
	working bool
	// 单独关闭该Coroutine的信号
	closeSign chan int

	runRWLock  sync.RWMutex
	workRWLock sync.RWMutex
}

func newCoroutine(pool *Pool) *Coroutine {

	coro := &Coroutine{}
	coro.pool = pool
	coro.closeSign = make(chan int)

	return coro
}

func (cor *Coroutine) start() {
	// 开启一个新的协程

	cor.setRunning(true)

	go func() {

		for cor.isRunning() {

			select {
			case task, ok := <-cor.pool.taskBenc:

				if !ok || cor.pool.IsClosed() {
					// Pool关闭了退出
					cor.setRunning(false)
					return
				}

				cor.setWorking(true)

				// 为waitingList减1
				cor.pool.waitingDecr()

				if !cor.isRunning() {
					// Pool没有关闭， 但该Coroutine关闭了
					// 因为上面已经为waitingList减1，无需再减
					return
				}

				// Coroutine将要取走一个任务，
				// // 因为上面已经为waitingList减1，无需再减
				// 为workingList加1
				cor.pool.workingIncr()

				task.Run()

				// 再次检查pool是否被关闭(因为可能task执行了一段时间，在这段时间内pool被关闭了)
				if cor.pool.IsClosed() {
					// pool被关闭了，退出
					cor.setRunning(false)
					return
				}

				// task 执行完，要为workingList减1
				cor.pool.workingDecr()

				// 再次检查Coroutine是否被关闭(因为可能task执行了一段时间，在这段时间内Coroutine被Pool关闭了)
				if !cor.isRunning() {
					// 如果关闭了，就不需要为waitingList加1了
					return
				}

				// 为waitingList加1
				cor.pool.waitingIncr()

				cor.setWorking(false)

			case <-cor.pool.exitAllChan:
				// pool池关闭的信号
				cor.setRunning(false)
				return
			case <-cor.closeSign:
				// 单独关闭该Coroutine的信号
				cor.pool.waitingDecr()
				cor.setRunning(false)
				return
			}

		}

	}()
}

// 获取Coroutine是否正在执行(与关闭对应)
func (cor *Coroutine) isRunning() bool {

	cor.runRWLock.RLock()
	result := cor.running
	cor.runRWLock.RUnlock()

	return result
}

// 设置Coroutine是否正在执行(与关闭对应)
func (cor *Coroutine) setRunning(running bool) {
	cor.runRWLock.Lock()
	cor.running = running
	cor.runRWLock.Unlock()
}

// 获取Coroutine是否正在执行Task
func (cor *Coroutine) isWorking() bool {

	cor.workRWLock.RLock()
	result := cor.working
	cor.workRWLock.RUnlock()

	return result
}

// 设置Coroutine是否正在执行Task
func (cor *Coroutine) setWorking(working bool) {
	cor.workRWLock.Lock()
	cor.working = working
	cor.workRWLock.Unlock()
}

func (cor *Coroutine) close() {
	close(cor.closeSign)
}

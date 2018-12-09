package gpool

import "sync"

// 协程
type Coroutine struct {

	// 分配得到的任务
	task *Task

	// 所属的池
	pool *Pool

	running bool

	runRWLock sync.RWMutex
}

func newCoroutine(pool *Pool) *Coroutine {

	coro := &Coroutine{}
	coro.pool = pool

	return coro
}

func (cor *Coroutine) start() {
	// 开启一个新的协程

	cor.setRunning(true)

	go func() {

		for cor.getRunning() {

			select {
			case task, ok := <-cor.pool.taskBenc:

				if !ok || cor.pool.IsClosed() {
					// Pool关闭了退出
					cor.setRunning(false)
					return
				}

				// 为waitingList减1
				cor.pool.waitingDecr()

				if !cor.getRunning() {
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
				if !cor.getRunning() {
					// 如果关闭了，就不需要为waitingList加1了
					return
				}

				// 为waitingList加1
				cor.pool.waitingIncr()

			case <-cor.pool.exitAllChan:
				cor.setRunning(false)
				return
			}

		}

	}()
}

func (cor *Coroutine) getRunning() bool {

	cor.runRWLock.RLock()
	result := cor.running
	cor.runRWLock.RUnlock()

	return result
}

func (cor *Coroutine) setRunning(running bool) {
	cor.runRWLock.Lock()
	cor.running = running
	cor.runRWLock.Unlock()
}

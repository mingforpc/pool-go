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

				if !cor.getRunning() {
					// Pool关闭该Coroutine，退出
					return
				}

				<-cor.pool.waitingList
				cor.pool.wrokingList <- working

				task.Run()

				if !cor.getRunning() {
					return
				}

				<-cor.pool.wrokingList
				cor.pool.waitingList <- waiting

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

package gpool

import (
	"errors"
	"fmt"
	"time"
)

var CLOSE_ERR = errors.New("Pool already closed!")

// 任务
type Task interface {
	Run()
}

// 监控Pool比例的Watcher
type CoroutineLimiter struct {
	timeChan <-chan time.Time

	pool *Pool

	step int // 每次遍历的个数

	stepMask int // 标记遍历位置
}

func NewCoroutineLimiter(timeChan <-chan time.Time, pool *Pool) CoroutineLimiter {
	return CoroutineLimiter{timeChan: timeChan, pool: pool, step: 20}
}

func (limiter CoroutineLimiter) Run() {

	for !limiter.pool.IsClosed() {
		select {
		case _, ok := <-limiter.timeChan:
			if !ok {
				// limiter.timeChan这个channel 关闭了，退出
				return
			}

			limiter.limit()

		case <-limiter.pool.exitAllChan:
			// pool的退出信号
			return
		}
	}

}

func (limiter *CoroutineLimiter) limit() {

	pool := limiter.pool
	currentRatio := pool.getCurrentRatio()

	fmt.Printf("current:[%f], ratio:[%f]\n", currentRatio, pool.ratio)

	if pool.GetCorosCount() > pool.initSize && currentRatio > pool.ratio {
		// 只有当Coroutine大于最小数才执行

		pool.lock.Lock()

		corosLen := len(pool.coros)

		for i := 0; i < limiter.step; i++ {

			if corosLen > limiter.stepMask {

				idleCoro := pool.coros[limiter.stepMask]

				if !idleCoro.isWorking() {
					idleCoro.close()

					if limiter.stepMask == corosLen-1 {
						pool.coros = pool.coros[:limiter.stepMask]
					} else {
						pool.coros = append(pool.coros[:limiter.stepMask], pool.coros[limiter.stepMask+1:]...)
					}

					limiter.stepMask++
					break
				} else {
					limiter.stepMask++
				}

			} else {
				limiter.stepMask = 0
			}

		}

		// limiter.stepMask += limiter.step

		pool.lock.Unlock()

	}

}

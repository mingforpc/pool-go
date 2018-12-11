package gpool

const DONE = 1

// 协程
type Coroutine struct {

	// 分配得到的任务
	task Task
}

func NewCoroutine(task Task) Coroutine {

	return Coroutine{task: task}

}

func (cor *Coroutine) Start() {

	task := cor.task

	val := task.runnable.Run()

	result := NewReuslt(DONE, val)

	task.done <- result

	close(task.done)
}

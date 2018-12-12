package gpool

import (
	"fmt"
)

// 协程
type Coroutine struct {

	// 分配得到的任务
	task Task
}

func NewCoroutine(task Task) Coroutine {

	return Coroutine{task: task}

}

func (cor *Coroutine) Start() {

	defer cor.corRecover()

	task := cor.task

	val := task.runnable.Run()

	result := NewReuslt(DONE, val)

	task.done <- result

	close(task.done)
}

func (cor *Coroutine) corRecover() {
	if err := recover(); err != nil {
		fmt.Println(err)
		result := NewReuslt(DONE, nil)

		cor.task.done <- result

		close(cor.task.done)
	}
}

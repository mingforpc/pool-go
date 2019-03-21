package gpool

import (
	"fmt"
)

// Goroutine 封装的go
type Goroutine struct {

	// 分配得到的任务
	task *Task
}

// NewGoroutine 创建一个NewGoroutine
func NewGoroutine(task *Task) Goroutine {

	return Goroutine{task: task}

}

// Start 开始执行
func (cor *Goroutine) Start() {

	defer cor.corRecover()

	task := cor.task

	val, err := task.runnable.Run(task.ctx, task.args)

	result := NewReuslt(TaskDone, val, err)

	task.setResult(&result)

	task.done <- &result
	close(task.done)

}

func (cor *Goroutine) corRecover() {
	if err := recover(); err != nil {

		result := NewReuslt(TaskExcept, nil, fmt.Errorf("%+v", err))

		cor.task.setResult(&result)
		cor.task.done <- &result

		close(cor.task.done)
	}
}

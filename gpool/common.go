package gpool

import (
	"errors"
)

const DONE = 1
const EXCEPT = -1

var ERR_POOL_COSED = errors.New("Pool closed")

// 执行的方法
type Runnable interface {
	Run() interface{}
}

// 返回值
type Result struct {
	status int
	runRes interface{}
}

func NewReuslt(status int, runRes interface{}) Result {
	return Result{status: status, runRes: runRes}
}

// 任务对象
type Task struct {
	runnable Runnable
	done     chan Result
	result   *Result
}

func NewTask(runnable Runnable) Task {

	task := Task{runnable: runnable, done: make(chan Result, 1)}
	return task
}

func (task *Task) Done() Result {
	result, _ := <-task.done
	task.result = &result
	return result
}

func (task *Task) IsDone() bool {
	if task.result == nil {
		return false
	} else {
		return true
	}
}

func (task *Task) Result() *Result {
	return task.result
}

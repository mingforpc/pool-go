package gpool

import (
	"errors"
	"sync"
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
	Status int
	RunRes interface{}
}

func NewReuslt(status int, runRes interface{}) Result {
	return Result{Status: status, RunRes: runRes}
}

// 任务对象
type Task struct {
	runnable Runnable
	done     chan *Result
	result   *Result

	resultLock sync.RWMutex
}

func NewTask(runnable Runnable) *Task {

	task := &Task{runnable: runnable, done: make(chan *Result, 1)}
	return task
}

func (task *Task) Done() *Result {
	result, _ := <-task.done

	return result
}

func (task *Task) IsDone() bool {
	done := false
	task.resultLock.RLock()
	if task.result != nil {
		done = true
	}
	task.resultLock.RUnlock()

	return done
}

func (task *Task) Result() *Result {

	task.resultLock.RLock()
	result := task.result
	task.resultLock.RUnlock()

	return result
}

func (task *Task) setResult(result *Result) {
	task.resultLock.Lock()
	task.result = result
	task.resultLock.Unlock()
}

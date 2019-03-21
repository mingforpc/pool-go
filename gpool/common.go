package gpool

import (
	"context"
	"errors"
	"sync"
)

// Task 的执行完成状态
const (
	TaskDone   = 1
	TaskExcept = -1
)

var ERR_POOL_COSED = errors.New("Pool closed")

// Runnable 执行的方法
type Runnable interface {
	Run(ctx context.Context, args ...interface{}) (interface{}, error)
}

type funcRunnable struct {
	runFun func(ctx context.Context, args ...interface{}) (interface{}, error)
}

func (f *funcRunnable) Run(ctx context.Context, args ...interface{}) (interface{}, error) {
	return f.runFun(ctx, args)
}

// Result 返回值
type Result struct {
	Status int
	RunRes interface{}
	Err    error
}

func NewReuslt(status int, runRes interface{}, err error) Result {
	return Result{Status: status, RunRes: runRes, Err: err}
}

// Task 任务对象
type Task struct {
	ctx      context.Context
	args     []interface{}
	runnable Runnable
	done     chan *Result
	result   *Result

	resultLock sync.RWMutex
}

// NewTask 创建一个task
func NewTask(runnable Runnable, args []interface{}) *Task {

	task := &Task{runnable: runnable, done: make(chan *Result, 1), args: args}
	return task
}

// Done 阻塞等待任务完成
func (task *Task) Done() *Result {
	result, _ := <-task.done

	return result
}

// IsDone 检查任务是否完成
func (task *Task) IsDone() bool {
	done := false
	task.resultLock.RLock()
	if task.result != nil {
		done = true
	}
	task.resultLock.RUnlock()

	return done
}

// Result 获得任务的返回结果
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

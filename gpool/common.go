package gpool

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

	task := Task{runnable: runnable, done: make(chan Result)}
	return task
}

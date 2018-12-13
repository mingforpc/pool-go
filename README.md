# pool-go

一个简单的goroutine封装，并没有进行goroutine的复用，只是通过channel进行goroutine的数量限制

## 使用方法

```go

...
type TSNoWait struct {
	i int
}

func NewTSNoWait(i int) TSNoWait {
	return TSNoWait{i: i}
}

func (ts TSNoWait) Run() interface{} {
	fmt.Println("sleep:" + strconv.Itoa(ts.i))
	time.Sleep(time.Duration(ts.i) * time.Second)
	return ts.i
}
...


// 指定最大的goroutine数量
p := gpool.NewPool(10)
p.Serv()

var wg sync.WaitGroup

wg.Add(1)
// Pool.Run(runnable Runnable)接受一个Runnable的Interface作为参数
// Runnable的Interface只需要实现 `Run() interface{}` 方法
// 返回的Task对象可以判断任务是否完成，和接受返回的`interface{}`对象
task, err := p.Run(NewTSNoWait(5))

// 返回task是否已经完成
task.IsDone()

// task.Done()会返回result对象，会阻塞直至task完成
result := task.Done()

// 任务结果的状态码
// DONE:表示完成
// EXCEPT: 表示出现异常了
result.Status

// ` Run() interface{} `中返回的interface{}结果
result.RunRes
```
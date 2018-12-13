package test

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/mingforpc/pool-go/gpool"
)

type TSWait struct {
	i  int
	wg *sync.WaitGroup
}

func NewTSWait(i int, wg *sync.WaitGroup) TSWait {
	return TSWait{i: i, wg: wg}
}

func (ts TSWait) Run() interface{} {
	fmt.Println("sleep:" + strconv.Itoa(ts.i))
	time.Sleep(time.Duration(ts.i) * time.Second)
	ts.wg.Done()
	return nil
}

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

type TSPainc struct {
	i int
}

func NewTSPainc(i int) TSPainc {
	return TSPainc{i: i}
}

func (ts TSPainc) Run() interface{} {
	fmt.Println("sleep:" + strconv.Itoa(ts.i))
	time.Sleep(time.Duration(ts.i) * time.Second)
	panic("抛出异常！")
}

// 测试Task.IsDone()函数，和暴露其中一个问题
func TestRuneOnce(t *testing.T) {
	p := gpool.NewPool(10)
	p.Serv()

	var wg sync.WaitGroup

	wg.Add(1)
	task, _ := p.Run(NewTSWait(1, &wg))

	done := task.IsDone()
	fmt.Printf("task.IsDone(): %t \n", done)
	if done == true {
		panic("Task 并没有完成！")
	}

	wg.Wait()

	// TODO: 这里需要sleep一下的原因是:
	// Runnable 的 wg.Done() 是会比 Coroutine 中的 `task.done <- &result`这句先执行的
	// 那样可能导致wg.Wait()不阻塞了，但是还没有通知到task，现在任务已经完成了
	// 所以当这个sleep注释掉时，这个case是通过不了的
	// time.Sleep(time.Second)

	done = task.IsDone()
	fmt.Printf("task.IsDone(): %t \n", done)
	if done == false {
		panic("Task 已经完成！")
	}

	p.Close()
}

// 测试Task.Done()函数
func TestTaskDone(t *testing.T) {

	p := gpool.NewPool(10)
	p.Serv()

	task, _ := p.Run(NewTSNoWait(5))

	done := task.IsDone()
	fmt.Printf("task.IsDone(): %t \n", done)
	if done == true {
		panic("Task 并没有完成！")
	}

	result := task.Done()

	if result.Status != gpool.DONE {
		panic("返回结果的Status异常！")
	}

	if result.RunRes != 5 {
		panic("返回结果异常！")
	}

	done = task.IsDone()
	fmt.Printf("task.IsDone(): %t \n", done)
	if done == false {
		panic("Task 已经完成！")
	}

	p.Close()

}

func TestMaxTask(t *testing.T) {
	p := gpool.NewPool(1000)
	p.Serv()

	var wg sync.WaitGroup
	wg.Add(20000)
	for i := 0; i < 20000; i++ {

		p.Run(NewTSWait(1, &wg))
	}

	wg.Wait()

	p.Close()
}

func TestClose(t *testing.T) {
	p := gpool.NewPool(10)
	p.Serv()

	task1, err := p.Run(NewTSNoWait(2))

	if err != nil {
		panic(err)
	}

	if task1 == nil {
		panic("Task1 不应该为nil！")
	}

	p.Close()

	task2, err := p.Run(NewTSNoWait(2))

	if err == nil {
		panic(err)
	}

	if task2 != nil {
		panic("Task2 应该为nil！")
	}
}

func TestPainc(t *testing.T) {
	p := gpool.NewPool(10)
	p.Serv()

	task, _ := p.Run(NewTSPainc(1))

	task.Done()

	rs := task.Result()

	if rs.Status != gpool.EXCEPT {
		panic("Result 的状态码应该为EXCEPT")
	}
}

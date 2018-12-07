package gpool

import (
	"errors"
)

var CLOSE_ERR = errors.New("Pool already closed!")

// 任务
type Task interface {
	Run()
}

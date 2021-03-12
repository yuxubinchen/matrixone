package routines

import (
	"matrixbase/pkg/vm/routines/task"
	"matrixbase/pkg/vm/routines/worker"
)

type Routines interface {
	Run()
	Stop()
	AddTask(task.Task)
}

type routines struct {
	cnt uint64
	num uint64
	ch  chan struct{}
	ws  []worker.Worker
}
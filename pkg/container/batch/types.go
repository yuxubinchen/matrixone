package batch

import (
	"matrixbase/pkg/container/vector"

	aio "github.com/traetox/goaio"
)

type Info struct {
	Alg int
	Wg  *WaitGroup
}

type WaitGroup struct {
	Ap *aio.AIO
	Id aio.RequestId
}

type Batch struct {
	Is    []Info
	Sels  []int64
	Attrs []string
	Vecs  []*vector.Vector
}
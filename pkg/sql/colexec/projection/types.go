package projection

import "matrixbase/pkg/sql/colexec/extend"

type Argument struct {
	Attrs []string
	Es    []extend.Extend
}
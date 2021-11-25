// Copyright 2021 Matrix Origin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ring

import (
	"matrixone/pkg/container/types"
	"matrixone/pkg/container/vector"
	"matrixone/pkg/vm/mheap"
)

type Ring interface {
	Count() int

	Size() int

	Dup() Ring
	Type() types.Type

	String() string

	Free(*mheap.Mheap)
	Grow(*mheap.Mheap) error

	SetLength(int)
	Shrink([]int64)

	Shuffle([]int64, *mheap.Mheap) error

	Eval([]int64) *vector.Vector

	Fill(int64, int64, int64, *vector.Vector)

	BulkFill(int64, []int64, *vector.Vector)

	Mul(int64, int64)
	Add(interface{}, int64, int64)
}
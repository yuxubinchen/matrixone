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

package min

import (
	"fmt"
	"math"
	"matrixone/pkg/container/nulls"
	"matrixone/pkg/container/ring"
	"matrixone/pkg/container/types"
	"matrixone/pkg/container/vector"
	"matrixone/pkg/encoding"
	"matrixone/pkg/vm/mheap"
)

func NewUInt64(typ types.Type) *UInt64Ring {
	return &UInt64Ring{Typ: typ}
}

func (r *UInt64Ring) String() string {
	return fmt.Sprintf("%v-%v", r.Vs, r.Ns)
}

func (r *UInt64Ring) Free(m *mheap.Mheap) {
	if r.Da != nil {
		mheap.Free(m, r.Da)
		r.Da = nil
		r.Vs = nil
		r.Ns = nil
	}
}

func (r *UInt64Ring) Count() int {
	return len(r.Vs)
}

func (r *UInt64Ring) Size() int {
	return cap(r.Da)
}

func (r *UInt64Ring) Dup() ring.Ring {
	return &UInt64Ring{
		Typ: r.Typ,
	}
}

func (r *UInt64Ring) Type() types.Type {
	return r.Typ
}

func (r *UInt64Ring) SetLength(n int) {
	r.Vs = r.Vs[:n]
	r.Ns = r.Ns[:n]
}

func (r *UInt64Ring) Shrink(sels []int64) {
	for i, sel := range sels {
		r.Vs[i] = r.Vs[sel]
		r.Ns[i] = r.Ns[sel]
	}
	r.Vs = r.Vs[:len(sels)]
	r.Ns = r.Ns[:len(sels)]
}

func (r *UInt64Ring) Shuffle(_ []int64, _ *mheap.Mheap) error {
	return nil
}

func (r *UInt64Ring) Grow(m *mheap.Mheap) error {
	n := len(r.Vs)
	if n == 0 {
		data, err := mheap.Alloc(m, 8*8)
		if err != nil {
			return err
		}
		r.Da = data
		r.Ns = make([]int64, 0, 8)
		r.Vs = encoding.DecodeUint64Slice(data)
	} else if n+1 >= cap(r.Vs) {
		data, err := mheap.Grow(m, r.Da, int64(n+1)*8)
		if err != nil {
			return err
		}
		mheap.Free(m, r.Da)
		r.Da = data
		r.Vs = encoding.DecodeUint64Slice(data)
	}
	r.Vs = r.Vs[:n+1]
	r.Vs[n] = math.MaxUint64
	r.Ns = append(r.Ns, 0)
	return nil
}

func (r *UInt64Ring) Fill(i int64, sel, _ int64, vec *vector.Vector) {
	if v := vec.Col.([]uint64)[sel]; v < r.Vs[i] {
		r.Vs[i] = v
	}
	if nulls.Contains(vec.Nsp, uint64(sel)) {
		r.Ns[i]++
	}
}

func (r *UInt64Ring) BulkFill(i int64, _ []int64, vec *vector.Vector) {
	vs := vec.Col.([]uint64)
	for _, v := range vs {
		if v < r.Vs[i] {
			r.Vs[i] = v
		}
	}
	r.Ns[i] += int64(nulls.Length(vec.Nsp))
}

func (r *UInt64Ring) Add(a interface{}, x, y int64) {
	ar := a.(*UInt64Ring)
	if r.Vs[x] > ar.Vs[y] {
		r.Vs[x] = ar.Vs[y]
	}
	r.Ns[x] += ar.Ns[y]
}

func (r *UInt64Ring) Mul(_, _ int64) {
}

func (r *UInt64Ring) Eval(zs []int64) *vector.Vector {
	defer func() {
		r.Da = nil
		r.Vs = nil
		r.Ns = nil
	}()
	nsp := new(nulls.Nulls)
	for i, z := range zs {
		if z-r.Ns[i] == 0 {
			nulls.Add(nsp, uint64(i))
		}
	}
	return &vector.Vector{
		Nsp:  nsp,
		Data: r.Da,
		Col:  r.Vs,
		Or:   false,
		Typ:  r.Typ,
	}
}
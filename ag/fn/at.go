// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fn

import (
	"github.com/nlpodyssey/spago/mat"
)

// At is an operator to obtain the i,j-th value of a matrix.
type At[T mat.DType, O Operand[T]] struct {
	x O
	i int
	j int
}

// NewAt returns a new At Function.
func NewAt[T mat.DType, O Operand[T]](x O, i int, j int) *At[T, O] {
	return &At[T, O]{x: x, i: i, j: j}
}

// Operands returns the list of operands.
func (r *At[T, O]) Operands() []O {
	return []O{r.x}
}

// Forward computes the output of the function.
func (r *At[T, O]) Forward() mat.Matrix[T] {
	return mat.NewScalar(r.x.Value().At(r.i, r.j))
}

// Backward computes the backward pass.
func (r *At[T, O]) Backward(gy mat.Matrix[T]) {
	if r.x.RequiresGrad() {
		dx := r.x.Value().ZerosLike()
		defer mat.ReleaseMatrix(dx)
		dx.Set(r.i, r.j, gy.Scalar())
		r.x.PropagateGrad(dx)
	}
}

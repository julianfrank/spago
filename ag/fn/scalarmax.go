// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fn

import (
	"github.com/nlpodyssey/spago/mat"
)

// ScalarMax is an operator to perform reduce-max function on a list of scalars.
// It gets the maximum element of the Operand x
type ScalarMax[T mat.DType, O Operand[T]] struct {
	xs       []O
	argmax   int
	operands []O
}

// NewScalarMax returns a new ScalarMax Function.
func NewScalarMax[T mat.DType, O Operand[T]](xs []O) *ScalarMax[T, O] {
	return &ScalarMax[T, O]{xs: xs}
}

// Operands returns the list of operands.
func (r *ScalarMax[T, O]) Operands() []O {
	return r.xs
}

// Forward computes the output of this function.
func (r *ScalarMax[T, O]) Forward() mat.Matrix[T] {
	var max T
	var argmax int
	for i, x := range r.xs {
		val := x.Value().Scalar()
		if val > max {
			max = val
			argmax = i
		}
	}
	r.argmax = argmax
	return mat.NewScalar(max)
}

// Backward computes the backward pass.
func (r *ScalarMax[T, O]) Backward(gy mat.Matrix[T]) {
	if !mat.IsScalar(gy) {
		panic("fn: the gradient had to be a scalar")
	}
	target := r.xs[r.argmax]
	if target.RequiresGrad() {
		target.PropagateGrad(gy)
	}
}
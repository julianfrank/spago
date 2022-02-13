// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package layernorm

import (
	"github.com/nlpodyssey/spago/ag"
	"github.com/nlpodyssey/spago/mat"
	"github.com/nlpodyssey/spago/nn"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestModel_Forward(t *testing.T) {
	t.Run("float32", testModelForward[float32])
	t.Run("float64", testModelForward[float64])
}

func testModelForward[T mat.DType](t *testing.T) {
	model := newTestModel[T]()
	g := ag.NewGraph[T](ag.WithMode[T](ag.Training))

	// == Forward
	x := g.NewVariable(mat.NewVecDense([]T{0.4, 0.8, -0.7, -0.5}), true)
	y := nn.ToNode[T](nn.Reify(model, g).Forward(x))

	assert.InDeltaSlice(t, []T{1.157863, 0.2, -0.561554, -0.444658}, y.Value().Data(), 1.0e-06)

	// == Backward
	y.PropagateGrad(mat.NewVecDense([]T{-1.0, -0.2, 0.4, 0.6}))
	g.BackwardAll()

	assert.InDeltaSlice(t, []T{-0.496261, 0.280677, -0.408772, 0.624355}, x.Grad().Data(), 1.0e-06)
	assert.InDeltaSlice(t, []T{-0.644658, -0.257863, -0.45126, -0.483493}, model.W.Grad().Data(), 1.0e-06)
	assert.InDeltaSlice(t, []T{-1.0, -0.2, 0.4, 0.6}, model.B.Grad().Data(), 1.0e-06)
}

func newTestModel[T mat.DType]() *Model[T] {
	model := New[T](4)
	model.W.Value().SetData([]T{0.4, 0.0, -0.3, 0.8})
	model.B.Value().SetData([]T{0.9, 0.2, -0.9, 0.2})
	return model
}
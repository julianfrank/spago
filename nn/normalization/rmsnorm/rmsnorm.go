// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rmsnorm implements the Root Mean Square Layer Normalization method.
//
// Reference: "Root Mean Square Layer Normalization" by Biao Zhang and Rico Sennrich (2019).
// (https://arxiv.org/pdf/1910.07467.pdf)
package rmsnorm

import (
	"encoding/gob"

	"github.com/nlpodyssey/spago/ag"
	"github.com/nlpodyssey/spago/mat"
	"github.com/nlpodyssey/spago/nn"
)

var _ nn.Model = &Model[float32]{}

// Model contains the serializable parameters.
type Model[T mat.DType] struct {
	nn.Module
	W nn.Param[T] `spago:"type:weights"`
	B nn.Param[T] `spago:"type:biases"`
}

func init() {
	gob.Register(&Model[float32]{})
	gob.Register(&Model[float64]{})
}

// New returns a new model with parameters initialized to zeros.
func New[T mat.DType](size int) *Model[T] {
	return &Model[T]{
		W: nn.NewParam[T](mat.NewEmptyVecDense[T](size)),
		B: nn.NewParam[T](mat.NewEmptyVecDense[T](size)),
	}
}

// Forward performs the forward step for each input node and returns the result.
func (m *Model[T]) Forward(xs ...ag.Node[T]) []ag.Node[T] {
	eps := ag.Constant[T](1e-10)
	ys := make([]ag.Node[T], len(xs))
	for i, x := range xs {
		rms := ag.Sqrt(ag.ReduceMean(ag.Square(x)))
		ys[i] = ag.Add[T](ag.Prod[T](ag.DivScalar(x, ag.AddScalar(rms, eps)), m.W), m.B)
	}
	return ys
}

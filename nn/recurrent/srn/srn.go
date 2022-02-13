// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srn

import (
	"encoding/gob"
	"github.com/nlpodyssey/spago/ag"
	"github.com/nlpodyssey/spago/mat"
	"github.com/nlpodyssey/spago/nn"
	"log"
)

var _ nn.Model[float32] = &Model[float32]{}

// Model contains the serializable parameters.
type Model[T mat.DType] struct {
	nn.BaseModel[T]
	W      nn.Param[T] `spago:"type:weights"`
	WRec   nn.Param[T] `spago:"type:weights"`
	B      nn.Param[T] `spago:"type:biases"`
	States []*State[T] `spago:"scope:processor"`
}

// State represent a state of the SRN recurrent network.
type State[T mat.DType] struct {
	Y ag.Node[T]
}

func init() {
	gob.Register(&Model[float32]{})
	gob.Register(&Model[float64]{})
}

// New returns a new model with parameters initialized to zeros.
func New[T mat.DType](in, out int) *Model[T] {
	return &Model[T]{
		W:    nn.NewParam[T](mat.NewEmptyDense[T](out, in)),
		WRec: nn.NewParam[T](mat.NewEmptyDense[T](out, out)),
		B:    nn.NewParam[T](mat.NewEmptyVecDense[T](out)),
	}
}

// SetInitialState sets the initial state of the recurrent network.
// It panics if one or more states are already present.
func (m *Model[T]) SetInitialState(state *State[T]) {
	if len(m.States) > 0 {
		log.Fatal("srn: the initial state must be set before any input")
	}
	m.States = append(m.States, state)
}

// Forward performs the forward step for each input node and returns the result.
func (m *Model[T]) Forward(xs ...ag.Node[T]) []ag.Node[T] {
	ys := make([]ag.Node[T], len(xs))
	for i, x := range xs {
		s := m.forward(x)
		m.States = append(m.States, s)
		ys[i] = s.Y
	}
	return ys
}

// LastState returns the last state of the recurrent network.
// It returns nil if there are no states.
func (m *Model[T]) LastState() *State[T] {
	n := len(m.States)
	if n == 0 {
		return nil
	}
	return m.States[n-1]
}

// y = tanh(w (dot) x + b + wRec (dot) yPrev)
func (m *Model[T]) forward(x ag.Node[T]) (s *State[T]) {
	g := m.Graph()
	s = new(State[T])
	yPrev := m.prev()
	s.Y = g.Tanh(g.Affine(m.B, m.W, x, m.WRec, yPrev))
	return
}

func (m *Model[T]) prev() (yPrev ag.Node[T]) {
	s := m.LastState()
	if s != nil {
		yPrev = s.Y
	}
	return
}
// Copyright 2021 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gmlp

import (
	"encoding/gob"

	"github.com/nlpodyssey/spago/mat"
	"github.com/nlpodyssey/spago/nn"
	"github.com/nlpodyssey/spago/nn/activation"
	"github.com/nlpodyssey/spago/nn/linear"
	"github.com/nlpodyssey/spago/nn/sgu"
	"github.com/nlpodyssey/spago/nn/stack"
)

var _ nn.Model = &Block[float32]{}

// Block is the core model of the gMLP.
type Block[T mat.DType] struct {
	nn.Module
	*stack.Model[T]
}

// BlockConfig provides configuration parameters for a single Block of the gMLP Model.
type BlockConfig struct {
	Dim        int
	DimFF      int
	SeqLen     int
	Activation activation.Name
}

func init() {
	gob.Register(&Block[float32]{})
	gob.Register(&Block[float64]{})
}

// NewBlock returns a new Block.
func NewBlock[T mat.DType](config BlockConfig) *Block[T] {
	return &Block[T]{
		Model: stack.New[T](
			linear.New[T](config.Dim, config.DimFF),
			activation.New[T](activation.GELU),
			sgu.New(sgu.Config[T]{
				Dim:        config.DimFF,
				DimSeq:     config.SeqLen,
				InitEps:    1e-3,
				Activation: config.Activation,
			}),
			linear.New[T](config.DimFF/2, config.Dim),
		),
	}
}

// Copyright 2020 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fn

import (
	"github.com/nlpodyssey/spago/mat"
	"github.com/nlpodyssey/spago/mat/matutils"
	matsort "github.com/nlpodyssey/spago/mat/sort"
	"sort"
)

// SparseMax function implementation, based on https://github.com/gokceneraslan/SparseMax.torch
type SparseMax[T mat.DType] struct {
	x Operand[T]
	y mat.Matrix[T] // initialized during the forward pass, required by the backward pass
}

var _ Function[float32] = &SparseMax[float32]{}
var _ Function[float32] = &SparseMaxLoss[float32]{}

// NewSparseMax returns a new SparseMax Function.
func NewSparseMax[T mat.DType](x Operand[T]) *SparseMax[T] {
	return &SparseMax[T]{x: x}
}

// Forward computes the output of the function.
func (s *SparseMax[T]) Forward() mat.Matrix[T] {
	s.y = mat.NewVecDense(sparseMax(translateInput(s.x.Value().Data())))
	return s.y
}

// Backward computes the backward pass.
func (s *SparseMax[T]) Backward(gy mat.Matrix[T]) {
	if s.x.RequiresGrad() {
		var nzSum T = 0.0
		var nzCount T = 0.0
		s.y.DoVecNonZero(func(i int, _ T) {
			nzSum += gy.AtVec(i)
			nzCount++
		})
		nzSum = nzSum / nzCount

		gx := s.x.Value().ZerosLike()
		defer mat.ReleaseMatrix(gx)
		s.y.DoVecNonZero(func(i int, _ T) {
			gx.SetVec(i, gy.AtVec(i)-nzSum)
		})

		s.x.PropagateGrad(gx)
	}
}

// translateInput translates the input by max for numerical stability
func translateInput[T mat.DType](v []T) []T {
	maximum := max(v)
	translated := make([]T, len(v))
	for i := range v {
		translated[i] = v[i] - maximum
	}
	return translated
}

func sparseMaxCommon[T mat.DType](v []T) (zs, bounds, cumSumInput []T, tau T) {
	zs = make([]T, len(v))
	copy(zs, v)

	// Sort zs in descending order.
	sort.Sort(sort.Reverse(matsort.DTSlice[T](zs)))

	bounds = make([]T, len(zs))
	for i := range bounds {
		bounds[i] = 1 + T(i+1)*zs[i]
	}

	cumSumInput = make([]T, len(zs))
	matutils.CumSum(cumSumInput, zs)

	k := -1
	tau = 0.0
	for i := range zs {
		if bounds[i] > cumSumInput[i] {
			if k < (i + 1) {
				k = i + 1
			}
			tau += zs[i]
		}
	}
	tau = (tau - 1) / T(k)

	return zs, bounds, cumSumInput, tau
}

func sparseMax[T mat.DType](v []T) []T {
	zs, _, _, tau := sparseMaxCommon(v)

	//Reuses zs to avoid allocating new slice
	for i := range zs {
		zs[i] = mat.Max(0.0, v[i]-tau)
	}
	return zs
}

// SparseMaxLoss function implementation, based on https://github.com/gokceneraslan/SparseMax.torch
type SparseMaxLoss[T mat.DType] struct {
	x   Operand[T]
	tau T             // computed during the forward pass
	y   mat.Matrix[T] // computed during forward pass
}

// NewSparseMaxLoss returns a new SparseMaxLoss Function.
func NewSparseMaxLoss[T mat.DType](x Operand[T]) *SparseMaxLoss[T] {
	return &SparseMaxLoss[T]{x: x}
}

// sparseMaxLoss computes the sparseMax loss function and returns
// the loss and the tau parameter (needed by backward)
func sparseMaxLoss[T mat.DType](v []T) ([]T, T) {
	zs, bounds, cumSumInput, tau := sparseMaxCommon(v)
	var regTerm T = 0.0
	tauSquared := tau * tau

	for i := range zs {
		if bounds[i] > cumSumInput[i] {
			regTerm += zs[i]*zs[i] - tauSquared
		}
	}
	regTerm = regTerm*0.5 + 0.5

	// Reuse zs to avoid allocating a new slice
	for i := range zs {
		zs[i] = v[i] - regTerm
	}
	return zs, tau
}

// Forward computes the output of the function.
func (s *SparseMaxLoss[T]) Forward() mat.Matrix[T] {
	output, tau := sparseMaxLoss(s.x.Value().Data())
	s.y = mat.NewVecDense(output)
	s.tau = tau
	return s.y
}

// Backward computes the backward pass.
func (s *SparseMaxLoss[T]) Backward(gy mat.Matrix[T]) {
	if s.x.RequiresGrad() {
		tau := s.tau
		gySum := gy.Sum()

		sparseMax := s.x.Value().Apply(func(_, _ int, v T) T {
			return mat.Max(0, v-tau) * gySum
		})
		defer mat.ReleaseMatrix(sparseMax)

		gx := gy.Sub(sparseMax)
		defer mat.ReleaseMatrix(gx)

		s.x.PropagateGrad(gx)
	}
}
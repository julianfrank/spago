// Copyright 2022 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sort

import (
	"github.com/nlpodyssey/spago/mat"
	"sort"
)

// DTSlice attaches the methods of sort.Interface to []T, sorting in increasing order
// (not-a-number values are treated as less than other values).
type DTSlice[T mat.DType] []T

func (p DTSlice[_]) Len() int           { return len(p) }
func (p DTSlice[_]) Less(i, j int) bool { return p[i] < p[j] || isNaN(p[i]) && !isNaN(p[j]) }
func (p DTSlice[_]) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// isNaN is a copy of math.IsNaN.
func isNaN[T mat.DType](f T) bool {
	return f != f
}

// Sort is a convenience method.
func (p DTSlice[_]) Sort() { sort.Sort(p) }

// Code from `https://stackoverflow.com/questions/31141202/get-the-indices-of-the-array-after-sorting-in-golang`

// Slice is an extension of sort.Interface which keeps track of the
// original indices of the elements of a slice after sorting.
type Slice struct {
	sort.Interface
	Indices []int
}

// Swap swaps the elements with indexes i and j.
func (s Slice) Swap(i, j int) {
	s.Interface.Swap(i, j)
	s.Indices[i], s.Indices[j] = s.Indices[j], s.Indices[i]
}

// NewSlice returns a new Slice.
func NewSlice(n sort.Interface) *Slice {
	s := &Slice{Interface: n, Indices: make([]int, n.Len())}
	for i := range s.Indices {
		s.Indices[i] = i
	}
	return s
}

// NewIntSlice returns a new Slice for the given sequence of int values.
func NewIntSlice(n ...int) *Slice {
	return NewSlice(sort.IntSlice(n))
}

// NewDTSlice returns a new Slice for the given sequence of float32 values.
func NewDTSlice[T mat.DType](n ...T) *Slice {
	return NewSlice(DTSlice[T](n))
}

// NewStringSlice returns a new Slice for the given sequence of string values.
func NewStringSlice(n ...string) *Slice {
	return NewSlice(sort.StringSlice(n))
}
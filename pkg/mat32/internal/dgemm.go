// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"runtime"
	"sync"

	"github.com/nlpodyssey/spago/pkg/mat32/internal/asm/f32"
)

const (
	mLT0 = "blas: m < 0"
	nLT0 = "blas: n < 0"
	kLT0 = "blas: k < 0"

	badTranspose = "blas: illegal transpose"

	badLdA = "blas: bad leading dimension of A"
	badLdB = "blas: bad leading dimension of B"
	badLdC = "blas: bad leading dimension of C"

	shortA = "blas: insufficient length of a"
	shortB = "blas: insufficient length of b"
	shortC = "blas: insufficient length of c"

	// NoTrans is a Transpose option.
	NoTrans Transpose = 'N'
	// Trans is a Transpose option.
	Trans Transpose = 'T'
	// ConjTrans is a Transpose option.
	ConjTrans Transpose = 'C'
)

// Transpose specifies the transposition operation of a matrix.
type Transpose byte

// Dgemm performs one of the matrix-matrix operations
//  C = alpha * A * B + beta * C
//  C = alpha * Aᵀ * B + beta * C
//  C = alpha * A * Bᵀ + beta * C
//  C = alpha * Aᵀ * Bᵀ + beta * C
// where A is an m×k or k×m dense matrix, B is an n×k or k×n dense matrix, C is
// an m×n matrix, and alpha and beta are scalars. tA and tB specify whether A or
// B are transposed.
//gocyclo:ignore
func Dgemm(tA, tB Transpose, m, n, k int, alpha float32, a []float32, lda int, b []float32, ldb int, beta float32, c []float32, ldc int) {
	switch tA {
	default:
		panic(badTranspose)
	case NoTrans, Trans, ConjTrans:
	}
	switch tB {
	default:
		panic(badTranspose)
	case NoTrans, Trans, ConjTrans:
	}
	if m < 0 {
		panic(mLT0)
	}
	if n < 0 {
		panic(nLT0)
	}
	if k < 0 {
		panic(kLT0)
	}
	aTrans := tA == Trans || tA == ConjTrans
	if aTrans {
		if lda < max(1, m) {
			panic(badLdA)
		}
	} else {
		if lda < max(1, k) {
			panic(badLdA)
		}
	}
	bTrans := tB == Trans || tB == ConjTrans
	if bTrans {
		if ldb < max(1, k) {
			panic(badLdB)
		}
	} else {
		if ldb < max(1, n) {
			panic(badLdB)
		}
	}
	if ldc < max(1, n) {
		panic(badLdC)
	}

	// Quick return if possible.
	if m == 0 || n == 0 {
		return
	}

	// For zero matrix size the following slice length checks are trivially satisfied.
	if aTrans {
		if len(a) < (k-1)*lda+m {
			panic(shortA)
		}
	} else {
		if len(a) < (m-1)*lda+k {
			panic(shortA)
		}
	}
	if bTrans {
		if len(b) < (n-1)*ldb+k {
			panic(shortB)
		}
	} else {
		if len(b) < (k-1)*ldb+n {
			panic(shortB)
		}
	}
	if len(c) < (m-1)*ldc+n {
		panic(shortC)
	}

	// Quick return if possible.
	if (alpha == 0 || k == 0) && beta == 1 {
		return
	}

	// scale c
	if beta != 1 {
		if beta == 0 {
			for i := 0; i < m; i++ {
				ctmp := c[i*ldc : i*ldc+n]
				for j := range ctmp {
					ctmp[j] = 0
				}
			}
		} else {
			for i := 0; i < m; i++ {
				ctmp := c[i*ldc : i*ldc+n]
				for j := range ctmp {
					ctmp[j] *= beta
				}
			}
		}
	}

	dgemmParallel(aTrans, bTrans, m, n, k, a, lda, b, ldb, c, ldc, alpha)
}

func dgemmParallel(aTrans, bTrans bool, m, n, k int, a []float32, lda int, b []float32, ldb int, c []float32, ldc int, alpha float32) {
	// dgemmParallel computes a parallel matrix multiplication by partitioning
	// a and b into sub-blocks, and updating c with the multiplication of the sub-block
	// In all cases,
	// A = [ 	A_11	A_12 ... 	A_1j
	//			A_21	A_22 ...	A_2j
	//				...
	//			A_i1	A_i2 ...	A_ij]
	//
	// and same for B. All of the submatrix sizes are blockSize×blockSize except
	// at the edges.
	//
	// In all cases, there is one dimension for each matrix along which
	// C must be updated sequentially.
	// Cij = \sum_k Aik Bki,	(A * B)
	// Cij = \sum_k Aki Bkj,	(Aᵀ * B)
	// Cij = \sum_k Aik Bjk,	(A * Bᵀ)
	// Cij = \sum_k Aki Bjk,	(Aᵀ * Bᵀ)
	//
	// This code computes one {i, j} block sequentially along the k dimension,
	// and computes all of the {i, j} blocks concurrently. This
	// partitioning allows Cij to be updated in-place without race-conditions.
	// Instead of launching a goroutine for each possible concurrent computation,
	// a number of worker goroutines are created and channels are used to pass
	// available and completed cases.
	//
	// http://alexkr.com/docs/matrixmult.pdf is a good reference on matrix-matrix
	// multiplies, though this code does not copy matrices to attempt to eliminate
	// cache misses.

	maxKLen := k
	parBlocks := blocks(m, blockSize) * blocks(n, blockSize)
	if parBlocks < minParBlock {
		// The matrix multiplication is small in the dimensions where it can be
		// computed concurrently. Just do it in serial.
		DgemmSerial(aTrans, bTrans, m, n, k, a, lda, b, ldb, c, ldc, alpha)
		return
	}

	// workerLimit acts a number of maximum concurrent workers,
	// with the limit set to the number of procs available.
	workerLimit := make(chan struct{}, runtime.GOMAXPROCS(0))

	// wg is used to wait for all
	var wg sync.WaitGroup
	wg.Add(parBlocks)
	defer wg.Wait()

	for i := 0; i < m; i += blockSize {
		for j := 0; j < n; j += blockSize {
			workerLimit <- struct{}{}
			go func(i, j int) {
				defer func() {
					wg.Done()
					<-workerLimit
				}()

				leni := blockSize
				if i+leni > m {
					leni = m - i
				}
				lenj := blockSize
				if j+lenj > n {
					lenj = n - j
				}

				cSub := sliceView32(c, ldc, i, j, leni, lenj)

				// Compute A_ik B_kj for all k
				for k := 0; k < maxKLen; k += blockSize {
					lenk := blockSize
					if k+lenk > maxKLen {
						lenk = maxKLen - k
					}
					var aSub, bSub []float32
					if aTrans {
						aSub = sliceView32(a, lda, k, i, lenk, leni)
					} else {
						aSub = sliceView32(a, lda, i, k, leni, lenk)
					}
					if bTrans {
						bSub = sliceView32(b, ldb, j, k, lenj, lenk)
					} else {
						bSub = sliceView32(b, ldb, k, j, lenk, lenj)
					}
					DgemmSerial(aTrans, bTrans, leni, lenj, lenk, aSub, lda, bSub, ldb, cSub, ldc, alpha)
				}
			}(i, j)
		}
	}
}

// DgemmSerial is serial matrix multiply
func DgemmSerial(aTrans, bTrans bool, m, n, k int, a []float32, lda int, b []float32, ldb int, c []float32, ldc int, alpha float32) {
	switch {
	case !aTrans && !bTrans:
		dgemmSerialNotNot(m, n, k, a, lda, b, ldb, c, ldc, alpha)
		return
	case aTrans && !bTrans:
		dgemmSerialTransNot(m, n, k, a, lda, b, ldb, c, ldc, alpha)
		return
	case !aTrans && bTrans:
		dgemmSerialNotTrans(m, n, k, a, lda, b, ldb, c, ldc, alpha)
		return
	case aTrans && bTrans:
		dgemmSerialTransTrans(m, n, k, a, lda, b, ldb, c, ldc, alpha)
		return
	default:
		panic("unreachable")
	}
}

// dgemmSerial where neither a nor b are transposed
func dgemmSerialNotNot(m, n, k int, a []float32, lda int, b []float32, ldb int, c []float32, ldc int, alpha float32) {
	// This style is used instead of the literal [i*stride +j]) is used because
	// approximately 5 times faster as of go 1.3.
	for i := 0; i < m; i++ {
		ctmp := c[i*ldc : i*ldc+n]
		for l, v := range a[i*lda : i*lda+k] {
			tmp := alpha * v
			if tmp != 0 {
				f32.AxpyUnitary(tmp, b[l*ldb:l*ldb+n], ctmp)
			}
		}
	}
}

// dgemmSerial where neither a is transposed and b is not
func dgemmSerialTransNot(m, n, k int, a []float32, lda int, b []float32, ldb int, c []float32, ldc int, alpha float32) {
	// This style is used instead of the literal [i*stride +j]) is used because
	// approximately 5 times faster as of go 1.3.
	for l := 0; l < k; l++ {
		btmp := b[l*ldb : l*ldb+n]
		for i, v := range a[l*lda : l*lda+m] {
			tmp := alpha * v
			if tmp != 0 {
				ctmp := c[i*ldc : i*ldc+n]
				f32.AxpyUnitary(tmp, btmp, ctmp)
			}
		}
	}
}

// dgemmSerial where neither a is not transposed and b is
func dgemmSerialNotTrans(m, n, k int, a []float32, lda int, b []float32, ldb int, c []float32, ldc int, alpha float32) {
	// This style is used instead of the literal [i*stride +j]) is used because
	// approximately 5 times faster as of go 1.3.
	for i := 0; i < m; i++ {
		atmp := a[i*lda : i*lda+k]
		ctmp := c[i*ldc : i*ldc+n]
		for j := 0; j < n; j++ {
			ctmp[j] += alpha * f32.DotUnitary(atmp, b[j*ldb:j*ldb+k])
		}
	}
}

// dgemmSerial where both are transposed
func dgemmSerialTransTrans(m, n, k int, a []float32, lda int, b []float32, ldb int, c []float32, ldc int, alpha float32) {
	// This style is used instead of the literal [i*stride +j]) is used because
	// approximately 5 times faster as of go 1.3.
	for l := 0; l < k; l++ {
		for i, v := range a[l*lda : l*lda+m] {
			tmp := alpha * v
			if tmp != 0 {
				ctmp := c[i*ldc : i*ldc+n]
				f32.AxpyInc(tmp, b[l:], ctmp, uintptr(n), uintptr(ldb), 1, 0, 0)
			}
		}
	}
}

func sliceView32(a []float32, lda, i, j, r, c int) []float32 {
	return a[i*lda+j : (i+r-1)*lda+j+c]
}

// [SD]gemm behavior constants. These are kept here to keep them out of the
// way during single precision code genration.
const (
	blockSize   = 64 // b x b matrix
	minParBlock = 4  // minimum number of blocks needed to go parallel
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// blocks returns the number of divisions of the dimension length with the given
// block size.
func blocks(dim, bsize int) int {
	return (dim + bsize - 1) / bsize
}
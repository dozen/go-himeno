package main

import (
	"testing"
)

//プリプロセス
const MIMAX = 0
const MJMAX = 0
const MKMAX = 0

var (
	p    [MIMAX][MJMAX][MKMAX]float32
	a    [4][MIMAX][MJMAX][MKMAX]float32
	b    [3][MIMAX][MJMAX][MKMAX]float32
	c    [3][MIMAX][MJMAX][MKMAX]float32
	bnd  [MIMAX][MJMAX][MKMAX]float32
	wrk1 [MIMAX][MJMAX][MKMAX]float32
	wrk2 [MIMAX][MJMAX][MKMAX]float32

	pSlice    [][][]float32
	aSlice    [][][][]float32
	bSlice    [][][][]float32
	cSlice    [][][][]float32
	bndSlice  [][][]float32
	wrk1Slice [][][]float32
	wrk2Slice [][][]float32

	imax = MIMAX - 1
	jmax = MJMAX - 1
	kmax = MKMAX - 1

	omega  = float32(0.8)
	target = 60.0
)

func initArray() {
	var i, j, k int
	p = [MIMAX][MJMAX][MKMAX]float32{}
	a = [4][MIMAX][MJMAX][MKMAX]float32{}
	b = [3][MIMAX][MJMAX][MKMAX]float32{}
	c = [3][MIMAX][MJMAX][MKMAX]float32{}
	bnd = [MIMAX][MJMAX][MKMAX]float32{}
	wrk1 = [MIMAX][MJMAX][MKMAX]float32{}
	wrk2 = [MIMAX][MJMAX][MKMAX]float32{}

	for i = 0; i < imax; i++ {
		for j = 0; j < jmax; j++ {
			for k = 0; k < kmax; k++ {
				a[0][i][j][k] = 1.0
				a[1][i][j][k] = 1.0
				a[2][i][j][k] = 1.0
				a[3][i][j][k] = 1.0 / 6.0
				b[0][i][j][k] = 0.0
				b[1][i][j][k] = 0.0
				b[2][i][j][k] = 0.0
				c[0][i][j][k] = 1.0
				c[1][i][j][k] = 1.0
				c[2][i][j][k] = 1.0
				p[i][j][k] = (float32)(i*i) / (float32)((imax-1)*(imax-1))
				wrk1[i][j][k] = 0.0
				bnd[i][j][k] = 1.0
			}
		}
	}
}

func initSlice() {
	var i, j, k int

	aSlice = make([][][][]float32, 4)
	bSlice = make([][][][]float32, 3)
	cSlice = make([][][][]float32, 3)

	pSlice = make([][][]float32, MIMAX, MIMAX)
	aSlice[0] = make([][][]float32, MIMAX, MIMAX)
	aSlice[1] = make([][][]float32, MIMAX, MIMAX)
	aSlice[2] = make([][][]float32, MIMAX, MIMAX)
	aSlice[3] = make([][][]float32, MIMAX, MIMAX)
	bSlice[0] = make([][][]float32, MIMAX, MIMAX)
	bSlice[1] = make([][][]float32, MIMAX, MIMAX)
	bSlice[2] = make([][][]float32, MIMAX, MIMAX)
	cSlice[0] = make([][][]float32, MIMAX, MIMAX)
	cSlice[1] = make([][][]float32, MIMAX, MIMAX)
	cSlice[2] = make([][][]float32, MIMAX, MIMAX)
	bndSlice = make([][][]float32, MIMAX, MIMAX)
	wrk1Slice = make([][][]float32, MIMAX, MIMAX)
	wrk2Slice = make([][][]float32, MIMAX, MIMAX)

	for i = 0; i < imax; i++ {
		pSlice[i] = make([][]float32, MJMAX, MJMAX)
		aSlice[0][i] = make([][]float32, MJMAX, MJMAX)
		aSlice[1][i] = make([][]float32, MJMAX, MJMAX)
		aSlice[2][i] = make([][]float32, MJMAX, MJMAX)
		aSlice[3][i] = make([][]float32, MJMAX, MJMAX)
		bSlice[0][i] = make([][]float32, MJMAX, MJMAX)
		bSlice[1][i] = make([][]float32, MJMAX, MJMAX)
		bSlice[2][i] = make([][]float32, MJMAX, MJMAX)
		cSlice[0][i] = make([][]float32, MJMAX, MJMAX)
		cSlice[1][i] = make([][]float32, MJMAX, MJMAX)
		cSlice[2][i] = make([][]float32, MJMAX, MJMAX)
		bndSlice[i] = make([][]float32, MJMAX, MJMAX)
		wrk1Slice[i] = make([][]float32, MJMAX, MJMAX)
		wrk2Slice[i] = make([][]float32, MJMAX, MJMAX)
		for j = 0; j < jmax; j++ {
			pSlice[i][j] = make([]float32, MKMAX, MKMAX)
			aSlice[0][i][j] = make([]float32, MKMAX, MKMAX)
			aSlice[1][i][j] = make([]float32, MKMAX, MKMAX)
			aSlice[2][i][j] = make([]float32, MKMAX, MKMAX)
			aSlice[3][i][j] = make([]float32, MKMAX, MKMAX)
			bSlice[0][i][j] = make([]float32, MKMAX, MKMAX)
			bSlice[1][i][j] = make([]float32, MKMAX, MKMAX)
			bSlice[2][i][j] = make([]float32, MKMAX, MKMAX)
			cSlice[0][i][j] = make([]float32, MKMAX, MKMAX)
			cSlice[1][i][j] = make([]float32, MKMAX, MKMAX)
			cSlice[2][i][j] = make([]float32, MKMAX, MKMAX)
			bndSlice[i][j] = make([]float32, MKMAX, MKMAX)
			wrk1Slice[i][j] = make([]float32, MKMAX, MKMAX)
			wrk2Slice[i][j] = make([]float32, MKMAX, MKMAX)
			for k = 0; k < kmax; k++ {
				aSlice[0][i][j][k] = 1.0
				aSlice[1][i][j][k] = 1.0
				aSlice[2][i][j][k] = 1.0
				aSlice[3][i][j][k] = 1.0 / 6.0
				bSlice[0][i][j][k] = 0.0
				bSlice[1][i][j][k] = 0.0
				bSlice[2][i][j][k] = 0.0
				cSlice[0][i][j][k] = 1.0
				cSlice[1][i][j][k] = 1.0
				cSlice[2][i][j][k] = 1.0
				pSlice[i][j][k] = (float32)(i*i) / (float32)((imax-1)*(imax-1))
				wrk1Slice[i][j][k] = 0.0
				bndSlice[i][j][k] = 1.0
			}
		}
	}
}

func jacobi() float32 {
	var i, j, k int
	var gosa, s0, ss float32
	gosa = 0.0

	for i = 1; i < imax-1; i++ {
		for j = 1; j < jmax-1; j++ {
			for k = 1; k < kmax-1; k++ {
				s0 = a[0][i][j][k]*p[i+1][j][k] +
					a[1][i][j][k]*p[i][j+1][k] +
					a[2][i][j][k]*p[i][j][k+1] +
					b[0][i][j][k]*(p[i+1][j+1][k]-p[i+1][j-1][k]-p[i-1][j+1][k]+p[i-1][j-1][k]) +
					b[1][i][j][k]*(p[i][j+1][k+1]-p[i][j-1][k+1]-p[i][j+1][k-1]+p[i][j-1][k-1]) +
					b[2][i][j][k]*(p[i+1][j][k+1]-p[i-1][j][k+1]-p[i+1][j][k-1]+p[i-1][j][k-1]) +
					c[0][i][j][k]*p[i-1][j][k] + c[1][i][j][k]*p[i][j-1][k] + c[2][i][j][k]*p[i][j][k-1] +
					wrk1[i][j][k]

				ss = (s0*a[3][i][j][k] - p[i][j][k]) * bnd[i][j][k]

				gosa += ss * ss
				/* gosa= (gosa > ss*ss) ? a : b; */
				wrk2[i][j][k] = p[i][j][k] + omega*ss
			}
		}
	}

	for i = 1; i < imax-1; i++ {
		for j = 1; j < jmax-1; j++ {
			for k = 1; k < kmax-1; k++ {
				p[i][j][k] = wrk2[i][j][k]
			}
		}
	}

	return gosa
}

func jacobiSlice() float32 {
	var i, j, k int
	var gosa, s0, ss float32
	gosa = 0.0

	for i = 1; i < imax-1; i++ {
		for j = 1; j < jmax-1; j++ {
			for k = 1; k < kmax-1; k++ {
				s0 = aSlice[0][i][j][k]*pSlice[i+1][j][k] +
					aSlice[1][i][j][k]*pSlice[i][j+1][k] +
					aSlice[2][i][j][k]*pSlice[i][j][k+1] +
					bSlice[0][i][j][k]*(pSlice[i+1][j+1][k]-pSlice[i+1][j-1][k]-pSlice[i-1][j+1][k]+pSlice[i-1][j-1][k]) +
					bSlice[1][i][j][k]*(pSlice[i][j+1][k+1]-pSlice[i][j-1][k+1]-pSlice[i][j+1][k-1]+pSlice[i][j-1][k-1]) +
					bSlice[2][i][j][k]*(pSlice[i+1][j][k+1]-pSlice[i-1][j][k+1]-pSlice[i+1][j][k-1]+pSlice[i-1][j][k-1]) +
					cSlice[0][i][j][k]*pSlice[i-1][j][k] + cSlice[1][i][j][k]*pSlice[i][j-1][k] + cSlice[2][i][j][k]*pSlice[i][j][k-1] +
					wrk1Slice[i][j][k]

				ss = (s0*aSlice[3][i][j][k] - pSlice[i][j][k]) * bndSlice[i][j][k]

				gosa += ss * ss
				/* gosa= (gosa > ss*ss) ? a : b; */
				wrk2Slice[i][j][k] = pSlice[i][j][k] + omega*ss
			}
		}
	}

	for i = 1; i < imax-1; i++ {
		for j = 1; j < jmax-1; j++ {
			for k = 1; k < kmax-1; k++ {
				pSlice[i][j][k] = wrk2Slice[i][j][k]
			}
		}
	}

	return gosa
}

func BenchmarkJacobiArray(b *testing.B) {
	initArray()
	b.ResetTimer()
	gosa := float32(0)
	for i := 0; i < b.N; i++ {
		gosa = jacobi()
	}
	b.Logf("gosa: %v", gosa)
}

func BenchmarkJacobiSlice(b *testing.B) {
	initSlice()
	b.ResetTimer()
	gosa := float32(0)
	for i := 0; i < b.N; i++ {
		gosa = jacobiSlice()
	}
	b.Logf("gosa: %v", gosa)
}

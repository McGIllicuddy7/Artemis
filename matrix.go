package main

import (
	"math/cmplx"
	"strconv"
)

type matrix_t struct {
	data   []complex128
	height int
	width  int
}
type vector_t []complex128

func vector_add(v0 vector_t, v1 vector_t) vector_t {
	out := make([]complex128, len(v0))
	for i := 0; i < len(v0); i++ {
		out[i] = v0[i] + v1[i]
	}
	return out
}
func vector_sub(v0 vector_t, v1 vector_t) vector_t {
	out := make([]complex128, len(v0))
	for i := 0; i < len(v0); i++ {
		out[i] = v0[i] - v1[i]
	}
	return out
}
func vector_scale(v0 vector_t, s complex128) vector_t {
	out := make(vector_t, len(v0))
	for i := 0; i < len(v0); i++ {
		out[i] = v0[i] * s
	}
	return out
}
func vector_dot(v0 vector_t, v1 vector_t) complex128 {
	out := complex128(0)
	for i := 0; i < len(v0); i++ {
		out += v0[i] * v1[i]
	}
	return out
}
func vector_length(v0 vector_t) complex128 {
	return cmplx.Sqrt(vector_dot(v0, v0))
}
func (this *matrix_t) get(x int, y int) complex128 {
	return this.data[y*this.width+x]
}
func (this *matrix_t) set(x int, y int, v complex128) {
	this.data[y*this.width+x] = v
}
func (this *matrix_t) to_string() string {
	out := ""
	for j := 0; j < this.height; j++ {
		for i := 0; i < this.width; i++ {
			out += strconv.FormatComplex(this.get(i, j), 'G', 4, 128) + " "
		}
		out += "\n"
	}
	return out
}
func (this *matrix_t) swap_rows(r0 int, r1 int) {
	for i := 0; i < this.width; i++ {
		tmp0 := this.get(i, r0)
		tmp1 := this.get(i, r1)
		this.set(i, r1, tmp0)
		this.set(i, r0, tmp1)
	}
}

// adds r0 to r1 scaled by s
func (this *matrix_t) add_rows(r0 int, r1 int, s complex128) {
	for i := 0; i < this.width; i++ {
		tmp0 := this.get(i, r0) * s
		this.set(i, r1, tmp0+this.get(i, r1))
	}
}
func (this *matrix_t) sub_rows(r0 int, r1 int, s complex128) {
	for i := 0; i < this.width; i++ {
		tmp0 := this.get(i, r0) * s
		this.set(i, r1, this.get(i, r1)-tmp0)
	}
}
func (this *matrix_t) scale_row(r0 int, s complex128) {
	for i := 0; i < this.width; i++ {
		tmp0 := this.get(i, r0) * s
		this.set(i, r0, tmp0)
	}
}
func matrix_row_reduce(matrx matrix_t) matrix_t {
	mtrx := matrix_t{make([]complex128, len(matrx.data)), matrx.height, matrx.width}
	for i := 0; i < len(matrx.data); i++ {
		mtrx.data[i] = matrx.data[i]
	}
	for i := 0; i < mtrx.width; i++ {
		r := i
		degen := false
		for mtrx.get(i, r) == 0 {
			r++
			if r >= mtrx.height {
				degen = true
				break
			}
		}
		if degen {
			continue
		}
		if r != i {
			mtrx.swap_rows(r, i)
		}
		v := mtrx.get(i, i)
		mtrx.scale_row(i, 1/v)
		for j := 0; j < mtrx.height; j++ {
			if j == i {
				continue
			}
			mlt := mtrx.get(i, j)
			mtrx.sub_rows(i, j, mlt)
		}
	}
	return mtrx
}
func (this *matrix_t) determinant() complex128 {
	mtrx := matrix_t{make([]complex128, len(this.data)), this.height, this.width}
	out := complex(1, 0)
	for i := 0; i < len(this.data); i++ {
		mtrx.data[i] = this.data[i]
	}
	for i := 0; i < mtrx.width; i++ {
		r := i
		degen := false
		for mtrx.get(i, r) == 0 {
			r++
			if r >= mtrx.height {
				degen = true
				break
			}
		}
		if degen {
			continue
		}
		if r != i {
			mtrx.swap_rows(r, i)
			out *= -1
		}
		v := mtrx.get(i, i)
		mtrx.scale_row(i, 1/v)
		out *= v
		for j := 0; j < mtrx.height; j++ {
			if j == i {
				continue
			}
			mlt := mtrx.get(i, j)
			if mlt == 0 {
				continue
			}
			mtrx.sub_rows(i, j, mlt)
		}
	}
	for i := 0; i < this.width; i++ {
		out *= mtrx.get(i, i)
	}
	return out
}

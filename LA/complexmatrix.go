package La

import (
	"artemis/utils"
	"math/rand"
)

type MatrixComplex struct {
	data   []complex128
	height int
	width  int
}

func (tmat *MatrixComplex) Get(x int, y int) complex128 {
	return tmat.data[y*tmat.width+x]
}
func (tmat *MatrixComplex) Set(x int, y int, v complex128) {
	tmat.data[y*tmat.width+x] = v
}
func (tmat *MatrixComplex) NumCols() int {
	return tmat.width
}
func (tmat *MatrixComplex) NumRows() int {
	return tmat.width
}
func (tmat *MatrixComplex) Clone() MatrixComplex {
	out := MatrixComplex{make([]complex128, len(tmat.data)), tmat.height, tmat.width}
	copy(out.data, tmat.data)
	return out
}
func (tmat *MatrixComplex) ToString() string {
	out := ""
	out_strs := make([]string, 0, tmat.height*tmat.width)
	for j := 0; j < tmat.height; j++ {
		for i := 0; i < tmat.width; i++ {
			out_strs = append(out_strs, utils.FormatComplex(tmat.Get(i, j)))
		}
	}
	max := 0
	for i := 0; i < len(out_strs); i++ {
		if len(out_strs[i]) > max {
			max = len(out_strs[i])
		}
	}
	for i := 0; i < len(out_strs); i++ {
		if out_strs[i][0] != '-' {
			out_strs[i] = " " + out_strs[i]
		}
		for len(out_strs[i]) < max {
			out_strs[i] += " "
		}
	}
	for j := 0; j < tmat.height; j++ {
		for i := 0; i < tmat.width; i++ {
			out += out_strs[j*tmat.width+i]
			if i < tmat.width-1 {
				out += " "
			}
		}
		out += "\n"
	}
	return out
}
func (tmat *MatrixComplex) SwapRows(r0 int, r1 int) {
	for i := 0; i < tmat.width; i++ {
		tmp0 := tmat.Get(i, r0)
		tmp1 := tmat.Get(i, r1)
		tmat.Set(i, r1, tmp0)
		tmat.Set(i, r0, tmp1)
	}
}
func (tmat *MatrixComplex) AddRows(r0 int, r1 int, s complex128) {
	for i := 0; i < tmat.width; i++ {
		tmp0 := tmat.Get(i, r0) * s
		tmat.Set(i, r1, tmp0+tmat.Get(i, r1))
	}
}
func (tmat *MatrixComplex) SubRows(r0 int, r1 int, s complex128) {
	for i := 0; i < tmat.width; i++ {
		tmp0 := tmat.Get(i, r0) * s
		tmat.Set(i, r1, tmat.Get(i, r1)-tmp0)
	}
}
func (tmat *MatrixComplex) ScaleRow(r0 int, s complex128) {
	for i := 0; i < tmat.width; i++ {
		tmp0 := tmat.Get(i, r0) * s
		tmat.Set(i, r0, tmp0)
	}
}
func ComplexIdentity(n int) MatrixComplex {
	out := MatrixComplex{make([]complex128, n*n), n, n}
	for i := 0; i < n; i++ {
		out.Set(i, i, 1)
	}
	return out
}
func ComplexMatrixAdd(m0 MatrixComplex, m1 MatrixComplex) MatrixComplex {
	if m0.height != m1.height || m0.width != m1.width {
		panic("error adding matrices without the same dimension")
	}
	out := MatrixComplex{make([]complex128, m0.height*m0.width), m0.height, m0.width}
	for i := 0; i < m0.height; i++ {
		for j := 0; j < m0.width; j++ {
			out.Set(j, i, (m0.Get(j, i) + m1.Get(j, i)))
		}
	}
	return out
}
func ComplexMatrixSub(m0 MatrixComplex, m1 MatrixComplex) MatrixComplex {
	if m0.height != m1.height || m0.width != m1.width {
		panic("error adding matrices without the same dimension")
	}
	out := MatrixComplex{make([]complex128, m0.height*m0.width), m0.height, m0.width}
	for i := 0; i < m0.height; i++ {
		for j := 0; j < m0.width; j++ {
			out.Set(j, i, (m1.Get(j, i) - m0.Get(j, i)))
		}
	}
	return out
}
func ComplexMatrixScale(m0 MatrixComplex, s complex128) MatrixComplex {
	out := MatrixComplex{make([]complex128, m0.height*m0.width), m0.height, m0.width}
	for i := 0; i < len(m0.data); i++ {
		out.data[i] = m0.data[i] * s
	}
	return out
}
func ComplexMatrixRowReduce(matrx MatrixComplex) MatrixComplex {
	mtrx := matrx.Clone()
	for i := 0; i < mtrx.width; i++ {
		r := i
		degen := false
		for mtrx.Get(i, r) == 0 {
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
			mtrx.SwapRows(r, i)
		}
		v := mtrx.Get(i, i)
		mtrx.ScaleRow(i, 1/v)
		for j := 0; j < mtrx.height; j++ {
			if j == i {
				continue
			}
			mlt := mtrx.Get(i, j)
			mtrx.SubRows(i, j, mlt)
		}
	}
	return mtrx
}
func ComplexMatrixFromInts(slice [][]int) MatrixComplex {
	height := len(slice)
	width := len(slice[0])
	out := MatrixComplex{make([]complex128, height*width), height, width}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			out.Set(x, y, complex(float64(slice[y][x]), 0))
		}
	}
	return out
}
func ComplexMatrixPairRowReduce(source MatrixComplex, target MatrixComplex) (MatrixComplex, MatrixComplex) {
	mtrx := source.Clone()
	out := target.Clone()
	for i := 0; i < mtrx.width; i++ {
		r := i
		degen := false
		for mtrx.Get(i, r) == 0 {
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
			out.SwapRows(r, i)
			mtrx.SwapRows(r, i)
		}
		v := mtrx.Get(i, i)
		mtrx.ScaleRow(i, 1/v)
		out.ScaleRow(1, 1/v)
		for j := 0; j < mtrx.height; j++ {
			if j == i {
				continue
			}
			mlt := mtrx.Get(i, j)
			mtrx.SubRows(i, j, mlt)
			out.SubRows(i, j, mlt)
		}
	}
	return mtrx, out
}
func (tmat *MatrixComplex) Determinant() complex128 {
	mtrx := MatrixComplex{make([]complex128, len(tmat.data)), tmat.height, tmat.width}
	out := complex128(1)
	copy(mtrx.data, tmat.data)
	for i := 0; i < mtrx.width; i++ {
		r := i
		degen := false
		for mtrx.Get(i, r) == 0 {
			r++
			if r >= mtrx.height {
				degen = true
				break
			}
		}
		if degen {
			return 0
		}
		if r != i {
			mtrx.SwapRows(r, i)
			out *= -1
		}
		v := mtrx.Get(i, i)
		mtrx.ScaleRow(i, 1/v)
		out = v * out
		for j := 0; j < mtrx.height; j++ {
			if j == i {
				continue
			}
			mlt := mtrx.Get(i, j)
			if mlt == 0 {
				continue
			}
			mtrx.SubRows(i, j, mlt)
		}
	}
	return out
}
func (tmat *MatrixComplex) Inverse() MatrixComplex {
	_, out := ComplexMatrixPairRowReduce(*tmat, ComplexIdentity(tmat.width))
	return out
}
func RandomComplexMatrix(height int, width int) MatrixComplex {
	var out MatrixComplex
	out.data = make([]complex128, height*width)
	out.height = height
	out.width = width
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			out.Set(x, y, complex(float64(rand.Int31()%10-5), 0))
		}
	}
	return out
}
func (tmat *MatrixComplex) ToUpperTriangular() MatrixComplex {
	mtrx := tmat.Clone()
	for i := 0; i < mtrx.width; i++ {
		r := i
		degen := false
		for mtrx.Get(i, r) == 0 {
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
			mtrx.SwapRows(r, i)
		}
		v := mtrx.Get(i, i)
		mtrx.ScaleRow(i, 1/v)
		for j := r; j < mtrx.height; j++ {
			if j == i {
				continue
			}
			mlt := mtrx.Get(i, j)
			mtrx.SubRows(i, j, mlt)
		}
	}
	for i := 0; i < mtrx.width; i++ {
		r := i
		degen := false
		for mtrx.Get(i, r) == 0 {
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
			mtrx.SwapRows(r, i)
		}
		v := mtrx.Get(i, i)
		mtrx.ScaleRow(i, 1/v)
		for j := r; j < mtrx.height; j++ {
			if j == i {
				continue
			}
			mlt := mtrx.Get(i, j)
			mtrx.SubRows(i, j, mlt)
		}
	}
	return mtrx
}

func (tmat *MatrixComplex) Solve(values Vector) Vector {
	mtrx := tmat.ToUpperTriangular()
	definedSymbols := make([]bool, mtrx.width)
	for i := 0; i < len(definedSymbols); i++ {
		definedSymbols[i] = false
	}
	symbolTable := make([]complex128, mtrx.width)
	for i := 0; i < len(symbolTable); i++ {
		symbolTable[i] = 0
	}
	for y := mtrx.height - 1; y >= 0; y-- {
		idx := 0
		for utils.ComplexNearlyEqual(mtrx.Get(idx, y), 0) || definedSymbols[idx] {
			idx++
			if idx >= mtrx.width-1 {
				break
			}
		}
		if idx == -1 {
			continue
		}
		for i := idx + 1; i < mtrx.width; i++ {
			if !utils.ComplexNearlyEqual(mtrx.Get(i, y), 0) && !definedSymbols[i] {
				symbolTable[i] = 1
				definedSymbols[i] = true
			}
		}
		total := complex128(0)
		for x := 0; x < mtrx.width; x++ {
			if x != idx {
				total += mtrx.Get(x, y) * symbolTable[x]
			}
		}
		symbolTable[idx] = -total + values[y]
		definedSymbols[idx] = true
	}
	return symbolTable
}
func (tmat *MatrixComplex) MultByVector(v Vector) Vector {
	out := make(Vector, len(v))
	for i := 0; i < tmat.height; i++ {
		total := complex128(0)
		for j := 0; j < tmat.width; j++ {
			total += tmat.Get(j, i) * v[j]
		}
		out[i] = total
	}
	return out
}

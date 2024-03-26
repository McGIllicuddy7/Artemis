package La

import (
	fr "artemis/fractions"
	"math/rand"
)

type Matrix struct {
	data   []fr.Fraction
	height int
	width  int
}

func (tmat *Matrix) Get(x int, y int) fr.Fraction {
	return tmat.data[y*tmat.width+x]
}
func (tmat *Matrix) Set(x int, y int, v fr.Fraction) {
	tmat.data[y*tmat.width+x] = v
}
func (tmat *Matrix) NumCols() int {
	return tmat.width
}
func (tmat *Matrix) NumRows() int {
	return tmat.height
}
func (tmat *Matrix) Clone() Matrix {
	out := Matrix{make([]fr.Fraction, len(tmat.data)), tmat.height, tmat.width}
	copy(out.data, tmat.data)
	return out
}
func (tmat *Matrix) ToString() string {
	out := ""
	out_strs := make([]string, 0, tmat.height*tmat.width)
	for j := 0; j < tmat.height; j++ {
		for i := 0; i < tmat.width; i++ {
			out_strs = append(out_strs, tmat.Get(i, j).ToString())
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
func (tmat *Matrix) SwapRows(r0 int, r1 int) {
	for i := 0; i < tmat.width; i++ {
		tmp0 := tmat.Get(i, r0)
		tmp1 := tmat.Get(i, r1)
		tmat.Set(i, r1, tmp0)
		tmat.Set(i, r0, tmp1)
	}
}

// adds r0 to r1 scaled by s
func (tmat *Matrix) AddRows(r0 int, r1 int, s fr.Fraction) {
	for i := 0; i < tmat.width; i++ {
		tmp0 := fr.Mult(tmat.Get(i, r0), s)
		tmat.Set(i, r1, fr.Add(tmp0, tmat.Get(i, r1)))
	}
}

// subtracts r0 from r1
func (tmat *Matrix) SubRows(r0 int, r1 int, s fr.Fraction) {
	for i := 0; i < tmat.width; i++ {
		tmp0 := fr.Mult(tmat.Get(i, r0), s)
		tmat.Set(i, r1, fr.Sub(tmat.Get(i, r1), tmp0))
	}
}
func (tmat *Matrix) ScaleRow(r0 int, s fr.Fraction) {
	for i := 0; i < tmat.width; i++ {
		tmp0 := fr.Mult(tmat.Get(i, r0), s)
		tmat.Set(i, r0, tmp0)
	}
}
func Identity(n int) Matrix {
	out := Matrix{make([]fr.Fraction, n*n), n, n}
	for i := 0; i < n; i++ {
		out.Set(i, i, fr.FromInt(1))
	}
	return out
}
func MatrixAdd(m0 Matrix, m1 Matrix) Matrix {
	if m0.height != m1.height || m0.width != m1.width {
		panic("error adding matrices without the same dimension")
	}
	out := Matrix{make([]fr.Fraction, m0.height*m0.width), m0.height, m0.width}
	for i := 0; i < m0.height; i++ {
		for j := 0; j < m0.width; j++ {
			out.Set(j, i, fr.Add(m0.Get(j, i), m1.Get(j, i)))
		}
	}
	return out
}
func MatrixSub(m0 Matrix, m1 Matrix) Matrix {
	if m0.height != m1.height || m0.width != m1.width {
		panic("error subtracting matrices without the same dimension")
	}
	out := Matrix{make([]fr.Fraction, m0.height*m0.width), m0.height, m0.width}
	for i := 0; i < m0.height; i++ {
		for j := 0; j < m0.width; j++ {
			out.Set(j, i, fr.Sub(m1.Get(j, i), m0.Get(j, i)))
		}
	}
	return out
}
func MatrixScale(m0 Matrix, s fr.Fraction) Matrix {
	out := Matrix{make([]fr.Fraction, m0.height*m0.width), m0.height, m0.width}
	for i := 0; i < len(m0.data); i++ {
		out.data[i] = fr.Mult(m0.data[i], s)
	}
	return out
}
func MatrixRowReduce(matrx Matrix) Matrix {
	mtrx := matrx.Clone()
	for i := 0; i < mtrx.width; i++ {
		r := i
		degen := false
		for fr.Equals(mtrx.Get(i, r), fr.FromInt(0)) {
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
		mtrx.ScaleRow(i, fr.Recip(v))
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
func MatrixFromInts(slice [][]int) Matrix {
	height := len(slice)
	width := len(slice[0])
	out := Matrix{make([]fr.Fraction, height*width), height, width}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			out.Set(x, y, fr.FromInt(slice[y][x]))
		}
	}
	return out
}
func MatrixPairRowReduce(source Matrix, target Matrix) (Matrix, Matrix) {
	mtrx := source.Clone()
	out := target.Clone()
	for i := 0; i < mtrx.width; i++ {
		r := i
		degen := false
		for fr.Equals(mtrx.Get(i, r), fr.FromInt(0)) {
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
		mtrx.ScaleRow(i, fr.Recip(v))
		out.ScaleRow(1, fr.Recip(v))
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
func (tmat *Matrix) Determinant() fr.Fraction {
	mtrx := Matrix{make([]fr.Fraction, len(tmat.data)), tmat.height, tmat.width}
	out := fr.FromInt(1)
	copy(mtrx.data, tmat.data)
	for i := 0; i < mtrx.width; i++ {
		r := i
		degen := false
		for fr.Equals(mtrx.Get(i, r), fr.FromInt(0)) {
			r++
			if r >= mtrx.height {
				degen = true
				break
			}
		}
		if degen {
			return fr.FromInt(0)
		}
		if r != i {
			mtrx.SwapRows(r, i)
			out = fr.Scale(out, -1)
		}
		v := mtrx.Get(i, i)
		mtrx.ScaleRow(i, fr.Recip(v))
		out = fr.Mult(out, v)
		for j := 0; j < mtrx.height; j++ {
			if j == i {
				continue
			}
			mlt := mtrx.Get(i, j)
			if fr.Equals(mlt, fr.FromInt(0)) {
				continue
			}
			mtrx.SubRows(i, j, mlt)
		}
	}
	return out
}
func (tmat *Matrix) Inverse() Matrix {
	_, out := MatrixPairRowReduce(*tmat, Identity(tmat.width))
	return out
}
func RandomMatrix(height int, width int) Matrix {
	var out Matrix
	out.data = make([]fr.Fraction, height*width)
	out.height = height
	out.width = width
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			out.Set(x, y, fr.FromInt(int(rand.Int31()%10-5)))
		}
	}
	return out
}
func (tmat *Matrix) ToUpperTriangular() Matrix {
	mtrx := tmat.Clone()
	for i := 0; i < mtrx.width; i++ {
		r := i
		degen := false
		for fr.Equals(mtrx.Get(i, r), fr.FromInt(0)) {
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
		mtrx.ScaleRow(i, fr.Recip(v))
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
func (tmat *Matrix) Solve(values Vector) Vector {
	mtrx := tmat.Clone()
	vals := values.Clone()
	for i := 0; i < mtrx.width; i++ {
		r := i
		degen := false
		for fr.Equals(mtrx.Get(i, r), fr.FromInt(0)) {
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
			vals.Swap(r, i)
		}
		v := mtrx.Get(i, i)
		recip := fr.Recip(v)
		mtrx.ScaleRow(i, recip)
		vals[i] *= recip.ToComplex()
		for j := r; j < mtrx.height; j++ {
			if j == i {
				continue
			}
			mlt := mtrx.Get(i, j)
			mtrx.SubRows(i, j, mlt)
			vals[j] -= vals[i] * mlt.ToComplex()
		}
	}
	symbolTable := make(Vector, tmat.width)
	definedSymbols := make([]bool, tmat.width)
	for i := 0; i < len(symbolTable); i++ {
		symbolTable[i] = 0
		definedSymbols[i] = false
	}
	for y := tmat.height - 1; y >= 0; y-- {
		syms := make([]int, 0)
		for x := 0; x < tmat.width; x++ {
			if !(fr.Equals(mtrx.Get(x, y), fr.FromInt(0))) {
				syms = append(syms, x)
			}
		}
		undefined := make([]int, 0)
		for i := 0; i < len(syms); i++ {
			if !definedSymbols[syms[i]] {
				undefined = append(undefined, syms[i])
			}
		}
		for i := len(undefined) - 1; i > 0; i-- {
			idx := undefined[i]
			symbolTable[idx] = 1
			definedSymbols[idx] = true
		}
		if len(undefined) > 0 {
			//println(undefined[0])
			newSym := complex128(0)
			for i := 1; i < len(syms); i++ {
				newSym -= symbolTable[syms[i]] * tmat.Get(syms[i], y).ToComplex()
			}
			newSym += vals[y]
			newSym /= mtrx.Get(undefined[0], y).ToComplex()
			symbolTable[undefined[0]] = newSym
			definedSymbols[undefined[0]] = true
		}
	}
	return symbolTable
}
func (tmat *Matrix) MultByVector(v Vector) Vector {
	out := make(Vector, len(v))
	for i := 0; i < tmat.height; i++ {
		total := complex128(0)
		for j := 0; j < tmat.width; j++ {
			total += tmat.Get(j, i).ToComplex() * v[j]
		}
		out[i] = total
	}
	return out
}
func (tmat *Matrix) ToComplex() MatrixComplex {
	out := MatrixComplex{make([]complex128, tmat.height*tmat.width), tmat.height, tmat.width}
	for i := 0; i < out.width*out.height; i++ {
		out.data[i] = tmat.data[i].ToComplex()
	}
	return out
}

func ZeroMatrix(height int, width int) Matrix {
	out := Matrix{make([]fr.Fraction, height*width), height, width}
	return out
}

package La

import (
	al "artemis/algebra"
	fr "artemis/fractions"
	"artemis/utils"
)

type PolyMatrix struct {
	data   []al.Polynomial
	height int
	width  int
}

func (tmat *PolyMatrix) Get(x int, y int) al.Polynomial {
	return tmat.data[y*tmat.width+x]
}
func (tmat *PolyMatrix) Set(x int, y int, v al.Polynomial) {
	tmat.data[y*tmat.height+x] = v
}
func (tmat *PolyMatrix) ToString() string {
	strings := make([]string, tmat.height*tmat.width)
	for i := 0; i < len(tmat.data); i++ {
		strings[i] = tmat.data[i].ToString()
	}
	strings = utils.NormalizeStrlens(strings)
	out := ""
	for y := 0; y < tmat.height; y++ {
		for x := 0; x < tmat.width; x++ {
			out += strings[y*tmat.height+x]
			if x < tmat.width-1 {
				out += " "
			}
		}
		out += "\n"
	}
	return out
}
func (tmat *Matrix) ToEigenMatrix() PolyMatrix {
	out := PolyMatrix{make([]al.Polynomial, tmat.height*tmat.width), tmat.height, tmat.width}
	for y := 0; y < tmat.height; y++ {
		for x := 0; x < tmat.width; x++ {
			var v al.Polynomial
			if x == y {
				v = al.CompPoly(tmat.Get(x, y), fr.NewFrac(-1, 1), 1)
			} else {
				v = al.CompPoly(tmat.Get(x, y), fr.NewFrac(0, 1), 0)
			}
			out.Set(x, y, v)
		}
	}
	return out
}
func (tmat *Matrix) ToPolyMatrix() PolyMatrix {
	out := PolyMatrix{make([]al.Polynomial, tmat.height*tmat.width), tmat.height, tmat.width}
	for y := 0; y < tmat.height; y++ {
		for x := 0; x < tmat.width; x++ {
			v := al.CompPoly(tmat.Get(x, y), fr.NewFrac(0, 1), 0)
			out.Set(x, y, v)
		}
	}
	return out
}
func (tmat *PolyMatrix) ToMatrix() Matrix {
	out := Matrix{make([]fr.Fraction, tmat.height*tmat.width), tmat.height, tmat.width}
	for y := 0; y < tmat.height; y++ {
		for x := 0; x < tmat.width; x++ {
			v := tmat.Get(x, y).ZeroCoef()
			out.Set(x, y, v)
		}
	}
	return out
}
func (tmat *PolyMatrix) elimRowCollumn(idx int) PolyMatrix {
	out := PolyMatrix{make([]al.Polynomial, (tmat.width-1)*(tmat.height-1)), tmat.height - 1, tmat.width - 1}
	for y := 1; y < tmat.height; y++ {
		dy := y - 1
		for x := 0; x < tmat.width; x++ {
			dx := x
			if dx == idx {
				continue
			}
			if x > idx {
				dx--
			}
			out.Set(dx, dy, tmat.Get(x, y))
		}
	}
	return out
}

func (tmat PolyMatrix) CharacteristicPolynomial() al.Polynomial {
	if tmat.width == 2 && tmat.height == 2 {
		a := tmat.Get(0, 0)
		b := tmat.Get(0, 1)
		c := tmat.Get(1, 0)
		d := tmat.Get(1, 1)
		ad := al.PolynomialMult(a, d)
		bc := al.PolynomialMult(b, c)
		ret := al.PolynomialSub(ad, bc)
		return ret
	}
	var out al.Polynomial
	for i := 0; i < tmat.width; i++ {
		tmp := tmat.elimRowCollumn(i)
		m := tmat.Get(i, 0)
		det := tmp.CharacteristicPolynomial()
		mdet := al.PolynomialMult(m, det)
		if i%2 == 0 {
			out = al.PolynomialAdd(out, mdet)
		} else {
			mdet = al.PolynomialScale(mdet, fr.NewFrac(-1, 1))
			out = al.PolynomialAdd(out, mdet)
		}
	}
	return out
}
func (tmat *Matrix) EigenValues() []complex128 {
	eigen := tmat.ToEigenMatrix()
	poly := eigen.CharacteristicPolynomial()
	return poly.FindZeros()
}
func (tmat *Matrix) EigenVectors() []Vector {
	eigens := tmat.EigenValues()
	out := make([]Vector, 0)
	for i := 0; i < len(eigens); i++ {
		mat := ComplexMatrixSub(tmat.ToComplex(), ComplexMatrixScale(ComplexIdentity(tmat.height), eigens[i]))
		tmp := mat.Solve(ZeroVector(tmat.height))
		if !VectorEqual(tmat.MultByVector(tmp), ZeroVector(tmat.height)) {
			println("failed")
		}
		out = append(out, tmp)
	}
	return out
}

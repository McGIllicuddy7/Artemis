package La

import (
	"artemis/utils"
)

type Vector []complex128

func ZeroVector(dim int) Vector {
	out := make(Vector, dim)
	for i := 0; i < dim; i++ {
		out[i] = 0
	}
	return out
}
func (vec *Vector) Swap(a int, b int) {
	tmp := (*vec)[a]
	(*vec)[a] = (*vec)[b]
	(*vec)[b] = tmp
}
func (vec *Vector) Clone() Vector {
	out := make(Vector, len(*vec))
	for i := 0; i < len(out); i++ {
		out[i] = (*vec)[i]
	}
	return out
}
func (vec Vector) ToString() string {
	out := ""
	out += "<"
	for k := 0; k < len(vec); k++ {
		out += utils.FormatComplex(vec[k])
		if k != len(vec)-1 {
			out += ", "
		}
	}
	out += ">"
	return out
}
func VectorEqual(a Vector, b Vector) bool {
	for i := 0; i < len(a); i++ {
		if !utils.ComplexNearlyEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}
func (vec *Vector) Reverse() {
	for i := 0; i < len(*vec)/2; i++ {
		tmp := (*vec)[i]
		idx := len(*vec) - i - 1
		(*vec)[i] = (*vec)[idx]
		(*vec)[idx] = tmp
	}
}

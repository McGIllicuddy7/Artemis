package utils

import (
	"cmp"
	"fmt"
	"math"
	"math/cmplx"
	"math/rand"
)

func partition[T any](slice []T, cmp_func func(T, T) int, start int, end int) int {
	pivot := slice[end]
	p := start - 1
	for i := start; i < end; i++ {
		if cmp_func(slice[i], pivot) < 0 {
			p++
			tmp := slice[i]
			slice[i] = slice[p]
			slice[p] = tmp
		}
	}
	tmp := slice[p+1]
	slice[p+1] = slice[end]
	slice[end] = tmp
	return p + 1
}
func quick_sort[T any](slice []T, cmp_func func(T, T) int, start int, end int) {
	if start < end {
		part := partition[T](slice, cmp_func, start, end)
		quick_sort[T](slice, cmp_func, start, part-1)
		quick_sort[T](slice, cmp_func, part+1, end)
	}
}
func SortInplace[T any](slice []T, cmp_func func(T, T) int) {
	quick_sort[T](slice, cmp_func, 0, len(slice)-1)
}
func SortCopy[T any](slice []T, cmp_func func(T, T) int) []T {
	out := make([]T, len(slice))
	copy(out, slice)
	SortInplace[T](out, cmp_func)
	return out
}

func TrivSortInplace[T cmp.Ordered](slice []T) {
	SortInplace[T](slice, cmp.Compare[T])
}
func TrivSortCopy[T cmp.Ordered](slice []T) {
	SortCopy[T](slice, cmp.Compare[T])
}
func IsSorted[T any](slice []T, cmp_func func(T, T) int) bool {
	for i := 0; i < len(slice)-1; i++ {
		if cmp_func(slice[i], slice[i+1]) > 0 {
			return false
		}
	}
	return true
}
func TrivIsSorted[T cmp.Ordered](slice []T) bool {
	for i := 0; i < len(slice)-1; i++ {
		if cmp.Compare[T](slice[i], slice[i+1]) > 0 {
			return false
		}
	}
	return true
}
func FormatFloat64(f float64) string {
	if math.Floor(f) == f {
		return fmt.Sprintf("%d", int(f))
	} else {
		return fmt.Sprintf("%f", f)
	}
}
func FormatComplex(c complex128) string {
	out := ""
	if real(c) != 0 {
		out += FormatFloat64(real(c))
		if imag(c) != 0 {
			out += "+"
		}
	}
	if imag(c) != 0 {
		out += FormatFloat64(imag(c))
		out += "i"
	} else {
		if real(c) == 0 {
			out = "0"
		}
	}
	return out
}
func normalize_strlen(str string, length int) string {
	for len(str) < length {
		str += " "
	}
	return str
}
func NormalizeStrlens(in []string) []string {
	max := 0
	for i := 0; i < len(in); i++ {
		if len(in[i]) > max {
			max = len(in[i])
		}
	}
	out := make([]string, len(in))
	for i := 0; i < len(in); i++ {
		out[i] = normalize_strlen(in[i], max)
	}
	return out
}
func RandomComplex() complex128 {
	const max = 10000000
	return complex(float64(rand.Int63()%(max*2)-max), float64(rand.Int63()%(max*2)-max))
}
func Contains[T comparable](slice []T, value T) bool {
	for i := 0; i < len(slice); i++ {
		if slice[i] == value {
			return true
		}
	}
	return false
}
func ComplexNearlyEqual(v0 complex128, v1 complex128) bool {
	v := v0 - v1
	return cmplx.Abs(v) < 0.1
}

package main

import (
	"fmt"
	"os"
	"strconv"
)

type matrix_t struct {
	data   []int
	height int
	width  int
}

func (this *matrix_t) get(x int, y int) int {
	return this.data[y*this.width+x]
}
func (this *matrix_t) set(x int, y int, v int) {
	this.data[y*this.width+x] = v
}
func (this *matrix_t) to_string() string {
	out := ""
	for j := 0; j < this.height; j++ {
		for i := 0; i < this.width; i++ {
			out += strconv.Itoa(this.get(i, j)) + " "
		}
		out += "\n"
	}
	return out
}
func determinant(matrix matrix_t) float64 {
	if matrix.height != matrix.width {
		fmt.Fprint(os.Stderr, "error cannot take determinant of non square matrix\n")
		os.Exit(1)
	}
	if matrix.width == 2 {
		a := float64(matrix.get(0, 0))
		b := float64(matrix.get(1, 0))
		c := float64(matrix.get(0, 1))
		d := float64(matrix.get(1, 1))
		return a*d - b*c
	}
	out := float64(0)
	for i := 0; i < matrix.width; i++ {
		var m matrix_t
		m.height = matrix.height - 1
		m.width = matrix.width - 1
		m.data = make([]int, m.width*m.width)
		for x := 0; x < matrix.width; x++ {
			for y := 0; y < matrix.height; y++ {
				xv := x
				yv := y
				if x == i || y == i {
					continue
				}
				if x > i {
					xv--
				}
				if y > i {
					yv--
				}
				m.set(xv, yv, matrix.get(x, y))
			}
		}
		mlt := float64(1)
		if i%2 == 1 {
			mlt = -1
		}
		out += mlt * determinant(m)
	}
	return out
}

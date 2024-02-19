package main

import "strconv"

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
			out += strconv.Itoa(i) + " "
		}
		out += "\n"
	}
	return out
}

package main

import (
	LA "artemis/LA"
	"math/rand"
)

type Quiver struct {
	points     []vertex_t
	num_points int
}

func (q *Quiver) MutateAt(point int) Quiver {
	return Quiver{mutate(q.points, q.num_points, point), q.num_points}
}
func (q *Quiver) ToMatrix() LA.Matrix {
	return make_matrix_from_quiver(q.points, q.num_points)
}
func AddCost(quiv Quiver) int {
	cost := 0
	for i := 0; i < quiv.num_points; i++ {
		for j := 0; j < quiv.num_points; j++ {
			if quiv.points[i].edges[j] > 0 {
				cost += quiv.points[i].edges[j]
			}
		}
	}
	return cost
}
func MaxCost(quiv Quiver) int {
	cost := 0
	for i := 0; i < quiv.num_points; i++ {
		for j := 0; j < quiv.num_points; j++ {
			if quiv.points[i].edges[j] > cost {
				cost = quiv.points[i].edges[j]
			}
		}
	}
	return cost
}
func Cost(quiv Quiver) int {
	return AddCost(quiv)
}
func RandomQuiver(dim int, max int) Quiver {
	out := Quiver{make([]vertex_t, 64), dim}
	for i := 0; i < dim; i++ {
		out.points[i].location.X = rand.Float32()*500 + 250
		out.points[i].location.Y = rand.Float32()*500 + 250
		for j := 0; j < dim; j++ {
			if j == i {
				continue
			}
			r := rand.Int()%(max*2) - max
			out.points[i].edges[j] = r
			out.points[j].edges[i] = -r
		}
	}
	return out
}

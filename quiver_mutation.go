package main

import (
	fr "artemis/Fractions"
	La "artemis/LA"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// structure of a vertex. the edges are stored just as integers
type vertex_t struct {
	edges    [MAX_VERTICES]int
	location rl.Vector2
}

// stores a number to be added to a pair of edges in a mutation
// so that mutations don't self interfere
type mutation_event_t struct {
	start int
	end   int
	value int
}

func returns_to_self(idx int, target_idx int, visited []bool, quiver []vertex_t) bool {
	if visited[idx] {
		return false
	}
	visited[idx] = true
	if quiver[idx].edges[target_idx] > 0 {
		return true
	}
	for i := 0; i < len(quiver[idx].edges); i++ {
		if !visited[i] && quiver[idx].edges[i] > 0 {
			tmp := returns_to_self(i, target_idx, visited, quiver)
			if tmp {
				return true
			}
		}
	}
	return false
}
func is_cyclic(quiver []vertex_t) bool {
	for i := 0; i < len(quiver); i++ {
		visited := make([]bool, len(quiver))
		if returns_to_self(i, i, visited, quiver) {
			return true
		}
	}
	return false
}

// to make the matrix graph
func make_matrix_from_quiver(quiver []vertex_t, num int) La.Matrix {
	out := La.ZeroMatrix(num, num)
	for i := 0; i < num; i++ {
		for j := 0; j < num; j++ {
			out.Set(i, j, fr.FromInt(quiver[i].edges[j]))
		}
	}
	return out
}

func create_vertex(location rl.Vector2, vertices *[]vertex_t, num_vertices *int) {
	if *num_vertices+1 > MAX_VERTICES {
		return
	}
	var v vertex_t
	for i := 0; i < MAX_VERTICES; i++ {
		v.edges[i] = 0
	}
	v.location = location
	(*vertices)[*num_vertices] = v
	(*num_vertices)++
}

// links a pair of vertices by incrementing the edges for a and decrementing them for b
func vertex_link(vertices []vertex_t, a int, b int) {
	vertices[a].edges[b]++
	vertices[b].edges[a]--
}

/*
mutation: changing quivers to different quivers
mutate at vertex y:
step 1: for every x->y->z add an arrow x->z do this for every path in the original you don't have to do it recursively, number of arrows from x->y times number of arrows from y->z is how many arrows you add, mutation is local its the stuff directly connected to y
step 2: reverse all arrows of the form x->y or x<-y once again this is local
step 3: if you end up witn arrows pointing in opposite directions each pair of opposites is deleted, this is I BELIEVE global
*/
func new_mutation_event(start int, end int, value int) mutation_event_t {
	var out mutation_event_t
	out.start = start
	out.end = end
	out.value = value
	return out
}
func mutate_inline(vertices []vertex_t, num_vertices int, a int) {
	edges := vertices[a].edges
	mutations := make([]mutation_event_t, 4096)
	eventque_len := 0
	// step one
	for i := 0; i < num_vertices; i++ {
		if i == a {
			continue
		}
		for j := 0; j < num_vertices; j++ {
			if vertices[i].edges[a] > 0 {
				mutations[eventque_len] = new_mutation_event(i, j, edges[j]*vertices[i].edges[a])
				eventque_len++
			}
		}
	}
	for i := 0; i < eventque_len; i++ {
		vertices[mutations[i].start].edges[mutations[i].end] += mutations[i].value
		vertices[mutations[i].end].edges[mutations[i].start] -= mutations[i].value
	}
	for i := 0; i < num_vertices; i++ {
		tmp1 := vertices[i].edges[a]
		tmp2 := vertices[a].edges[i]
		vertices[a].edges[i] = tmp1
		vertices[i].edges[a] = tmp2
	}
}
func mutate(in_vertices []vertex_t, num_vertices int, a int) []vertex_t {
	vertices := make([]vertex_t, len(in_vertices))
	copy(vertices, in_vertices)
	edges := vertices[a].edges
	mutations := make([]mutation_event_t, 4096)
	eventque_len := 0
	// step one
	for i := 0; i < num_vertices; i++ {
		if i == a {
			continue
		}
		for j := 0; j < num_vertices; j++ {
			if vertices[i].edges[a] > 0 {
				mutations[eventque_len] = new_mutation_event(i, j, edges[j]*vertices[i].edges[a])
				eventque_len++
			}
		}
	}
	for i := 0; i < eventque_len; i++ {
		vertices[mutations[i].start].edges[mutations[i].end] += mutations[i].value
		vertices[mutations[i].end].edges[mutations[i].start] -= mutations[i].value
	}
	for i := 0; i < num_vertices; i++ {
		tmp1 := vertices[i].edges[a]
		tmp2 := vertices[a].edges[i]
		vertices[a].edges[i] = tmp1
		vertices[i].edges[a] = tmp2
	}
	return vertices
}
func Sanitize(in_vertices []vertex_t, num_vertices int) {
	for i := 0; i < len(in_vertices); i++ {
		for j := 0; j < len(in_vertices); j++ {
			if i >= num_vertices || j >= num_vertices {
				in_vertices[i].edges[j] = 0
			}
		}
	}
}

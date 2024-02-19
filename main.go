package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const SCREEN_HEIGHT = 900
const SCREEN_WIDTH = 900
const MAX_VERTICES = 64
const MAX_EDGES = 128
const CMD_LEN = 128

type vertex_t struct {
	edges     [MAX_VERTICES]int
	num_edges int
	location  rl.Vector2
}
type mutation_event_t struct {
	start int
	end   int
	value int
}

func make_matrix_from_quiver(quiver []vertex_t, num int) matrix_t {
	out := matrix_t{make([]int, num*num), num, num}
	for i := 0; i < num; i++ {
		for j := 0; j < num; j++ {
			out.set(i, j, quiver[i].edges[j])
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
func mutate(vertices []vertex_t, num_vertices int, a int) {
	edges := vertices[a].edges
	mutations := make([]mutation_event_t, 4096)
	var eventque_len int
	// step one
	for i := 0; i < num_vertices; i++ {
		if i == a {
			continue
		}
		for j := 0; j < num_vertices; j++ {
			if j == a || j == i {
				continue
			}
			if vertices[i].edges[a] > 0 {
				mutations[eventque_len] = new_mutation_event(i, j, edges[j]*vertices[i].edges[a])
				eventque_len++
			} else {
				mutations[eventque_len] = new_mutation_event(j, i, edges[i]*vertices[j].edges[a])
				eventque_len++
			}
		}
	}
	for i := 0; i < eventque_len; i++ {
		vertices[mutations[i].start].edges[mutations[i].end] += mutations[i].value
	}
	//step 2
	for i := 0; i < num_vertices; i++ {
		if edges[i] > 0 || vertices[i].edges[a] > 0 {
			tmp := edges[i]
			tmp2 := vertices[i].edges[a]
			edges[i] = tmp2
			vertices[i].edges[a] = tmp

		}
	}
	for i := 0; i < num_vertices-1; i++ {
		for j := i + 1; j < num_vertices; j++ {
			tmp_i := vertices[i].edges[j]
			tmp_j := vertices[j].edges[i]
			vertices[i].edges[j] = tmp_i - tmp_j
			vertices[j].edges[i] = tmp_j - tmp_i
		}
	}
}
func main() {
	cmd := ""
	vertices := make([]vertex_t, MAX_VERTICES)
	num_vertices := 0
	rl.SetTraceLogLevel(rl.LogError)
	rl.InitWindow(SCREEN_WIDTH, SCREEN_HEIGHT, "Artemis")
	rl.SetTargetFPS(120)
	for !rl.WindowShouldClose() {
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) && rl.GetMousePosition().Y < SCREEN_HEIGHT-20 {
			create_vertex(rl.GetMousePosition(), &vertices, &num_vertices)
		}
		if rl.IsKeyPressed(rl.KeyEnter) {
			cmd_execute(&cmd, &vertices, &num_vertices)
		} else {
			cmd_parse(&cmd)
		}
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		mat := make_matrix_from_quiver(vertices, num_vertices)
		str := mat.to_string()
		rl.DrawText(str, 100, 100, 12, rl.White)
		textbox := rl.NewColor(60, 60, 60, 255)
		rl.DrawRectangle(0, SCREEN_HEIGHT-20, SCREEN_WIDTH, 20, textbox)
		rl.DrawText(string(cmd), 20, SCREEN_HEIGHT-20, 14, rl.Black)
		for i := 0; i < num_vertices; i++ {
			rl.DrawCircleV(vertices[i].location, 5, rl.Gray)
			for j := 0; j < num_vertices; j++ {
				if vertices[i].edges[j] > 0 && i != j {
					tmp_buff := make([]byte, 0)
					tmp_buff = append(tmp_buff, []byte(fmt.Sprintf("%d", vertices[i].edges[j]))...)
					dv := rl.Vector2Subtract(vertices[j].location, vertices[i].location)
					theta := math.Atan(float64(dv.Y / dv.X))
					theta0 := theta + math.Pi/12
					dx0 := float32(math.Cos(theta0) * 25)
					dy0 := float32(math.Sin(theta0) * 25)
					theta1 := theta - math.Pi/12
					dx1 := float32(math.Cos(theta1) * 25)
					dy1 := float32(math.Sin(theta1) * 25)
					if vertices[j].location.X < vertices[i].location.X {
						dx0 *= -1
						dy0 *= -1
						dx1 *= -1
						dy1 *= -1
					}
					l0 := rl.Vector2Subtract(vertices[j].location, rl.NewVector2(dx0, dy0))
					l1 := rl.Vector2Subtract(vertices[j].location, rl.NewVector2(dx1, dy1))
					rl.DrawLineEx(vertices[i].location, vertices[j].location, 2, rl.Red)
					rl.DrawLineEx(vertices[j].location, l0, 2, rl.Red)
					rl.DrawLineEx(vertices[j].location, l1, 2, rl.Red)
					theta2 := theta - math.Pi/2
					v := rl.NewVector2(float32(math.Cos(theta2)*15), float32(math.Cos(theta2)*15))
					avg := rl.Vector2Scale(rl.Vector2Add(vertices[j].location, vertices[i].location), 0.5)
					v = rl.Vector2Add(v, avg)
					rl.DrawText(string(tmp_buff), int32(v.X), int32(v.Y), 12, rl.White)
				}
			}
			l := fmt.Sprintf("%d", i)
			rl.DrawText(l, int32(vertices[i].location.X+SCREEN_WIDTH/100.0), int32(vertices[i].location.Y-SCREEN_HEIGHT/100.0), 14, rl.White)
		}
		rl.EndDrawing()
	}
}

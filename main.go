package main

import (
	fr "artemis/Fractions"
	La "artemis/LA"
	"artemis/utils"
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const SCREEN_HEIGHT = 900
const SCREEN_WIDTH = 900
const MAX_VERTICES = 64

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
func main() {
	/*
		mat := matrix_t{make([]int, 3, 3), 3, 3}
		mat.data = []int{0, -1, 1, 1, 0, -1, -1, 1, 0}
		println(mat.to_string())
		det := determinant(mat)
		println()
		println(det)
	*/
	cmd := ""
	vertices := make([]vertex_t, MAX_VERTICES)
	num_vertices := 0
	//graphics initalization
	rl.SetTraceLogLevel(rl.LogError)
	rl.InitWindow(SCREEN_WIDTH, SCREEN_HEIGHT, "Artemis")
	rl.SetTargetFPS(120)
	mat := make_matrix_from_quiver(vertices, num_vertices)
	eigens := mat.EigenValues()
	_ = mat.EigenVectors()
	//program loop
	for !rl.WindowShouldClose() {
		//input handling
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) && rl.GetMousePosition().Y < SCREEN_HEIGHT-20 {
			create_vertex(rl.GetMousePosition(), &vertices, &num_vertices)
			mat = make_matrix_from_quiver(vertices, num_vertices)
			eigens = mat.EigenValues()
			_ = mat.EigenVectors()
		}
		if rl.IsKeyPressed(rl.KeyEnter) {
			cmd_execute(&cmd, &vertices, &num_vertices)
			mat = make_matrix_from_quiver(vertices, num_vertices)
			eigens = mat.EigenValues()
			_ = mat.EigenVectors()
		} else {
			cmd_parse(&cmd)
		}
		//rendering
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		//rendering the matrix
		str := mat.ToString()
		rl.DrawText(str, 20, 800, 16, rl.White)
		//rendering the terminal
		textbox := rl.NewColor(60, 60, 60, 255)
		rl.DrawRectangle(0, SCREEN_HEIGHT-20, SCREEN_WIDTH, 20, textbox)
		rl.DrawText(string(cmd), 20, SCREEN_HEIGHT-20, 14, rl.Black)
		arrow_color := rl.Red
		//rendering the points
		for i := 0; i < num_vertices; i++ {
			rl.DrawCircleV(vertices[i].location, 5, rl.Gray)
			for j := 0; j < num_vertices; j++ {
				if vertices[i].edges[j] > 0 && i != j {
					//this nightmare of trig renders arrows I don't understand it either
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
					rl.DrawLineEx(vertices[i].location, vertices[j].location, 2, arrow_color)
					rl.DrawLineEx(vertices[j].location, l0, 2, arrow_color)
					rl.DrawLineEx(vertices[j].location, l1, 2, arrow_color)
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
		cyc_msg := "is not cyclic"
		if is_cyclic(vertices) {
			cyc_msg = "is cyclic"
		}
		rl.DrawText(cyc_msg, 20, 20, 16, rl.RayWhite)
		rl.DrawText(fmt.Sprintf("determinant is %s", mat.Determinant().ToString()), 20, 40, 16, rl.RayWhite)
		rl.DrawText("eigen values: ", 600, 80, 16, rl.RayWhite)
		for i := 0; i < len(eigens); i++ {
			msg := utils.FormatComplex(eigens[i])
			if i != len(eigens)-1 {
				msg += ", "
			}
			rl.DrawText(msg, int32(600), int32(32+80+i*32), 16, rl.RayWhite)
		}
		/*
			for i := 0; i < len(vecs); i++ {
				for j := 0; j < len(vecs[i]); j++ {
					rl.DrawText(fmt.Sprintf("%s ", utils.FormatComplex(eigens[i])), int32(300+80*i), int32(630+30*j), 20, rl.RayWhite)
				}
			}
		*/
		rl.EndDrawing()
	}
}

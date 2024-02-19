package main

import (
	"fmt"
	"math"
	"os"
	"strconv"

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
type matrix_t struct {
	data   []int
	height int
	width  int
}

func make_matrix_from_quiver() {

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
func cmd_append_character(cmd *string, c string) {
	*cmd += c
}
func strip_white_space_left(str string) string {
	cmd_v := make([]byte, len(str))
	for i := 0; i < len(str); i++ {
		cmd_v[i] = str[i]
	}
	l := len(str)
	for i := 0; i < l; i++ {
		if str[i] != ' ' {
			break
		} else {
			cmd_v = cmd_v[1:]
		}
	}
	return string(cmd_v)
}
func next_int(str *string) int {
	*str = strip_white_space_left(*str)
	buff := ""
	counter := 0
	for i := 0; i < len(*str); i++ {
		if (*str)[i] == ' ' {
			break
		}
		buff += string((*str)[i])
		counter++
	}
	(*str) = (*str)[counter:]
	v, _ := strconv.Atoi(string(buff))
	return v
}
func vertex_link(vertices []vertex_t, a int, b int) {
	vertices[a].edges[b]++
	vertices[b].edges[a]--
}

type mutation_event_t struct {
	start int
	end   int
	value int
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
func cmd_execute(cmd *string, vertices *[]vertex_t, num_vertices *int) {
	old_vertices := make([]vertex_t, MAX_VERTICES)
	old_num_vertices := 0
	if len(*cmd) < 3 {
		return
	}
	cmd_v := strip_white_space_left(*cmd)
	if cmd_v[0] == 'h' {
		cmd_v = (fmt.Sprintf(string(cmd_v), "h for help, ln a b to create an edge from a to b, lnr a b to remove an edge from a to b, mut a to morph a,shift+delete to clear terminal exit to exit, "))
	}
	if cmd_v[0] == 'e' && cmd_v[1] == 'x' && cmd_v[2] == 'i' && cmd_v[3] == 't' {
		os.Exit(0)
	}
	if cmd_v[0] == 'l' && cmd_v[1] == 'n' && cmd_v[2] != 'r' {
		for i := 0; i < len(old_vertices); i++ {
			old_vertices[i] = (*vertices)[i]
		}
		old_num_vertices = *num_vertices
		cmd_v = cmd_v[2:]
		cmd_v = strip_white_space_left(cmd_v)
		a := next_int(&cmd_v)
		cmd_v = strip_white_space_left(cmd_v)
		b := next_int(&cmd_v)
		vertex_link(*vertices, a, b)
		*cmd = ""
		goto done
	}
	if cmd_v[0] == 'l' && cmd_v[1] == 'n' && cmd_v[2] == 'r' {
		for i := 0; i < len(old_vertices); i++ {
			old_vertices[i] = (*vertices)[i]
		}
		old_num_vertices = *num_vertices
		cmd_v = cmd_v[2:]
		cmd_v = strip_white_space_left(cmd_v)
		a := next_int(&cmd_v)
		cmd_v = strip_white_space_left(cmd_v)
		b := next_int(&cmd_v)
		vertex_link(*vertices, b, a)
		*cmd = ""
		goto done
	}
	if cmd_v[0] == 'm' && cmd_v[1] == 'u' && cmd_v[2] == 't' {
		for i := 0; i < len(old_vertices); i++ {
			old_vertices[i] = (*vertices)[i]
		}
		old_num_vertices = *num_vertices
		cmd_v = cmd_v[3:]
		cmd_v = strip_white_space_left(cmd_v)
		a := next_int(&cmd_v)
		mutate(*vertices, *num_vertices, a)
		*cmd = ""
		goto done
	}
	if cmd_v[0] == 'u' && cmd_v[1] == 'n' && cmd_v[2] == 'd' && cmd_v[3] == 'o' {
		for i := 0; i < len(old_vertices); i++ {
			old_vertices[i] = (*vertices)[i]
		}
		*num_vertices = old_num_vertices
		*cmd = ""
		goto done
	}
done:
	return
}
func cmd_parse(cmd *string) {
	if (rl.IsKeyDown(rl.KeyDelete) || rl.IsKeyPressed(rl.KeyBackspace)) && rl.IsKeyDown(rl.KeyLeftShift) {
		*cmd = ""
		return
	}
	if rl.IsKeyPressed(rl.KeyDelete) || rl.IsKeyPressed(rl.KeyBackspace) && len(*cmd) > 0 {
		*cmd = (*cmd)[:len(*cmd)-1]
		goto done
	}
	if rl.IsKeyPressed(rl.KeySpace) {
		cmd_append_character(cmd, " ")
		goto done
	}
	if rl.IsKeyPressed(rl.KeyA) {
		cmd_append_character(cmd, "a")
		goto done
	}
	if rl.IsKeyPressed(rl.KeyB) {
		cmd_append_character(cmd, "b")
		goto done
	}
	if rl.IsKeyPressed(rl.KeyC) {
		cmd_append_character(cmd, "c")
		goto done
	}
	if rl.IsKeyPressed(rl.KeyD) {
		cmd_append_character(cmd, "d")
		goto done
	}
	if rl.IsKeyPressed(rl.KeyE) {
		cmd_append_character(cmd, "e")
		goto done
	}
	if rl.IsKeyPressed(rl.KeyF) {
		cmd_append_character(cmd, "f")
		goto done
	}
	if rl.IsKeyPressed(rl.KeyG) {
		cmd_append_character(cmd, "g")
		goto done
	}
	if rl.IsKeyPressed(rl.KeyH) {
		cmd_append_character(cmd, "h")
		goto done
	}
	if rl.IsKeyPressed(rl.KeyI) {
		cmd_append_character(cmd, "i")
		goto done
	}
	if rl.IsKeyPressed(rl.KeyJ) {
		cmd_append_character(cmd, "j")
		goto done
	}
	if rl.IsKeyPressed(rl.KeyK) {
		cmd_append_character(cmd, "k")
		goto done
	}
	if rl.IsKeyPressed(rl.KeyL) {
		cmd_append_character(cmd, "l")
		goto done
	}
	if rl.IsKeyPressed(rl.KeyM) {
		cmd_append_character(cmd, "m")
		goto done
	}
	if rl.IsKeyPressed(rl.KeyN) {
		cmd_append_character(cmd, "n")
		goto done
	}
	if rl.IsKeyPressed(rl.KeyO) {
		cmd_append_character(cmd, "o")
		goto done
	}
	if rl.IsKeyPressed(rl.KeyP) {
		cmd_append_character(cmd, "p")
		goto done
	}
	if rl.IsKeyPressed(rl.KeyQ) {
		cmd_append_character(cmd, "q")
		goto done
	}
	if rl.IsKeyPressed(rl.KeyR) {
		cmd_append_character(cmd, "r")
		goto done
	}
	if rl.IsKeyPressed(rl.KeyS) {
		cmd_append_character(cmd, "s")
		goto done
	}
	if rl.IsKeyPressed(rl.KeyT) {
		cmd_append_character(cmd, "t")
		goto done
	}
	if rl.IsKeyPressed(rl.KeyU) {
		cmd_append_character(cmd, "u")
		goto done
	}
	if rl.IsKeyPressed(rl.KeyV) {
		cmd_append_character(cmd, "v")
		goto done
	}
	if rl.IsKeyPressed(rl.KeyW) {
		cmd_append_character(cmd, "w")
		goto done
	}
	if rl.IsKeyPressed(rl.KeyX) {
		cmd_append_character(cmd, "x")
		goto done
	}
	if rl.IsKeyPressed(rl.KeyY) {
		cmd_append_character(cmd, "y")
		goto done
	}
	if rl.IsKeyPressed(rl.KeyZ) {
		cmd_append_character(cmd, "z")
		goto done
	}
	if rl.IsKeyPressed(rl.KeyZero) {
		cmd_append_character(cmd, "0")
		goto done
	}
	if rl.IsKeyPressed(rl.KeyOne) {
		cmd_append_character(cmd, "1")
		goto done
	}
	if rl.IsKeyPressed(rl.KeyTwo) {
		cmd_append_character(cmd, "2")
		goto done
	}
	if rl.IsKeyPressed(rl.KeyThree) {
		cmd_append_character(cmd, "3")
		goto done
	}
	if rl.IsKeyPressed(rl.KeyFour) {
		cmd_append_character(cmd, "4")
		goto done
	}
	if rl.IsKeyPressed(rl.KeyFive) {
		cmd_append_character(cmd, "5")
		goto done
	}
	if rl.IsKeyPressed(rl.KeySix) {
		cmd_append_character(cmd, "6")
		goto done
	}
	if rl.IsKeyPressed(rl.KeySeven) {
		cmd_append_character(cmd, "7")
		goto done
	}
	if rl.IsKeyPressed(rl.KeyEight) {
		cmd_append_character(cmd, "8")
		goto done
	}
	if rl.IsKeyPressed(rl.KeyNine) {
		cmd_append_character(cmd, "9")
		goto done
	}
done:
	return
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
		matrix := make([]int, MAX_VERTICES*MAX_VERTICES)
		for i := 0; i < num_vertices; i++ {
			for j := 0; j < num_vertices; j++ {
				matrix[i*MAX_VERTICES+j] = vertices[i].edges[j]
			}
		}
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		textbox := rl.NewColor(60, 60, 60, 255)
		rl.DrawRectangle(0, SCREEN_HEIGHT-20, SCREEN_WIDTH, 20, textbox)
		rl.DrawText(string(cmd), 20, SCREEN_HEIGHT-20, 14, rl.Black)
		for i := 0; i < num_vertices; i++ {
			rl.DrawCircleV(vertices[i].location, 5, rl.Gray)
			for j := 0; j < num_vertices; j++ {
				if vertices[i].edges[j] > 0 {
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

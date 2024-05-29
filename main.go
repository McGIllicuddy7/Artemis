package main

import (
	La "artemis/LA"
	"artemis/utils"
	"fmt"
	"math"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const SCREEN_HEIGHT = 900
const SCREEN_WIDTH = 900
const MAX_VERTICES = 64
const OFFSET = 0

func main() {
	cmd := ""
	Q, err := LoadQuiver("q.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load %s\n", err.Error())
		Q = RandomQuiver(4, 3)
	}
	vertices := Q.points
	num_vertices := Q.num_points
	println(num_vertices)
	//graphics initalization
	rl.SetTraceLogLevel(rl.LogError)
	rl.InitWindow(SCREEN_WIDTH, SCREEN_HEIGHT, "Artemis")
	rl.SetTargetFPS(120)
	mat := make_matrix_from_quiver(vertices, num_vertices)
	eigens := mat.EigenValues()
	_ = mat.EigenVectors()
	t_error := 0.0
	poly_str := ""
	mlt_str := ""
	//program loop
	for !rl.WindowShouldClose() {
		//input handling
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) && rl.GetMousePosition().Y < SCREEN_HEIGHT-20 {
			create_vertex(rl.GetMousePosition(), &vertices, &num_vertices)
			mat = make_matrix_from_quiver(vertices, num_vertices)
			//eigens = mat.EigenValues()
			//_ = mat.EigenVectors()
			//poly := mat.ToEigenMatrix().CharacteristicPolynomial()
			//poly_str = poly.ToString()
			mlt_str = ""
		}
		if rl.IsKeyPressed(rl.KeyEnter) {
			old := make_matrix_from_quiver(vertices, num_vertices)
			err := cmd_execute(&cmd, &vertices, &num_vertices)
			Sanitize(vertices, num_vertices)
			mat = make_matrix_from_quiver(vertices, num_vertices)
			//eigens = mat.EigenValues()
			//_ = mat.EigenVectors()
			if !err {
				t_error = 1.0
				cmd = ""
			}
			//poly := mat.ToEigenMatrix().CharacteristicPolynomial()
			//poly_str = poly.ToString()
			if mat.NumCols() == old.NumCols() && mat.NumRows() == old.NumRows() {
				_, tmp := La.MatrixPairRowReduce(old, mat)
				mlt_str = tmp.ToString()
			} else {
				mlt_str = ""
			}
		} else {
			cmd_parse(&cmd)
		}
		//rendering
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		//rendering the matrix
		str := mat.ToString()
		rl.DrawText(str, 20, 800, 16, rl.White)
		rl.DrawText(mlt_str, 200, 800, 16, rl.White)
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
			l := fmt.Sprintf("%d", i+OFFSET)
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
		if t_error > 0.01 {
			t_error -= float64(rl.GetFrameTime())
			rl.DrawText("error: improper input", 600, 40, 16, rl.Red)
		}
		rl.DrawText("char poly: "+poly_str, 400, 700, 16, rl.White)
		quiv := Quiver{vertices, num_vertices}
		c := fmt.Sprintf("%d", Cost(quiv))
		rl.DrawText("cost: "+c, 700, 800, 16, rl.RayWhite)
		for i := 0; i < num_vertices; i++ {
			m := fmt.Sprintf("cost after mutation at %d: %d", i, Cost(quiv.MutateAt(i)))
			rl.DrawText(m, 400, int32(800+i*20), 16, rl.RayWhite)
		}
		rl.EndDrawing()
	}
	//SaveQuiver("q.txt", Quiver{vertices, num_vertices})
}

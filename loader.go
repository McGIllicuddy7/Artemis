package main

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func ParseQuiver(quiver string) (Quiver, error) {
	out := Quiver{make([]vertex_t, 64), 0}
	rowstrs := strings.SplitAfter(string(quiver), "\n")
	rows := make([][]string, 0)
	for i := 0; i < len(rowstrs); i++ {
		tmp := strings.Split(rowstrs[i], " ")
		row := make([]string, 0)
		for j := 0; j < len(tmp); j++ {
			if !strings.ContainsAny(tmp[j], "1234567890") || tmp[j] == "" {
				continue
			}
			row = append(row, strings.TrimSpace(tmp[j]))
			//println("<", row[len(row)-1], ">")
		}
		rows = append(rows, row)
	}
	out.num_points = len(rows[0])
	for i := 0; i < out.num_points; i++ {
		for j := 0; j < out.num_points; j++ {
			a, err := strconv.Atoi(rows[i][j])
			if err != nil {
				return out, errors.New(err.Error() + fmt.Sprintf(" num points: %d", out.num_points))
			}
			out.points[j].edges[i] = a
		}
	}
	for i := 0; i < len(out.points); i++ {
		l := float32(len(out.points)) - 1
		if l == 0 {
			l = 1
		}
		loc := rl.Vector2Zero()
		loc.X = float32(250*math.Cos(float64(float32(i)*100*math.Pi/l)) + 500)
		loc.Y = float32(250*math.Sin(float64(float32(i)*100*math.Pi/l)) + 500)
		out.points[i].location = loc
	}
	return out, nil
}
func LoadQuiver(filename string) (Quiver, error) {
	out := Quiver{make([]vertex_t, 64), 0}
	f, err := os.Open(filename)
	if err != nil {
		return out, errors.New("failed to open")
	}
	strs, err := io.ReadAll(f)
	if err != nil {
		return out, errors.New("failed to read")
	}
	f.Close()
	out, err = ParseQuiver(string(strs))
	if err != nil {
		return out, err
	}
	if out.num_points == 0 {
		return out, errors.New("no points")
	}
	return out, err
}
func WriteQuiver(file io.Writer, q Quiver) {
	m := q.ToMatrix()
	_, _ = io.WriteString(file, m.ToString())
}
func SaveQuiver(filename string, q Quiver) {
	f, err := os.Create(filename)
	if err != nil {
		return
	}
	m := q.ToMatrix()
	_, _ = io.WriteString(f, m.ToString())
}

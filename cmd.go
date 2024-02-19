package main

import (
	"fmt"
	"os"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

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
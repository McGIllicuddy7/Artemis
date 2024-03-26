package autopsy

import (
	"fmt"
	"sync/atomic"
)

type autopsy struct {
	info []string
}

var aut autopsy
var key atomic.Bool

func Init() {
	aut = autopsy{make([]string, 0)}
	key.Store(false)
}
func key_wait() {
	for key.Load() == true {

	}
	key.Store(true)
}
func key_release() {
	key.Store(false)
}
func Reset() {
	key_wait()
	aut = autopsy{make([]string, 0)}
	key_release()
}
func Store(msg string) {
	key_wait()
	aut.info = append(aut.info, msg)
	key_release()
}
func Dump() {
	for i := 0; i < len(aut.info); i++ {
		fmt.Printf("%s\n", aut.info[i])
	}
	clear(aut.info)
}

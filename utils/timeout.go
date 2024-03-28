package utils

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

var should_die atomic.Bool
var exisiting_timeout = false

type timeoutStack struct {
	t     time.Time
	stack []string
}

var tstack timeoutStack

func TimeoutPush(name string) {
	tstack.stack = append(tstack.stack, name)
}
func TimeoutPop() {
	tstack.stack = tstack.stack[0 : len(tstack.stack)-1]
}
func timeout() {
	for {
		if time.Since(tstack.t).Milliseconds() > 10000 {
			fmt.Printf("timed out during %s", tstack.stack[len(tstack.stack)-1])
			os.Exit(0)
		}
		if should_die.Load() {
			should_die.Store(false)
			return
		}
	}
}
func tkill() {
	should_die.Store(true)
	for should_die.Load() {

	}
}
func TimeoutStart() {
	should_die.Store(false)
	tstack.t = time.Now()
	tstack.stack = make([]string, 0)
	if exisiting_timeout {
		tkill()
	}
	go timeout()
	exisiting_timeout = true
}

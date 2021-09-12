// varcs generates vectorized arcs and is based on the example here:
// https://gitlab.com/gcmc/gcmc/blob/master/example/varcs.gcmc
//
// Usage:
//   go run examples/varcs/main.go > varcs.gcode
package main

import (
	"fmt"

	. "github.com/gmlewis/go-gcode/gcode"
)

const (
	safeHeight   = 1
	cuttingDepth = -1
)

func main() {
	g := gcmc()
	fmt.Printf("%v\n", g)
}

func gcmc() *GCode {
	g := New()
	g.Prologue = true
	g.Epilogue = true

	g.Feedrate(600)

	return g
}

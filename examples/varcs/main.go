// varcs generates vectorized arcs and is based on the example here:
// https://gitlab.com/gcmc/gcmc/blob/master/example/varcs.gcmc
//
// Usage:
//   go run examples/varcs/main.go > varcs.gcode
package main

import (
	"fmt"

	. "github.com/gmlewis/go-gcode/gcode"
	"github.com/gmlewis/go-gcode/utils"
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
	radius := 15.0
	flushItAt(g, XY(-20, 100), utils.VArcCW(XY(10, 15), radius, nil))
	flushItAt(g, XY(-20, 60), utils.VArcCW(XY(10, 15), -radius, nil))
	flushItAt(g, XY(20, 100), utils.VArcCCW(XY(10, 15), radius, nil))
	flushItAt(g, XY(20, 60), utils.VArcCCW(XY(10, 15), -radius, nil))

	return g
}

func flushItAt(g *GCode, v Tuple, vl []Tuple) {
	oldPos := g.Position()
	g.GotoXY(v)
	g.MoveXY(v.Offset(vl...)...)
	g.GotoXYZ(oldPos)
}

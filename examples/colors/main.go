// colors generates a design based on the example here:
// https://gitlab.com/gcmc/gcmc/blob/master/example/colors.gcmc
//
// Usage:
//   go run examples/colors/main.go > colors.gcode
package main

import (
	"fmt"
	"math"

	. "github.com/gmlewis/go-gcode/gcode"
)

var (
	box = []Tuple{
		XY(-1, -1),
		XY(-1, 1),
		XY(1, 1),
		XY(1, -1),
	}
)

func main() {
	g := gcmc()
	fmt.Printf("%v", g)
}

func gcmc() *GCode {
	g := New(UseGeneric)

	drawbox(g, 0, 10)

	for r := 1.0; r <= 30.0; r += 1.0 {
		drawbox(g, 3.0*(r-1.0), 50)
	}

	drawbox(g, 0, 20)
	drawbox(g, 0, 25)

	for _, b := range box {
		drawboxOffset(g, 0, 2.5, XY(10*b[0], 10*b[1]))
	}

	return g
}

func drawbox(g *GCode, ang, scl float64) {
	drawboxOffset(g, ang, scl, XY(0, 0))
}

func drawboxOffset(g *GCode, ang, scl float64, offset Tuple) {
	newBox := RotationZ(ang*math.Pi/180).
		Scale(scl, scl, 0).
		Translate(offset.X(), offset.Y(), offset.Z()).
		Transform(box...)
	g.GotoXY(LastXY(newBox))
	g.MoveXY(newBox...)
}

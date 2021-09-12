// bezier generates a design based on the example here:
// https://gitlab.com/gcmc/gcmc/blob/master/example/bezier.gcmc
//
// Usage:
//   go run examples/bezier/main.go > bezier.gcode
package main

import (
	"fmt"
	"math"

	. "github.com/gmlewis/go-gcode/gcode"
	"github.com/gmlewis/go-gcode/utils"
)

func main() {
	g := gcmc()
	fmt.Printf("%v", g)
}

func gcmc() *GCode {
	g := New()
	g.Prologue = true
	g.Epilogue = true

	g.Feedrate(500)

	square := []Tuple{
		XY(-1, -1),
		XY(1, -1),
		XY(1, 1),
		XY(-1, 1),
	}

	cp := []Tuple{
		XY(-2, 1), XY(-2, -1),
		XY(-1, -2), XY(1, -2),
		XY(2, -1), XY(2, 1),
		XY(1, 2), XY(-1, 2),
	}

	square = Scaling(8, 8, 0).Transform(square...)
	cp = Scaling(4.4, 4.4, 0).Transform(cp...)
	factor := 1.25

	for j := 0; j < 30; j++ {
		g.GotoXY(LastXY(square))
		for i := 0; i < len(square); i++ {
			last := i - 1
			if last < 0 {
				last = len(square) - 1
			}
			g.MoveXY(utils.VBezier3(square[last], cp[2*i], cp[2*i+1], square[i], nil)...)
		}
		cp = Scaling(factor, factor, 0).RotateZ(4 * math.Pi / 180).Transform(cp...)
		factor *= 0.99
	}

	return g
}

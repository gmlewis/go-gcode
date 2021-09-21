// involute-gear generates two gears and is based on the example here:
// https://gitlab.com/gcmc/gcmc/blob/master/example/involute-gear.gcmc
//
// Usage:
//   go run examples/involute-gear/main.go > gear.gcode
package main

import (
	"fmt"

	. "github.com/gmlewis/go-gcode/gcode"
	"github.com/gmlewis/go-gcode/utils"
)

const (
	HD = 6                       // Gear center-hole diameter (mm)
	N  = 9                       // Number of teeth
	PA = 20                      // Pressure angle (deg)
	D  = 100                     // Pitch diameter (mm)
	P  = float64(N) / float64(D) // Diametral pitch (teeth/mm)
)

func main() {
	g := gcmc()
	fmt.Printf("%v", g)
}

func gcmc() *GCode {
	g := New(UseGeneric)

	g.Feedrate(600)

	// First gear
	hole(g, X(D/2), HD/2)
	trace(g, utils.GearP(g, N, PA, P), X(D/2))

	// Second gear
	hole(g, X(-D/2), HD/2)
	trace(g, utils.GearP(g, N, PA, P), X(-D/2))

	return g
}

// Trace a path at given offset.
func trace(g *GCode, path []Tuple, offset Tuple) {
	g.GotoXY(path[len(path)-1].Add(offset))
	g.MoveXY(offset.Offset(path...)...)
}

// Make a hole at center point with given radius.
func hole(g *GCode, point Tuple, radius float64) {
	g.GotoXY(point.Sub(X(radius)))
	g.CircleCWRel(X(radius), nil)
}

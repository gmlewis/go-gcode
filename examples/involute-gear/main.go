// involute-gear generates two gears and is based on the example here:
// https://gitlab.com/gcmc/gcmc/blob/master/example/involute-gear.gcmc
package main

import (
	"fmt"

	. "github.com/gmlewis/go-gcode/gcode"
)

const (
	HD = 6     // Gear center-hole diameter (mm)
	N  = 9     // Number of teeth
	PA = 20    // Pressure angle (deg)
	D  = 100   // Pitch diameter (mm)
	P  = N / D // Diametral pitch (teeth/mm)
)

func main() {
	g := New()

	// First gear
	hole(g, Point(D/2, 0, 0), HD/2)
	trace(g, gearP(N, PA, P), Point(D/2, 0, 0))

	// Second gear
	hole(g, Point(-D/2, 0, 0), HD/2)
	trace(g, gearP(N, PA, P), Point(-D/2, 0, 0))

	fmt.Printf("%v\n", g)
}

// Trace a path at given offset.
func trace(g *GCode, path []Tuple, offset Tuple) {
	g.GotoXY(path[len(path)-1].Add(offset))
	g.MoveXY(offset.Offset(path...)...)
}

// Make a hole at center point with given radius.
func hole(g *GCode, point Tuple, radius float64) {
	g.GotoXY(point.Sub(Point(radius, 0, 0)))
	g.CircleCWRel(Point(radius, 0, 0))
}

func gearP(numTeeth int, pressureAngleDeg float64, diametralPitch float64) []Tuple {
	return nil
}

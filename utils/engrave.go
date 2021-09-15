package utils

import (
	. "github.com/gmlewis/go-gcode/gcode"
)

// Engrave takes care of pen-up/down handling for tracing text.
//
// Input:
// - g     GCode design
// - vs    Tuples (possibly from Typeset function)
// - zup   Pen-up height
// - zdown Pen-down height
//
// The engrave() function assumes that the Z-coordinate has not been altered or
// scaled as returned from the typeset() function. The Z-coordinate indicates
// pen-up/down, where 0.0 means pen-down and larger than 0.0 means pen-up (1.0
// is returned from the typeset() function). The pen movement is always in a
// single vector, as in: [-, -, penpos].
func Engrave(g *GCode, vs []Tuple, zUp, zDown float64) {
	for _, v := range vs {
		up := v.Z() > 0.0
		if up {
			g.GotoXYZ(XYZ(v.X(), v.Y(), zUp))
		} else {
			g.MoveXYZ(XYZ(v.X(), v.Y(), zDown))
		}
	}
}

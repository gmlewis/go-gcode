package utils

import (
	. "github.com/gmlewis/go-gcode/gcode"
)

// TracePath traces a path.
// Operation:
// - Goto last path entry
// - Move to cutting depth z
// - Move to all path entries
// - Return to orginal Z position
// Optional dwells for dw seconds at each path entry.
//
// Input:
// z    - cutting depth
// dw   - dwell at each point if >= 0
// path - vectorlist containing XY points
//
// Return value: none
func TracePath(g *GCode, z, dw float64, path ...Tuple) {
	if len(path) == 0 {
		return
	}
	g.Comment("-- tracepath at Z=", z, " --")
	oldZ := g.Position().Z()
	g.GotoXY(path[0])
	g.MoveZ(Z(z))
	if dw >= 0.0 {
		g.Dwell(dw)
	}
	for _, p := range path {
		g.MoveXY(p)
		if dw >= 0.0 {
			g.Dwell(dw)
		}
	}
	g.MoveXY(path[0])
	if dw >= 0.0 {
		g.Dwell(dw)
	}
	g.MoveZ(Z(oldZ))
	g.Comment("-- tracepath end --")
}

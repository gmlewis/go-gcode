package gcode

import "fmt"

// CircleCWRel perform a clockwise (cw) circle with radius length(centerPoint)
// and where centerPoint is the center point of the circle.
// The specified centerPoint is a relative position.
// The non-active plane coordinate may be used to create a helical movement.
func (g *GCode) CircleCWRel(centerPoint Tuple) *GCode {
	pos := g.Position()
	g.steps = append(g.steps, &Step{
		s:   fmt.Sprintf("G2 X%.8f Y%.8f I%.8f J%.8f", pos.X(), pos.Y(), centerPoint.X(), centerPoint.Y()),
		pos: pos,
	})
	return g
}

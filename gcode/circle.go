package gcode

import "fmt"

// ArcOptions represents options for the arc and circle methods.
type ArcOptions struct {
	Turns int
}

// ArcCCW performs a counter clockwise arc from the current position
// to endpoint with given radius.
// The arc will be the shortest angular movement with positive radius and
// the largest angular movement with negative radius.
// Optional turns sets the number of turns to perform.
func (g *GCode) ArcCCW(endPoint Tuple, radius float64, opts *ArcOptions) *GCode {
	pos := g.Position()
	s := fmt.Sprintf("G3 X%.8f Y%.8f Z%.8f J%.8f K%.8f", pos.X(), pos.Y(), pos.Z(), endPoint.X(), endPoint.Y())
	g.steps = append(g.steps, &Step{s: s, pos: pos})
	return g
}

// ArcCW performs a clockwise arc from the current position
// to endpoint with given radius.
// The arc will be the shortest angular movement with positive radius and
// the largest angular movement with negative radius.
// Optional turns sets the number of turns to perform.
func (g *GCode) ArcCW(endPoint Tuple, radius float64, opts *ArcOptions) *GCode {
	pos := g.Position()
	s := fmt.Sprintf("G2 X%.8f Y%.8f Z%.8f I%.8f K%.8f", pos.X(), pos.Y(), pos.Z(), endPoint.X(), endPoint.Y())
	g.steps = append(g.steps, &Step{s: s, pos: pos})
	return g
}

// CircleCWRel performs a clockwise (cw) circle with radius length(centerPoint)
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

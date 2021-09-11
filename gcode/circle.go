package gcode

// CircleCWRel perform a clockwise (cw) circle with radius length(centerPoint)
// and where centerPoint is the center point of the circle.
// The specified centerPoint is a relative position.
// The non-active plane coordinate may be used to create a helical movement.
func (g *GCode) CircleCWRel(centerPoint Tuple) *GCode {
	return g
}

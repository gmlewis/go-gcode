package gcode

func (g *GCode) home(opCode string, p Tuple) {
	g.steps = append(g.steps, &Step{s: opCode, pos: p})
	g.hasMoved = true
}

// HomeX homes the X axis.
func (g *GCode) HomeX() *GCode {
	pos := g.Position()
	newPos := XYZ(0, pos.Y(), pos.Z())
	g.home("G28 X", newPos)
	return g
}

// HomeY homes the Y axis.
func (g *GCode) HomeY() *GCode {
	pos := g.Position()
	newPos := XYZ(pos.X(), 0, pos.Z())
	g.home("G28 Y", newPos)
	return g
}

// HomeZ homes the Z axis.
func (g *GCode) HomeZ() *GCode {
	pos := g.Position()
	newPos := XYZ(pos.X(), pos.Y(), 0)
	g.home("G28 Z", newPos)
	return g
}

// HomeXY homes the X and Y axes.
func (g *GCode) HomeXY() *GCode {
	pos := g.Position()
	newPos := XYZ(0, 0, pos.Z())
	g.home("G28 X Y", newPos)
	return g
}

// HomeYZ homes the Y and Z axes.
func (g *GCode) HomeYZ() *GCode {
	pos := g.Position()
	newPos := XYZ(pos.X(), 0, 0)
	g.home("G28 Y Z", newPos)
	return g
}

// HomeXZ homes the X and Z axes.
func (g *GCode) HomeXZ() *GCode {
	pos := g.Position()
	newPos := XYZ(0, pos.Y(), 0)
	g.home("G28 X Z", newPos)
	return g
}

// HomeXYZ homes the X, Y, and Z axes.
func (g *GCode) HomeXYZ() *GCode {
	g.home("G28", XYZ(0, 0, 0))
	return g
}

package gcode

// PlaneT represents a construction plane.
type PlaneT string

const (
	PlaneXY PlaneT = "XY" // G17
	PlaneXZ PlaneT = "XZ" // G18
	PlaneYZ PlaneT = "YZ" // G19
)

// Plane sets the current construction plane.
func (g *GCode) Plane(p PlaneT) *GCode {
	var s string
	switch p {
	case PlaneXY:
		s = "G17"
	case PlaneXZ:
		s = "G18"
	case PlaneYZ:
		s = "G19"
	}

	g.steps = append(g.steps, &Step{
		s:   s,
		pos: g.Position(),
	})
	return g
}

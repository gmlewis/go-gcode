package gcode

// Pathmode sets exact pathmode.
func (g *GCode) Pathmode(exact bool) *GCode {
	// TODO: support G64.
	if exact {
		s := "G61"
		step := &Step{s: s, pos: g.Position()}
		g.steps = append(g.steps, step)
	}
	return g
}

package gcode

import "fmt"

// Dwell inserts a dwell command.
func (g *GCode) Dwell(dw float64) *GCode {
	s := fmt.Sprintf("G4 P%.8f", dw)
	step := &Step{s: s, pos: g.Position()}
	g.steps = append(g.steps, step)
	return g
}

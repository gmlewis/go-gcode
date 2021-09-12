package gcode

import "fmt"

// Feedrate sets the feedrate (F) to rate.
// The rate is interpreted following the setting of the Feedmode function.
func (g *GCode) Feedrate(rate float64) *GCode {
	g.steps = append(g.steps, &Step{
		s:   fmt.Sprintf("F%.8f", rate),
		pos: g.Position(),
	})
	return g
}
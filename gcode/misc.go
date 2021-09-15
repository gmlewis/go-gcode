package gcode

import (
	"fmt"
	"math"
)

// Feedrate sets the feedrate (F) to rate.
// The rate is interpreted following the setting of the Feedmode function.
func (g *GCode) Feedrate(rate float64) *GCode {
	g.steps = append(g.steps, &Step{
		s:   fmt.Sprintf("F%.8f", rate),
		pos: g.Position(),
	})
	return g
}

// Float returns a pointer to a float64 which is useful for options.
func Float(v float64) *float64 {
	return &v
}

// MilToMM converts mil (1/1000th of an inch) to millimeters.
func MilToMM(v float64) float64 {
	return 0.0254 * v
}

// ToRad converts degrees to radians.
func ToRad(a float64) float64 {
	return a * math.Pi / 180
}

// ToDeg converts radius to degrees.
func ToDeg(a float64) float64 {
	return a * 180 / math.Pi
}

// Package gcode provides methods used to generate G-Code.
package gcode

import "strings"

// GCode represents a G-Code design.
type GCode struct {
	steps []*Step
}

// New returns a new gcode design.
func New() *GCode {
	return &GCode{}
}

// String converts the design to a string.
func (g *GCode) String() string {
	var lines []string
	for _, step := range g.steps {
		lines = append(lines, step.s)
	}
	return strings.Join(lines, "\n")
}

// Step represents a step in the GCode.
type Step struct {
	s   string
	pos Tuple // position after performing the step.
}

// Position returns the current tool position.
func (g *GCode) Position() Tuple {
	if g == nil || len(g.steps) == 0 {
		return Point(0, 0, 0)
	}
	return g.steps[len(g.steps)-1].pos
}

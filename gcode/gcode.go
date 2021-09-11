// Package gcode provides methods used to generate G-Code.
package gcode

import (
	"fmt"
	"strings"
	"time"
)

// GCode represents a G-Code design.
type GCode struct {
	steps []*Step
}

// New returns a new gcode design.
func New() *GCode {
	const timeFmt = "2006-01-02 15:04:05"
	now := time.Now().Local()
	return &GCode{
		steps: []*Step{
			{s: fmt.Sprintf(prologue, now.Format(timeFmt))},
		},
	}
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

var prologue = `(go-gcode compiled code, do not change)
(%v)
(-- prologue begin --)
G17 ( Use XY plane )
G21 ( Use mm )
G40 ( Cancel cutter radius compensation )
G49 ( Cancel tool length compensation )
G54 ( Default coordinate system )
G80 ( Cancel canned cycle )
G90 ( Use absolute distance mode )
G94 ( Units Per Minute feed rate mode )
G64 ( Enable path blending for best speed )
(-- prologue end --)`

// Package gcode provides methods used to generate G-Code.
package gcode

import "fmt"

// GCode represents a G-Code design.
type GCode struct {
}

// New returns a new gcode design.
func New() *GCode {
	return &GCode{}
}

// String converts the design to a string.
func (g *GCode) String() string {
	return fmt.Sprintf("yo")
}

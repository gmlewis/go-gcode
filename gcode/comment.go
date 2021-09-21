package gcode

import (
	"fmt"
	"strings"
)

// Comment adds a comment to the G-Code using the provided args.
func (g *GCode) Comment(args ...interface{}) *GCode {
	var parts []string
	for _, arg := range args {
		parts = append(parts, fmt.Sprintf("%v", arg))
	}
	s := fmt.Sprintf(g.commentFmt, strings.Join(parts, ""))

	step := &Step{s: s, pos: g.Position()}
	g.steps = append(g.steps, step)
	return g
}

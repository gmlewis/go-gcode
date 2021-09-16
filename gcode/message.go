package gcode

import "fmt"

// Message displays the provided message.
func (g *GCode) Message(message string) *GCode {
	s := fmt.Sprintf("M117 %v", message)
	return g.sendOpCode(s)
}

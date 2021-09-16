package gcode

import "fmt"

// FanSpeed sets the named fan (0...) to the given speed (0-255).
func (g *GCode) FanSpeed(fanNum, speed int) *GCode {
	s := fmt.Sprintf("M106 P%v S%v", fanNum, speed)
	return g.sendOpCode(s)
}

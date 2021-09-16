package gcode

import "fmt"

// SpindleOnCW turns the spindle on in the clockwise direction
// with the given speed.
func (g *GCode) SpindleOnCW(rpm float64) *GCode {
	s := fmt.Sprintf("M3 S%v", rpm)
	return g.sendOpCode(s)
}

// SpindleOnCCW turns the spindle on in the counter-clockwise direction
// with the given speed.
func (g *GCode) SpindleOnCCW(rpm float64) *GCode {
	s := fmt.Sprintf("M4 S%v", rpm)
	return g.sendOpCode(s)
}

// SpindleOff turns the spindle off.
func (g *GCode) SpindleOff() *GCode { return g.sendOpCode("M5") }

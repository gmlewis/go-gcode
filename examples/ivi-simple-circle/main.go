// ivi-simple-circle generates a simple design to cut a hole
// on the IVI Closed-Loop 3D Printer/CNC/Laser-Engraver: https://ivi3d.com
//
// Usage:
//   go run examples/ivi-simple-circle/main.go > ivi-simple-circle.gcode
package main

import (
	"fmt"

	. "github.com/gmlewis/go-gcode/gcode"
	"github.com/gmlewis/go-gcode/utils"
)

const (
	spindleDiam = 20.4         // mm
	toolDiam    = 1000.0 / 8.0 // mils
	cutStep     = 1.1
	cutZ        = -1.0
)

var (
	home   = XYZ(0, 0, 10)
	offset = XY(0, -4)
)

func main() {
	g := gcmc()
	fmt.Printf("%v", g)
}

func gcmc() *GCode {
	g := New(UseIVI)

	g.Feedrate(300)

	targetRadius := 0.5 * spindleDiam
	toolRadius := 0.5 * MilToMM(toolDiam)

	g.GotoXYZ(home)
	utils.CCHole(g, offset, targetRadius, toolRadius, cutStep, cutZ)

	g.GotoXYZ(home)

	return g
}
// ivi-spool-cut generates a simple design to cut two slits in a spool
// on the IVI Closed-Loop 3D Printer/CNC/Laser-Engraver: https://ivi3d.com
//
// Usage:
//   go run examples/ivi-spool-cut/main.go > ivi-spool-cut.gcode
package main

import (
	"fmt"

	. "github.com/gmlewis/go-gcode/gcode"
)

const (
	spindleDiam   = 22.0 // mm
	wireFinalDiam = 48.0 // mm

	cutStep  = 0.5
	cutZ     = -4.0
	feedrate = 100
	safeZ    = 5.0
)

var (
	home   = XYZ(0, 0, 100)
	offset = XY(0, 1.7)
)

func main() {
	g := gcmc()
	fmt.Printf("%v", g)
}

func gcmc() *GCode {
	g := New(UseIVI)

	g.GotoZ(Z(safeZ))

	doCut := func(fromX, toX, speed float64) {
		g.GotoXY(offset.Add(X(fromX)))
		g.MoveZWithF(speed, Z(cutZ))
		g.MoveX(offset.Add(X(toX)))
		g.MoveZ(Z(safeZ))
	}

	doCut(0.5*spindleDiam, 0.5*wireFinalDiam, feedrate)
	doCut(0.5*spindleDiam, 0.5*wireFinalDiam, 3*feedrate)
	doCut(0.5*spindleDiam, 0.5*wireFinalDiam, 5*feedrate)

	doCut(-0.5*spindleDiam, -0.5*wireFinalDiam, feedrate)
	doCut(-0.5*spindleDiam, -0.5*wireFinalDiam, 3*feedrate)
	doCut(-0.5*spindleDiam, -0.5*wireFinalDiam, 5*feedrate)

	g.GotoXYZWithF(3000, home)

	return g
}

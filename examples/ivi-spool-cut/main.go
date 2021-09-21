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
	feedrate = 80
	safeZ    = 5.0
)

var (
	home   = XYZ(0, 0, 10)
	offset = XY(0, -5)
)

func main() {
	g := gcmc()
	fmt.Printf("%v", g)
}

func gcmc() *GCode {
	g := New(UseIVI)

	g.GotoZ(Z(safeZ))
	g.GotoXY(offset.Add(X(0.5 * spindleDiam)))
	g.MoveZWithF(feedrate, Z(cutZ))
	g.MoveX(offset.Add(X(0.5 * wireFinalDiam)))
	g.MoveX(offset.Add(X(0.5 * spindleDiam)))
	g.MoveZ(Z(safeZ))

	g.GotoXY(offset.Add(X(-0.5 * spindleDiam)))
	g.MoveZ(Z(cutZ))
	g.MoveX(offset.Add(X(-0.5 * wireFinalDiam)))
	g.MoveX(offset.Add(X(-0.5 * spindleDiam)))
	g.MoveZ(Z(safeZ))

	g.GotoXYZ(home)

	return g
}

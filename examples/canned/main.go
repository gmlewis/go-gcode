// canned generates a design based on the example here:
// https://gitlab.com/gcmc/gcmc/blob/master/example/canned.gcmc
//
// Usage:
//   go run examples/canned/main.go > canned.gcode
package main

import (
	"fmt"

	. "github.com/gmlewis/go-gcode/gcode"
	"github.com/gmlewis/go-gcode/utils"
)

func main() {
	g := gcmc()
	fmt.Printf("%v", g)
}

func gcmc() *GCode {
	g := New()
	g.Prologue = true
	g.Epilogue = true

	msg := "MSG,"
	homePos := Point(0, 0, 0)
	initPos := Point(-1, -1, 10)

	retract := 5.0 // R-plane level
	incr := 2.2    // Peck increment
	// oldZ := 1.0    // Return to initial Z in canned cycle (set to 0 for start at R-plane)
	// dwl := 0.2     // dwell time

	// Some holes...
	// They may include full XYZ coordinates for each point, however, specifying
	// only the coordinates that differ from the previous is enough, except for the
	// first one, which must include a Z-coordinate.
	lst := []Tuple{
		Point(10, 0, -1),
		X(15),
		X(20),
		X(25),
		Z(5),
		X(20),
		XZ(15, -2),
		X(10),
	}

	g.Feedrate(400)

	g.GotoXYZ(homePos)
	g.MoveXYZ(homePos)

	g.Comment(msg, "Canned peck drill without return-to-Z")
	g.GotoXYZ(initPos.Add(Y(0)))
	utils.CannedDrillPeck(g, retract, incr, false, lst...)

	g.Comment(msg, "Canned peck drill with return-to-Z")
	g.GotoXYZ(initPos.Add(Y(10)))
	utils.CannedDrillPeck(g, retract, incr, true, Translation(0, 10, 0).Transform(lst...)...)

	g.Comment(msg, "Canned drill without return-to-Z")
	g.GotoXYZ(initPos.Add(Y(20)))
	utils.CannedDrill(g, retract, -1, false, Translation(0, 20, 0).Transform(lst...)...)

	g.Comment(msg, "Canned drill with return-to-Z")
	g.GotoXYZ(initPos.Add(Y(30)))
	utils.CannedDrill(g, retract, -1, true, Translation(0, 30, 0).Transform(lst...)...)

	g.Comment(msg, "Canned drill dwell without return-to-Z")
	g.GotoXYZ(initPos.Add(Y(40)))
	utils.CannedDrill(g, retract, 0.5, false, Translation(0, 40, 0).Transform(lst...)...)

	g.Comment(msg, "Canned drill dwell with return-to-Z")
	g.GotoXYZ(initPos.Add(Y(50)))
	utils.CannedDrill(g, retract, 0.5, true, Translation(0, 50, 0).Transform(lst...)...)

	g.GotoXY(XY(homePos.X(), homePos.Y()))
	g.GotoXYZ(homePos)

	return g
}
// tool-compensate generates a design based on the example here:
// https://gitlab.com/gcmc/gcmc/blob/master/example/tool-compensate.gcmc
//
// Usage:
//   go run examples/tool-compensate/main.go > tool-compensate.gcode
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
	g := New(UseGeneric)

	g.Feedrate(600)

	safeZ := 5.0 // Safe Z-level
	cutZ := -1.0 // Cutting level
	home := XYZ(0.0, 0.0, safeZ)
	tw := 4.0       // Tool width
	tw2 := tw / 2.0 // Half tool-width for compensation

	g.GotoXYZ(home)
	g.MoveXYZ(home) // So LinuxCNC will show the following rapids

	path = Scaling(10, 10, 0).Transform(path...) // Scale to something visibly useful
	trace(g, path, Z(safeZ))                     // Show a rapid path to see the difference

	path = Z(cutZ).Offset(path...) // Set the cutting depth for all points
	run := func(opts utils.TPCOptions) {
		utils.TracePathComp(g, tw2, opts, path...)
	}

	// Choose your test trace by uncommenting the line(s) you want to see
	// run(utils.TPCLeft | utils.TPCOldZ | utils.TPCArcIn | utils.TPCArcOut | utils.TPCClosed)
	// run(utils.TPCRight | utils.TPCOldZ | utils.TPCArcIn | utils.TPCArcOut | utils.TPCClosed)
	run(utils.TPCRight | utils.TPCOldZ | utils.TPCArcIn | utils.TPCArcOut | utils.TPCQuiet)
	// run(utils.TPCLeft | utils.TPCOldZ | utils.TPCArcIn | utils.TPCArcOut | utils.TPCQuiet)
	// run(utils.TPCRight | utils.TPCOldZ)
	// run(utils.TPCLeft | utils.TPCOldZ)

	HD := 6.0           // Gear center-hole diameter
	N := 9              // Number of teeth
	PA := 20.0          // Pressure angle
	D := 100.0          // Pitch diameter
	P := float64(N) / D // Diametral pitch

	utils.CCHole(g, XY(120, 0), HD, tw2, tw2/2.0, cutZ)
	gearpath := XY(120, 0).Offset(utils.GearP(g, N, PA, P)...)
	trace(g, gearpath, XY(0, 0))           // Show the original path as a rapid
	gearpath = Z(cutZ).Offset(gearpath...) // Set cutting depth
	utils.TracePathComp(g, tw2, utils.TPCRight|utils.TPCOldZ|utils.TPCArcIn|utils.TPCArcOut|utils.TPCClosed, gearpath...)

	g.GotoXYZ(home)

	return g
}

func trace(g *GCode, p []Tuple, offset Tuple) {
	g.GotoXY(offset.Offset(p...)...)
	g.GotoXY(offset.Offset(p[0])...)
}

var path = []Tuple{
	XY(4, 3), XY(3.5, 3.5), XY(3, 4), // Entry with co-linear point
	XY(2.7, 4), XY(2.5, 3), XY(2.3, 4), // inside marginal entry

	XY(2, 4), XY(1, 5),

	XY(0.2, 5), XY(0.1, 6), XY(0, 3), XY(-0.1, 4), // Sharp both directions
	XY(-1, 4), XY(-1.1, 6), XY(-1.2, 5), // Outside "horns"
	XY(-2.0, 5), XY(-2.1, 6), XY(-2, 4),

	XY(-3, 4), XY(-4, 3),

	XY(-4, 2.5), XY(-5, 3), XY(-4, 2.5), // Path reversal angled
	XY(-4, 2), XY(-3, 2), XY(-4, 2), // Path reversal
	XY(-4, 1.5), XY(-5, 1), XY(-4, 1.5), // Path reversal angled
	XY(-4, 1), XY(-3, 0), XY(-4, 0), // Inside angle
	XY(-4, -1), XY(-5, -1.05), XY(-4, -1.1), // Symmetric outside
	XY(-4, -2), XY(-5, -3), XY(-4, -2.1), // Symmetric outside angled

	XY(-4, -3), XY(-3, -4),

	XY(-2, -4), XY(-1.9, -2), XY(-1.8, -3), // Inside "horns"
	XY(-1.0, -3), XY(-0.9, -2), XY(-1, -4),
	XY(-0, -4), XY(0, -5), XY(0.1, -2), XY(0.2, -3), // Sharp both directions

	XY(1, -3), XY(2, -4),

	XY(2.3, -4), XY(2.5, -5), XY(2.7, -4), // Outside marginal entry

	XY(3, -4), XY(4, -3),

	XY(4, -2.1), XY(3, -2.5), XY(4, -2), // Symmetric inside angled
	XY(4, -1.1), XY(3, -1.05), XY(4, -1), // Symmetric inside
	XY(4, 0), XY(5, 0), XY(4, 1), // Outside angle
	XY(4, 1.5), XY(3, 1), XY(4, 1.5), // Path reversal
	XY(4, 2), XY(5, 2), XY(4, 2), // Path reversal

	XY(4, 2.5), // Exit point (below entry)
}

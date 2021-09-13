package utils

import (
	"log"

	"github.com/gmlewis/go-gcode/gcode"
)

// Based on:
// https://gitlab.com/gcmc/gcmc/blob/master/example/cc_hole.inc.gcmc

// Hole milling example
// --------------------
// Mill a hole in continuous curvature movements. Not a single straight line is
// required to mill an arbitrary large hole from any size milling bit.
// Continuous curvature milling reduces the stress on the mill, bit and object
// by preventing any jerking.
//
// A hole is milled at a given center and depth with a target radius. The
// milling-bit radius and the cutting step define how many turning cycles are
// required to finish the hole. The mill is retracted with a helical move back
// to the center and starting Z-position.
func CCHole(g *gcode.GCode, center gcode.Tuple, targetRadius, toolRadius, cutStep, cutZ float64) {
	if targetRadius <= 0.0 {
		log.Fatal("targetRadius must be positive")
	}
	if toolRadius <= 0.0 {
		log.Fatal("toolRadius must be positive")
	}
	if targetRadius <= toolRadius {
		log.Fatalf("hole targetRadius (%.8f) must be larger than toolRadius (%.8f)", targetRadius, toolRadius)
	}
	if cutStep <= 0.0 {
		log.Fatal("cutStep must be positive")
	}
	if cutStep > 2.0*toolRadius {
		log.Printf("WARNING: cutStep is larger than twice the toolRadius, not all material will be removed.")
	} else if cutStep == 2.0*toolRadius {
		log.Printf("WARNING: cutSteo is exactly twice the toolRadius, material may be left at the inner edge")
	}

	oldZ := g.Position().Z()

	g.Comment("-- CCHole center=", center, " targetRadius=", targetRadius, " toolRadius=", toolRadius, " cutStep=", cutStep, " cutZ=", cutZ, " --")

	g.GotoXYZ(gcode.XYZ(center.X(), center.Y(), oldZ))
	g.MoveZ(gcode.Z(cutZ))

	r := toolRadius
	n := 1
	dir := -1.0
	var p float64
	for r < targetRadius {
		if targetRadius-r >= cutStep {
			p = float64(2*n-1) * cutStep
			r += cutStep
		} else {
			p = float64(2*n-2)*cutStep + targetRadius - r
			r += targetRadius - r
		}
		g.ArcCWRel(gcode.XY(0, dir*p), 0.5*p, nil)
		g.CircleCW(gcode.XYZ(center.X(), center.Y(), cutZ), nil)
		n++
		dir = -dir
	}

	g.ArcCWRel(
		gcode.XYZ(0, dir*(targetRadius-toolRadius), oldZ-cutZ),
		0.5*(targetRadius-toolRadius), nil)

	g.Comment("-- end CCHole --")
}

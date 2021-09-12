// ball-in-cube generates a design based on the example here:
// https://gitlab.com/gcmc/gcmc/blob/master/example/ball-in-cube.gcmc
//
// Usage:
//   go run examples/ball-in-cube/main.go > ball-in-cube.gcode
package main

import (
	"fmt"
	"math"

	. "github.com/gmlewis/go-gcode/gcode"
)

const (
	boxSize    = 50.0      // mm
	cutterRad  = 3.0 / 2.0 // mm
	boxBarSize = 3.0       // size of box bars
	angInc     = 5.0       // deg
	safeZ      = 10.0      // safe height above object
)

func main() {
	g := gcmc()
	fmt.Printf("%v", g)
}

func gcmc() *GCode {
	g := New()
	g.Prologue = true
	g.Epilogue = true

	cRad := math.Sqrt2 * (boxSize/2.0 + cutterRad)

	g.Feedrate(250)
	g.MoveXY(XY(0, 0))

	// FIXME:
	// --- Warning ---
	// The space above the ball is not carved out. The cutter is unlikely to have
	// enough vertical free-room to remove material. The current path would jam the
	// cutter into the remaining material and that would end badly. A separate
	// routine is required to remove enough material from the block to let the
	// cutter do a proper job.
	// --- Warning ---

	ang := angInc
	for {
		radians := ang * math.Pi / 180.0
		xy := math.Sin(radians) * cRad * 0.5    // Corner point
		zz := -(1.0 - math.Cos(radians)) * cRad // Depth at corner
		rr := math.Cos(radians) * cRad          // Adjusted radius
		g.Feedrate(75)
		g.MoveXY(XY(xy, xy))
		g.MoveZ(Z(zz))
		g.Feedrate(100)
		g.Plane(PlaneXZ)
		g.ArcCW(Point(-xy, xy, zz), rr, nil)
		g.Plane(PlaneYZ)
		g.ArcCCW(Point(-xy, -xy, zz), rr, nil)
		g.Plane(PlaneXZ)
		g.ArcCCW(Point(xy, -xy, zz), rr, nil)
		g.Plane(PlaneYZ)
		g.ArcCW(Point(xy, xy, zz), rr, nil)
		ang += angInc
		if ang >= 55.0 || xy >= 0.5*boxSize-cutterRad-boxBarSize {
			break
		}
	}

	g.MoveZ(Z(safeZ))

	g.Feedrate(75)

	// Cut the inside bar
	xy := 0.5*boxSize - boxBarSize - cutterRad
	g.GotoXY(XY(xy, xy))

	for zz := -cutterRad; zz >= -boxBarSize; zz -= cutterRad {
		g.MoveZ(Z(zz))
		g.MoveXY(XY(-xy, xy))
		g.MoveXY(XY(-xy, -xy))
		g.MoveXY(XY(xy, -xy))
		g.MoveXY(XY(xy, xy))
	}
	g.MoveZ(Z(safeZ))

	// Cut the outside bar
	xy = 0.5*boxSize + boxBarSize + cutterRad
	g.GotoXY(XY(xy, xy))

	for zz := -cutterRad; zz >= -boxBarSize; zz -= cutterRad {
		g.MoveZ(Z(zz))
		g.MoveXY(XY(-xy, xy))
		g.MoveXY(XY(-xy, -xy))
		g.MoveXY(XY(xy, -xy))
		g.MoveXY(XY(xy, xy))
	}
	g.MoveZ(Z(safeZ))
	g.GotoXY(XY(0, 0))

	return g
}

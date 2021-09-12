// wheels generates wheel paths and is based on the example here:
// https://gitlab.com/gcmc/gcmc/blob/master/example/wheels.gcmc
//
// Usage:
//   go run examples/wheels/main.go > wheels.gcode
package main

import (
	"fmt"
	"math"

	. "github.com/gmlewis/go-gcode/gcode"
)

const (
	safeHeight   = 1
	cuttingDepth = -1
)

func main() {
	g := gcmc()
	fmt.Printf("%v", g)
}

func gcmc() *GCode {
	g := New()
	g.Prologue = true
	g.Epilogue = true

	svec := Point(5, 5, 1)
	wheels := []*Wheel{
		createWheel(10, 1, 0),
		createWheel(5, 7, 0),
		createWheel(3.333, -17, 90),
	}

	g.Feedrate(60)
	goAtSafeHeight(g, 0, 0)
	cutPath(g, wheels, 0, 0.01, 360, cuttingDepth, svec)
	goAtSafeHeight(g, 0, 0)

	return g
}

// Wheel represents a wheel.
type Wheel struct {
	radius float64
	speed  float64
	phase  float64
}

func createWheel(radius, speed, phase float64) *Wheel {
	return &Wheel{radius: radius, speed: speed, phase: phase}
}

func calcPoint(wheels []*Wheel, angleDeg float64) Tuple {
	var x, y float64
	angle := angleDeg * math.Pi / 180.0
	for _, w := range wheels {
		at := w.speed*angle + w.phase
		x += w.radius * math.Cos(at)
		y += w.radius * math.Sin(at)
	}
	return XY(x, y)
}

func cutPath(g *GCode, wheels []*Wheel, start, inc, end, cdepth float64, scale Tuple) {
	for angle := start; angle <= end; angle += inc {
		if angle == start {
			// we should be at safe height, so move to cutting depth
			p := calcPoint(wheels, angle)
			g.GotoXY(XY(p.X()*scale.X(), p.Y()*scale.Y()))
			g.GotoZ(Z(cdepth * scale.Z()))
			continue
		}

		p := calcPoint(wheels, angle)
		g.MoveXY(XY(p.X()*scale.X(), p.Y()*scale.Y()))
	}
}

func goAtSafeHeight(g *GCode, x, y float64) {
	g.GotoZ(Z(safeHeight))
	g.GotoXYZ(Point(x, y, safeHeight))
}

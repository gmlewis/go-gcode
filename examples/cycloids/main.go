// cycloids generates a design based on the example here:
// https://gitlab.com/gcmc/gcmc/blob/master/example/cycloids.gcmc
//
// Usage:
//   go run examples/cycloids/main.go > cycloids.gcode
package main

import (
	"fmt"
	"math"

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

	g.Feedrate(2000)

	parms := []Tuple{
		XYZ(5, 3, 2),
		XYZ(10, 1, 2),
		XYZ(5, 3, 5),
		XYZ(10, 2, 1),
		XYZ(17, 5, 3),
		XYZ(7, 13, 5),
	}

	i := -2.5
	for _, p := range parms {
		maxW := p[0] + p[1] + p[2]
		hTroch := doChoid(p, 2, hypotrochoidPoint)
		eTroch := doChoid(p, 2, epitrochoidPoint)
		f := 30.0 / maxW
		utils.TracePath(g, -1, -1, Scaling(f, f, f).Translate(2.1*30*i, -30, 0).Transform(hTroch...)...)
		utils.TracePath(g, -1, -1, Scaling(f, f, f).Translate(2.1*30*i, 30, 0).Transform(eTroch...)...)
		i += 1.0
	}

	return g
}

// Create a hypto- or epi-trochoid path
//
// p[0] - R	- Fixed circle radius
// p[1] - r	- Rolling circle radius
// p[2] - d	- Tracking point distance from rolling circle center
// aStep- Angular step (in degrees)
func doChoid(p Tuple, aStep float64, f func(float64, Tuple) Tuple) []Tuple {
	var path []Tuple
	maxA := 360 * p[1]
	for a := 0.0; a < maxA; a += aStep {
		path = append(path, f(a, p))
	}
	return path
}

// Hypotrochoid - Circle trace inside a circle
// See: https://en.wikipedia.org/wiki/Hypotrochoid
//
// A hypotrochoid becomes a hypocycloid when d == r
//
// a	- Angle of rotation for the point (in degrees)
// p[0] - R	- Fixed circle radius
// p[1] - r	- Rolling circle radius
// p[2] - d	- Tracking point distance from rolling circle center
func hypotrochoidPoint(a float64, p Tuple) Tuple {
	phi := a * math.Pi / 180.0
	dr := p[0] - p[1]
	drOr := dr / p[1]
	x := dr*math.Cos(phi) + p[2]*math.Cos(phi*drOr)
	y := dr*math.Sin(phi) - p[2]*math.Sin(phi*drOr)
	return XY(x, y)
}

// Epitrochoid - Circle trace inside a circle
// See: https://en.wikipedia.org/wiki/Epitrochoid
//
// A epitrochoid becomes a epicycloid when d == r
//
// a	- Angle of rotation for the point (in degrees)
// p[0] - R	- Fixed circle radius
// p[1] - r	- Rolling circle radius
// p[2] - d	- Tracking point distance from rolling circle center
func epitrochoidPoint(a float64, p Tuple) Tuple {
	phi := a * math.Pi / 180.0
	dr := p[0] + p[1]
	drOr := dr / p[1]
	x := dr*math.Cos(phi) - p[2]*math.Cos(phi*drOr)
	y := dr*math.Sin(phi) - p[2]*math.Sin(phi*drOr)
	return XY(x, y)
}

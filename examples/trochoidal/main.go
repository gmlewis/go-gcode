// trochoidal generates a design based on the example here:
// https://gitlab.com/gcmc/gcmc/blob/master/example/trochoidal.gcmc
//
// Usage:
//   go run examples/trochoidal/main.go > trochoidal.gcode
package main

import (
	"fmt"
	"log"
	"math"

	. "github.com/gmlewis/go-gcode/gcode"
)

const (
	cutZ  = -1.0
	safeZ = 5.0
)

func main() {
	g := gcmc()
	fmt.Printf("%v", g)
}

func gcmc() *GCode {
	g := New()
	g.Prologue = true
	g.Epilogue = true

	home := XYZ(0.0, 0.0, safeZ)

	path := []Tuple{
		XY(0, 1),
		XY(2, 2),
		XY(0, 4),
		XY(-1, 2),
	}

	g.Feedrate(300.0)
	g.GotoZ(Z(safeZ))
	g.GotoXYZ(home)

	path = Scaling(25, 25, 0).Transform(path...)

	// Trochoidal high-speed milling of the outline
	g.Feedrate(3000.0) // *Really* high-speed milling
	for i, p := range path {
		trochoidMove(g, path[(i+len(path)-1)%len(path)], p, cutZ, 5.0, 2.0)
	}

	// Clean-cutting the object
	g.Feedrate(150.0)           // "Finishing the edge" speed
	g.Pathmode(true)            // Exact path mode so we hit the corners exactly
	g.MoveZ(Z(cutZ))            // We are at the "outside" of the path, reenter cutting depth
	g.MoveXY(path[len(path)-1]) // The first corner
	g.MoveXY(path...)           // Trace the object

	g.GotoZ(Z(safeZ))
	g.GotoXYZ(home)

	return g
}

// Trochoidal point calculation.
// See: https://en.wikipedia.org/wiki/Trochoid
func trochoidPoint(ang, a, b float64) Tuple {
	// The first part is the trochoid, the second part moves the first 180
	// degree point at a relative "0, 0" location so we can scale in any
	// way without having to do hard math
	return XY(a*ang-b*math.Sin(ang)-a*math.Pi, b-b*math.Cos(ang)-2*b)
}

// Perform a move from startpoint to endpoint using a trochoidal path.
// - Cutting at depth cutz (returns to old Z)
// - Trochoid radius as specified
// - Increment for each turn as specified
func trochoidMove(g *GCode, startpoint, endpoint Tuple, cutz, radius, increment float64) {
	a := increment / (2.0 * math.Pi)        // Trochoid step parameter
	ainc := math.Log10(radius) * ToRad(5.0) // Steps are logarithmic based on the radius to reduce small steps
	oldZ := g.Position().Z()
	vec := endpoint.Sub(startpoint) // Vector denoting path to move

	// Of we are not moving, it is an error
	if vec.Magnitude() <= 0.0 {
		log.Printf("trochoid move is not going anywhere")
		return
	}

	g.Comment("-- trochoid_move at ", cutz, " from ", startpoint, " to ", endpoint, " radius=", radius, " increment=", increment, " --")
	// Calculate the number of *whole* rotations, rounded up, we need to make
	n := 2.0 * math.Pi * math.Ceil(vec.Magnitude()/increment)

	// The path may be arbitrary angled, get the angle for rotating the trochoid
	rot := math.Atan2(vec[1], vec[0])

	// Go to the trochoid entry-point and move to cutting depth
	g.GotoXY(startpoint.Add(RotationZ(rot).Transform(trochoidPoint(0.0, a, radius))[0]))
	g.MoveZ(Z(cutZ))

	// Calculate each next point of the trochoid until we traversed the whole path to the endpoint
	for i := 0.0; i < n; i += ainc {
		g.MoveXY(startpoint.Add(RotationZ(rot).Transform(trochoidPoint(i, a, radius))[0]))
	}

	// Return to old Z so we will not bump into stuff
	g.GotoZ(Z(oldZ))
	g.Comment("-- trochoid_move end --")
}

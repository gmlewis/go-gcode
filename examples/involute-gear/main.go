// involute-gear generates two gears and is based on the example here:
// https://gitlab.com/gcmc/gcmc/blob/master/example/involute-gear.gcmc
//
// Usage:
//   go run examples/involute-gear/main.go > gear.gcode
package main

import (
	"fmt"
	"log"
	"math"

	. "github.com/gmlewis/go-gcode/gcode"
)

const (
	HD = 6                       // Gear center-hole diameter (mm)
	N  = 9                       // Number of teeth
	PA = 20                      // Pressure angle (deg)
	D  = 100                     // Pitch diameter (mm)
	P  = float64(N) / float64(D) // Diametral pitch (teeth/mm)
)

func main() {
	g := gcmc()
	fmt.Printf("%v\n", g)
}

func gcmc() *GCode {
	g := New()

	g.Feedrate(600)

	// First gear
	hole(g, X(D/2), HD/2)
	trace(g, gearP(g, N, PA, P), X(D/2))

	// Second gear
	hole(g, X(-D/2), HD/2)
	trace(g, gearP(g, N, PA, P), X(-D/2))

	return g
}

// Trace a path at given offset.
func trace(g *GCode, path []Tuple, offset Tuple) {
	g.GotoXY(path[len(path)-1].Add(offset))
	g.MoveXY(offset.Offset(path...)...)
}

// Make a hole at center point with given radius.
func hole(g *GCode, point Tuple, radius float64) {
	g.GotoXY(point.Sub(X(radius)))
	g.CircleCWRel(X(radius))
}

// Gear terms:
// N	- Number of Teeth
// Pa	- Pressure Angle
// D	- Pitch	Diameter	- D = N/P = Do - 2/P	(Gear radius at center of the teeth)
// P	- Diametral Pitch	- P = N/D
// p	- Circular Pitch	- p = pi() / P
// Db	- Base Diameter		- Db = D * cos(Pa)	(Bottom of teeth insertion)
// Dr	- Root Diameter		- Dr = D - 2b		(Bottom of tooth cutout)
// Do	- Outside Diameter	- Do = D + 2a
// a	- Addendum		- a = 1/P
// b	- Dedendum		- b = ht - a
// ht	- Whole Depth (Pa<20)	- 2.157/P
// ht	- Whole Depth (Pa>=20)	- 2.2/P + 0.05mm	(Total depth from outer dia to bottom)
// t	- Tooth Thickness	- t = pi()/(2*P)	(Thinckness at Pitch Diameter)

const angleStepDeg = 2 // Trace interval for curves, in degrees.

func toRad(a float64) float64 {
	return a * math.Pi / 180
}

// Point on involute curve at specified angle in degrees, see https://en.wikipedia.org/wiki/Involute
// Cartesian:
//	x = a * ( cos(t) + t * sin(t))
//	y = a * ( sin(t) - t * cos(t))
// Polar:
//	r   = a * sqrt(1 + t^2) = sqrt(a^2 + (a*t)^2)
//	phi = t - atan(t)
// where:
// - a = circle radius
// - t = angle (radians)
//
// For angle from circle radius: t^2 = (r/a)^2 - 1
func involutePoint(angle, radius float64) Tuple {
	angle = toRad(angle) // Multiplication must be in radians.
	cos := math.Cos(angle)
	sin := math.Sin(angle)
	return XY(cos+angle*sin, sin-angle*cos).MultScalar(radius)
}

func involuteAngle(radius, outrad float64) float64 {
	return toRad(math.Sqrt(math.Pow(outrad/radius, 2) - 1))
}

// Make a gear with:
// - nteeth		Number of teeth
// - pressureAngleDeg	Teeth contact pressure angle
// - diametralPitch	Diametral pitch (teeth/length)
//
// Return a vectorlist with outer points of the gear centered at [0,0]
func gearP(g *GCode, nteeth int, pressureAngleDeg float64, diametralPitch float64) []Tuple {
	// The routine gets in serious trouble if you make the pressure angle
	// too large or too small. Warn the user if such case occurs.
	if pressureAngleDeg > 24.6 {
		log.Printf("Pressure angle (%v) too large, cannot fit teeth inside the set outside diameter", pressureAngleDeg)
	}
	if pressureAngleDeg < 12.0 {
		log.Printf("Pressure angle (%v) too small, teeth may get stuck at pitch radius", pressureAngleDeg)
	}

	pitchDiameter := float64(nteeth) / diametralPitch
	baseDiameter := pitchDiameter * math.Cos(toRad(pressureAngleDeg))
	addendum := 1 / diametralPitch
	ht := 2.157 / diametralPitch
	dedendum := ht - addendum
	outsideDiameter := pitchDiameter + 2*addendum
	rootDiameter := baseDiameter - 2*dedendum
	workDiameter := outsideDiameter - 4*addendum

	var tooth []Tuple // The curve for one tooth

	log.Printf("nteeth=%v, pressureAngleDeg=%v, diametralPitch=%v", nteeth, pressureAngleDeg, diametralPitch)
	log.Printf("addendum=%v, dedendum=%v, ht=%v", addendum, dedendum, ht)
	log.Printf("pitchDiameter=%v", pitchDiameter)
	log.Printf("baseDiameter=%v", baseDiameter)
	log.Printf("outsideDiameter=%v", outsideDiameter)
	log.Printf("rootDiameter=%v", rootDiameter)
	log.Printf("workDiameter=%v", workDiameter)

	// Show the different diameters:
	hole(g, XY(0, 0), pitchDiameter/2)
	hole(g, XY(0, 0), baseDiameter/2)
	hole(g, XY(0, 0), outsideDiameter/2)
	hole(g, XY(0, 0), rootDiameter/2)
	hole(g, XY(0, 0), workDiameter/2)

	// Fillet radius is approx. Will not reach root exactly, but close enough.
	// Otherwise need to calculate intersection with root-circle.
	filletrad := (baseDiameter - rootDiameter) / 8

	// Center of the fillet arc, involute makes a ~240deg angle with fillet arc.
	// The fillet arc runs from the root to the working depth of the gear.
	center := RotationZ(toRad(60)).MultTuple(X(-filletrad)).Add(X(workDiameter / 2))

	// Trace the fillet arc from ~root-circle to working depth at involute arc starting Y-level
	var a float64
	for a = 180.0; a > 60.0; a -= angleStepDeg * 2.5 {
		r := toRad(a)
		tooth = append(tooth, XY(math.Cos(r), math.Sin(r)).MultScalar(filletrad).Add(center))
	}
	if a != 60.0 {
		// Add the last point if we did not reach the working depth
		r := toRad(60)
		tooth = append(tooth, XY(math.Cos(r), math.Sin(r)).MultScalar(filletrad).Add(center))
	}

	// Calculate the maximum involute angle to intersect at the outside radius
	maxA := involuteAngle(baseDiameter/2, outsideDiameter/2)

	// Trace the involute arc from the base up to outside radius
	for a = 0.0; a < maxA; a += angleStepDeg {
		tooth = append(tooth, involutePoint(a, baseDiameter/2))
	}
	if a != maxA {
		// Add the last point if we did not reach the outside radius
		tooth = append(tooth, involutePoint(maxA, baseDiameter/2))
	}

	// We now have one side of the tooth. Rotate to be at tooth-symmetry on X-axis
	tooth = RotationZ(toRad(-90 / float64(nteeth))).Transform(tooth...)

	// Remember how many point we have in a side
	ntooth := len(tooth)

	// Add the same curve mirrored to make the other side of the tooth
	// Coordinates reverse to have them all in one direction only.
	// Also add a point in the middle of the outside linear segment connecting
	// both sides of the tooth. This will help the caller to attach a
	// tool-compensated path at that point.
	mirror := []Tuple{X(-tooth[len(tooth)-1].X())}
	for i := len(tooth) - 1; i >= 0; i-- {
		mirror = append(mirror, XY(tooth[i].X(), -tooth[i].Y()))
	}
	tooth = append(tooth, mirror...)

	// Create all teeth of the gear by adding each tooth at correct angle.
	var gear []Tuple
	for i := 0; i < nteeth; i++ {
		a := toRad(float64(i) * 360 / float64(nteeth))
		gear = append(gear, RotationZ(a).Transform(tooth...)...)
	}

	// Return the gear with the gear points rotated by a tooth's side
	// point-count plus one for the intermediate point to have the middle
	// of the outside segment as entry-point into the path.
	// return tail(gear, ntooth+1) + head(gear, -ntooth-1)
	result := append([]Tuple{}, gear[ntooth+1:]...)
	result = append(result, gear[0:ntooth+1]...)
	return result
}

package utils

import (
	"log"
	"math"

	. "github.com/gmlewis/go-gcode/gcode"
)

// GearP makes a gear with:
// - nteeth		Number of teeth
// - pressureAngleDeg	Teeth contact pressure angle
// - diametralPitch	Diametral pitch (teeth/length)
//
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
//
// Return a vectorlist with outer points of the gear centered at [0,0]
func GearP(g *GCode, nteeth int, pressureAngleDeg float64, diametralPitch float64) []Tuple {
	// The routine gets in serious trouble if you make the pressure angle
	// too large or too small. Warn the user if such case occurs.
	if pressureAngleDeg > 24.6 {
		log.Printf("Pressure angle (%.8f) too large, cannot fit teeth inside the set outside diameter", pressureAngleDeg)
	}
	if pressureAngleDeg < 12.0 {
		log.Printf("Pressure angle (%.8f) too small, teeth may get stuck at pitch radius", pressureAngleDeg)
	}

	pitchDiameter := float64(nteeth) / diametralPitch
	baseDiameter := pitchDiameter * math.Cos(ToRad(pressureAngleDeg))
	addendum := 1 / diametralPitch
	ht := 2.157 / diametralPitch
	dedendum := ht - addendum
	outsideDiameter := pitchDiameter + 2*addendum
	rootDiameter := baseDiameter - 2*dedendum
	workDiameter := outsideDiameter - 4*addendum

	// log.Printf("nteeth=%v, pressureAngleDeg=%.8f, diametralPitch=%.8f", nteeth, pressureAngleDeg, diametralPitch)
	// log.Printf("addendum=%.8f, dedendum=%.8f, ht=%.8f", addendum, dedendum, ht)
	// log.Printf("pitchDiameter=%.8f", pitchDiameter)
	// log.Printf("baseDiameter=%.8f", baseDiameter)
	// log.Printf("outsideDiameter=%.8f", outsideDiameter)
	// log.Printf("rootDiameter=%.8f", rootDiameter)
	// log.Printf("workDiameter=%.8f", workDiameter)

	// Show the different diameters:
	// hole(g, XY(0, 0), pitchDiameter/2)
	// hole(g, XY(0, 0), baseDiameter/2)
	// hole(g, XY(0, 0), outsideDiameter/2)
	// hole(g, XY(0, 0), rootDiameter/2)
	// hole(g, XY(0, 0), workDiameter/2)

	tooth := rotatedHalfTooth(baseDiameter, outsideDiameter, rootDiameter, workDiameter, nteeth)

	// Remember how many points we have in a side.
	ntooth := len(tooth)

	// Add the same curve mirrored to make the other side of the tooth
	// Coordinates reverse to have them all in one direction only.
	// Also add a point in the middle of the outside linear segment connecting
	// both sides of the tooth. This will help the caller to attach a
	// tool-compensated path at that point.
	mirror := mirrorTooth(tooth)
	tooth = append(tooth, mirror...)

	// Create all teeth of the gear by adding each tooth at correct angle.
	var gear []Tuple
	for i := 0; i < nteeth; i++ {
		a := ToRad(float64(i) * 360.0 / float64(nteeth))
		gear = append(gear, RotationZ(a).Transform(tooth...)...)
	}

	// Return the gear with the gear points rotated by a tooth's side
	// point-count plus one for the intermediate point to have the middle
	// of the outside segment as entry-point into the path.
	// return tail(gear, ntooth+1) + head(gear, -ntooth-1)
	result := append([]Tuple{}, gear[ntooth:]...)
	result = append(result, gear[0:ntooth]...)
	return result
}

const angleStepDeg = 2 // Trace interval for curves, in degrees.

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
	angle = ToRad(angle) // Multiplication must be in radians.
	cos := math.Cos(angle)
	sin := math.Sin(angle)
	return XY(radius*(cos+angle*sin), radius*(sin-angle*cos))
}

func involuteAngle(radius, outrad float64) float64 {
	return ToDeg(math.Sqrt(math.Pow(outrad/radius, 2.0) - 1.0))
}

func mirrorTooth(tooth []Tuple) []Tuple {
	mirror := []Tuple{X(tooth[len(tooth)-1].X())}
	for i := len(tooth) - 1; i >= 0; i-- {
		mirror = append(mirror, XY(tooth[i].X(), -tooth[i].Y()))
	}
	return mirror
}

func rotatedHalfTooth(baseDiameter, outsideDiameter, rootDiameter, workDiameter float64, nteeth int) []Tuple {
	tooth := halfTooth(baseDiameter, outsideDiameter, rootDiameter, workDiameter)

	// We now have one side of the tooth. Rotate to be at tooth-symmetry on X-axis
	tooth = RotationZ(ToRad(-90.0 / float64(nteeth))).Transform(tooth...)

	return tooth
}

func halfTooth(baseDiameter, outsideDiameter, rootDiameter, workDiameter float64) []Tuple {
	var tooth []Tuple

	// Fillet radius is approx. Will not reach root exactly, but close enough.
	// Otherwise need to calculate intersection with root-circle.
	filletrad := (baseDiameter - rootDiameter) / 8

	// Center of the fillet arc, involute makes a ~240deg angle with fillet arc.
	// The fillet arc runs from the root to the working depth of the gear.
	center := RotationZ(ToRad(60)).MultTuple(X(-filletrad)).Add(XYZ(workDiameter/2, 0, 0))

	// Trace the fillet arc from ~root-circle to working depth at involute arc starting Y-level
	var a float64
	for a = 180.0; a > 60.0; a -= angleStepDeg * 2.5 {
		angle := ToRad(a)
		pt := XY(filletrad*math.Cos(angle)+center.X(), filletrad*math.Sin(angle)+center.Y())
		tooth = append(tooth, pt)
	}
	if a != 60.0 {
		// Add the last point if we did not reach the working depth
		angle := ToRad(60)
		pt := XY(filletrad*math.Cos(angle)+center.X(), filletrad*math.Sin(angle)+center.Y())
		tooth = append(tooth, pt)
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

	return tooth
}

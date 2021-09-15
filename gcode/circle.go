package gcode

import (
	"fmt"
	"log"
	"math"
)

// TurnsOption represents options for the arc and circle methods.
type TurnsOption struct {
	Turns int
}

type arcFnEnumT int

const (
	fnArcCW arcFnEnumT = iota
	fnArcCWRel
	fnArcCCW
	fnArcCCWRel
)

func (g *GCode) allArcs(endP Tuple, origRad float64, relative bool, ft arcFnEnumT, opCode string, opts *TurnsOption) *Step {
	radius := origRad
	if ft == fnArcCCW || ft == fnArcCCWRel {
		radius *= -1.0
	}

	vecep := endP
	vecab := vecep
	if !relative {
		vecab = vecep.Sub(g.Position())
	}
	var center Tuple
	var length float64
	// log.Printf("allArcs(endP=%v, radius=%v, relative=%v, ft=%v, opCode=%v): vecab=%v, pos=%v", endP, radius, relative, ft, opCode, vecab, g.Position())

	switch g.activePlane {
	default: // XY
		length = math.Sqrt(vecab.X()*vecab.X() + vecab.Y()*vecab.Y())
		if length < epsilon {
			log.Fatalf("distance between start and endPoint is zero (%v)", length)
		}
		var normal Tuple
		if radius < 0.0 {
			normal[0] = -0.5 * vecab.Y() / length
			normal[1] = 0.5 * vecab.X() / length
		} else {
			normal[0] = 0.5 * vecab.Y() / length
			normal[1] = -0.5 * vecab.X() / length
		}
		a := normal[0]*normal[0] + normal[1]*normal[1]
		b := 0.5*vecab[0]*normal[0] + 0.5*vecab[1]*normal[1]
		c := 0.25*vecab[0]*vecab[0] + 0.25*vecab[1]*vecab[1] - radius*radius
		d := b*b - 4.0*a*c
		// log.Printf("vecab=%v, length=%v, radius=%v, normal=%v, a=%v, b=%v, c=%v, d=%v", vecab, length, radius, normal, a, b, c, d)
		if d < 0.0 && math.Abs(d) < epsilon {
			d = 0.0
		}
		if d < 0.0 {
			log.Fatalf("allArcs(endP=%v, opCode=%q), (pos=%v)(A): radius %v is less than two times distance from start to end (D=%v)", endP, opCode, g.Position(), origRad, d)
		}
		// lambda := math.Sqrt(vecab.X()*vecab.X() + vecab.Z()*vecab.Z())
		lambda := (-b + math.Sqrt(d)) / (2.0 * a)
		center[0] = 0.5*vecab.X() + lambda*normal.X()
		center[1] = 0.5*vecab.Y() + lambda*normal.Y()
	case PlaneXZ:
		length = math.Sqrt(vecab.X()*vecab.X() + vecab.Z()*vecab.Z())
		if length < epsilon {
			log.Fatalf("distance between start and endPoint is zero (%v)", length)
		}
		var normal Tuple
		// In XZ we need the normal on the other side. Otherwise we'd get the wrong
		// behavior with respect to +/-radius for short/long anglular movement.
		if radius >= 0.0 {
			normal[0] = -0.5 * vecab.Z() / length
			normal[2] = 0.5 * vecab.X() / length
		} else {
			normal[0] = 0.5 * vecab.Z() / length
			normal[2] = -0.5 * vecab.X() / length
		}
		a := normal[0]*normal[0] + normal[2]*normal[2]
		b := 0.5*vecab[0]*normal[0] + 0.5*vecab[2]*normal[2]
		c := 0.25*vecab[0]*vecab[0] + 0.25*vecab[2]*vecab[2] - radius*radius
		d := b*b - 4.0*a*c
		if d < 0.0 && math.Abs(d) < epsilon {
			d = 0.0
		}
		if d < 0.0 {
			log.Fatalf("allArcs(endP=%v, opCode=%q), (pos=%v)(B): radius %v is less than two times distance from start to end (D=%v)", endP, opCode, g.Position(), origRad, d)
		}
		lambda := (-b + math.Sqrt(d)) / (2.0 * a)
		center[0] = 0.5*vecab.X() + lambda*normal.X()
		center[2] = 0.5*vecab.Z() + lambda*normal.Z()
	case PlaneYZ:
		length = math.Sqrt(vecab.Y()*vecab.Y() + vecab.Z()*vecab.Z())
		if length < epsilon {
			log.Fatalf("distance between start and endPoint is zero (%v)", length)
		}
		var normal Tuple
		if radius < 0.0 {
			normal[1] = -0.5 * vecab.Z() / length
			normal[2] = 0.5 * vecab.Y() / length
		} else {
			normal[1] = 0.5 * vecab.Z() / length
			normal[2] = -0.5 * vecab.Y() / length
		}
		a := normal[1]*normal[1] + normal[2]*normal[2]
		b := 0.5*vecab[1]*normal[1] + 0.5*vecab[2]*normal[2]
		c := 0.25*vecab[1]*vecab[1] + 0.25*vecab[2]*vecab[2] - radius*radius
		d := b*b - 4.0*a*c
		if d < 0.0 && math.Abs(d) < epsilon {
			d = 0.0
		}
		if d < 0.0 {
			log.Fatalf("allArcs(endP=%v, opCode=%q), (pos=%v)(C): radius %v is less than two times distance from start to end (D=%v)", endP, opCode, g.Position(), origRad, d)
		}
		lambda := (-b + math.Sqrt(d)) / (2.0 * a)
		center[1] = 0.5*vecab.Y() + lambda*normal.Y()
		center[2] = 0.5*vecab.Z() + lambda*normal.Z()
	}

	if math.Abs(radius)-(0.5*length) < -epsilon {
		log.Fatalf("allArcs(endP=%v, opCode=%q), (pos=%v)(D): radius %v is less than two times start-to-endpoint distance %v", endP, opCode, g.Position(), origRad, 0.5*length)
	}

	pos := g.Position()
	xyz := pos.Add(vecab)
	s := fmt.Sprintf("%v X%.8f Y%.8f", opCode, xyz.X(), xyz.Y())
	if math.Abs(xyz.Z()-pos.Z()) >= epsilon {
		s += fmt.Sprintf(" Z%.8f", xyz.Z())
	}

	if relative {
		pos = pos.Add(vecep)
	} else {
		pos = vecep
	}

	switch g.activePlane {
	default: // XY
		s += fmt.Sprintf(" I%.8f J%.8f", center.X(), center.Y())
	case PlaneXZ:
		s += fmt.Sprintf(" I%.8f K%.8f", center.X(), center.Z())
	case PlaneYZ:
		s += fmt.Sprintf(" J%.8f K%.8f", center.Y(), center.Z())
	}

	if opts != nil && opts.Turns > 0 {
		s += fmt.Sprintf(" P%v", opts.Turns)
	}

	return &Step{s: s, pos: pos}
}

// ArcCCW performs a counter clockwise arc from the current position
// to endpoint with given radius.
// The arc will be the shortest angular movement with positive radius and
// the largest angular movement with negative radius.
// Optional turns sets the number of turns to perform.
func (g *GCode) ArcCCW(endPoint Tuple, radius float64, opts *TurnsOption) *GCode {
	step := g.allArcs(endPoint, radius, false, fnArcCCW, "G3", opts)
	g.steps = append(g.steps, step)
	return g
}

// ArcCCWRel performs a counter clockwise arc from the current position
// to endpoint with given radius using relative offsets.
// The arc will be the shortest angular movement with positive radius and
// the largest angular movement with negative radius.
// Optional turns sets the number of turns to perform.
func (g *GCode) ArcCCWRel(endPoint Tuple, radius float64, opts *TurnsOption) *GCode {
	step := g.allArcs(endPoint, radius, true, fnArcCCWRel, "G3", opts)
	g.steps = append(g.steps, step)
	return g
}

// ArcCW performs a clockwise arc from the current position
// to endpoint with given radius.
// The arc will be the shortest angular movement with positive radius and
// the largest angular movement with negative radius.
// Optional turns sets the number of turns to perform.
func (g *GCode) ArcCW(endPoint Tuple, radius float64, opts *TurnsOption) *GCode {
	step := g.allArcs(endPoint, radius, false, fnArcCW, "G2", opts)
	g.steps = append(g.steps, step)
	return g
}

// ArcCWRel performs a clockwise arc from the current position
// to endpoint with given radius using relative offsets.
// The arc will be the shortest angular movement with positive radius and
// the largest angular movement with negative radius.
// Optional turns sets the number of turns to perform.
func (g *GCode) ArcCWRel(endPoint Tuple, radius float64, opts *TurnsOption) *GCode {
	step := g.allArcs(endPoint, radius, true, fnArcCWRel, "G2", opts)
	g.steps = append(g.steps, step)
	return g
}

// CircleCWRel performs a clockwise circle with radius length(centerPoint)
// and where centerPoint is the center point of the circle.
// The non-active plane coordinate may be used to create a helical movement.
// Optional turns sets the number of turns to perform.
func (g *GCode) CircleCW(centerPoint Tuple, opts *TurnsOption) *GCode {
	step := g.allCircles(centerPoint, false, fnCircleCW, "G2", opts)
	g.steps = append(g.steps, step)
	return g
}

// CircleCWRel performs a clockwise circle with radius length(centerPoint)
// and where centerPoint is the center point of the circle.
// The specified centerPoint is a relative position.
// The non-active plane coordinate may be used to create a helical movement.
// Optional turns sets the number of turns to perform.
func (g *GCode) CircleCWRel(centerPoint Tuple, opts *TurnsOption) *GCode {
	step := g.allCircles(centerPoint, true, fnCircleCWRel, "G2", opts)
	g.steps = append(g.steps, step)
	return g
}

// CircleCCWRel performs a clockwise circle with radius length(centerPoint)
// and where centerPoint is the center point of the circle.
// The non-active plane coordinate may be used to create a helical movement.
// Optional turns sets the number of turns to perform.
func (g *GCode) CircleCCW(centerPoint Tuple, opts *TurnsOption) *GCode {
	step := g.allCircles(centerPoint, false, fnCircleCCW, "G2", opts)
	g.steps = append(g.steps, step)
	return g
}

// CircleCCWRel performs a clockwise circle with radius length(centerPoint)
// and where centerPoint is the center point of the circle.
// The specified centerPoint is a relative position.
// The non-active plane coordinate may be used to create a helical movement.
// Optional turns sets the number of turns to perform.
func (g *GCode) CircleCCWRel(centerPoint Tuple, opts *TurnsOption) *GCode {
	step := g.allCircles(centerPoint, true, fnCircleCCWRel, "G2", opts)
	g.steps = append(g.steps, step)
	return g
}

type circleFnEnumT int

const (
	fnCircleCW circleFnEnumT = iota
	fnCircleCWRel
	fnCircleCCW
	fnCircleCCWRel
)

func (g *GCode) allCircles(arg0 Tuple, relative bool, ft circleFnEnumT, opCode string, opts *TurnsOption) *Step {
	endP := g.Position()
	var coor1, coor2 float64

	switch g.activePlane {
	default: // XY
		if relative {
			coor1 = arg0.X()
			coor2 = arg0.Y()
		} else {
			coor1 = arg0.X() - endP.X()
			coor2 = arg0.Y() - endP.Y()
		}
		endP[2] = arg0.Z()
	case PlaneXZ:
		if relative {
			coor1 = arg0.X()
			coor2 = arg0.Z()
		} else {
			coor1 = arg0.X() - endP.X()
			coor2 = arg0.Z() - endP.Z()
		}
		endP[1] = arg0.Y()
	case PlaneYZ:
		if relative {
			coor1 = arg0.Y()
			coor2 = arg0.Z()
		} else {
			coor1 = arg0.Y() - endP.Y()
			coor2 = arg0.Z() - endP.Z()
		}
		endP[0] = arg0.X()
	}

	if math.Sqrt(coor1*coor1+coor2*coor2) < epsilon {
		log.Fatal("radius is zero")
	}

	s := fmt.Sprintf("%v X%.8f Y%.8f", opCode, endP.X(), endP.Y())
	if math.Abs(endP.Z()-g.Position().Z()) >= epsilon {
		s += fmt.Sprintf(" Z%.8f", endP.Z())
	}
	pos := endP

	switch g.activePlane {
	default: // XY
		s += fmt.Sprintf(" I%.8f J%.8f", coor1, coor2)
	case PlaneXZ:
		s += fmt.Sprintf(" I%.8f K%.8f", coor1, coor2)
	case PlaneYZ:
		s += fmt.Sprintf(" J%.8f K%.8f", coor1, coor2)
	}

	if opts != nil && opts.Turns > 0 {
		s += fmt.Sprintf(" P%v", opts.Turns)
	}

	return &Step{s: s, pos: pos}
}

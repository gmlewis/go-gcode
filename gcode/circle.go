package gcode

import (
	"fmt"
	"log"
	"math"
)

// ArcOptions represents options for the arc and circle methods.
type ArcOptions struct {
	Turns int
}

type arcFnEnumT int

const (
	fnArcCW arcFnEnumT = iota
	fnArcCWRel
	fnArcCCW
	fnArcCCWRel
)

func (g *GCode) allArcs(endP Tuple, origRad float64, relative bool, ft arcFnEnumT, opCode string, opts *ArcOptions) *Step {
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
		if d < 0.0 && math.Abs(d) < epsilon {
			d = 0.0
		}
		if d < 0.0 {
			log.Fatalf("radius %v is less than two times distance from start to end (D=%v", origRad, d)
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
			log.Fatalf("radius %v is less than two times distance from start to end (D=%v", origRad, d)
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
			log.Fatalf("radius %v is less than two times distance from start to end (D=%v", origRad, d)
		}
		lambda := (-b + math.Sqrt(d)) / (2.0 * a)
		center[1] = 0.5*vecab.Y() + lambda*normal.Y()
		center[2] = 0.5*vecab.Z() + lambda*normal.Z()
	}

	if math.Abs(radius)-(0.5*length) < epsilon {
		log.Fatalf("radius %v is less than two times start-to-endpoint distance %v", origRad, 0.5*length)
	}

	pos := g.Position()
	xyz := pos.Add(vecab)
	s := fmt.Sprintf("%v X%.8f Y%.8f Z%.8f", opCode, xyz.X(), xyz.Y(), xyz.Z())
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
func (g *GCode) ArcCCW(endPoint Tuple, radius float64, opts *ArcOptions) *GCode {
	step := g.allArcs(endPoint, radius, false, fnArcCCW, "G3", opts)
	g.steps = append(g.steps, step)
	return g
}

// ArcCW performs a clockwise arc from the current position
// to endpoint with given radius.
// The arc will be the shortest angular movement with positive radius and
// the largest angular movement with negative radius.
// Optional turns sets the number of turns to perform.
func (g *GCode) ArcCW(endPoint Tuple, radius float64, opts *ArcOptions) *GCode {
	step := g.allArcs(endPoint, radius, false, fnArcCW, "G2", opts)
	g.steps = append(g.steps, step)
	return g
}

// CircleCWRel performs a clockwise (cw) circle with radius length(centerPoint)
// and where centerPoint is the center point of the circle.
// The specified centerPoint is a relative position.
// The non-active plane coordinate may be used to create a helical movement.
func (g *GCode) CircleCWRel(centerPoint Tuple) *GCode {
	pos := g.Position()
	g.steps = append(g.steps, &Step{
		s:   fmt.Sprintf("G2 X%.8f Y%.8f I%.8f J%.8f", pos.X(), pos.Y(), centerPoint.X(), centerPoint.Y()),
		pos: pos,
	})
	return g
}

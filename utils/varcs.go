package utils

import (
	"errors"
	"fmt"
	"log"
	"math"

	. "github.com/gmlewis/go-gcode/gcode"
)

const (
	defaultMaxL = 0.1
	defaultMaxA = 1 // degrees
)

// VOptions are options that alter the behavior of the arc and circle methods.
type VOptions struct {
	// Turns represents the number of turns.
	Turns int
	// MaxL is the maximum length.
	MaxL *float64
	// MaxA is the maximum angle (in degrees).
	MaxA *float64
	// ActPlane is the acting plane.
	ActPlane PlaneT
}

// GetMaxL gets the value of the MaxL option.
func (v *VOptions) GetMaxL() float64 {
	if v == nil || v.MaxL == nil {
		return defaultMaxL
	}
	return *v.MaxL
}

// GetMaxA gets the value of the MaxA option.
func (v *VOptions) GetMaxA() float64 {
	if v == nil || v.MaxA == nil {
		return defaultMaxA
	}
	return *v.MaxA
}

// VArcCW generates a clockwise arc.
func VArcCW(endPoint Tuple, radius float64, opts *VOptions) []Tuple {
	if opts == nil {
		opts = &VOptions{}
	}
	endPoint = planePtConvert(endPoint, opts.ActPlane)
	vl, err := genAll(false, false, endPoint, radius, opts.Turns, opts.GetMaxL(), opts.GetMaxA())
	if err != nil {
		log.Fatal(err)
	}

	toActivePlane(vl, opts.ActPlane)
	return vl
}

// VArcCCW generates a counter-clockwise arc.
func VArcCCW(endPoint Tuple, radius float64, opts *VOptions) []Tuple {
	if opts == nil {
		opts = &VOptions{}
	}
	endPoint = planePtConvert(endPoint, opts.ActPlane)
	vl, err := genAll(true, false, endPoint, radius, opts.Turns, opts.GetMaxL(), opts.GetMaxA())
	if err != nil {
		log.Fatal(err)
	}

	toActivePlane(vl, opts.ActPlane)
	return vl
}

// VCircleCW generates a clockwise circle.
func VCircleCW(centerPoint Tuple, opts *VOptions) []Tuple {
	if opts == nil {
		opts = &VOptions{}
	}
	centerPoint = planePtConvert(centerPoint, opts.ActPlane)
	radius := -math.Sqrt(centerPoint.X()*centerPoint.X() + centerPoint.Y()*centerPoint.Y())
	vl, err := genAll(false, true, centerPoint, radius, opts.Turns, opts.GetMaxL(), opts.GetMaxA())
	if err != nil {
		log.Fatal(err)
	}

	toActivePlane(vl, opts.ActPlane)
	return vl
}

// VCircleCCW generates a counter-clockwise circle.
func VCircleCCW(centerPoint Tuple, opts *VOptions) []Tuple {
	if opts == nil {
		opts = &VOptions{}
	}
	centerPoint = planePtConvert(centerPoint, opts.ActPlane)
	radius := math.Sqrt(centerPoint.X()*centerPoint.X() + centerPoint.Y()*centerPoint.Y())
	vl, err := genAll(true, true, centerPoint, radius, opts.Turns, opts.GetMaxL(), opts.GetMaxA())
	if err != nil {
		log.Fatal(err)
	}

	toActivePlane(vl, opts.ActPlane)
	return vl
}

func planePtConvert(pt Tuple, actPlane PlaneT) Tuple {
	switch actPlane {
	case PlaneXZ:
		return XYZ(pt.X(), pt.Z(), pt.Y())
	case PlaneYZ:
		return XYZ(pt.Y(), pt.Z(), pt.X())
	default:
		return pt
	}
}

// genAll calculates a vectorlist for all arcs CW and CCW in the XY plane.
// The arc runs from [0, 0] to endpoint (epIn) with given radius. The shortest
// arc route is taken on positive radius and the longest route on negative
// radius. The turns parameter adds a number of full turns in the CW or CCW
// directions.
// The arc is vectorized with straight lines with intervals with a maximum
// angle of maxA (in degrees) or a maximum length of maxL.
// The value resulting in the maximum number of steps is used.
// Circles are denoted with the epIn argument pointing to the center point of
// the circle. The rotation is always full in the desired direction.
func genAll(ccw, isCircle bool, epIn Tuple, radius float64, turns int, maxL, maxA float64) ([]Tuple, error) {
	if radius == 0 {
		return nil, errors.New("radius must not be zero")
	}
	if turns < 0 {
		return nil, errors.New("turns must not be negative")
	}
	if maxL <= 0 {
		return nil, errors.New("maxL must be positive")
	}

	origRadius := radius // for error messages
	if ccw {
		radius *= -1
	}

	var sp, ep, cp Tuple
	if !isCircle {
		// Arcs need to find the center point.
		ep = XY(epIn.X(), epIn.Y())
		// The normal of the sp <-> ep vector pointing to the center.
		normal := XY(ep.Y(), -ep.X()).Normalize()
		if radius < 0 {
			normal = normal.MultScalar(-1)
			radius *= -1
		}
		b := 0.5 * ep.Dot(normal)
		c := 0.25*ep.Dot(ep) - radius*radius
		d := b*b - 4*c
		if d < 0 {
			return nil, fmt.Errorf("radius (%.8f) is less than twice the distance from start to end (D=%.f)", origRadius, d)
		}
		cp = ep.MultScalar(0.5).Add(normal.MultScalar(-b + math.Sqrt(d)*0.5))
	} else {
		// Full circles have center point as argument and the endpoint
		// is always the same as the start point.
		ep = sp
		cp = XY(-ep.X(), -ep.Y())
	}

	aStart := math.Atan2(-cp.Y(), -cp.X())
	var aEnd float64
	switch {
	case !isCircle:
		aEnd = math.Atan2(ep.Y()-cp.Y(), ep.X()-cp.X())
	case ccw:
		aEnd = aStart + 2*math.Pi
	default:
		aEnd = aStart - 2*math.Pi
	}

	for aStart < 0 || aEnd < 0 {
		aStart += 2 * math.Pi
		aEnd += 2 * math.Pi
	}

	if ccw && aEnd < aStart {
		aEnd += 2 * math.Pi
	} else if !ccw && aStart < aEnd {
		aStart += 2 * math.Pi
	}

	// Should this be done above?
	if ccw {
		aEnd += float64(turns) * 2 * math.Pi
	} else {
		aEnd -= float64(turns) * 2 * math.Pi
	}

	aTot := aEnd - aStart
	var aStep float64
	var nStep int

	sa := int(math.Ceil(math.Abs(aTot) / (2.0 * math.Asin(0.5*maxL/cp.Magnitude()))))
	sl := int(math.Ceil(math.Abs(aTot) / (maxA * math.Pi / 180.0)))
	if sa <= 1 {
		sa = 2
	}
	if sl <= 1 {
		sl = 2
	}

	if sa > sl {
		aStep = aTot / float64(sa)
		nStep = sa
	} else {
		aStep = aTot / float64(sl)
		nStep = sl
	}

	var arcv []Tuple
	sp[2] = epIn.Z()
	ang := aStart
	offs := XYZ(radius*math.Cos(ang), radius*math.Sin(ang), 0)
	for i := 0; i < nStep; i++ {
		ang = float64(i)*aStep + aStart
		x := radius*math.Cos(ang) - offs.X()
		y := radius*math.Sin(ang) - offs.Y()
		z := float64(i)*epIn.Z()/float64(nStep) - offs.Z()
		arcv = append(arcv, XYZ(x, y, z))
	}

	if isCircle {
		arcv = append(arcv, sp)
	} else {
		arcv = append(arcv, epIn)
	}

	return arcv, nil
}

func toActivePlane(vl []Tuple, actPlane PlaneT) {
	if actPlane == "" || actPlane == PlaneXY {
		return
	}

	if actPlane == PlaneXZ {
		for i, v := range vl {
			vl[i] = XYZ(v.X(), v.Z(), v.Y())
		}
		return
	}

	for i, v := range vl {
		vl[i] = XYZ(v.Z(), v.X(), v.Y())
	}
}

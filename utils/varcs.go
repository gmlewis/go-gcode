package utils

import (
	"errors"
	"fmt"
	"math"

	"github.com/gmlewis/go-gcode/gcode"
)

const (
	defaultMaxL = 0.1
	defaultMaxA = 1 // degrees
)

// Plane represents a major axis.
type Plane int

const (
	PlaneXY Plane = iota
	PlaneYZ
	PlaneXZ
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
	ActPlane Plane
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
func VArcCW(endPoint gcode.Tuple, radius float64, opts *VOptions) ([]gcode.Tuple, error) {
	if opts == nil {
		opts = &VOptions{}
	}
	endPoint = planePtConvert(endPoint, opts.ActPlane)
	vl, err := genAll(false, false, endPoint, radius, opts.Turns, opts.GetMaxL(), opts.GetMaxA())
	if err != nil {
		return nil, err
	}

	return vl, nil
}

// VArcCCW generates a counter-clockwise arc.
func VArcCCW(endPoint gcode.Tuple, radius float64, opts *VOptions) ([]gcode.Tuple, error) {
	if opts == nil {
		opts = &VOptions{}
	}
	endPoint = planePtConvert(endPoint, opts.ActPlane)
	vl, err := genAll(true, false, endPoint, radius, opts.Turns, opts.GetMaxL(), opts.GetMaxA())
	if err != nil {
		return nil, err
	}

	return vl, nil
}

// VCircleCW generates a clockwise circle.
func VCircleCW(centerPoint gcode.Tuple, opts *VOptions) ([]gcode.Tuple, error) {
	if opts == nil {
		opts = &VOptions{}
	}
	centerPoint = planePtConvert(centerPoint, opts.ActPlane)
	radius := -math.Sqrt(centerPoint.X()*centerPoint.X() + centerPoint.Y()*centerPoint.Y())
	vl, err := genAll(false, true, centerPoint, radius, opts.Turns, opts.GetMaxL(), opts.GetMaxA())
	if err != nil {
		return nil, err
	}

	return vl, nil
}

// VCircleCCW generates a counter-clockwise circle.
func VCircleCCW(centerPoint gcode.Tuple, opts *VOptions) ([]gcode.Tuple, error) {
	if opts == nil {
		opts = &VOptions{}
	}
	centerPoint = planePtConvert(centerPoint, opts.ActPlane)
	radius := math.Sqrt(centerPoint.X()*centerPoint.X() + centerPoint.Y()*centerPoint.Y())
	vl, err := genAll(true, true, centerPoint, radius, opts.Turns, opts.GetMaxL(), opts.GetMaxA())
	if err != nil {
		return nil, err
	}

	return vl, nil
}

func planePtConvert(pt gcode.Tuple, actPlane Plane) gcode.Tuple {
	switch actPlane {
	case PlaneXZ:
		return gcode.Point(pt.X(), pt.Z(), pt.Y())
	case PlaneYZ:
		return gcode.Point(pt.Y(), pt.Z(), pt.X())
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
func genAll(ccw, circ bool, epIn gcode.Tuple, radius float64, turns int, maxL, maxA float64) ([]gcode.Tuple, error) {
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

	var sp, ep, cp gcode.Tuple
	if !circ {
		// Arcs need to find the center point.
		ep = gcode.XY(epIn.X(), epIn.Y())
		// The normal of the sp <-> ep vector pointing to the center.
		normal := gcode.XY(ep.Y(), -ep.X()).Normalize()
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
		cp = gcode.XY(-ep.X(), -ep.Y())
	}

	aStart := math.Atan2(-cp.Y(), -cp.X())
	var aEnd float64

	var arcv []gcode.Tuple
	return arcv, nil
}

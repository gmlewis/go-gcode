package utils

import (
	"math"

	"github.com/gmlewis/go-gcode/gcode"
)

const (
	defaultFlatness = 1e-4
	defaultMinL     = 0.1 // mm
)

// VBezierOptions represents options for VBezier3.
type VBezierOptions struct {
	Flatness float64
	MinL     float64
}

// Vectorize cubic Bezier curves by recursive subdivision
// using De Casteljau's algorithm
// Input:
// b0		: node left
// b1		: control point left
// b2		: control point right
// b3		: node right
// flatness	: resisidual 2-2*|cos(phi)| to determine direction of line
// minl		: minimum distance between points
// Output:
// Vectorlist of points excluding b0 and exactly ending at b3
func VBezier3(b0, b1, b2, b3 gcode.Tuple, opts *VBezierOptions) []gcode.Tuple {
	flatness, minL := defaultFlatness, defaultMinL
	if opts != nil {
		flatness = opts.Flatness
		minL = opts.MinL
	}

	return genBezier(b0, b1, b2, b3, flatness, minL)
}

func genBezier(b0, b1, b2, b3 gcode.Tuple, flatness, minL float64) []gcode.Tuple {
	l := b0.Add(b1).MultScalar(0.5)
	m := b1.Add(b2).MultScalar(0.5)
	r := b2.Add(b3).MultScalar(0.5)
	lm := l.Add(m).MultScalar(0.5)
	rm := r.Add(m).MultScalar(0.5)
	t := lm.Add(rm).MultScalar(0.5)

	// If the b0-top-b3 triangle is under the minimum length then return
	// this triangle. Extreme cases where the control points pull out the
	// curve to the sides is handled by ensuring that left and right
	// node-to-conrol points must adhere to the same minimum length.
	if (t.Sub(b0).Magnitude()+b3.Sub(t).Magnitude()) < minL && r.Sub(l).Magnitude() < minL {
		return []gcode.Tuple{t, b3}
	}

	// cos(Angles) as seen from both sides should sum to 2.0 within
	// tolerance to be co-linear but from opposite direction.
	cplv := b0.Sub(m).Normalize().Dot(lm.Sub(t).Normalize())
	cprv := b3.Sub(m).Normalize().Dot(t.Sub(rm).Normalize())

	// We're done if we are co-linear.
	// The angle test is scaled by the distance between the end-points
	// relative to the minL argument. This reduces the perpendicular
	// deviation from the curve to a minimum.
	if b3.Sub(b0).Magnitude()/minL*(2.0-math.Abs(cprv-cplv)) < flatness {
		return []gcode.Tuple{b3}
	}

	// Otherwise, return the points with bisected recursion.
	left := genBezier(b0, l, lm, t, flatness, minL)
	right := genBezier(t, rm, r, b3, flatness, minL)
	result := append([]gcode.Tuple{}, left[0:len(left)-1]...)
	result = append(result, t)
	result = append(result, right[0:len(right)-1]...)
	result = append(result, b3)

	return result
}

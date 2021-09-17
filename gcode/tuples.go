package gcode

import (
	"log"
	"math"
	"strconv"
	"strings"
)

const (
	epsilon = 1e-8
)

// Tuple is a 4-float vector.
type Tuple [4]float64

// X returns the X value of the Tuple.
func (t Tuple) X() float64 {
	if len(t) == 0 {
		return 0
	}
	return t[0]
}

// Y returns the Y value of the Tuple.
func (t Tuple) Y() float64 {
	if len(t) == 0 {
		return 0
	}
	return t[1]
}

// Z returns the Z value of the Tuple.
func (t Tuple) Z() float64 {
	if len(t) == 0 {
		return 0
	}
	return t[2]
}

// This is not a ray-tracer - do not use W.
// // W returns the W value of the Tuple.
// func (t Tuple) W() float64 {
// 	if len(t) == 0 {
// 		return 0
// 	}
// 	return t[3]
// }
//
// // IsPoint identifies the Tuple as a Point.
// func (t Tuple) IsPoint() bool {
// 	return len(t) > 3 && t[3] == 1.0
// }
//
// // IsVector identifies the Tuple as a Vector.
// func (t Tuple) IsVector() bool {
// 	return len(t) == 0 || t[3] == 0.0
// }
//
// // Point returns a new Tuple as a Point.
// func Point(x, y, z float64) Tuple {
// 	return Tuple{x, y, z, 1}
// }
//
// // Vector returns a new Tuple as a Vector.
// func Vector(x, y, z float64) Tuple {
// 	return Tuple{x, y, z, 0}
// }

// Equal tests if two Tuples are equal.
func (t Tuple) Equal(other Tuple) bool {
	return math.Abs(t.X()-other.X()) < epsilon &&
		math.Abs(t.Y()-other.Y()) < epsilon &&
		math.Abs(t.Z()-other.Z()) < epsilon
	// this is not a ray-tracer - don't compare w.
	// math.Abs(t.W()-other.W()) < epsilon
}

// Add adds two Tuples and returns a new one.
func (t Tuple) Add(other Tuple) Tuple {
	return Tuple{
		t.X() + other.X(),
		t.Y() + other.Y(),
		t.Z() + other.Z(),
		1, // this is not a ray-tracer - keep w=1 // t.W() + other.W(),
	}
}

// Offset adds an offset to multiple points.
func (t Tuple) Offset(points ...Tuple) []Tuple {
	var result []Tuple
	for _, p := range points {
		newPt := t.Add(p)
		newPt[3] = 1
		result = append(result, newPt)
	}
	return result
}

// Sub subtracts two Tuples and returns a new one.
func (t Tuple) Sub(other Tuple) Tuple {
	return Tuple{
		t.X() - other.X(),
		t.Y() - other.Y(),
		t.Z() - other.Z(),
		1, // this is not a ray-tracer - keep w=1 // t.W() - other.W(),
	}
}

// Negate negates a Tuple.
func (t Tuple) Negate() Tuple {
	return Tuple{
		-t.X(),
		-t.Y(),
		-t.Z(),
		1, // this is not a ray-tracer - keep w=1 // -t.W(),
	}
}

// MultScalar multiplies a tuple by a scalar.
func (t Tuple) MultScalar(f float64) Tuple {
	return Tuple{
		f * t.X(),
		f * t.Y(),
		f * t.Z(),
		1, // this is not a ray-tracer - keep w=1 // f * t.W(),
	}
}

// DivScalar divides a tuple by a scalar.
func (t Tuple) DivScalar(f float64) Tuple {
	return t.MultScalar(1 / f)
}

// Magnitude computes the magnitude or length of a vector (Tuple).
func (t Tuple) Magnitude() float64 {
	return math.Sqrt(
		t.X()*t.X() +
			t.Y()*t.Y() +
			t.Z()*t.Z())
	// this is not a ray-tracer - don't consider w. // t.W()*t.W())
}

// Normalize normalizes a vector to a unit vector (of length 1).
func (t Tuple) Normalize() Tuple {
	return t.DivScalar(t.Magnitude())
}

// Dot computes the dot product (aka "scalar product" or "inner product")
// of two vectors (Tuples). The dot product is the cosine of the angle
// between two unit vectors.
func (t Tuple) Dot(other Tuple) float64 {
	return t.X()*other.X() +
		t.Y()*other.Y() +
		t.Z()*other.Z() +
		t[3]*other[3] // Must consider W to perform translation in matrix multiplies!
}

// Cross computes the cross product of two vectors (order matters and this
// implements t cross other).
func (t Tuple) Cross(other Tuple) Tuple {
	return XYZ( // this is not a ray-tracer - don't use: Vector(
		t.Y()*other.Z()-t.Z()*other.Y(),
		t.Z()*other.X()-t.X()*other.Z(),
		t.X()*other.Y()-t.Y()*other.X(),
	)
}

// Reflect reflects a vector around a prodived normal vector.
func (t Tuple) Reflect(normal Tuple) Tuple {
	return t.Sub(normal.MultScalar(2.0 * t.Dot(normal)))
}

// Reverse makes a reversed copy of the slice of Tuples.
func Reverse(vs []Tuple) []Tuple {
	v := make([]Tuple, 0, len(vs))
	for i := len(vs) - 1; i >= 0; i-- {
		v = append(v, vs[i])
	}
	return v
}

// LastXY returns the last XY point from a slice and sets Z=0.
func LastXY(vs []Tuple) Tuple {
	if len(vs) == 0 {
		return XY(0, 0)
	}
	p := vs[len(vs)-1]
	return XY(p.X(), p.Y())
}

// X returns a Point with only X set.
func X(x float64) Tuple { return XYZ(x, 0, 0) }

// Y returns a Point with only Y set.
func Y(y float64) Tuple { return XYZ(0, y, 0) }

// Z returns a Point with only Z set.
func Z(z float64) Tuple { return XYZ(0, 0, z) }

// XY returns a Point with only X and Y set.
func XY(x, y float64) Tuple { return XYZ(x, y, 0) }

// YZ returns a Point with only Y and Z set.
func YZ(y, z float64) Tuple { return XYZ(0, y, z) }

// XZ returns a Point with only X and Z set.
func XZ(x, z float64) Tuple { return XYZ(x, 0, z) }

// XYZ returns a Point.
func XYZ(x, y, z float64) Tuple { return Tuple{x, y, z, 1} }

// Path returns a list of tuples, keeping track of the changes
// made at each step.
//
// They may include full XYZ coordinates for each point, however, specifying
// only the coordinates that differ from the previous is enough, except for the
// first one, which should include a Z-coordinate (missing coordinates will
// start at zero). A "-" or empty string means to keep the previous value.
func Path(ss ...string) []Tuple {
	var result []Tuple
	lastPos := XYZ(0, 0, 0)
	for _, s := range ss {
		lastPos = parseDiffs(s, lastPos)
		result = append(result, lastPos)
	}
	return result
}

func parseDiffs(s string, lastPos Tuple) Tuple {
	parts := strings.Split(s, ",")
	for i, part := range parts {
		p := strings.TrimSpace(part)
		if p == "" || p == "-" {
			continue
		}
		v, err := strconv.ParseFloat(p, 64)
		if err != nil {
			log.Fatal(err)
		}
		if i > 2 {
			log.Fatalf("too many fields in path: %q", s)
		}
		lastPos[i] = v
	}
	return lastPos
}

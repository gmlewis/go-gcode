package gcode

import "math"

// Translation returns a 4x4 translation matrix.
func Translation(x, y, z float64) M4 {
	return M4{
		Tuple{1, 0, 0, x},
		Tuple{0, 1, 0, y},
		Tuple{0, 0, 1, z},
		Tuple{0, 0, 0, 1},
	}
}

// Translate translates a 4x4 matrix and returns a new one.
func (m M4) Translate(x, y, z float64) M4 {
	t := Translation(x, y, z)
	return t.Mult(m)
}

// Scaling returns a 4x4 scaling matrix.
func Scaling(x, y, z float64) M4 {
	return M4{
		Tuple{x, 0, 0, 0},
		Tuple{0, y, 0, 0},
		Tuple{0, 0, z, 0},
		Tuple{0, 0, 0, 1},
	}
}

// Scale scales a 4x4 matrix and returns a new one.
func (m M4) Scale(x, y, z float64) M4 {
	t := Scaling(x, y, z)
	return t.Mult(m)
}

// RotationX returns a 4x4 rotation matrix clockwise about the X axis
// using the left-hand rule.
func RotationX(radians float64) M4 {
	c := math.Cos(radians)
	s := math.Sin(radians)
	return M4{
		Tuple{1, 0, 0, 0},
		Tuple{0, c, -s, 0},
		Tuple{0, s, c, 0},
		Tuple{0, 0, 0, 1},
	}
}

// RotateX rotates a 4x4 matrix clockwise about the X axis
// using the left-hand rule and returns a new 4x4 matrix.
func (m M4) RotateX(radians float64) M4 {
	t := RotationX(radians)
	return t.Mult(m)
}

// RotationY returns a 4x4 rotation matrix clockwise about the Y axis
// using the left-hand rule.
func RotationY(radians float64) M4 {
	c := math.Cos(radians)
	s := math.Sin(radians)
	return M4{
		Tuple{c, 0, s, 0},
		Tuple{0, 1, 0, 0},
		Tuple{-s, 0, c, 0},
		Tuple{0, 0, 0, 1},
	}
}

// RotateX rotates a 4x4 matrix clockwise about the Y axis
// using the left-hand rule and returns a new 4x4 matrix.
func (m M4) RotateY(radians float64) M4 {
	t := RotationY(radians)
	return t.Mult(m)
}

// RotationZ returns a 4x4 rotation matrix clockwise about the Z axis
// using the left-hand rule.
func RotationZ(radians float64) M4 {
	c := math.Cos(radians)
	s := math.Sin(radians)
	return M4{
		Tuple{c, -s, 0, 0},
		Tuple{s, c, 0, 0},
		Tuple{0, 0, 1, 0},
		Tuple{0, 0, 0, 1},
	}
}

// RotateX rotates a 4x4 matrix clockwise about the Z axis
// using the left-hand rule and returns a new 4x4 matrix.
func (m M4) RotateZ(radians float64) M4 {
	t := RotationZ(radians)
	return t.Mult(m)
}

// Shearing returns a 4x4 shearing matrix.
func Shearing(xy, xz, yx, yz, zx, zy float64) M4 {
	return M4{
		Tuple{1, xy, xz, 0},
		Tuple{yx, 1, yz, 0},
		Tuple{zx, zy, 1, 0},
		Tuple{0, 0, 0, 1},
	}
}

// Shear shears a 4x4 matrix and returns a new one.
func (m M4) Shear(xy, xz, yx, yz, zx, zy float64) M4 {
	t := Shearing(xy, xz, yx, yz, zx, zy)
	return t.Mult(m)
}

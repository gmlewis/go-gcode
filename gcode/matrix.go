package gcode

import "log"

// M4 is a 4x4 matrix.
type M4 [4]Tuple

// Get returns a value within the matrix.
func (m M4) Get(row, col int) float64 {
	return m[row][col]
}

// Equal tests if two matrices are equal.
func (m M4) Equal(other M4) bool {
	return m[0].Equal(other[0]) &&
		m[1].Equal(other[1]) &&
		m[2].Equal(other[2]) &&
		m[3].Equal(other[3])
}

// full4Dot computes the dot product (aka "scalar product" or "inner product")
// of two vectors (Tuples). The dot product is the cosine of the angle
// between two unit vectors.
// This version is needed for M4 full matrix multiples but should not be
// used for general dot products of XYZ tuples in GCode scripts, as the W
// factor messes up the dot products on points since they are actually vectors.
func (t Tuple) full4Dot(other Tuple) float64 {
	return t.X()*other.X() +
		t.Y()*other.Y() +
		t.Z()*other.Z() +
		t[3]*other[3] // Must consider W to perform translation in matrix multiplies!
}

// Mult multiplies two M4 matrices. Order is important.
func (m M4) Mult(other M4) M4 {
	oc := M4{other.Column(0), other.Column(1), other.Column(2), other.Column(3)}
	return M4{
		Tuple{m[0].full4Dot(oc[0]), m[0].full4Dot(oc[1]), m[0].full4Dot(oc[2]), m[0].full4Dot(oc[3])},
		Tuple{m[1].full4Dot(oc[0]), m[1].full4Dot(oc[1]), m[1].full4Dot(oc[2]), m[1].full4Dot(oc[3])},
		Tuple{m[2].full4Dot(oc[0]), m[2].full4Dot(oc[1]), m[2].full4Dot(oc[2]), m[2].full4Dot(oc[3])},
		Tuple{m[3].full4Dot(oc[0]), m[3].full4Dot(oc[1]), m[3].full4Dot(oc[2]), m[3].full4Dot(oc[3])},
	}
}

// MultTuple multiples an M4 matrix by a tuple.
func (m M4) MultTuple(other Tuple) Tuple {
	return Tuple{
		m[0].full4Dot(other),
		m[1].full4Dot(other),
		m[2].full4Dot(other),
		m[3].full4Dot(other),
	}
}

// Transform transforms all points by the given transform.
func (m M4) Transform(points ...Tuple) []Tuple {
	var result []Tuple
	for _, p := range points {
		newPt := m.MultTuple(p)
		newPt[3] = 1
		result = append(result, newPt)
	}
	return result
}

// M4Identity returns a 4x4 identity matrix.
func M4Identity() M4 {
	return M4{
		Tuple{1, 0, 0, 0},
		Tuple{0, 1, 0, 0},
		Tuple{0, 0, 1, 0},
		Tuple{0, 0, 0, 1},
	}
}

// Transpose transposes a 4x4 matrix.
func (m M4) Transpose() M4 {
	return M4{
		m.Column(0),
		m.Column(1),
		m.Column(2),
		m.Column(3),
	}
}

// Submatrix returns a 3x3 submatrix with a row and column removed from a 4x4 matrix.
func (m M4) Submatrix(row, col int) M3 {
	v := func(r, c int) float64 {
		if r >= row {
			r++
		}
		if c >= col {
			c++
		}
		return m[r][c]
	}
	return M3{
		Tuple{v(0, 0), v(0, 1), v(0, 2)},
		Tuple{v(1, 0), v(1, 1), v(1, 2)},
		Tuple{v(2, 0), v(2, 1), v(2, 2)},
	}
}

// Minor returns the determinant of a submatrix of a 4x4 matrix.
func (m M4) Minor(row, col int) float64 {
	return m.Submatrix(row, col).Determinant()
}

// Cofactor returns the cofactor of a submatrix of a 4x4 matrix.
func (m M4) Cofactor(row, col int) float64 {
	minor := m.Minor(row, col)
	if (row+col)%2 == 1 {
		minor = -minor
	}
	return minor
}

// Determinant returns the determinant of the 4x4 matrix.
func (m M4) Determinant() float64 {
	return m[0][0]*m.Cofactor(0, 0) + m[0][1]*m.Cofactor(0, 1) + m[0][2]*m.Cofactor(0, 2) + m[0][3]*m.Cofactor(0, 3)
}

// Invertible returns the invertibility of the 4x4 matrix.
func (m M4) Invertible() bool {
	return m.Determinant() != 0
}

// Inverse calculates the inverse of the 4x4 matrix.
func (m M4) Inverse() M4 {
	d := m.Determinant()
	if d == 0 {
		log.Fatalf("cannot take inverse of non-invertible matrix: %v", m)
	}

	v := func(row, col int) float64 {
		return m.Cofactor(col, row) / d // transpose happens here: row,col=>col,row
	}

	return M4{
		Tuple{v(0, 0), v(0, 1), v(0, 2), v(0, 3)},
		Tuple{v(1, 0), v(1, 1), v(1, 2), v(1, 3)},
		Tuple{v(2, 0), v(2, 1), v(2, 2), v(2, 3)},
		Tuple{v(3, 0), v(3, 1), v(3, 2), v(3, 3)},
	}
}

// Column returns a column of the matrix as a Tuple.
func (m M4) Column(col int) Tuple {
	return Tuple{m[0][col], m[1][col], m[2][col], m[3][col]}
}

// M3 is a 3x3 matrix.
type M3 [3]Tuple

// Get returns a value within the matrix.
func (m M3) Get(row, col int) float64 {
	return m[row][col]
}

// Equal tests if two matrices are equal.
func (m M3) Equal(other M3) bool {
	return m[0].Equal(other[0]) &&
		m[1].Equal(other[1]) &&
		m[2].Equal(other[2])
}

// Submatrix returns a 2x2 submatrix with a row and column removed from a 3x3 matrix.
func (m M3) Submatrix(row, col int) M2 {
	v := func(r, c int) float64 {
		if r >= row {
			r++
		}
		if c >= col {
			c++
		}
		return m[r][c]
	}
	return M2{
		Tuple{v(0, 0), v(0, 1)},
		Tuple{v(1, 0), v(1, 1)},
	}
}

// Minor returns the determinant of a submatrix of a 3x3 matrix.
func (m M3) Minor(row, col int) float64 {
	return m.Submatrix(row, col).Determinant()
}

// Cofactor returns the cofactor of a submatrix of a 3x3 matrix.
func (m M3) Cofactor(row, col int) float64 {
	minor := m.Minor(row, col)
	if (row+col)%2 == 1 {
		minor = -minor
	}
	return minor
}

// Determinant returns the determinant of the 3x3 matrix.
func (m M3) Determinant() float64 {
	return m[0][0]*m.Cofactor(0, 0) + m[0][1]*m.Cofactor(0, 1) + m[0][2]*m.Cofactor(0, 2)
}

// M2 is a 2x2 matrix.
type M2 [2]Tuple

// Get returns a value within the matrix.
func (m M2) Get(row, col int) float64 {
	return m[row][col]
}

// Equal tests if two matrices are equal.
func (m M2) Equal(other M2) bool {
	return m[0].Equal(other[0]) &&
		m[1].Equal(other[1])
}

// Determinant finds the determinant of a 2x2 matrix.
func (m M2) Determinant() float64 {
	return m[0][0]*m[1][1] - m[0][1]*m[1][0]
}

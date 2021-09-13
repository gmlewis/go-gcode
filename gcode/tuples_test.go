package gcode

import "testing"

func TestPath(t *testing.T) {
	got := Path(
		"10, 0, -1",
		"15",
		"20",
		"25",
		"-, 5",
		"20",
		"15, -, -2",
		"10",
	)

	want := []Tuple{
		XYZ(10, 0, -1),
		XYZ(15, 0, -1),
		XYZ(20, 0, -1),
		XYZ(25, 0, -1),
		XYZ(25, 5, -1),
		XYZ(20, 5, -1),
		XYZ(15, 5, -2),
		XYZ(10, 5, -2),
	}

	if len(got) != len(want) {
		t.Fatalf("got = %v tuples, want %v", len(got), len(want))
	}

	for i, p := range got {
		if !p.Equal(want[i]) {
			t.Errorf("i=%v: got %v, want %v", i, p, want[i])
		}
	}
}

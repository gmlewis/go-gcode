package gcode

import "testing"

func TestTranslation_Transform(t *testing.T) {
	lst := []Tuple{
		XYZ(0, 0, 0),
		XYZ(1, 2, 3),
	}

	got := Translation(10, 20, 30).Transform(lst...)
	want := []Tuple{
		XYZ(10, 20, 30),
		XYZ(11, 22, 33),
	}

	if len(got) != len(want) {
		t.Fatalf("Translation.Transform = %v points, want %v", len(got), len(want))
	}

	for i, g := range got {
		if !g.Equal(want[i]) {
			t.Errorf("Translation.Transform[%v] = %v, want %v", i, g, want[i])
		}
	}
}

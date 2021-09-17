package gcode

import "testing"

func TestMoveOrGoto(t *testing.T) {
	g := New()
	g.GotoXYZ(XYZ(-145, -30, -1))
	g.MoveZ(Z(0))
	g.MoveXYZ(XYZ(-139.5, 30, 0))
	g.MoveZ(Z(-1))

	got := g.String()
	want := `G0 X-145.00000000 Y-30.00000000 Z-1.00000000
G1 Z0.00000000
G1 X-139.50000000 Y30.00000000
G1 Z-1.00000000
`

	if got != want {
		t.Errorf("moveOrGoto =\n%v\nwant:\n%v", got, want)
	}
}

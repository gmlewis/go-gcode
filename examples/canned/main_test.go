package main

import (
	"strings"
	"testing"
)

func TestGCMC(t *testing.T) {
	g := gcmc()
	got := strings.Split(g.String(), "\n")
	want := strings.Split(gcmcOut, "\n")

	if len(got) != len(want) {
		t.Errorf("gcmc = %v lines, want %v", len(got), len(want))
	}

	for i, line := range got {
		if i == 1 {
			continue
		}
		if i < len(want) && line != want[i] {
			t.Fatalf("gcmc line #%v:\n%v\nwant:\n%v", i+1, line, want[i])
		}
	}
}

var gcmcOut = `(go-gcode compiled code, do not change)
(2021-09-11 08:33:05)
(-- prologue begin --)
G17 ( Use XY plane )
G21 ( Use mm )
G40 ( Cancel cutter radius compensation )
G49 ( Cancel tool length compensation )
G54 ( Default coordinate system )
G80 ( Cancel canned cycle )
G90 ( Use absolute distance mode )
G94 ( Units Per Minute feed rate mode )
G64 ( Enable path blending for best speed )
(-- prologue end --)
F400.00000000
G0 X0.00000000 Y0.00000000 Z0.00000000
(MSG,Canned peck drill without return-to-Z)
G0 X-1.00000000 Y-1.00000000 Z10.00000000
(-- canned_drill_peck R-plane=5 peck-increment=2.2 return-to-old-Z=false --)
G61
G0 X10.00000000 Y0.00000000
G0 Z5.00000000
G1 Z2.80000000
G0 Z5.00000000
G0 Z3.00000000
G1 Z0.60000000
G0 Z5.00000000
G0 Z0.80000000
G1 Z-1.00000000
G0 Z5.00000000
G0 X15.00000000
G1 Z2.80000000
G0 Z5.00000000
G0 Z3.00000000
G1 Z0.60000000
G0 Z5.00000000
G0 Z0.80000000
G1 Z-1.00000000
G0 Z5.00000000
G0 X20.00000000
G1 Z2.80000000
G0 Z5.00000000
G0 Z3.00000000
G1 Z0.60000000
G0 Z5.00000000
G0 Z0.80000000
G1 Z-1.00000000
G0 Z5.00000000
G0 X25.00000000
G1 Z2.80000000
G0 Z5.00000000
G0 Z3.00000000
G1 Z0.60000000
G0 Z5.00000000
G0 Z0.80000000
G1 Z-1.00000000
G0 Z5.00000000
G0 Y5.00000000
G1 Z2.80000000
G0 Z5.00000000
G0 Z3.00000000
G1 Z0.60000000
G0 Z5.00000000
G0 Z0.80000000
G1 Z-1.00000000
G0 Z5.00000000
G0 X20.00000000
G1 Z2.80000000
G0 Z5.00000000
G0 Z3.00000000
G1 Z0.60000000
G0 Z5.00000000
G0 Z0.80000000
G1 Z-1.00000000
G0 Z5.00000000
G0 X15.00000000
G1 Z2.80000000
G0 Z5.00000000
G0 Z3.00000000
G1 Z0.60000000
G0 Z5.00000000
G0 Z0.80000000
G1 Z-1.60000000
G0 Z5.00000000
G0 Z-1.40000000
G1 Z-2.00000000
G0 Z5.00000000
G0 X10.00000000
G1 Z2.80000000
G0 Z5.00000000
G0 Z3.00000000
G1 Z0.60000000
G0 Z5.00000000
G0 Z0.80000000
G1 Z-1.60000000
G0 Z5.00000000
G0 Z-1.40000000
G1 Z-2.00000000
G0 Z5.00000000
(-- end canned_drill_peck --)
(MSG,Canned peck drill with return-to-Z)
G0 X-1.00000000 Y9.00000000 Z10.00000000
(-- canned_drill_peck R-plane=5 peck-increment=2.2 return-to-old-Z=true --)
G61
G0 X10.00000000 Y10.00000000
G0 Z5.00000000
G1 Z2.80000000
G0 Z5.00000000
G0 Z3.00000000
G1 Z0.60000000
G0 Z5.00000000
G0 Z0.80000000
G1 Z-1.00000000
G0 Z5.00000000
G0 Z10.00000000
G0 X15.00000000
G0 Z5.00000000
G1 Z2.80000000
G0 Z5.00000000
G0 Z3.00000000
G1 Z0.60000000
G0 Z5.00000000
G0 Z0.80000000
G1 Z-1.00000000
G0 Z5.00000000
G0 Z10.00000000
G0 X20.00000000
G0 Z5.00000000
G1 Z2.80000000
G0 Z5.00000000
G0 Z3.00000000
G1 Z0.60000000
G0 Z5.00000000
G0 Z0.80000000
G1 Z-1.00000000
G0 Z5.00000000
G0 Z10.00000000
G0 X25.00000000
G0 Z5.00000000
G1 Z2.80000000
G0 Z5.00000000
G0 Z3.00000000
G1 Z0.60000000
G0 Z5.00000000
G0 Z0.80000000
G1 Z-1.00000000
G0 Z5.00000000
G0 Z10.00000000
G0 Y15.00000000
G0 Z5.00000000
G1 Z2.80000000
G0 Z5.00000000
G0 Z3.00000000
G1 Z0.60000000
G0 Z5.00000000
G0 Z0.80000000
G1 Z-1.00000000
G0 Z5.00000000
G0 Z10.00000000
G0 X20.00000000
G0 Z5.00000000
G1 Z2.80000000
G0 Z5.00000000
G0 Z3.00000000
G1 Z0.60000000
G0 Z5.00000000
G0 Z0.80000000
G1 Z-1.00000000
G0 Z5.00000000
G0 Z10.00000000
G0 X15.00000000
G0 Z5.00000000
G1 Z2.80000000
G0 Z5.00000000
G0 Z3.00000000
G1 Z0.60000000
G0 Z5.00000000
G0 Z0.80000000
G1 Z-1.60000000
G0 Z5.00000000
G0 Z-1.40000000
G1 Z-2.00000000
G0 Z5.00000000
G0 Z10.00000000
G0 X10.00000000
G0 Z5.00000000
G1 Z2.80000000
G0 Z5.00000000
G0 Z3.00000000
G1 Z0.60000000
G0 Z5.00000000
G0 Z0.80000000
G1 Z-1.60000000
G0 Z5.00000000
G0 Z-1.40000000
G1 Z-2.00000000
G0 Z5.00000000
G0 Z10.00000000
(-- end canned_drill_peck --)
(MSG,Canned drill without return-to-Z)
G0 X-1.00000000 Y19.00000000
(-- canned_drill R-plane=5 dwelling=-1 return-to-old-Z=false --)
G61
G0 X10.00000000 Y20.00000000
G0 Z5.00000000
G1 Z-1.00000000
G0 Z5.00000000
G0 X15.00000000
G1 Z-1.00000000
G0 Z5.00000000
G0 X20.00000000
G1 Z-1.00000000
G0 Z5.00000000
G0 X25.00000000
G1 Z-1.00000000
G0 Z5.00000000
G0 Y25.00000000
G1 Z-1.00000000
G0 Z5.00000000
G0 X20.00000000
G1 Z-1.00000000
G0 Z5.00000000
G0 X15.00000000
G1 Z-2.00000000
G0 Z5.00000000
G0 X10.00000000
G1 Z-2.00000000
G0 Z5.00000000
(-- end canned_drill --)
(MSG,Canned drill with return-to-Z)
G0 X-1.00000000 Y29.00000000 Z10.00000000
(-- canned_drill R-plane=5 dwelling=-1 return-to-old-Z=true --)
G61
G0 X10.00000000 Y30.00000000
G0 Z5.00000000
G1 Z-1.00000000
G0 Z10.00000000
G0 X15.00000000
G0 Z5.00000000
G1 Z-1.00000000
G0 Z10.00000000
G0 X20.00000000
G0 Z5.00000000
G1 Z-1.00000000
G0 Z10.00000000
G0 X25.00000000
G0 Z5.00000000
G1 Z-1.00000000
G0 Z10.00000000
G0 Y35.00000000
G0 Z5.00000000
G1 Z-1.00000000
G0 Z10.00000000
G0 X20.00000000
G0 Z5.00000000
G1 Z-1.00000000
G0 Z10.00000000
G0 X15.00000000
G0 Z5.00000000
G1 Z-2.00000000
G0 Z10.00000000
G0 X10.00000000
G0 Z5.00000000
G1 Z-2.00000000
G0 Z10.00000000
(-- end canned_drill --)
(MSG,Canned drill dwell without return-to-Z)
G0 X-1.00000000 Y39.00000000
(-- canned_drill R-plane=5 dwelling=0.5 return-to-old-Z=false --)
G61
G0 X10.00000000 Y40.00000000
G0 Z5.00000000
G1 Z-1.00000000
G4 P0.50000000
G0 Z5.00000000
G0 X15.00000000
G1 Z-1.00000000
G4 P0.50000000
G0 Z5.00000000
G0 X20.00000000
G1 Z-1.00000000
G4 P0.50000000
G0 Z5.00000000
G0 X25.00000000
G1 Z-1.00000000
G4 P0.50000000
G0 Z5.00000000
G0 Y45.00000000
G1 Z-1.00000000
G4 P0.50000000
G0 Z5.00000000
G0 X20.00000000
G1 Z-1.00000000
G4 P0.50000000
G0 Z5.00000000
G0 X15.00000000
G1 Z-2.00000000
G4 P0.50000000
G0 Z5.00000000
G0 X10.00000000
G1 Z-2.00000000
G4 P0.50000000
G0 Z5.00000000
(-- end canned_drill --)
(MSG,Canned drill dwell with return-to-Z)
G0 X-1.00000000 Y49.00000000 Z10.00000000
(-- canned_drill R-plane=5 dwelling=0.5 return-to-old-Z=true --)
G61
G0 X10.00000000 Y50.00000000
G0 Z5.00000000
G1 Z-1.00000000
G4 P0.50000000
G0 Z10.00000000
G0 X15.00000000
G0 Z5.00000000
G1 Z-1.00000000
G4 P0.50000000
G0 Z10.00000000
G0 X20.00000000
G0 Z5.00000000
G1 Z-1.00000000
G4 P0.50000000
G0 Z10.00000000
G0 X25.00000000
G0 Z5.00000000
G1 Z-1.00000000
G4 P0.50000000
G0 Z10.00000000
G0 Y55.00000000
G0 Z5.00000000
G1 Z-1.00000000
G4 P0.50000000
G0 Z10.00000000
G0 X20.00000000
G0 Z5.00000000
G1 Z-1.00000000
G4 P0.50000000
G0 Z10.00000000
G0 X15.00000000
G0 Z5.00000000
G1 Z-2.00000000
G4 P0.50000000
G0 Z10.00000000
G0 X10.00000000
G0 Z5.00000000
G1 Z-2.00000000
G4 P0.50000000
G0 Z10.00000000
(-- end canned_drill --)
G0 X0.00000000 Y0.00000000
G0 Z0.00000000
(-- epilogue begin --)
M30 (-- epilogue end --)
`
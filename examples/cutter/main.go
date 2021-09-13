// cutter generates a design based on the example here:
// https://gitlab.com/gcmc/gcmc/blob/master/example/cutter.gcmc
//
// Usage:
//   go run examples/cutter/main.go > cutter.gcode
package main

import (
	"fmt"

	. "github.com/gmlewis/go-gcode/gcode"
)

const (
	Feedrate  = 250
	TouchZ    = -1.0
	SafeZ     = 5.0
	PredrillZ = -4.0
	TW        = 3.0
	TW2       = TW / 2.0
	Dwell     = -1 // >= 0 means dwell enable
)

var (
	Origin = XYZ(0, 0, 5)
)

func main() {
	g := gcmc()
	fmt.Printf("%v", g)
}

func gcmc() *GCode {
	g := New()
	g.Prologue = true
	g.Epilogue = true

	g.Comment("-- feed        ", Feedrate, " --")
	g.Comment("-- origin      ", Origin, " --")
	g.Comment("-- touchz      ", TouchZ, " --")
	g.Comment("-- safez       ", SafeZ, " --")
	g.Comment("-- predrillz   ", PredrillZ, " --")
	g.Comment("-- toolwidth   ", TW, " --")
	g.Comment("-- toolwidth/2 ", TW2, " --")

	holes := []Tuple{
		XY(MilToMM(275), MilToMM(155)),  // Mount
		XY(MilToMM(275), MilToMM(555)),  // InA
		XY(MilToMM(275), MilToMM(1485)), // InB
		XY(MilToMM(275), MilToMM(1785)), // Mount
	}

	return g
}

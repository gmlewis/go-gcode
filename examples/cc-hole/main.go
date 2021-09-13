// cc-hole generates a design based on the example here:
// https://gitlab.com/gcmc/gcmc/blob/master/example/cc_hole.gcmc
//
// Usage:
//   go run examples/cc-hole/main.go > cc-hole.gcode
package main

import (
	"fmt"

	. "github.com/gmlewis/go-gcode/gcode"
	"github.com/gmlewis/go-gcode/utils"
)

func main() {
	g := gcmc()
	fmt.Printf("%v", g)
}

func gcmc() *GCode {
	g := New()
	g.Prologue = true
	g.Epilogue = true

	g.Feedrate(600)

	home := XYZ(0, 0, 10)
	g.GotoXYZ(home)
	utils.CCHole(g, XY(10.0, 25.0), 4.0, 3.0, 1.1, -1.0)
	utils.CCHole(g, XY(25.0, 25.0), 7.0, 3.0, 0.6, -1.0)
	utils.CCHole(g, XY(55.0, 25.0), 17.0, 3.0, 3.0, -1.0)
	utils.CCHole(g, XY(100.0, 25.0), 26.0, 3.0, 5.0, -1.0)

	g.GotoXYZ(home)

	return g
}

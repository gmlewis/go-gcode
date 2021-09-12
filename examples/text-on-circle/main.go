// text-on-circle generates text and is based on the example here:
// https://gitlab.com/gcmc/gcmc/blob/master/example/text-on-circle.gcmc
//
// Usage:
//   go run examples/text-on-circle/main.go > text-on-circle.gcode
package main

import (
	"fmt"

	_ "github.com/gmlewis/go-fonts/fonts/allura_regular"
	_ "github.com/gmlewis/go-fonts/fonts/amerikasans"
	_ "github.com/gmlewis/go-fonts/fonts/freesans"
	_ "github.com/gmlewis/go-fonts/fonts/freeserif"
	_ "github.com/gmlewis/go-fonts/fonts/freeserifbold"
	_ "github.com/gmlewis/go-fonts/fonts/freeserifbolditalic"
	_ "github.com/gmlewis/go-fonts/fonts/freeserifitalic"
	_ "github.com/gmlewis/go-fonts/fonts/grandhotel_regular"
	. "github.com/gmlewis/go-gcode/gcode"
	"github.com/gmlewis/go-gcode/utils"
)

const (
	sf = 10
)

func main() {
	g := gcmc()
	fmt.Printf("%v", g)
}

func gcmc() *GCode {
	g := New()
	g.Prologue = true
	g.Epilogue = true

	g.Feedrate(400)

	vl := utils.Typeset("Text ", "freesans")
	vl = append(vl, LastXY(vl).Offset(utils.Typeset("and ", "amerikasans")...)...)
	vl = append(vl, LastXY(vl).Offset(utils.Typeset("fonts ", "allura_regular")...)...)
	vl = append(vl, LastXY(vl).Offset(utils.Typeset("going ", "grandhotel_regular")...)...)
	vl = append(vl, LastXY(vl).Offset(utils.Typeset("Round ", "freeserif")...)...)
	vl = append(vl, LastXY(vl).Offset(utils.Typeset("and ", "freeserifbold")...)...)
	vl = append(vl, LastXY(vl).Offset(utils.Typeset("Round ", "freeserifitalic")...)...)
	vl = append(vl, LastXY(vl).Offset(utils.Typeset("<dizzy>. ", "freeserifbolditalic")...)...)

	utils.Engrave(g, vl, 1, 0)

	return g
}

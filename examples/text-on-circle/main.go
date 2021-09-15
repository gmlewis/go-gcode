// text-on-circle generates text and is based on the example here:
// https://gitlab.com/gcmc/gcmc/blob/master/example/text-on-circle.gcmc
//
// Usage:
//   go run examples/text-on-circle/main.go > text-on-circle.gcode
package main

import (
	"fmt"
	"math"

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

func main() {
	g := gcmc()
	fmt.Printf("%v", g)
}

func gcmc() *GCode {
	g := New()
	g.Prologue = true
	g.Epilogue = true

	g.Feedrate(400)

	genText(g, 20)

	genInsetText(g, 16)

	return g
}

func genText(g *GCode, sf float64) {
	vl := utils.Typeset("Text ", "freesans")
	vl = append(vl, LastXY(vl).Offset(utils.Typeset("and ", "amerikasans")...)...)
	vl = append(vl, LastXY(vl).Offset(utils.Typeset("fonts ", "allura_regular")...)...)
	vl = append(vl, LastXY(vl).Offset(utils.Typeset("going ", "grandhotel_regular")...)...)
	vl = append(vl, LastXY(vl).Offset(utils.Typeset("Round ", "freeserif")...)...)
	vl = append(vl, LastXY(vl).Offset(utils.Typeset("and ", "freeserifbold")...)...)
	vl = append(vl, LastXY(vl).Offset(utils.Typeset("Round ", "freeserifitalic")...)...)
	vl = append(vl, LastXY(vl).Offset(utils.Typeset("<dizzy>. ", "freeserifbolditalic")...)...)

	vl = Scaling(sf/(2*math.Pi), sf, 1).Transform(vl...)

	// The end-point's X is the actual size of the string, which is the circumference.
	circ := LastXY(vl).X()

	var rotvl []Tuple
	for _, v := range vl {
		angle := 0.5*math.Pi - 2*math.Pi*v.X()/circ
		r := v.Y() + circ
		x := r * math.Cos(angle)
		y := r * math.Sin(angle)
		rotvl = append(rotvl, XYZ(x, y, v.Z()))
	}

	utils.Engrave(g, rotvl, 1, 0)
}

func genInsetText(g *GCode, sf float64) {
	vl := utils.Typeset("Blá Blå Blà Blæ Blø", "freesans")

	vl = Scaling(sf/math.Pi, sf, 1).Transform(vl...)

	circ := LastXY(vl).X()

	var rotvl []Tuple
	for _, v := range vl {
		angle := math.Pi - math.Pi*v.X()/circ
		r := v.Y() + circ
		x := r * math.Cos(angle)
		y := r * math.Sin(angle)
		rotvl = append(rotvl, XYZ(x, y, v.Z()))
	}

	// Engrave on a different engraving plane.
	utils.Engrave(g, rotvl, 5, 1)
}

// text-fonts generates a design based on the example here:
// https://gitlab.com/gcmc/gcmc/blob/master/example/text-fonts.gcmc
//
// Usage:
//   go run examples/text-fonts/main.go > text-fonts.gcode
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
	text = "Hello World! All your Glyphs are belong to us."
	sf   = 2.5
)

func main() {
	g := gcmc()
	fmt.Printf("%v", g)
}

func gcmc() *GCode {
	g := New(UseGeneric)

	g.Feedrate(400)

	vl := utils.Typeset(text, "freesans")
	vl = Scaling(sf, sf, 0).Translate(0, 18*sf, 0).Transform(vl...)

	utils.Engrave(g, vl, 1, 0)

	f := func(x, y float64) M4 {
		return Scaling(sf, sf, 0).Translate(x, y, 0)
	}
	fs := func(x, y float64) M4 {
		return Scaling(sf, sf, 0).Shear(0.15, 0, 0, 0, 0, 0).Translate(x, y, 0)
	}
	t := func(xfm M4, fontName string) {
		utils.Engrave(g, xfm.Transform(utils.Typeset(text, fontName)...), 1, 0)
	}

	t(f(0, 16*sf), "amerikasans")
	t(fs(0, 14*sf), "allura_regular")
	t(fs(0, 12*sf), "grandhotel_regular")
	t(f(0, 10*sf), "freeserif")
	t(f(0, 8*sf), "freeserifbold")
	t(f(0, 6*sf), "freeserifitalic")
	t(f(0, 4*sf), "freeserifbolditalic")
	t(f(0, 2*sf), "allura_regular")
	t(f(0, 0), "grandhotel_regular")

	return g
}

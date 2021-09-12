package utils

import (
	"log"

	"github.com/gmlewis/go-fonts/fonts"
	. "github.com/gmlewis/go-gcode/gcode"
)

// Typeset creates a vector list from the given string
// using the provided go-font.
func Typeset(msg, fontName string) []Tuple {
	render, err := fonts.Text(0, 0, 1, 1, msg, fontName, nil)
	if err != nil {
		log.Fatal(err)
	}

	var offX, offY float64
	if len(render.Info) > 0 {
		gi := render.Info[0]
		offX, offY = gi.X, gi.Y
	}

	var vs []Tuple
	for _, poly := range render.Polygons {
		pts := poly.Pts
		if len(pts) == 0 {
			continue
		}

		vs = append(vs, Point(pts[0][0]-offX, pts[0][1]-offY, 1))
		for _, pt := range pts {
			vs = append(vs, Point(pt[0]-offX, pt[1]-offY, 0))
		}
		lastPt := pts[len(pts)-1]
		vs = append(vs, Point(lastPt[0]-offX, lastPt[1]-offY, 1))
	}

	if len(render.Info) > 0 {
		gi := render.Info[len(render.Info)-1]
		vs = append(vs, Point(gi.X+gi.Width-offX, gi.Y-offY, 1))
	}

	return vs
}

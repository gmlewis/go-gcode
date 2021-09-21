// edm-pro-box generates a design based on the example here:
// https://gitlab.com/gcmc/gcmc/blob/master/example/edm-pro-box.gcmc
//
// Usage:
//   go run examples/edm-pro-box/main.go > edm-pro-box.gcode
package main

import (
	"fmt"
	"math"

	. "github.com/gmlewis/go-gcode/gcode"
)

func main() {
	g := gcmc()
	fmt.Printf("%v", g)
}

func gcmc() *GCode {
	g := New(UseGeneric)

	g.Feedrate(500)
	bottomPlate(g, XY(0, 0))
	sidePlate(g, XY(113, 150), false)
	sidePlate(g, XY(164, 0), true)
	topPlate(g, XY(187, 0))
	frontPlate(g, XY(113, 162))
	backPlate(g, XY(0, 159))

	return g
}

var (
	wall  = 3.0              // Plexiglas thickness (mm)
	caseW = 104.0 + 2.0*wall // Case size  (mm)
	caseL = 150.0 + 2.0*wall
	caseH = 20.0                   // Front height  (mm)
	angle = 10.0 * math.Pi / 180.0 // Top angled by this

	pcbholeW   = MilToMM(3545) // Mounting hole spacing on PCB
	pcbholeL   = MilToMM(5515)
	pcbholeRad = 1.55         // PCB mounting hole radius (mm)
	pcbswRad   = 4.75         // Switch hole radius (mm)
	pcbswOfsW  = MilToMM(985) // Offset for switch hole from top mounting holes
	pcbswOfsL  = MilToMM(515)

	screwholeRad = 1.55       // Screw mount through the plate (mm)
	screwRetract = wall / 2.0 // Hole's center from the edge
	footholeRad  = 2.125      // Rubber foot mounting hole (mm)
	footRetract  = 20.0       // Rubber foot movement from the edges (mm)

	// Special path to cut a nut-spacing and screw mount into the plate (mm)
	screwCut = []Tuple{
		XY(-1.50, 0.0),
		XY(-1.50, 4.0),
		XY(-2.75, 4.0),
		XY(-2.75, 6.5),
		XY(-1.50, 6.5),
		XY(-1.50, 8.5),
		XY(1.50, 8.5),
		XY(1.50, 6.5),
		XY(2.75, 6.5),
		XY(2.75, 4.0),
		XY(1.50, 4.0),
		XY(1.50, 0.0),
	}
)

func hole(g *GCode, center Tuple, radius float64) {
	g.GotoXY(XY(center.X()-radius, center.Y()))
	g.CircleCWRel(X(radius), &TurnsOption{Turns: 0})
}

func trace(g *GCode, vl ...Tuple) {
	if len(vl) == 0 {
		return
	}
	g.GotoXY(vl[0])
	g.MoveXY(vl...)
	g.MoveXY(vl[0])
}

func bottomPlate(g *GCode, offset Tuple) {
	bp := []Tuple{
		XY(0, 0),
		XY(caseW, 0),
		XY(caseW, caseL),
		XY(0, caseL),
	}
	bpS := []Tuple{
		XY(screwRetract, caseL/2.0),
		XY(caseW/2.0, screwRetract),
		XY(caseW-screwRetract, caseL/2.0),
		XY(caseW/2.0, caseL-screwRetract),
	}
	corners := []Tuple{XY(1, 1), XY(-1, 1), XY(-1, -1), XY(1, -1)}

	for i := 0; i < len(bp); i++ {
		hole(g, bpS[i].Add(offset), screwholeRad)
		hole(g, bp[i].Add(corners[i].MultScalar(footRetract).Add(offset)), footholeRad)
	}

	trace(g, offset.Offset(bp...)...)
}

func frontPlate(g *GCode, offset Tuple) {
	tp := append([]Tuple{XY(0, 0), XY(caseW/2.0+screwCut[0][0], 0)},
		XY(caseW/2.0, 0).Offset(screwCut...)...)
	tp = append(tp, []Tuple{
		XY(caseW/2.0, 0).Add(LastXY(screwCut)),
		XY(caseW, 0),
		XY(caseW, caseH-wall*math.Sin(angle)),
		XY(0, caseH-wall*math.Sin(angle)),
	}...)

	hole(g, XY(caseW-screwRetract, caseH/2.0).Add(offset), screwholeRad)
	hole(g, XY(screwRetract, caseH/2.0).Add(offset), screwholeRad)

	trace(g, offset.Offset(tp...)...)
}

// Hole for DC power plug
func powerHole(g *GCode, offset Tuple) {
	g.GotoXY(XY(-5.9, -2.48).Add(offset))
	g.MoveXY(XY(-5.9, 2.48).Add(offset))
	g.ArcCW(XY(5.9, 2.48).Add(offset), 6.4, nil)
	g.MoveXY(XY(5.9, -2.48).Add(offset))
	g.ArcCW(XY(-5.9, -2.48).Add(offset), 6.4, nil)
}

func backPlate(g *GCode, offset Tuple) {
	h := caseH + (caseL-2*wall)*math.Tan(angle)
	tp := append([]Tuple{XY(0, 0), XY(caseW/2.0+screwCut[0][0], screwCut[0][1])},
		XY(caseW/2.0, 0).Offset(screwCut...)...)
	tp = append(tp, []Tuple{
		XY(caseW, 0),
		XY(caseW, h),
		XY(0, h),
	}...)

	hole(g, XY(caseW-screwRetract, h/2.0).Add(offset), screwholeRad)
	powerHole(g, XY(0.75*caseW, 0.4*h).Add(offset))
	hole(g, XY(screwRetract, h/2.0).Add(offset), screwholeRad)

	trace(g, offset.Offset(tp...)...)
}

func sidePlatePath() []Tuple {
	l := caseL - 2.0*wall
	h := caseH + l*math.Tan(angle)
	tl := l / math.Cos(angle)
	tv := XY(-math.Cos(angle), -math.Sin(angle))
	shr := RotationZ(0.5 * math.Pi).Transform(screwCut...)
	shl := RotationZ(-0.5 * math.Pi).Transform(screwCut...)
	sht := RotationZ(math.Pi + angle).Transform(screwCut...)
	result := append([]Tuple{
		XY(0, 0),
		XY(l/2.0, 0).Add(screwCut[0]),
	}, XY(l/2.0, 0).Offset(screwCut...)...)
	result = append(result, []Tuple{
		XY(l, 0),
		XY(l, h/2).Add(shr[0]),
	}...)
	result = append(result, XY(l, h/2).Offset(shr...)...)
	result = append(result, []Tuple{
		XY(l, h),
		XY(l, h).Add(tv.MultScalar(tl * 0.25)).Add(sht[0]),
	}...)
	result = append(result, XY(l, h).Add(tv.MultScalar(tl*0.25)).Offset(sht...)...)
	result = append(result, tv.MultScalar(tl*0.75).Add(XY(l, h)).Add(sht[0]))
	result = append(result, XY(l, h).Add(tv.MultScalar(tl*0.75)).Offset(sht...)...)
	result = append(result, []Tuple{
		XY(0, caseH),
		XY(0, caseH/2).Add(shl[0]),
	}...)
	result = append(result, XY(0, caseH/2).Offset(shl...)...)
	return result
}

// Left/right side plates
func sidePlate(g *GCode, offset Tuple, m bool) {
	sp := sidePlatePath()
	if m {
		sp = RotationZ(0.5*math.Pi-angle).Scale(-1, 1, 0).Transform(sp...)
	} else {
		sp = RotationZ(-0.5 * math.Pi).Transform(sp...)
	}

	trace(g, offset.Offset(sp...)...)
}

func topPlate(g *GCode, offset Tuple) {
	l := caseL / math.Cos(angle)
	l2 := (caseL - 2*wall) / math.Cos(angle)

	tp := []Tuple{
		XY(0, 0),
		XY(caseW, 0),
		XY(caseW, l),
		XY(0, l),
	}

	pbw := (caseW - pcbholeW) / 2.0
	pbl := (caseL - pcbholeL) / 2.0
	hole(g, XY(pbw, pbl).Add(offset), pcbholeRad)
	hole(g, XY(caseW-pbw, pbl).Add(offset), pcbholeRad)
	hole(g, XY(caseW-screwRetract, l2*0.25+wall).Add(offset), screwholeRad)
	hole(g, XY(caseW-screwRetract, l2*0.75+wall).Add(offset), screwholeRad)
	hole(g, XY(caseW-pbw, caseL-pbl).Add(offset), pcbholeRad)
	hole(g, XY(pbw, caseL-pbl).Add(offset), pcbholeRad)
	hole(g, XY(screwRetract, l2*0.75+wall).Add(offset), screwholeRad)
	hole(g, XY(screwRetract, l2*0.25+wall).Add(offset), screwholeRad)
	hole(g, XY(pbw+pcbswOfsW, pbl+pcbswOfsL).Add(offset), pcbswRad)
	hole(g, XY(caseW-pbw-pcbswOfsW, pbl+pcbswOfsL).Add(offset), pcbswRad)

	trace(g, offset.Offset(tp...)...)
}

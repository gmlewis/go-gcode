// cutter generates a design based on the example here:
// https://gitlab.com/gcmc/gcmc/blob/master/example/cutter.gcmc
//
// Usage:
//   go run examples/cutter/main.go > cutter.gcode
package main

import (
	"fmt"

	. "github.com/gmlewis/go-gcode/gcode"
	"github.com/gmlewis/go-gcode/utils"
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

	padPath = []Tuple{
		XY(TW2, MilToMM(-75)+TW2), XY(TW2, MilToMM(75)-TW2),
		XY(MilToMM(550)-TW2, MilToMM(75)-TW2), XY(MilToMM(550)-TW2, MilToMM(-75)+TW2),
	}
)

func main() {
	g := gcmc()
	fmt.Printf("%v", g)
}

func gcmc() *GCode {
	g := New(UseGeneric)

	g.Comment("-- feed        ", Feedrate, " --")
	g.Comment("-- origin      ", Origin, " --")
	g.Comment("-- touchz      ", TouchZ, " --")
	g.Comment("-- safez       ", SafeZ, " --")
	g.Comment("-- predrillz   ", PredrillZ, " --")
	g.Comment("-- toolwidth   ", TW, " --")
	g.Comment("-- toolwidth/2 ", TW2, " --")

	// Left-side holes
	holes := []Tuple{
		XY(MilToMM(275), MilToMM(155)),  // Mount
		XY(MilToMM(275), MilToMM(555)),  // InA
		XY(MilToMM(275), MilToMM(1485)), // InB
		XY(MilToMM(275), MilToMM(1785)), // Mount
	}

	// Press holes
	for x := 0.0; x < 4.0; x += 1.0 {
		holes = append(holes,
			XY(MilToMM(895+x*2025), MilToMM(335)),
			XY(MilToMM(895+x*2025), MilToMM(1605)))
	}

	holes = append(holes,
		// Right-side holes
		XY(MilToMM(7575), MilToMM(1605)), // Mount
		XY(MilToMM(7575), MilToMM(1200)), // OutA
		XY(MilToMM(7575), MilToMM(970)),  // OutB
		XY(MilToMM(7575), MilToMM(740)),  // OutC
		XY(MilToMM(7575), MilToMM(335)),  // Mount
		// Alignment holes
		XY(MilToMM(6495), MilToMM(80)),
		XY(MilToMM(6495), MilToMM(1860)),
		XY(MilToMM(2365), MilToMM(1860)),
		XY(MilToMM(2365), MilToMM(80)),
	)

	vPath := []Tuple{
		XY(MilToMM(280), MilToMM(1080)), XY(MilToMM(1870), MilToMM(1255)), XY(MilToMM(1990), MilToMM(1375)), XY(MilToMM(2460)-TW2, MilToMM(1375)),
		XY(MilToMM(2460)-TW2, MilToMM(1185)), XY(MilToMM(1990), MilToMM(1185)), XY(MilToMM(1870), MilToMM(1155)), XY(MilToMM(470), MilToMM(980)),
		XY(MilToMM(470), MilToMM(960)), XY(MilToMM(1870), MilToMM(785)), XY(MilToMM(1990), MilToMM(755)), XY(MilToMM(2460)-TW2, MilToMM(755)),
		XY(MilToMM(2460)-TW2, MilToMM(565)), XY(MilToMM(1990), MilToMM(565)), XY(MilToMM(1870), MilToMM(685)), XY(MilToMM(280), MilToMM(860)),
	}

	vPathDrill := []Tuple{
		XY(MilToMM(2410), MilToMM(585)),
		XY(MilToMM(2410), MilToMM(735)),
		XY(MilToMM(2410), MilToMM(1205)),
		XY(MilToMM(2410), MilToMM(1355)),
	}

	g.Feedrate(Feedrate)
	g.GotoXYZ(Origin)
	g.GotoZ(Z(SafeZ))

	for _, hole := range holes {
		g.GotoXY(hole)
		g.MoveZ(Z(TouchZ))
		doDwell(g)
		g.MoveZ(Z(SafeZ))
	}

	putPad(g, XY(MilToMM(1560), MilToMM(80)))
	putPad(g, XY(MilToMM(1560), MilToMM(1860)))

	// Component area
	ePath := erode(vPath, TW*2.0/3.0)
	g.GotoXY(vPath[0])
	g.MoveZ(Z(-1))
	utils.TracePath(g, -1, Dwell, vPath...)
	g.MoveXY(ePath[0])
	utils.TracePath(g, -1, Dwell, ePath...)
	g.MoveXY(vPath[0])
	g.MoveZ(Z(-2))
	utils.TracePath(g, -2, Dwell, vPath...)
	g.MoveXY(ePath[0])
	utils.TracePath(g, -2, Dwell, ePath...)
	g.MoveZ(Z(SafeZ))

	for _, hole := range vPathDrill {
		g.GotoXY(hole)
		doDwell(g)
		g.MoveZ(Z(PredrillZ))
		doDwell(g)
		g.MoveZ(Z(SafeZ))
	}

	putPad(g, XY(MilToMM(6620), MilToMM(970)))

	g.MoveZ(Z(SafeZ))
	g.GotoXYZ(Origin)

	return g
}

func doDwell(g *GCode) {
	if Dwell >= 0.0 {
		g.Dwell(Dwell)
	}
}

func erode(srcPath []Tuple, width float64) []Tuple {
	result := make([]Tuple, 0, len(srcPath))
	n := len(srcPath)
	for i, pc := range srcPath {
		pp := srcPath[(i-1+n)%n]
		pn := srcPath[(i+1)%n]
		v1 := pp.Sub(pc).Normalize()
		v2 := pn.Sub(pc).Normalize()
		bisect := v1.Add(v2).Normalize().MultScalar(width)
		newPoint := bisect.Add(pc)
		if i > 0 {
			crossP := v1[0]*v2[1] - v1[1]*v2[0]
			if crossP < 0.0 {
				newPoint = pc.Sub(bisect)
			}
		}
		newPoint[3] = 1
		result = append(result, newPoint)
	}
	return result
}

func putPad(g *GCode, offset Tuple) {
	g.Comment("-- putpad at offset=", offset, " --")
	g.GotoXY(padPath[0].Add(offset))
	g.MoveZ(Z(-1))
	utils.TracePath(g, -1, Dwell, offset.Offset(padPath...)...)
	g.MoveZ(Z(-2))
	utils.TracePath(g, -2, Dwell, offset.Offset(padPath...)...)
	g.MoveZ(Z(-1))

	g.GotoXY(XY(MilToMM(50), 0).Add(offset))
	doDwell(g)
	g.MoveZ(Z(PredrillZ))
	doDwell(g)
	g.MoveZ(Z(SafeZ))
}

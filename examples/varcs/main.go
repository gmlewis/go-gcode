// varcs generates vectorized arcs and is based on the example here:
// https://gitlab.com/gcmc/gcmc/blob/master/example/varcs.gcmc
//
// Usage:
//   go run examples/varcs/main.go > varcs.gcode
package main

import (
	"fmt"
	"math"

	. "github.com/gmlewis/go-gcode/gcode"
	"github.com/gmlewis/go-gcode/utils"
)

const (
	safeHeight   = 1
	cuttingDepth = -1
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

	simpleArcs(g)

	spiral(g)

	wavyArcs(g)

	simpleEllipse(g)

	angledEllipse(g)

	return g
}

func simpleArcs(g *GCode) {
	radius := 15.0
	flushItAt(g, XY(-20, 100), utils.VArcCW(XY(10, 15), radius, nil))
	flushItAt(g, XY(-20, 60), utils.VArcCW(XY(10, 15), -radius, nil))
	flushItAt(g, XY(20, 100), utils.VArcCCW(XY(10, 15), radius, nil))
	flushItAt(g, XY(20, 60), utils.VArcCCW(XY(10, 15), -radius, nil))
}

func flushItAt(g *GCode, v Tuple, vl []Tuple) {
	oldPos := g.Position()
	g.GotoXY(v)
	g.MoveXY(v.Offset(vl...)...)
	g.GotoXYZ(oldPos)
}

func spiral(g *GCode) {
	radius := 20.0
	center := XY(60, 100)
	zMove := -20.0
	turns := 5
	cp := center.Normalize().MultScalar(radius)
	cp[2] = zMove

	spiral := utils.VCircleCW(cp, &utils.VOptions{Turns: turns, MaxL: Float(1), MaxA: Float(3)})
	spiral = center.Sub(XY(cp.X(), cp.Y())).Offset(spiral...)

	g.GotoXY(spiral[len(spiral)-1])
	g.GotoZ(Z(0))
	g.MoveXYZ(spiral...)
	g.GotoZ(Z(0))
	g.GotoXY(XY(0, 0))
}

func wavyArcs(g *GCode) {
	waveLength := 10.0
	rr := 4.0
	up := utils.VArcCW(X(0.5*waveLength), rr, nil)
	dn := utils.VArcCCW(X(0.5*waveLength), rr, nil)

	period := append(up, X(0.5*waveLength).Offset(dn...)...)
	wave := append([]Tuple{}, period...)
	for i := 0; i < 4; i++ {
		wave = append(wave, X(float64(i)*waveLength).Offset(period...)...)
	}

	zWave := Reverse(RotationX(90 * math.Pi / 180).Transform(wave...))

	for i := 0; i < 12; i++ {
		g.MoveXYZ(wave...)
		g.MoveXYZ(zWave...)
		g.GotoXY(XY(0, 0))
		wave = RotationZ(30 * math.Pi / 180).Transform(wave...)
		zWave = RotationZ(30 * math.Pi / 180).Transform(zWave...)
	}

	g.GotoZ(Z(0))
}

func simpleEllipse(g *GCode) {
	majorRadius := 25.0
	minorRadius := 15.0
	ell := utils.VCircleCW(XY(majorRadius, 0), nil)
	ell = Scaling(1, minorRadius/majorRadius, 0).Translate(60, 50, 0).Transform(ell...)

	g.GotoXY(ell[len(ell)-1])
	g.MoveXY(ell...)
	g.GotoXY(XY(0, 0))
	g.GotoX(Z(0))
}

func angledEllipse(g *GCode) {
	majorRadius := 25.0
	minorRadius := 15.0
	angle := 30.0 * math.Pi / 180.0

	// Center point angle must be transformed with the major/minor axes ratio so
	// that the entry/exit point of the ellipse is on the bounding box' left side:
	// - take a unit vector
	// - rotate to desired ellipse angle
	// - scale according to radii ratio
	// - find the resulting angle
	x, y := math.Cos(angle), math.Sin(angle)*minorRadius/majorRadius
	cpa := math.Atan2(y, x)
	center := XY(majorRadius*math.Cos(-cpa), majorRadius*math.Sin(-cpa))

	ell := utils.VCircleCW(center, nil)
	ell = Scaling(1, minorRadius/majorRadius, 0).RotateZ(angle).Translate(60, 0, 0).Transform(ell...)

	g.GotoXY(ell[len(ell)-1])
	g.MoveXY(ell...)
	g.GotoXY(XY(0, 0))
	g.GotoX(Z(0))
}

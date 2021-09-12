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
	fmt.Printf("%v\n", g)
}

func gcmc() *GCode {
	g := New()
	g.Prologue = true
	g.Epilogue = true

	g.Feedrate(600)

	simpleArcs(g)

	spiral(g)

	wavyArcs(g)

	return g
}

func simpleArcs(g *GCode) {
	radius := 15.0
	flushItAt(g, XY(-20, 100), utils.VArcCW(XY(10, 15), radius, nil))
	flushItAt(g, XY(-20, 60), utils.VArcCW(XY(10, 15), -radius, nil))
	flushItAt(g, XY(20, 100), utils.VArcCCW(XY(10, 15), radius, nil))
	flushItAt(g, XY(20, 60), utils.VArcCCW(XY(10, 15), -radius, nil))
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

func flushItAt(g *GCode, v Tuple, vl []Tuple) {
	oldPos := g.Position()
	g.GotoXY(v)
	g.MoveXY(v.Offset(vl...)...)
	g.GotoXYZ(oldPos)
}

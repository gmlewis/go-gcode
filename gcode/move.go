package gcode

import (
	"fmt"
	"math"
	"strings"
)

const (
	forceX = 1 << iota
	forceY
	forceZ

	forceXY  = forceX | forceY
	forceYZ  = forceY | forceZ
	forceXZ  = forceX | forceZ
	forceXYZ = forceX | forceY | forceZ
)

func (g *GCode) genChangedXYZ(opCode string, p Tuple, force int) string {
	pos := g.Position()
	var parts []string
	if (!g.hasMoved && (force&forceX) != 0) || math.Abs(p.X()-pos.X()) >= epsilon {
		parts = append(parts, fmt.Sprintf("X%.8f", p.X()))
	}
	if (!g.hasMoved && (force&forceY) != 0) || math.Abs(p.Y()-pos.Y()) >= epsilon {
		parts = append(parts, fmt.Sprintf("Y%.8f", p.Y()))
	}
	if (!g.hasMoved && (force&forceZ) != 0) || math.Abs(p.Z()-pos.Z()) >= epsilon {
		parts = append(parts, fmt.Sprintf("Z%.8f", p.Z()))
	}
	if len(parts) == 0 {
		return ""
	}
	s := fmt.Sprintf("%v %v", opCode, strings.Join(parts, " "))
	return s
}

// moveOrGo optimizes the movement to only include the
// axes that have changed since the last moveOrGo.
// As a special case, for the very first move/goto command,
// force the output of all the mentioned axes, even if 0.
func (g *GCode) moveOrGo(opCode string, p Tuple, force int) {
	s := g.genChangedXYZ(opCode, p, force)
	if s == "" {
		return
	}
	p[3] = 1
	g.steps = append(g.steps, &Step{s: s, pos: p})
	g.hasMoved = true
}

// GotoX performs one or more rapid move(s) on the X axis.
func (g *GCode) GotoX(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		newPos := XYZ(p.X(), pos.Y(), pos.Z())
		g.moveOrGo("G0", newPos, forceX)
	}
	return g
}

// GotoY performs one or more rapid move(s) on the Y axis.
func (g *GCode) GotoY(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		newPos := XYZ(pos.X(), p.Y(), pos.Z())
		g.moveOrGo("G0", newPos, forceY)
	}
	return g
}

// GotoZ performs one or more rapid move(s) on the Z axis.
func (g *GCode) GotoZ(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		newPos := XYZ(pos.X(), pos.Y(), p.Z())
		g.moveOrGo("G0", newPos, forceZ)
	}
	return g
}

// GotoXY performs one or more rapid move(s) on the XY axes.
func (g *GCode) GotoXY(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		newPos := XYZ(p.X(), p.Y(), pos.Z())
		g.moveOrGo("G0", newPos, forceXY)
	}
	return g
}

// GotoYZ performs one or more rapid move(s) on the YZ axes.
func (g *GCode) GotoYZ(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		newPos := XYZ(pos.X(), p.Y(), p.Z())
		g.moveOrGo("G0", newPos, forceYZ)
	}
	return g
}

// GotoXZ performs one or more rapid move(s) on the XZ axes.
func (g *GCode) GotoXZ(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		newPos := XYZ(p.X(), pos.Y(), p.Z())
		g.moveOrGo("G0", newPos, forceXZ)
	}
	return g
}

// GotoXYZ performs one or more rapid move(s) on the XYZ axes.
func (g *GCode) GotoXYZ(ps ...Tuple) *GCode {
	for _, p := range ps {
		g.moveOrGo("G0", p, forceXYZ)
	}
	return g
}

// GotoXYZWithF performs one or more move(s) on the XYZ axes using the provided feed-rate.
func (g *GCode) GotoXYZWithF(feedrate float64, ps ...Tuple) *GCode {
	for i, p := range ps {
		g.moveOrGo("G0", p, forceXYZ)
		if i == 0 {
			lastStep := len(g.steps) - 1
			g.steps[lastStep].s += fmt.Sprintf(" F%v", feedrate)
		}
	}
	return g
}

// MoveX performs one or more move(s) on the X axis at the current feed-rate.
func (g *GCode) MoveX(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		newPos := XYZ(p.X(), pos.Y(), pos.Z())
		g.moveOrGo("G1", newPos, forceX)
	}
	return g
}

// MoveY performs one or more move(s) on the Y axis at the current feed-rate.
func (g *GCode) MoveY(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		newPos := XYZ(pos.X(), p.Y(), pos.Z())
		g.moveOrGo("G1", newPos, forceY)
	}
	return g
}

// MoveZ performs one or more move(s) on the Z axis at the current feed-rate.
func (g *GCode) MoveZ(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		newPos := XYZ(pos.X(), pos.Y(), p.Z())
		g.moveOrGo("G1", newPos, forceZ)
	}
	return g
}

// MoveZWithF performs one or more move(s) on the Z axis using the provided feed-rate.
func (g *GCode) MoveZWithF(feedrate float64, ps ...Tuple) *GCode {
	pos := g.Position()
	for i, p := range ps {
		newPos := XYZ(pos.X(), pos.Y(), p.Z())
		g.moveOrGo("G1", newPos, forceZ)
		if i == 0 {
			lastStep := len(g.steps) - 1
			g.steps[lastStep].s += fmt.Sprintf(" F%v", feedrate)
		}
	}
	return g
}

// MoveXY performs one or more move(s) on the XY axes at the current feed-rate.
func (g *GCode) MoveXY(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		newPos := XYZ(p.X(), p.Y(), pos.Z())
		g.moveOrGo("G1", newPos, forceXY)
	}
	return g
}

// MoveYZ performs one or more move(s) on the YZ axes at the current feed-rate.
func (g *GCode) MoveYZ(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		newPos := XYZ(pos.X(), p.Y(), p.Z())
		g.moveOrGo("G1", newPos, forceYZ)
	}
	return g
}

// MoveXZ performs one or more move(s) on the XZ axes at the current feed-rate.
func (g *GCode) MoveXZ(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		newPos := XYZ(p.X(), pos.Y(), p.Z())
		g.moveOrGo("G1", newPos, forceXZ)
	}
	return g
}

// MoveXYZ performs one or more move(s) on the XYZ axes at the current feed-rate.
func (g *GCode) MoveXYZ(ps ...Tuple) *GCode {
	for _, p := range ps {
		g.moveOrGo("G1", p, forceXYZ)
	}
	return g
}

// MoveXYZRel performs one or more move(s) on the XYZ axes at the current feed-rate
// using relative offsets.
func (g *GCode) MoveXYZRel(ps ...Tuple) *GCode {
	for _, p := range ps {
		pos := g.Position()
		g.moveOrGo("G1", pos.Add(p), forceXYZ)
	}
	return g
}

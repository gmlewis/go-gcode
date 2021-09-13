package gcode

import (
	"fmt"
	"math"
	"strings"
)

// moveOrGo optimizes the movement to only include the
// axes that have changed since the last moveOrGo.
func (g *GCode) moveOrGo(opCode string, p Tuple) {
	pos := g.Position()
	var parts []string
	if math.Abs(p.X()-pos.X()) >= epsilon {
		parts = append(parts, fmt.Sprintf("X%.8f", p.X()))
	}
	if math.Abs(p.Y()-pos.Y()) >= epsilon {
		parts = append(parts, fmt.Sprintf("Y%.8f", p.Y()))
	}
	if math.Abs(p.Z()-pos.Z()) >= epsilon {
		parts = append(parts, fmt.Sprintf("Z%.8f", p.Z()))
	}
	if len(parts) == 0 {
		return
	}
	s := fmt.Sprintf("%v %v", opCode, strings.Join(parts, " "))
	g.steps = append(g.steps, &Step{s: s, pos: p})
}

// GotoX performs one or more rapid move(s) on the X axis.
func (g *GCode) GotoX(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		newPos := Point(p.X(), pos.Y(), pos.Z())
		g.moveOrGo("G0", newPos)
	}
	return g
}

// GotoY performs one or more rapid move(s) on the Y axis.
func (g *GCode) GotoY(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		newPos := Point(pos.X(), p.Y(), pos.Z())
		g.moveOrGo("G0", newPos)
	}
	return g
}

// GotoZ performs one or more rapid move(s) on the Z axis.
func (g *GCode) GotoZ(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		newPos := Point(pos.X(), pos.Y(), p.Z())
		g.moveOrGo("G0", newPos)
	}
	return g
}

// GotoXY performs one or more rapid move(s) on the XY axes.
func (g *GCode) GotoXY(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		newPos := Point(p.X(), p.Y(), pos.Z())
		g.moveOrGo("G0", newPos)
	}
	return g
}

// GotoYZ performs one or more rapid move(s) on the YZ axes.
func (g *GCode) GotoYZ(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		newPos := Point(pos.X(), p.Y(), p.Z())
		g.moveOrGo("G0", newPos)
	}
	return g
}

// GotoXZ performs one or more rapid move(s) on the XZ axes.
func (g *GCode) GotoXZ(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		newPos := Point(p.X(), pos.Y(), p.Z())
		g.moveOrGo("G0", newPos)
	}
	return g
}

// GotoXYZ performs one or more rapid move(s) on the XYZ axis.
func (g *GCode) GotoXYZ(ps ...Tuple) *GCode {
	for _, p := range ps {
		g.moveOrGo("G0", p)
	}
	return g
}

// MoveX performs one or more move(s) on the X axis at the current feed-rate.
func (g *GCode) MoveX(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		newPos := Point(p.X(), pos.Y(), pos.Z())
		g.moveOrGo("G1", newPos)
	}
	return g
}

// MoveY performs one or more move(s) on the Y axis at the current feed-rate.
func (g *GCode) MoveY(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		newPos := Point(pos.X(), p.Y(), pos.Z())
		g.moveOrGo("G1", newPos)
	}
	return g
}

// MoveZ performs one or more move(s) on the Z axis at the current feed-rate.
func (g *GCode) MoveZ(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		newPos := Point(pos.X(), pos.Y(), p.Z())
		g.moveOrGo("G1", newPos)
	}
	return g
}

// MoveXY performs one or more move(s) on the XY axes at the current feed-rate.
func (g *GCode) MoveXY(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		newPos := Point(p.X(), p.Y(), pos.Z())
		g.moveOrGo("G1", newPos)
	}
	return g
}

// MoveYZ performs one or more move(s) on the YZ axes at the current feed-rate.
func (g *GCode) MoveYZ(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		newPos := Point(pos.X(), p.Y(), p.Z())
		g.moveOrGo("G1", newPos)
	}
	return g
}

// MoveXZ performs one or more move(s) on the XZ axes at the current feed-rate.
func (g *GCode) MoveXZ(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		newPos := Point(p.X(), pos.Y(), p.Z())
		g.moveOrGo("G1", newPos)
	}
	return g
}

// MoveXYZ performs one or more move(s) on the XYZ axes at the current feed-rate.
func (g *GCode) MoveXYZ(ps ...Tuple) *GCode {
	for _, p := range ps {
		g.moveOrGo("G1", p)
	}
	return g
}

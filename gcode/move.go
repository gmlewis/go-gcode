package gcode

import "fmt"

// GotoX performs one or more rapid move(s) on the X axis.
func (g *GCode) GotoX(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		pos = Point(p.X(), pos.Y(), pos.Z())
		g.steps = append(g.steps, &Step{
			s:   fmt.Sprintf("G0 X%.6f", p.X()),
			pos: pos,
		})
	}
	return g
}

// GotoY performs one or more rapid move(s) on the Y axis.
func (g *GCode) GotoY(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		pos = Point(pos.X(), p.Y(), pos.Z())
		g.steps = append(g.steps, &Step{
			s:   fmt.Sprintf("G0 Y%.6f", p.Y()),
			pos: pos,
		})
	}
	return g
}

// GotoZ performs one or more rapid move(s) on the Z axis.
func (g *GCode) GotoZ(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		pos = Point(pos.X(), pos.Y(), p.Z())
		g.steps = append(g.steps, &Step{
			s:   fmt.Sprintf("G0 Z%.6f", p.Z()),
			pos: pos,
		})
	}
	return g
}

// GotoXY performs one or more rapid move(s) on the XY axes.
func (g *GCode) GotoXY(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		pos = Point(p.X(), p.Y(), pos.Z())
		g.steps = append(g.steps, &Step{
			s:   fmt.Sprintf("G0 X%.6f Y%.6f", p.X(), p.Y()),
			pos: pos,
		})
	}
	return g
}

// GotoYZ performs one or more rapid move(s) on the YZ axes.
func (g *GCode) GotoYZ(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		pos = Point(pos.X(), p.Y(), p.Z())
		g.steps = append(g.steps, &Step{
			s:   fmt.Sprintf("G0 Y%.6f Z%.6f", p.Y(), p.Z()),
			pos: pos,
		})
	}
	return g
}

// GotoXZ performs one or more rapid move(s) on the XZ axes.
func (g *GCode) GotoXZ(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		pos = Point(p.X(), pos.Y(), p.Z())
		g.steps = append(g.steps, &Step{
			s:   fmt.Sprintf("G0 X%.6f Z%.6f", p.X(), p.Z()),
			pos: pos,
		})
	}
	return g
}

// GotoXYZ performs one or more rapid move(s) on the XYZ axis.
func (g *GCode) GotoXYZ(ps ...Tuple) *GCode {
	for _, p := range ps {
		g.steps = append(g.steps, &Step{
			s:   fmt.Sprintf("G0 X%.6f Y%.6f Z%.6f", p.X(), p.Y(), p.Z()),
			pos: p,
		})
	}
	return g
}

// MoveX performs one or more move(s) on the X axis at the current feed-rate.
func (g *GCode) MoveX(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		pos = Point(p.X(), pos.Y(), pos.Z())
		g.steps = append(g.steps, &Step{
			s:   fmt.Sprintf("G1 X%.6f", p.X()),
			pos: pos,
		})
	}
	return g
}

// MoveY performs one or more move(s) on the Y axis at the current feed-rate.
func (g *GCode) MoveY(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		pos = Point(pos.X(), p.Y(), pos.Z())
		g.steps = append(g.steps, &Step{
			s:   fmt.Sprintf("G1 Y%.6f", p.Y()),
			pos: pos,
		})
	}
	return g
}

// MoveZ performs one or more move(s) on the Z axis at the current feed-rate.
func (g *GCode) MoveZ(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		pos = Point(pos.X(), pos.Y(), p.Z())
		g.steps = append(g.steps, &Step{
			s:   fmt.Sprintf("G1 Z%.6f", p.Z()),
			pos: pos,
		})
	}
	return g
}

// MoveXY performs one or more move(s) on the XY axes at the current feed-rate.
func (g *GCode) MoveXY(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		pos = Point(p.X(), p.Y(), pos.Z())
		g.steps = append(g.steps, &Step{
			s:   fmt.Sprintf("G1 X%.6f Y%.6f", p.X(), p.Y()),
			pos: pos,
		})
	}
	return g
}

// MoveYZ performs one or more move(s) on the YZ axes at the current feed-rate.
func (g *GCode) MoveYZ(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		pos = Point(pos.X(), p.Y(), p.Z())
		g.steps = append(g.steps, &Step{
			s:   fmt.Sprintf("G1 Y%.6f Z%.6f", p.Y(), p.Z()),
			pos: pos,
		})
	}
	return g
}

// MoveXZ performs one or more move(s) on the XZ axes at the current feed-rate.
func (g *GCode) MoveXZ(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		pos = Point(p.X(), pos.Y(), p.Z())
		g.steps = append(g.steps, &Step{
			s:   fmt.Sprintf("G1 X%.6f Z%.6f", p.X(), p.Z()),
			pos: pos,
		})
	}
	return g
}

// MoveXYZ performs one or more move(s) on the XYZ axes at the current feed-rate.
func (g *GCode) MoveXYZ(ps ...Tuple) *GCode {
	for _, p := range ps {
		g.steps = append(g.steps, &Step{
			s:   fmt.Sprintf("G1 X%.6f Y%.6f Z%.6f", p.X(), p.Y(), p.Z()),
			pos: p,
		})
	}
	return g
}

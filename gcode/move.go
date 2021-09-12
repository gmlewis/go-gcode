package gcode

import "fmt"

// GotoX performs one or more rapid move(s) on the X axis.
func (g *GCode) GotoX(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		newPos := Point(p.X(), pos.Y(), pos.Z())
		if newPos.Equal(pos) {
			continue
		}
		pos = newPos
		g.steps = append(g.steps, &Step{
			s:   fmt.Sprintf("G0 X%.8f", p.X()),
			pos: pos,
		})
	}
	return g
}

// GotoY performs one or more rapid move(s) on the Y axis.
func (g *GCode) GotoY(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		newPos := Point(pos.X(), p.Y(), pos.Z())
		if newPos.Equal(pos) {
			continue
		}
		pos = newPos
		g.steps = append(g.steps, &Step{
			s:   fmt.Sprintf("G0 Y%.8f", p.Y()),
			pos: pos,
		})
	}
	return g
}

// GotoZ performs one or more rapid move(s) on the Z axis.
func (g *GCode) GotoZ(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		newPos := Point(pos.X(), pos.Y(), p.Z())
		if newPos.Equal(pos) {
			continue
		}
		pos = newPos
		g.steps = append(g.steps, &Step{
			s:   fmt.Sprintf("G0 Z%.8f", p.Z()),
			pos: pos,
		})
	}
	return g
}

// GotoXY performs one or more rapid move(s) on the XY axes.
func (g *GCode) GotoXY(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		newPos := Point(p.X(), p.Y(), pos.Z())
		if newPos.Equal(pos) {
			continue
		}
		pos = newPos
		g.steps = append(g.steps, &Step{
			s:   fmt.Sprintf("G0 X%.8f Y%.8f", p.X(), p.Y()),
			pos: pos,
		})
	}
	return g
}

// GotoYZ performs one or more rapid move(s) on the YZ axes.
func (g *GCode) GotoYZ(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		newPos := Point(pos.X(), p.Y(), p.Z())
		if newPos.Equal(pos) {
			continue
		}
		pos = newPos
		g.steps = append(g.steps, &Step{
			s:   fmt.Sprintf("G0 Y%.8f Z%.8f", p.Y(), p.Z()),
			pos: pos,
		})
	}
	return g
}

// GotoXZ performs one or more rapid move(s) on the XZ axes.
func (g *GCode) GotoXZ(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		newPos := Point(p.X(), pos.Y(), p.Z())
		if newPos.Equal(pos) {
			continue
		}
		pos = newPos
		g.steps = append(g.steps, &Step{
			s:   fmt.Sprintf("G0 X%.8f Z%.8f", p.X(), p.Z()),
			pos: pos,
		})
	}
	return g
}

// GotoXYZ performs one or more rapid move(s) on the XYZ axis.
func (g *GCode) GotoXYZ(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		if p.Equal(pos) {
			continue
		}
		pos = p
		g.steps = append(g.steps, &Step{
			s:   fmt.Sprintf("G0 X%.8f Y%.8f Z%.8f", p.X(), p.Y(), p.Z()),
			pos: pos,
		})
	}
	return g
}

// MoveX performs one or more move(s) on the X axis at the current feed-rate.
func (g *GCode) MoveX(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		newPos := Point(p.X(), pos.Y(), pos.Z())
		if newPos.Equal(pos) {
			continue
		}
		pos = newPos
		g.steps = append(g.steps, &Step{
			s:   fmt.Sprintf("G1 X%.8f", p.X()),
			pos: pos,
		})
	}
	return g
}

// MoveY performs one or more move(s) on the Y axis at the current feed-rate.
func (g *GCode) MoveY(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		newPos := Point(pos.X(), p.Y(), pos.Z())
		if newPos.Equal(pos) {
			continue
		}
		pos = newPos
		g.steps = append(g.steps, &Step{
			s:   fmt.Sprintf("G1 Y%.8f", p.Y()),
			pos: pos,
		})
	}
	return g
}

// MoveZ performs one or more move(s) on the Z axis at the current feed-rate.
func (g *GCode) MoveZ(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		newPos := Point(pos.X(), pos.Y(), p.Z())
		if newPos.Equal(pos) {
			continue
		}
		pos = newPos
		g.steps = append(g.steps, &Step{
			s:   fmt.Sprintf("G1 Z%.8f", p.Z()),
			pos: pos,
		})
	}
	return g
}

// MoveXY performs one or more move(s) on the XY axes at the current feed-rate.
func (g *GCode) MoveXY(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		newPos := Point(p.X(), p.Y(), pos.Z())
		if newPos.Equal(pos) {
			continue
		}
		pos = newPos
		g.steps = append(g.steps, &Step{
			s:   fmt.Sprintf("G1 X%.8f Y%.8f", p.X(), p.Y()),
			pos: pos,
		})
	}
	return g
}

// MoveYZ performs one or more move(s) on the YZ axes at the current feed-rate.
func (g *GCode) MoveYZ(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		newPos := Point(pos.X(), p.Y(), p.Z())
		if newPos.Equal(pos) {
			continue
		}
		pos = newPos
		g.steps = append(g.steps, &Step{
			s:   fmt.Sprintf("G1 Y%.8f Z%.8f", p.Y(), p.Z()),
			pos: pos,
		})
	}
	return g
}

// MoveXZ performs one or more move(s) on the XZ axes at the current feed-rate.
func (g *GCode) MoveXZ(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		newPos := Point(p.X(), pos.Y(), p.Z())
		if newPos.Equal(pos) {
			continue
		}
		pos = newPos
		g.steps = append(g.steps, &Step{
			s:   fmt.Sprintf("G1 X%.8f Z%.8f", p.X(), p.Z()),
			pos: pos,
		})
	}
	return g
}

// MoveXYZ performs one or more move(s) on the XYZ axes at the current feed-rate.
func (g *GCode) MoveXYZ(ps ...Tuple) *GCode {
	pos := g.Position()
	for _, p := range ps {
		if p.Equal(pos) {
			continue
		}
		pos = p
		g.steps = append(g.steps, &Step{
			s:   fmt.Sprintf("G1 X%.8f Y%.8f Z%.8f", p.X(), p.Y(), p.Z()),
			pos: pos,
		})
	}
	return g
}

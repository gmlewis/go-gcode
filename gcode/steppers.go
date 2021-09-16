package gcode

func (g *GCode) sendOpCode(opCode string) *GCode {
	pos := g.Position()
	g.steps = append(g.steps, &Step{s: opCode, pos: pos})
	return g
}

// EnableX enables the X stepper motor.
func (g *GCode) EnableX() *GCode { return g.sendOpCode("M17 X") }

// EnableY enables the Y stepper motor.
func (g *GCode) EnableY() *GCode { return g.sendOpCode("M17 Y") }

// EnableZ enables the Z stepper motor.
func (g *GCode) EnableZ() *GCode { return g.sendOpCode("M17 Z") }

// EnableXY enables the X and Y stepper motors.
func (g *GCode) EnableXY() *GCode { return g.sendOpCode("M17 X Y") }

// EnableYZ enables the Y and Z stepper motors.
func (g *GCode) EnableYZ() *GCode { return g.sendOpCode("M17 Y Z") }

// EnableXZ enables the X and Z stepper motors.
func (g *GCode) EnableXZ() *GCode { return g.sendOpCode("M17 X Z") }

// EnableXYZ enables the X, Y, and Z stepper motors.
func (g *GCode) EnableXYZ() *GCode { return g.sendOpCode("M17") }

// DisableX disables the X stepper motor.
func (g *GCode) DisableX() *GCode { return g.sendOpCode("M18 X") }

// DisableY disables the Y stepper motor.
func (g *GCode) DisableY() *GCode { return g.sendOpCode("M18 Y") }

// DisableZ disables the Z stepper motor.
func (g *GCode) DisableZ() *GCode { return g.sendOpCode("M18 Z") }

// DisableXY disables the X and Y stepper motors.
func (g *GCode) DisableXY() *GCode { return g.sendOpCode("M18 X Y") }

// DisableYZ disables the Y and Z stepper motors.
func (g *GCode) DisableYZ() *GCode { return g.sendOpCode("M18 Y Z") }

// DisableXZ disables the X and Z stepper motors.
func (g *GCode) DisableXZ() *GCode { return g.sendOpCode("M18 X Z") }

// DisableXYZ disables the X, Y, and Z stepper motors.
func (g *GCode) DisableXYZ() *GCode { return g.sendOpCode("M18") }

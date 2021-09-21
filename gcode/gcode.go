// Package gcode provides methods used to generate G-Code.
package gcode

import (
	"fmt"
	"strings"
	"time"
)

// GCode represents a G-Code design.
type GCode struct {
	noHeader   bool
	prologue   string
	epilogue   string
	commentFmt string

	activePlane PlaneT
	hasMoved    bool
	steps       []*Step
}

// Option represents various options for generating GCode.
type Option string

const (
	CommentsUseSemicolons Option = "CommentsUseSemicolons"
	NoHeader              Option = "NoHeader"
	UseIVI                Option = "UseIVI"
	UseGeneric            Option = "UseGeneric"
)

// New returns a new gcode design.
func New(opts ...Option) *GCode {
	g := &GCode{activePlane: PlaneXY, commentFmt: "(%v)"}

	for _, opt := range opts {
		switch opt {
		case CommentsUseSemicolons:
			g.commentFmt = ";%v"
		case NoHeader:
			g.noHeader = true
		case UseIVI:
			g.prologue = iviPrologue
			g.epilogue = iviEpilogue
			g.commentFmt = ";%v"
		case UseGeneric:
			g.prologue = genericPrologue
			g.epilogue = genericEpilogue
		}
	}

	return g
}

// String converts the design to a string.
func (g *GCode) String() string {
	var lines []string
	if !g.noHeader {
		const timeFmt = "2006-01-02 15:04:05"
		now := time.Now().Local()
		lines = append(lines, fmt.Sprintf(g.commentFmt+"\n"+g.commentFmt, identifier, now.Format(timeFmt)))
	}
	if g.prologue != "" {
		lines = append(lines, g.prologue)
	}
	for _, step := range g.steps {
		lines = append(lines, step.s)
	}
	if g.epilogue != "" {
		lines = append(lines, g.epilogue)
	}
	return strings.Join(lines, "\n") + "\n"
}

// Step represents a step in the GCode.
type Step struct {
	s   string
	pos Tuple // position after performing the step.
}

// Position returns the current tool position (defaulting to home 0,0,0).
func (g *GCode) Position() Tuple {
	if g == nil || len(g.steps) == 0 {
		return XYZ(0, 0, 0)
	}
	return g.steps[len(g.steps)-1].pos
}

const identifier = `go-gcode compiled code, do not change`

var genericPrologue = `(-- prologue begin --)
G17 ( Use XY plane )
G21 ( Use mm )
G40 ( Cancel cutter radius compensation )
G49 ( Cancel tool length compensation )
G54 ( Default coordinate system )
G80 ( Cancel canned cycle )
G90 ( Use absolute distance mode )
G94 ( Units Per Minute feed rate mode )
G64 ( Enable path blending for best speed )
(-- prologue end --)`

var genericEpilogue = `(-- epilogue begin --)
M30 (-- epilogue end --)`

var iviPrologue = `;-- prologue begin --
G17 ;Use XY plane
G21 ;Use mm
G90 ;Use absolute distance mode
G0 Z10.00 F400
G0 X-1.19 Y-11.57 F400
G0 Z5.00 F400
M3 P100 ;Start spindle clockwise, 100% Power
G0 Z5.00 F3000
;-- prologue end --`

var iviEpilogue = `;-- epilogue begin --
M5 ;Stop spindle
;-- epilogue end --`

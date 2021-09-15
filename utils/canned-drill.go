package utils

import (
	"log"

	. "github.com/gmlewis/go-gcode/gcode"
)

const (
	cannedDrillPeckClearance = 0.1 // mm
)

// Canned drilling cycle with/without dwell at bottom.
//
// Input:
// - retractZ: scalar
//      Defines the Z-coordinate to retract to (R-plane).
// - dw: scalar
//      Time to dwell at the bottom of the hole. If negative, no dwelling is
//      performed.
// - oldZ: bool
//      True assures that each cycle returns to the original
//      Z-position. False indicates to remain at the R-plane.
// - holes: vectorlist
//      List of coordinates to drill. Must include at least one entry and the
//      first vector in the list must include a Z-coordinate to define the
//      drilling depth. Each subsequent vector must include at least one X or Y
//      coordinate or possibly both. Each vector may include a Z-coordinate to
//      define a new drilling depth.
func CannedDrill(g *GCode, retractZ, dw float64, oldZ bool, holes ...Tuple) {
	prevZ := g.Position().Z()

	g.Comment("-- canned_drill R-plane=", retractZ, " dwelling=", dw, " return-to-old-Z=", oldZ, " --")
	g.Pathmode(true)

	if prevZ < retractZ {
		prevZ = retractZ
		g.GotoZ(Z(retractZ))
	}

	for _, v := range holes {
		zDrill := v.Z()
		g.GotoXY(v)
		if zDrill >= retractZ {
			// warning(pfx, "drilling at ", head(v, 2), " to depth ", zdrill, " is higher than retract-plane (", retractz, "), skipping")
			continue
		}

		g.GotoZ(Z(retractZ))
		g.MoveZ(Z(zDrill))

		if dw >= 0.0 {
			g.Dwell(dw)
		}

		if oldZ {
			g.GotoZ(Z(prevZ))
		} else {
			g.GotoZ(Z(retractZ))
		}
	}

	if oldZ {
		if prevZ > retractZ {
			g.GotoZ(Z(prevZ))
		} else if prevZ < retractZ {
			// warning(pfx, "oldz return requested, but oldZ (", prevZ,") is below retract-plane (", retractZ,"), staying at retract-plane");
		}
	}

	g.Comment("-- end canned_drill --")
}

// Canned drilling cycle with peck.
//
// Input:
// - retractZ: scalar
//      Defines the Z-coordinate to retract to (R-plane).
// - delta: scalar
//      Incremental drill depth for each peck cycle. The value of delta must be
//      larger than 0.0.
// - oldZ: bool
//      True assures that each cycle returns to the original
//      Z-position. False indicates to remain at the R-plane.
// - holes: vectorlist
//      List of coordinates to drill. Must include at least one entry and the
//      first vector in the list must include a Z-coordinate to define the
//      drilling depth. Each subsequent vector must include at least one X or Y
//      coordinate or possibly both. Each vector may include a Z-coordinate to
//      define a new drilling depth.
func CannedDrillPeck(g *GCode, retractZ, delta float64, oldZ bool, holes ...Tuple) {
	if delta <= 0.0 {
		log.Fatal("delta must be > 0")
	}

	prevZ := g.Position().Z()

	g.Comment("-- canned_drill_peck R-plane=", retractZ, " peck-increment=", delta, " return-to-old-Z=", oldZ, " --")
	g.Pathmode(true)

	clearance := 0.1 * delta
	if clearance > 2.0*cannedDrillPeckClearance {
		clearance = 2.0 * cannedDrillPeckClearance
	} else if clearance < cannedDrillPeckClearance {
		clearance = cannedDrillPeckClearance
	}

	if prevZ < retractZ {
		prevZ = retractZ
		g.GotoZ(Z(retractZ))
	}

	for _, v := range holes {
		zDrill := v.Z()
		g.GotoXY(v)
		if zDrill >= retractZ {
			// warning(pfx, "drilling at ", head(v, 2), " to depth ", zdrill, " is higher than retract-plane (", retractz, "), skipping")
			continue
		}

		g.GotoZ(Z(retractZ))

		if retractZ-delta >= zDrill {
			g.MoveZ(Z(retractZ - delta))
			g.GotoZ(Z(retractZ))
		} else {
			g.MoveZ(Z(zDrill))
			if oldZ {
				g.GotoZ(Z(prevZ))
			} else {
				g.GotoZ(Z(retractZ))
			}
			continue
		}

		var zPos float64
		for zPos = retractZ - 2.0*delta; zPos > zDrill; zPos -= delta {
			g.GotoZ(Z(zPos + delta + clearance))
			g.MoveZ(Z(zPos))
			g.GotoZ(Z(retractZ))
		}

		zPos += delta
		if zPos > zDrill {
			g.GotoZ(Z(zPos + clearance))
			g.MoveZ(Z(zDrill))
			g.GotoZ(Z(retractZ))
		}

		if oldZ {
			g.GotoZ(Z(prevZ))
		}
	}

	if oldZ {
		if prevZ > retractZ {
			g.GotoZ(Z(prevZ))
		} else if prevZ < retractZ {
			// warning(pfx, "oldz return requested, but oldz (", prevz,") is below retract-plane (", retractz,"), staying at retract-plane");
		}
	}

	g.Comment("-- end canned_drill_peck --")
}

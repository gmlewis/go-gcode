package utils

import (
	"log"
	"math"

	"github.com/gmlewis/go-gcode/gcode"
)

const (
	epsilon = 1e-8
)

// TPCOptions control how TracePathComp behaves.
type TPCOptions int

const (
	TPCArcIn  TPCOptions = 1 << iota // Arc in to the path
	TPCArcOut                        // Arc out of the path
	TPCClosed                        // Interpret path as being closed (path[0] and path[-1] are connected)
	TPCKeepZ                         // Don't move on Z-axis
	TPCLeft                          // Trace at left side of path
	TPCOldZ                          // Return to entry Z-coordinate
	TPCQuiet                         // Don't warn on unreachable inside corners
	TPCRight                         // Trace at right side of path (default)
)

// TracePathComp traces a path with tool compensation.
// Operation:
// - Goto last path entry
// - Move to cutting depth z
// - Move to all path entries
// - Return to orginal Z position
// Optional dwells for dw seconds at each path entry.
//
// Input:
// - width	Distance from path to trace at
// - path	Path to trace. The first coordinate should include a cutting
//		Z-level which is moved to before entering the path. Changing
//		the Z-level is allowed at any point as long as each XY position
//		is unique.
//
// Return value: none
func TracePathComp(g *gcode.GCode, width float64, flags TPCOptions, path ...gcode.Tuple) {
	if len(path) == 0 {
		return
	}
	if width <= 0.0 {
		log.Fatal("width must be positive")
	}

	prevZ := g.Position().Z()
	side := 1.0
	dirComment := "right"
	if flags&TPCLeft > 0 {
		side = -1.0
		dirComment = "left"
	}

	if flags&TPCKeepZ > 0 {
		z := path[0].Z()
		for i := range path {
			path[i][2] = z
		}
	}

	lastPath := path[len(path)-1]
	if math.Abs(path[0].X()-lastPath.X()) < epsilon &&
		math.Abs(path[0].Y()-lastPath.Y()) < epsilon &&
		math.Abs(path[0].Z()-lastPath.Z()) >= epsilon {
		// A manually closed path but start and end are at different Z
		log.Printf("First and last point only differ by Z coordinate, deleting last point and closing path")
		path = path[0 : len(path)-1]
		flags |= TPCClosed
	}

	normal, dir := calcDirs(side, path)

	// Find all sharp inside corners and remove them as we cannot enter them.
	// This may still leave a "reversible", but that we can handle in the main loop.

	var warnRemove int
	// TODO!!!

	if warnRemove > 0 && !(flags&TPCQuiet > 0) {
		log.Printf("Removed %v unreachable internal corner(s)", warnRemove)
	}

	// Start the trace of the path

	g.Comment("-- tracepath_comp at ", dirComment, " side at width=", width, " --")

	if flags&TPCClosed > 0 {
		// Check if the entry will collide with the exit
		lastDir := dir[len(dir)-1]
		crossP := side * tpcCrossProd(lastDir, dir[0])
		dotP := lastDir.Dot(dir[0])
		if crossP < 0.0 && dotP > 0.0 {
			log.Printf("Path entry and exit collides with exit and entry, use a >=180 degree entry angle point to prevent")
		}
	}

	// Entry into first segment
	if flags&TPCArcIn > 0 {
		p := path[0].Add(normal[0].Sub(dir[0]).MultScalar(2.0 * width))
		g.GotoXY(p)  // Rapid without Z
		g.MoveXYZ(p) // Move with Z
		// Arc into the first segment's start point
		p = path[0].Add(normal[0].MultScalar(width))
		if side > 0.0 {
			g.ArcCW(p, width, nil)
		} else {
			g.ArcCCW(p, width, nil)
		}
	} else {
		p := path[0].Add(normal[0].MultScalar(2.0 * width))
		g.GotoXY(p)  // Rapid without Z
		g.MoveXYZ(p) // Move with Z
	}

	// A closed path ends at the first point, loop once more
	npath := len(path)
	n := npath
	if flags&TPCClosed > 0 {
		n++
	}

	var i int
	//TODO!!!

	// Exit the path
	i--
	i = i % npath
	if flags&TPCArcOut > 0 {
		p := normal[i-1].Add(dir[i-1]).MultScalar(width)
		if side > 0.0 {
			g.ArcCWRel(p, width, nil)
		} else {
			g.ArcCCWRel(p, width, nil)
		}
	} else {
		p := normal[i-1].MultScalar(width)
		g.MoveXYZRel(p)
	}

	// Return to old Z if requested
	if flags&TPCOldZ > 0 {
		g.GotoZ(gcode.Z(prevZ))
	}
	g.Comment("-- tracepath_comp end --")
}

func calcDirs(side float64, path []gcode.Tuple) (npath int, path []gcode.Tuple, normal []gcode.Tuple, dir []gcode.Tuple) {
	npath = len(path)
	for i := 0; i < npath; i++ {
		pp := path[i]
		pc := path[(i+1)%npath]
		if pp.Equal(pc) {
			path = append(path[:i], path[i+1:]...)
			i--
			npath--
			continue
		}
		pc[2] = 0
		pp[2] = 0
		if pp.Equal(pc) {
			log.Printf("TracePathComp: Points %v and %v have same XY but different Z, deleting point %v", i, (i+1)%npath, i)
			path = append(path[:i], path[i+1:]...)
			i--
			npath--
			continue
		}
		dir = append(dir, pc.Sub(pp).Normalize())
		normal = append(normal, XY(side*dir[i].Y(), -side*dir[i].X()))
	}
	return npath, path, normal, dir
}

func length2D(v gcode.Tuple) float64 {
	return math.Sqrt(v.X()*v.X() + v.Y()*v.Y())
}

// Cross product divided by length returning sin(angle)
func tpcCrossProd(v1, v2 gcode.Tuple) float64 {
	crossP := v1[0]*v2[1] - v1[1]*v2[0]
	l := length2D(v1) * length2D(v2)
	if l == 0.0 {
		return 0.0
	}
	return crossP / l
}

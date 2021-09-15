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

	var normal, dir []gcode.Tuple
	path, normal, dir = calcDirs(side, path)

	var warnRemove int
	// Find all sharp inside corners and remove them as we cannot enter them.
	// This may still leave a "reversible", but that we can handle in the main loop.
	npath := len(path)
	for i := 1; i < npath; i++ {
		crossP := tpcCrossProd(dir[i-1], dir[i])
		dotP := dir[i-1].Dot(dir[i])
		if math.Abs(crossP) < epsilon && math.Abs(dotP+1) < epsilon { // crossP==0 && dotP==-1
			// Reversal of the path
			if tpcCrossProd(dir[(i+npath-2)%npath], dir[i-1])*side > 0.0 {
				// Reversal on wrong side of path, delete point and ensure
				// that any duplicate gets removed by recalculating the lot.
				path = append(path[:i], path[i+1:]...)
				path, normal, dir = calcDirs(side, path)
				npath = len(path)
				i--
				warnRemove++
			} // else the reversal is on the correct side and we may walk around.
		} else if crossP*side < 0.0 && dotP < 0.0 {
			// An internal angle detected... see if we can fit into it
			tmp := normal[i].Add(normal[i-1])
			crossP = side * tpcCrossProd(normal[i], tmp) // sin(angle/2)
			dotP = normal[i].Dot(tmp.Normalize())        // cos(angle/2)
			tmp = path[i].Sub(path[i-1])
			li := gcode.XY(tmp.X(), tmp.Y()).Magnitude() // Entry segment into sharp internal edge
			tmp = path[(i+1)%npath].Sub(path[i])
			lo := gcode.XY(tmp.X(), tmp.Y()).Magnitude() // Exit segment into sharp internal edge
			lm := width * crossP / dotP                  // Bisected calculated move distance
			if lm <= li && lm <= lo {
				continue // we can fit
			}
			warnRemove++
			if li > lo {
				// Entry segment is longer, project exit onto entry segment
				var tag float64
				if !normal[(i+1)%npath].Equal(dir[i]) {
					// Bottom of the pit has not a 90 degree turn to the top exit
					// Tag lifts bottom to 90 degree
					tag = lo * (dir[i].Dot(normal[(i+1)%npath]))
				} else {
					// Bottom already exits at 90 degrees
					tag = lo
				}
				// Move the next point onto the entry segment
				tmpDot := normal[(i+1)%npath].Dot(dir[i-1].Negate())
				path[(i+1)%npath] = path[i].Sub(dir[i-1].MultScalar(tag / tmpDot))
				// and delete this point
				path = append(path[:i], path[i+1:]...)
				dir = append(dir[:i], dir[i+1:]...)
				normal = append(normal[:i], normal[i+1:]...)
				npath--
				tpcRecalcDir(side, path, dir, normal, i-1)
				tpcRecalcDir(side, path, dir, normal, i)
				tpcRecalcDir(side, path, dir, normal, (i+1)%npath)
				i--
			} else if li < lo {
				// Exit segment is longer, project entry onto exit segment
				var tag float64
				if !normal[(i+npath-2)%npath].Equal(dir[i-1].Negate()) {
					// Bottom of the pit has not a 90 degree turn to the top entry
					// Tag lifts bottom to 90 degree
					tag = li * dir[i-1].Negate().Dot(normal[(i+npath-2)%npath])
				} else {
					// Bottom already entered at 90 degrees
					tag = li
				}
				// Move the bottom pit point up along the exit segment
				dotP = normal[(i+npath-2)%npath].Dot(dir[i])
				if math.Abs(dotP) < epsilon {
					dotP = 1.0
				}
				path[i] = path[i].Add(dir[i].MultScalar(tag / dotP))
				// and delete the previous point
				path = append(path[:i-1], path[i:]...)
				dir = append(dir[:i-1], dir[i:]...)
				normal = append(normal[:i-1], normal[i:]...)
				npath--
				tpcRecalcDir(side, path, dir, normal, (i+npath-2)%npath)
				tpcRecalcDir(side, path, dir, normal, i-1)
				tpcRecalcDir(side, path, dir, normal, i)
				i--
			} else {
				// Both entry and exit are same length, remove the point
				path = append(path[:i], path[i+1:]...)
				dir = append(dir[:i], dir[i+1:]...)
				normal = append(normal[:i], normal[i+1:]...)
				tpcRecalcDir(side, path, dir, normal, i-1)
				tpcRecalcDir(side, path, dir, normal, i)
				i--
			}
		}
	}

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
		p := path[0].Add(normal[0].MultScalar(2).Sub(dir[0]).MultScalar(width))
		// log.Printf("path[0]=%v, normal[0]=%v, dir[0]=%v, p=%v", path[0], normal[0], dir[0], p)
		g.GotoXY(p)  // Rapid without Z
		g.MoveXYZ(p) // Move with Z
		// Arc into the first segment's start point
		p = path[0].Add(normal[0].MultScalar(width))
		if side > 0.0 {
			// log.Printf("GML1: g.ArcCW(p=%v, width=%v, nil)", p, width)
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
	npath = len(path)
	n := npath
	if flags&TPCClosed > 0 {
		n++
	}

	var i int
	for i = 1; i < n; i++ {
		j := i % npath
		crossP := tpcCrossProd(dir[i-1], dir[j])
		dotP := dir[j-1].Dot(dir[j])
		if math.Abs(crossP) < epsilon {
			// Co-linear or 180 degree turn
			if dotP >= 0.0 {
				if i < npath-1 {
					// Co-linear, delete the point
					path = append(path[:i], path[i+1:]...)
					dir = append(dir[:i], dir[i+1:]...)
					normal = append(normal[:i], normal[i+1:]...)
					i--
					n--
					npath--
				} else {
					// Don't delete the last entry for closure
					g.MoveXYZ(path[j].Add(normal[j].MultScalar(width)))
				}
			} else {
				// 180 degree turn; wrong side 180'ies have already been deleted
				// Move to end of segment
				g.MoveXYZ(path[j].Sub(normal[j].MultScalar(width)))
				if i < n-1 { // Only if not last
					// Arc with 180 degrees
					if side > 0.0 {
						g.ArcCCWRel(normal[j].MultScalar(width*2.0), width, nil)
					} else {
						g.ArcCWRel(normal[j].MultScalar(width*2.0), width, nil)
					}
				}
			}
			continue
		}

		if crossP*side < 0.0 {
			// Inside angle move
			crossP = tpcCrossProd(normal[j], normal[j].Add(normal[j-1])) // sin(angle/2)
			dotP = normal[j].Dot(normal[j].Add(normal[j-1]).Normalize()) // cos(angle/2)
			// End at the projected direction of the next segment
			g.MoveXYZ(path[j].Add(normal[j].Add(dir[j].MultScalar(side * crossP / dotP)).MultScalar(width)))
		} else {
			// Outside angle move
			g.MoveXYZ(path[j].Add(normal[j-1].MultScalar(width)))
			if i < n-1 { // Only is not last
				// Arc around the angle
				if side > 0.0 {
					g.ArcCCW(path[j].Add(normal[j].MultScalar(width)), width, nil)
				} else {
					g.ArcCW(path[j].Add(normal[j].MultScalar(width)), width, nil)
				}
			}
		}
	}

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

func calcDirs(side float64, path []gcode.Tuple) (_ []gcode.Tuple, normal []gcode.Tuple, dir []gcode.Tuple) {
	npath := len(path)
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
		normal = append(normal, gcode.XY(side*dir[i].Y(), -side*dir[i].X()))
	}
	return path, normal, dir
}

func tpcRecalcDir(side float64, path, dir, normal []gcode.Tuple, idx int) {
	pp := gcode.XY(path[idx].X(), path[idx].Y())
	nxt := (idx + 1) % len(path)
	pc := gcode.XY(path[nxt].X(), path[nxt].Y())
	dir[idx] = pc.Sub(pp).Normalize()
	normal[idx] = gcode.XY(side*dir[idx].Y(), -side*dir[idx].X())
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

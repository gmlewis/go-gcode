// floret-vogel generates a design based on the example here:
// https://gitlab.com/gcmc/gcmc/blob/master/example/floret-vogel.gcmc
//
// Usage:
//   go run examples/floret-vogel/main.go > floret-vogel.gcode
package main

import (
	"fmt"
	"math"
	"sort"

	. "github.com/gmlewis/go-gcode/gcode"
	"github.com/gmlewis/go-gcode/utils"
)

func main() {
	g := gcmc()
	fmt.Printf("%v", g)
}

func gcmc() *GCode {
	g := New(UseGeneric)

	safeZ := 1.0   // Safe Z-level
	drillZ := -3.0 // Drilling depth
	g.Feedrate(600)
	g.GotoZ(Z(safeZ))

	c := 4.0
	ga := math.Pi * (3.0 - math.Sqrt(5))

	var list []Tuple
	for n := 1.0; n <= 500.0; n += 1.0 {
		r := c * math.Sqrt(n)
		p := n * ga
		list = append(list, XYZ(r*math.Sin(p), r*math.Cos(p), drillZ))
	}

	sList := XY(90, -35).Offset(sortList(list)...)
	bList := XY(0, 145).Offset(binningSort(list, 20)...)
	list = XY(-90, -35).Offset(list...)

	utils.CannedDrill(g, safeZ, -1, false, list...)
	utils.CannedDrill(g, safeZ, -1, false, sList...)
	utils.CannedDrill(g, safeZ, -1, false, bList...)

	return g
}

// sortlist() - Not-so-efficient sort-by-length routine
//
// Take a list and sort the vectors based on length between them. Starts at the
// last entry of the list and determines the closest neighbor for each next
// point. The routine runs in O(n^2), so don't give it a million points...
//
// The sorting is not optimal because it is _not_ a complete "shortest-path"
// algorithm. However, it does reduce the actual trajectory considerably.
func sortList(list []Tuple) []Tuple {
	lst := make([]*Tuple, len(list))
	for i, v := range list { // Make a copy of list
		t := XYZ(v.X(), v.Y(), v.Z())
		lst[i] = &t
	}

	res := []Tuple{*lst[len(lst)-1]}
	lst = lst[0 : len(lst)-1]
	numLeft := len(lst)
	for numLeft > 0 {
		p := res[len(res)-1] // point to measure from
		tag := -1
		var length float64
		for i, v := range lst {
			if v == nil {
				continue
			}
			d := p.Sub(*v).Magnitude()
			if tag < 0 || d < length {
				tag = i
				length = d
			}
		}
		res = append(res, *lst[tag])
		lst[tag] = nil
		numLeft--
	}
	return res
}

// binningsort() - Sort with Y-binning to band the path
//
// Makes bands of similar Y-coordinates by binning them. Each band's
// X-coordinates are then sorted and added to the result in left/right
// alternating direction.
func binningSort(lst []Tuple, nbins int) []Tuple {
	var res []Tuple

	miny, maxy := math.Inf(1), math.Inf(-1)
	for _, v := range lst {
		if v[1] < miny {
			miny = v[1]
		}
		if v[1] > maxy {
			maxy = v[1]
		}
	}

	// Create the bins based on the span of Y-coordinates and the number of bins requested
	bands := make([][]Tuple, nbins)
	dy := maxy - miny
	for _, v := range lst {
		j := int(float64(nbins) * (v[1] - miny) / dy)
		if j >= nbins {
			j = nbins - 1
		}
		bands[j] = append(bands[j], v)
	}

	for i := 0; i < nbins; i++ {
		sort.Slice(bands[i], func(a, b int) bool { return bands[i][a][0] < bands[i][b][0] })
		if i&1 != 0 {
			res = append(res, bands[i]...)
		} else {
			res = append(res, Reverse(bands[i])...)
		}
	}

	return res
}

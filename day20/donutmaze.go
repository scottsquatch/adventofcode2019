package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/scottsquatch/adventofcode2019/utils/datastructures"
	"github.com/scottsquatch/adventofcode2019/utils/fileutils"
	"github.com/scottsquatch/adventofcode2019/utils/geometry"
)

func isupperchar(r rune) bool {
	return r >= 'A' && r <= 'Z'
}

func getSpaceMap(str string) (spacemap [][]string, points []geometry.Point) {
	for y, row := range strings.Split(str, "\n") {
		var spacemaprow []string
		for x, r := range row {
			cell := string(r)
			spacemaprow = append(spacemaprow, cell)

			if cell == "." {
				points = append(points, geometry.Point{X: x, Y: y})
			}
		}
		spacemap = append(spacemap, spacemaprow)
	}

	return spacemap, points
}

func getPortals(sm [][]string) map[string][]geometry.Point {
	portals := make(map[string][]geometry.Point)
	// first row portal
	for y := 0; y < len(sm)-1; y++ {
		for x := 0; x < len(sm[y])-1; x++ {
			cell := sm[y][x]
			r := rune(cell[0])
			if isupperchar(r) {
				switch {
				case isupperchar(rune(sm[y][x+1][0])):
					str := cell + sm[y][x+1]
					var p geometry.Point
					if x > 0 && sm[y][x-1] == "." {
						p = geometry.Point{Y: y, X: x - 1}
					} else {
						p = geometry.Point{Y: y, X: x + 2}
					}

					portals[str] = append(portals[str], p)
				case isupperchar(rune(sm[y+1][x][0])):
					str := cell + sm[y+1][x]
					var p geometry.Point
					if y > 0 && sm[y-1][x] == "." {
						p = geometry.Point{Y: y - 1, X: x}
					} else {
						p = geometry.Point{Y: y + 2, X: x}
					}

					portals[str] = append(portals[str], p)
				}
			}
		}
	}

	return portals
}

func runPartA() {
	spacemap, points := getSpaceMap(fileutils.ReadFile(os.Args[1]))
	portals := getPortals(spacemap)

	g := datastructures.NewPointGraph(points)

	// Add edges to graph
	for _, u := range points {
		for _, dp := range []geometry.Point{geometry.Point{X: 0, Y: 1}, geometry.Point{X: 0, Y: -1}, geometry.Point{X: 1, Y: 0}, geometry.Point{X: -1, Y: 0}} {
			v := geometry.Point{X: u.X + dp.X, Y: u.Y + dp.Y}
			if spacemap[v.Y][v.X] == "." {
				g.AddEdge(u, v)
			}
		}
	}

	for k, ps := range portals {
		if k != "AA" && k != "ZZ" {
			g.AddEdge(ps[0], ps[1])
		}
	}

	start, end := portals["AA"][0], portals["ZZ"][0]
	type queueitem struct {
		pos   geometry.Point
		steps int
	}
	queue := []queueitem{queueitem{start, 0}}
	var next queueitem
	visited := make(map[geometry.Point]bool)
	for len(queue) > 0 {
		next, queue = queue[0], queue[1:]

		if next.pos == end {
			break
		}

		for _, v := range g.Adj(next.pos) {
			if !visited[v] {
				visited[v] = true
				queue = append(queue, queueitem{v, next.steps + 1})
			}
		}
	}

	fmt.Printf("Minimum steps to target %d\n", next.steps)
}

func runPartB() {
	spacemap, points := getSpaceMap(fileutils.ReadFile(os.Args[1]))
	portals := getPortals(spacemap)

	g := datastructures.NewPointGraph(points)

	// Add edges to graph
	for _, u := range points {
		for _, dp := range []geometry.Point{geometry.Point{X: 0, Y: 1}, geometry.Point{X: 0, Y: -1}, geometry.Point{X: 1, Y: 0}, geometry.Point{X: -1, Y: 0}} {
			v := geometry.Point{X: u.X + dp.X, Y: u.Y + dp.Y}
			if spacemap[v.Y][v.X] == "." {
				g.AddEdge(u, v)
			}
		}
	}

	outer, inner := make(map[geometry.Point]geometry.Point), make(map[geometry.Point]geometry.Point)
	pointToPortal := make(map[geometry.Point]string)
	for k, ps := range portals {
		if k != "AA" && k != "ZZ" {
			g.AddEdge(ps[0], ps[1])
			if ps[0].Y == 2 || ps[0].X == 2 || ps[0].Y == len(spacemap)-4 || ps[0].X == len(spacemap[ps[0].Y])-4 {
				pointToPortal[ps[0]] = "-" + k
				pointToPortal[ps[1]] = "+" + k
				outer[ps[0]] = ps[1]
				inner[ps[1]] = ps[0]
			} else {
				pointToPortal[ps[1]] = "-" + k
				pointToPortal[ps[0]] = "+" + k
				outer[ps[1]] = ps[0]
				inner[ps[0]] = ps[1]
			}
		} else {
			pointToPortal[ps[0]] = "-" + k
		}
	}

	start, end := portals["AA"][0], portals["ZZ"][0]
	type state struct {
		pos   geometry.Point
		level int
	}
	type queueitem struct {
		state state
		steps int
	}
	target, begin := state{end, 0}, state{start, 0}
	queue := []queueitem{queueitem{begin, 0}}
	var next queueitem
	pathTo := make(map[state]state)
	pathTo[begin] = begin
	for len(queue) > 0 {
		next, queue = queue[0], queue[1:]

		if next.state == target {
			break
		}

		otoin, visitingouter := outer[next.state.pos]
		itoout, visitinginner := inner[next.state.pos]
		for _, v := range g.Adj(next.state.pos) {
			s := state{v, next.state.level}
			skip := false
			if visitingouter && v == otoin {
				if next.state.level == 0 {
					skip = true
				} else {
					s.level--
				}
			} else if visitinginner && v == itoout {
				s.level++
			} else if (v == start || v == end) && next.state.level > 0 {
				skip = true
			}

			_, visited := pathTo[s]
			if !skip && !visited {
				pathTo[s] = next.state
				queue = append(queue, queueitem{s, next.steps + 1})
			}
		}
	}
	fmt.Printf("Minimum steps to target (recursive maze) %d\n", next.steps)
}

func main() {
	runPartA()
	runPartB()
}

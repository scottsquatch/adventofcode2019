package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/scottsquatch/adventofcode2019/utils/fileutils"

	"github.com/scottsquatch/adventofcode2019/utils/geometry"
)

func runPartA() {
	g := [][]string{}
	for _, row := range strings.Split(fileutils.ReadFile(os.Args[1]), "\n") {
		gRow := []string{}
		for _, cell := range row {
			gRow = append(gRow, string(cell))
		}
		g = append(g, gRow)
	}

	fmt.Println(shortestSteps(g, 0))
}

func printCave(g [][]string) {
	for _, row := range g {
		for _, cell := range row {
			fmt.Print(cell)
		}
		fmt.Println()
	}
}

const (
	wall   = "#"
	ground = "."
	us     = "@"
	keys   = "abcdefghijklmnopqrstuvwxyz"
	doors  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func shortestSteps(g [][]string, steps int) int {
	keys, doors, start, paths := getLandmarks(g)
	fmt.Printf("K: %v, D: %v, S: %v, P: %v\n", keys, doors, start, paths)
	if len(keys) == 0 {
		return steps
	}

	for k, p := range keys {
		newG := cloneCave(g)
		newG[p.Y][p.X] = us
		newG[start.Y][start.X] = ground
		q, found := doors[strings.ToUpper(k)]
		if found {
			newG[q.Y][q.X] = ground
		}

		subPaths := getPaths(newG)
		for _, path := range subPaths {
			paths = append(paths, append([]string{k}, path...))
		}

	}

	return paths
}

func cloneCave(g [][]string) [][]string {
	newG := make([][]string, len(g))

	for i, row := range g {
		newG[i] = make([]string, len(row))
		copy(newG[i], row)
	}

	return newG
}

func getLandmarks(g [][]string) (keys, doors map[string]geometry.Point, startpoint geometry.Point, paths map[string][]geometry.Point) {
	keys, doors = make(map[string]geometry.Point), make(map[string]geometry.Point)
	for y, row := range g {
		for x, cell := range row {
			if cell == us {
				startpoint = geometry.Point{X: x, Y: y}
			} else if isDoor(cell) {
				doors[cell] = geometry.Point{X: x, Y: y}
			}
		}
	}

	pathTo := make(map[geometry.Point]*geometry.Point)
	queue := []geometry.Point{startpoint}
	pathTo[startpoint] = nil
	var p geometry.Point
	for len(queue) > 0 {
		p, queue = queue[0], queue[1:]

		for _, q := range getNeighbors(g, p) {
			_, visited := pathTo[q]
			if !visited {
				pathTo[q] = &p

				cell := g[q.Y][q.X]
				if isKey(cell) {
					keys[cell] = q
				} else if cell == ground {
					queue = append(queue, q)
				}
			}
		}
	}

	paths = make(map[string][]geometry.Point)
	for k, p := range keys {
		pth := []geometry.Point{}
		next := pathTo[p]
		for pathTo != nil {
			pth = append([]geometry.Point{*next}, pth...)
			next = pathTo[*next]
		}
		paths[k] = pth
	}

	return keys, doors, startpoint, paths
}

func getNeighbors(g [][]string, p geometry.Point) []geometry.Point {
	directions := []geometry.Point{geometry.Point{X: p.X, Y: p.Y - 1}, geometry.Point{X: p.X, Y: p.Y + 1}, geometry.Point{X: p.X - 1, Y: p.Y}, geometry.Point{X: p.X + 1, Y: p.Y}}
	neighbors := []geometry.Point{}
	for _, q := range directions {
		if isValid(g, q) {
			cell := g[q.Y][q.X]
			if cell == ground || isKey(cell) {
				neighbors = append(neighbors, q)
			}
		}
	}

	return neighbors
}

func isValid(g [][]string, p geometry.Point) bool {
	return p.Y >= 0 && p.Y < len(g) && p.X >= 0 && p.X < len(g[p.Y])
}

func isKey(cell string) bool {
	return strings.Contains(keys, cell)
}

func isDoor(cell string) bool {
	return strings.Contains(doors, cell)
}

func runPartB() {

}

func main() {
	runPartA()
	runPartB()
}

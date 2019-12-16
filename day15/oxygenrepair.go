package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/scottsquatch/adventofcode2019/utils/fileutils"

	"github.com/scottsquatch/adventofcode2019/common"

	"github.com/scottsquatch/adventofcode2019/utils/geometry"
)

const (
	north int64 = 1
	south int64 = 2
	west  int64 = 3
	east  int64 = 4
)

const (
	blocked int64 = 0
	moved   int64 = 1
	found   int64 = 2
)

const (
	wall       = "#"
	floor      = "."
	oxygentank = "O"
	droid      = "D"
)

func getPos(direction int64, curr geometry.Point) geometry.Point {
	switch direction {
	case north:
		return geometry.Point{X: curr.X, Y: curr.Y - 1}
	case south:
		return geometry.Point{X: curr.X, Y: curr.Y + 1}
	case west:
		return geometry.Point{X: curr.X - 1, Y: curr.Y}
	case east:
		return geometry.Point{X: curr.X + 1, Y: curr.Y}
	default:
		panic(fmt.Sprintf("direction %d is not valid\n", direction))
	}
}

func getOpposite(direction int64) int64 {
	switch direction {
	case north:
		return south
	case south:
		return north
	case west:
		return east
	case east:
		return west
	default:
		panic(fmt.Sprintf("direction %d is not valid\n", direction))
	}
}

func printMap(m map[geometry.Point]string) {
	// Find bounding box
	maxX, minX, maxY, minY := 0, 0, 0, 0

	for p := range m {
		if p.X > maxX {
			maxX = p.X
		} else if p.X < minX {
			minX = p.X
		}

		if p.Y > maxY {
			maxY = p.Y
		} else if p.Y < minY {
			minY = p.Y
		}
	}

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			s, f := m[geometry.Point{X: x, Y: y}]

			if !f {
				fmt.Print(" ")
			} else if x == 0 && y == 0 {
				fmt.Print("D")
			} else {
				fmt.Print(s)
			}
		}

		fmt.Println()
	}
}

var dirs = []int64{north, south, west, east}

func bfs(start, end geometry.Point, sysmap map[geometry.Point]string) []geometry.Point {
	pathTo := make(map[geometry.Point]geometry.Point)
	queue := []geometry.Point{start}
	cur := geometry.Point{}
	for cur != end {
		cur, queue = queue[0], queue[1:]

		for _, d := range dirs {
			p := getPos(d, cur)
			_, visited := pathTo[p]
			if sysmap[p] != wall && !visited {
				pathTo[p] = cur

				queue = append(queue, p)
			}
		}
	}

	path := []geometry.Point{end}
	cur = pathTo[cur]
	for cur != start {
		path = append([]geometry.Point{cur}, path...)
		cur = pathTo[cur]
	}
	path = append([]geometry.Point{cur}, path...)

	return path
}

func runPartA() (map[geometry.Point]string, geometry.Point) {
	input, output := make(chan int64), make(chan int64)
	program := common.NewIntCodeProgram(fileutils.ReadFile(os.Args[1]))
	pc := common.NewComputer(input, output)
	go pc.Run(program)
	oxysystemmap := make(map[geometry.Point]string)
	oxysystemmap[geometry.Point{}] = floor
	cur := geometry.Point{}
	var oxy geometry.Point

	// Was having trouble getting maze solving algorithms to work, settled on a brute force
	// way. not guaranteed to work but should work for our input
	for i := 0; i < 1000000; i++ {
		// Random walk
		dir := rand.Int63n(4) + 1
		for oxysystemmap[getPos(dir, cur)] == wall {
			dir = rand.Int63n(4) + 1
		}

		input <- dir
		r := <-output
		p := getPos(dir, cur)

		switch r {
		case found:
			oxy = cur
			oxysystemmap[p] = oxygentank
			cur = p
		case moved:
			oxysystemmap[p] = floor
			cur = p
		case blocked:
			oxysystemmap[p] = wall
		}
	}

	printMap(oxysystemmap)
	fmt.Println(oxy)

	fmt.Printf("Fewest number of moveset commands: %d\n", len(bfs(geometry.Point{}, oxy, oxysystemmap)))

	return oxysystemmap, oxy
}

func turnLeft(facing int64) int64 {
	switch facing {
	case north:
		return west
	case south:
		return east
	case west:
		return south
	case east:
		return north
	default:
		panic(fmt.Sprintf("direction of %d is not valid", facing))
	}
}

func turnRight(facing int64) int64 {
	return getOpposite(turnLeft(facing))
}

func runPartB(sysmap map[geometry.Point]string, oxy geometry.Point) {
	type queueitem struct {
		pos    geometry.Point
		minute int
	}

	queue := []queueitem{queueitem{oxy, 1}}
	var cur queueitem
	maxminutes := 0
	for len(queue) != 0 {
		cur, queue = queue[0], queue[1:]
		if cur.minute > maxminutes {
			maxminutes = cur.minute
		}
		for _, d := range dirs {
			p := getPos(d, cur.pos)
			if sysmap[p] == floor {
				sysmap[p] = oxygentank
				queue = append(queue, queueitem{p, cur.minute + 1})
			}
		}
	}

	fmt.Printf("Oxygen fills area after %d minutes\n", maxminutes)
}

func main() {
	m, o := runPartA()
	runPartB(m, o)
}

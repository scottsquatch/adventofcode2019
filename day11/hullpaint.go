package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/scottsquatch/adventofcode2019/common"
	"github.com/scottsquatch/adventofcode2019/utils/fileutils"
	"github.com/scottsquatch/adventofcode2019/utils/geometry"
)

type direction int

const (
	up    direction = 0
	down  direction = 1
	left  direction = 2
	right direction = 3
)

const (
	turnLeft  int64 = 0
	turnRight int64 = 1
)

func runPartA() {
	program := common.NewIntCodeProgram(fileutils.ReadFile(os.Args[1]))
	robot := NewHullPaintingRobot(*program)
	painted := make(map[geometry.Point]int64)
	position := geometry.Point{X: 0, Y: 0}
	facing := up
	mux := &sync.Mutex{}
	paint := func(color int64) {
		mux.Lock()
		c, found := painted[position]
		if !found || c != color {
			painted[position] = color
		}
		mux.Unlock()
	}
	move := func(action int64) {
		switch action {
		case turnLeft:
			switch facing {
			case up:
				facing = left
				position.X--
			case down:
				facing = right
				position.X++
			case left:
				facing = down
				position.Y++
			case right:
				facing = up
				position.Y--
			}
		case turnRight:
			switch facing {
			case up:
				facing = right
				position.X++
			case down:
				facing = left
				position.X--
			case left:
				facing = up
				position.Y--
			case right:
				facing = down
				position.Y++
			}
		}
	}
	camera := func() int64 {
		mux.Lock()
		color, found := painted[position]
		mux.Unlock()
		if found {
			return color
		}

		return 0
	}
	robot.Start(paint, move, camera)

	fmt.Printf("Number of panels that are painted at least once: %d\n", len(painted))
}

func runPartB() {
	program := common.NewIntCodeProgram(fileutils.ReadFile(os.Args[1]))
	robot := NewHullPaintingRobot(*program)
	painted := make(map[geometry.Point]int64)
	position := geometry.Point{X: 0, Y: 0}
	painted[position] = 1
	facing := up
	mux := &sync.Mutex{}
	paint := func(color int64) {
		mux.Lock()
		c, found := painted[position]
		if !found || c != color {
			painted[position] = color
		}
		mux.Unlock()
	}
	move := func(action int64) {
		switch action {
		case turnLeft:
			switch facing {
			case up:
				facing = left
				position.X--
			case down:
				facing = right
				position.X++
			case left:
				facing = down
				position.Y++
			case right:
				facing = up
				position.Y--
			}
		case turnRight:
			switch facing {
			case up:
				facing = right
				position.X++
			case down:
				facing = left
				position.X--
			case left:
				facing = up
				position.Y--
			case right:
				facing = down
				position.Y++
			}
		}
	}
	camera := func() int64 {
		mux.Lock()
		color, found := painted[position]
		mux.Unlock()
		if found {
			return color
		}

		return 0
	}
	robot.Start(paint, move, camera)

	// Print out
	min := geometry.Point{X: 123123123, Y: 1232142414}
	max := geometry.Point{X: -12312312, Y: -123123123}
	// Get Bounding Box
	for p := range painted {
		if p.X < min.X {
			min.X = p.X
		}
		if p.X > max.X {
			max.X = p.X
		}
		if p.Y < min.Y {
			min.Y = p.Y
		}
		if p.Y > max.Y {
			max.Y = p.Y
		}
	}
	for y := min.Y; y <= max.Y; y++ {
		for x := min.X; x <= max.X; x++ {
			p := geometry.Point{X: x, Y: y}

			c, found := painted[p]
			if found && c == 1 {
				fmt.Print("*")
			} else {
				fmt.Print(" ")
			}
		}

		fmt.Println()
	}
}

func main() {
	runPartA()
	runPartB()
}

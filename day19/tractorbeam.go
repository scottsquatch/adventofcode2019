package main

import (
	"fmt"
	"os"

	"github.com/scottsquatch/adventofcode2019/utils/geometry"

	"github.com/scottsquatch/adventofcode2019/common"
	"github.com/scottsquatch/adventofcode2019/utils/fileutils"
)

func runPartA() {
	affectedcells := 0
	spacemap := make([][]string, 50)
	for y := 0; y < 50; y++ {
		spacemap[y] = make([]string, 50)
		for x := 0; x < 50; x++ {
			c := calldroid(x, y)
			switch c {
			case 0:
				spacemap[y][x] = "."
			case 1:
				spacemap[y][x] = "#"
				affectedcells++
			default:
				panic(fmt.Sprintf("Output of %d is not valid\n", c))
			}
		}
	}

	printmap(spacemap)
	fmt.Printf("Number of cells affected by tractor beam %d\n", affectedcells)
}

func calldroid(x, y int) int {
	in, out := make(chan int64), make(chan int64)
	defer func() { close(in); close(out) }()
	go func() {
		software := common.NewIntCodeProgram(fileutils.ReadFile(os.Args[1]))
		drone := common.NewComputer(in, out)
		drone.Run(software)
	}()
	in <- int64(x)
	in <- int64(y)
	c := <-out

	return int(c)
}

func printmap(m [][]string) {
	for _, row := range m {
		for _, c := range row {
			fmt.Print(c)
		}
		fmt.Println()
	}
}
func runPartB() {
	N := 100
	beamstartx := 1
	var topleft geometry.Point
	c := 0
	for y := N; c == 0; y++ {
		// find start of beam
		c := calldroid(beamstartx, y)
		for c == 0 {
			beamstartx++
			c = calldroid(beamstartx, y)
		}

		c = calldroid(beamstartx+N-1, y-(N-1))
		if c == 1 {
			topleft = geometry.Point{X: beamstartx, Y: y - (N - 1)}
			break
		}
	}

	fmt.Printf("topleft: %v. Answer is %d\n", topleft, topleft.X*10000+topleft.Y)
}

func main() {
	runPartA()
	runPartB()
}

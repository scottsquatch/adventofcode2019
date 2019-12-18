package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/scottsquatch/adventofcode2019/common"
	"github.com/scottsquatch/adventofcode2019/utils/fileutils"
)

func runPartA() [][]string {
	input, output := make(chan int64), make(chan int64)
	program := common.NewIntCodeProgram(fileutils.ReadFile(os.Args[1]))
	robot := common.NewComputer(input, output)

	go func() {
		robot.Run(program)
		close(output)
	}()

	var str strings.Builder
	for c := range output {
		str.WriteString(string(int(c)))
	}

	var floor [][]string
	for _, r := range strings.Split(str.String(), "\n") {
		var row []string
		for _, c := range r {
			row = append(row, string(c))
		}
		floor = append(floor, row)
	}

	sum := 0
	for y, r := range floor {
		for x := range r {
			if x > 0 && x < len(floor[y])-1 && y > 0 && y < len(floor)-3 {
				if floor[y][x] == "#" && floor[y+1][x] == "#" && floor[y-1][x] == "#" && floor[y][x-1] == "#" && floor[y][x+1] == "#" {
					sum += y * x
				}
			}
		}
	}

	fmt.Printf("Sum of alignment paramters is %d\n", sum)

	return floor
}

func runPartB(floor [][]string) {
	traverseMaze(floor)

	input, output := make(chan int64, 1), make(chan int64)
	program := common.NewIntCodeProgram(fileutils.ReadFile(os.Args[1]))
	program.ProgramData[0] = 2
	robot := common.NewComputer(input, output)

	go func() {
		robot.Run(program)
		close(output)
	}()

	go func() {

		// Path compression determined manually
		movements := "A,A,B,B,C,B,C,B,C,A\n"
		a := "L,10,L,10,R,6\n"
		b := "R,12,L,12,L,12\n"
		c := "L,6,L,10,R,12,R,12\n"
		feed := "n\n"

		inst := []string{movements, a, b, c, feed}
		for _, i := range inst {
			for _, c := range i {
				input <- int64(c)
			}
		}

	}()
	for c := range output {
		if c >= 0 && c <= 127 {
			fmt.Print(string(int(c)))
		} else {
			fmt.Printf("Spacedust collected: %d\n", c)
		}
	}
}

func traverseMaze(floor [][]string) {
	var x, y, dir int
	for i, r := range floor {
		for j, c := range r {
			switch c {
			case "^":
				dir = north
				x, y = j, i
				floor[y][x] = "@"
				break
			case "<":
				dir = east
				floor[y][x] = "@"
				x, y = j, i
				break
			case ">":
				dir = west
				floor[y][x] = "@"
				x, y = j, i
				break
			case "v", "V":
				dir = south
				floor[y][x] = "@"
				x, y = j, i
				break
			}
		}
	}

	// Find path through maze using Left-Hand method
	steps := 0
	var directions []string
	for hasUnvisitedScaffolding(floor) {
		if canGoForward(floor, x, y, dir) {
			floor[y][x] = "@"
			steps++
			x, y = move(x, y, dir)
		} else if canGoLeft(floor, x, y, dir) {
			floor[y][x] = "@"
			directions = append(directions, strconv.FormatInt(int64(steps), 10), "L")
			steps = 0
			dir = turnLeft(dir)
		} else if canGoRight(floor, x, y, dir) {
			floor[y][x] = "@"
			directions = append(directions, strconv.FormatInt(int64(steps), 10), "R")
			steps = 0
			dir = turnRight(dir)
		} else {
			directions = append(directions, strconv.FormatInt(int64(steps), 10), "L", "L")
			steps = 0
			dir = turnLeft(turnLeft(dir))
		}
	}

	N := len(directions)
	fmt.Printf("Directions for robot %v\n", directions[1:N-2])
}

func printMap(m [][]string) {
	for _, r := range m {
		for _, c := range r {
			fmt.Print(c)
		}
		fmt.Println()
	}
}

func move(x, y, direction int) (int, int) {

	switch direction {
	case north:
		y--
	case south:
		y++
	case east:
		x++
	case west:
		x--
	}

	return x, y
}

func canGoForward(floor [][]string, x, y, direction int) bool {
	newX, newY := move(x, y, direction)

	return newY >= 0 && newY < len(floor) && newX >= 0 && newX < len(floor[newY]) && (floor[newY][newX] == "#" || floor[newY][newX] == "@")
}

func canGoLeft(floor [][]string, x, y, direction int) bool {
	newX, newY := move(x, y, turnLeft(direction))

	return newY >= 0 && newY < len(floor) && newX >= 0 && newX < len(floor[newY]) && (floor[newY][newX] == "#" || floor[newY][newX] == "@")
}

func canGoRight(floor [][]string, x, y, direction int) bool {
	newX, newY := move(x, y, turnRight(direction))
	return newY >= 0 && newY < len(floor) && newX >= 0 && newX < len(floor[newY]) && (floor[newY][newX] == "#" || floor[newY][newX] == "@")
}

func turnLeft(direction int) int {
	switch direction {
	case north:
		return west
	case south:
		return east
	case east:
		return north
	case west:
		return south
	default:
		panic(fmt.Sprintf("direction of %d is not valid\n", direction))
	}
}

func turnRight(direction int) int {
	switch direction {
	case north:
		return east
	case south:
		return west
	case east:
		return south
	case west:
		return north
	default:
		panic(fmt.Sprintf("direction of %d is not valid\n", direction))
	}
}

const (
	north = 0
	south = 1
	east  = 2
	west  = 3
)

func hasUnvisitedScaffolding(m [][]string) bool {
	for _, r := range m {
		for _, c := range r {
			if c == "#" {
				return true
			}
		}
	}

	return false
}

func main() {
	f := runPartA()
	runPartB(f)
}

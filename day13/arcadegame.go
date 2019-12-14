package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/scottsquatch/adventofcode2019/common"
	"github.com/scottsquatch/adventofcode2019/utils/fileutils"
)

func printScreen(screen [][]string) {
	for _, row := range screen {
		for _, p := range row {
			fmt.Print(p)
		}
		fmt.Println()
	}
}

var mu sync.Mutex

func runPartA() {
	game := common.NewIntCodeProgram(fileutils.ReadFile(os.Args[1]))
	joystick := make(chan int64)
	output := make(chan int64, 1)
	cpu := common.NewComputer(joystick, output)

	blocks := 0
	go func() {
		cpu.Run(game)
		close(output)
	}()

	var x, y int
	mode := 0
	screen := make([][]string, 1)
	for i := range output {
		switch mode {
		case 0:
			x = int(i)
		case 1:
			y = int(i)
		default:
			if y >= len(screen) {
				newScreen := make([][]string, y+1)
				copy(newScreen, screen)
				screen = newScreen
			}
			if x >= len(screen[y]) {
				newRow := make([]string, x+1)
				copy(newRow, screen[y])
				screen[y] = newRow
			}

			screen[y][x] = getTile(i)

			if i == 2 {
				blocks++
			}
		}

		mode = (mode + 1) % 3
	}

	printScreen(screen)
	fmt.Printf("Amount of blocks in map %d\n", blocks)
}

func getTile(tileID int64) string {
	switch tileID {
	case 0:
		return " "
	case 1:
		return "|"
	case 2:
		return "B"
	case 3:
		return "-"
	case 4:
		return "*"
	}

	panic(fmt.Sprintf("tileID of %d is not valid\n", tileID))
}

func runPartB() {
	game := common.NewIntCodeProgram(fileutils.ReadFile(os.Args[1]))
	joystick := make(chan int64, 1)
	output := make(chan int64, 1)
	cpu := common.NewComputer(joystick, output)
	game.ProgramData[0] = 2

	go func() {
		cpu.Run(game)
		close(output)
	}()

	ballPos := make([]int, 2)
	paddlePos := make([]int, 2)
	var x, y int
	mode := 0
	score := 0
	screen := make([][]string, 1)
	for i := range output {
		switch mode {
		case 0:
			x = int(i)
		case 1:
			y = int(i)
		default:
			if x == -1 && y == 0 {
				score = int(i)
			} else {
				if y >= len(screen) {
					newScreen := make([][]string, y+1)
					copy(newScreen, screen)
					screen = newScreen
				}
				if x >= len(screen[y]) {
					newRow := make([]string, x+1)
					copy(newRow, screen[y])
					screen[y] = newRow
				}

				screen[y][x] = getTile(i)

				if i == 4 {
					ballPos[0], ballPos[1] = y, x
					if paddlePos[1] < x {
						joystick <- 1
					} else if paddlePos[1] > x {
						joystick <- (-1)
					} else {
						joystick <- 0
					}
				} else if i == 3 {
					paddlePos[0], paddlePos[1] = y, x
				}
			}
		}
		mode = (mode + 1) % 3

	}
	fmt.Printf("Final score: %d\n", score)
}

func main() {
	runPartA()
	runPartB()
}

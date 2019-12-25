package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/scottsquatch/adventofcode2019/utils/geometry"

	"github.com/scottsquatch/adventofcode2019/utils/fileutils"
)

func encodestate(statemap [][]string) int {
	state := 0
	cellnum := 0
	for _, row := range statemap {
		for _, cell := range row {
			if cell == "#" {
				state |= (1 << cellnum)
			}
			cellnum++
		}
	}

	return state
}

func runPartA() {
	// Enocde the state into an integer where a 1 in the nth bit means the nth cell has a bug
	statemap := make(map[int]bool)

	var bugmap [][]string
	for _, line := range strings.Split(fileutils.ReadFile(os.Args[1]), "\r\n") {
		var maprow []string
		for _, r := range line {
			maprow = append(maprow, string(r))
		}
		bugmap = append(bugmap, maprow)
	}

	state := encodestate(bugmap)
	minutes := 0
	for !statemap[state] {
		statemap[state] = true
		// Run simulation
		nextmap := make([][]string, len(bugmap))
		for y := range bugmap {
			nextmap[y] = make([]string, len(bugmap[y]))
			for x := range bugmap[y] {
				cell := bugmap[y][x]
				adjacentbugs := getnumadjacentbugs(bugmap, x, y)
				nextcell := bugmap[y][x]
				if cell == "#" {
					switch adjacentbugs {
					case 1:
						nextcell = "#"
					default:
						nextcell = "."
					}
				} else if cell == "." && (adjacentbugs == 1 || adjacentbugs == 2) {
					nextcell = "#"
				}

				nextmap[y][x] = nextcell
			}
		}

		bugmap = nextmap
		state = encodestate(bugmap)
		minutes++
	}

	fmt.Printf("Found duplicate state after %d minutes\n", minutes)
	printmap(bugmap)
	fmt.Printf("Biodiversity %d\n", calculatebiodiversity(state))
}

func calculatebiodiversity(state int) int64 {
	power := 1
	var bd int64 = 0
	for state > 0 {
		if state&1 == 1 {
			bd += int64(power)
		}

		power <<= 1
		state >>= 1
	}

	return bd
}

func printmap(m [][]string) {
	for _, row := range m {
		for _, cell := range row {
			fmt.Print(cell)
		}
		fmt.Println()
	}
}

func getnumadjacentbugs(m [][]string, x, y int) int {
	num := 0
	for _, dp := range []geometry.Point{geometry.Point{X: 0, Y: 1}, geometry.Point{X: 0, Y: -1}, geometry.Point{X: 1, Y: 0}, geometry.Point{X: -1, Y: 0}} {
		p := geometry.Point{X: x + dp.X, Y: y + dp.Y}
		if p.Y >= 0 && p.Y < len(m) && p.X >= 0 && p.X < len(m[p.Y]) && m[p.Y][p.X] == "#" {
			num++
		}
	}

	return num
}

func runPartB() {
	bugmap := make([][][]string, 1)
	baseidx := 0
	for _, line := range strings.Split(fileutils.ReadFile(os.Args[1]), "\r\n") {
		var maprow []string
		for _, r := range line {
			maprow = append(maprow, string(r))
		}
		bugmap[baseidx] = append(bugmap[baseidx], maprow)
	}

	for i := 0; i < 200; i++ {
		if numbugs(bugmap[0]) > 0 {
			bugmap = append([][][]string{makeemptygrid(5, 5)}, bugmap...)
			baseidx++
		}
		if numbugs(bugmap[len(bugmap)-1]) > 0 {
			bugmap = append(bugmap, makeemptygrid(5, 5))
		}
		// Run simulation
		nextmap := make([][][]string, len(bugmap))
		for z, levelmap := range bugmap {
			nextmap[z] = make([][]string, len(bugmap[z]))
			for y, row := range levelmap {
				nextmap[z][y] = make([]string, len(bugmap[z][y]))
				for x, cell := range row {
					if x != 2 || y != 2 {

						adjacentbugs := getreecursivenumadjacentbugs(bugmap, x, y, z, baseidx)
						nextcell := bugmap[z][y][x]
						if cell == "#" {
							switch adjacentbugs {
							case 1:
								nextcell = "#"
							default:
								nextcell = "."
							}
						} else if cell == "." && (adjacentbugs == 1 || adjacentbugs == 2) {
							nextcell = "#"
						}

						nextmap[z][y][x] = nextcell
					}
				}
			}
		}

		bugmap = nextmap
		// for z, bugmaplevel := range bugmap {
		// 	fmt.Printf("\nLevel: %d\n", z-baseidx)
		// 	printmap(bugmaplevel)
		// }
		fmt.Printf("Number of bugs %d\n", numbugsalllevels(bugmap))
	}
}

func numbugsalllevels(m [][][]string) int {
	num := 0
	for _, ml := range m {
		num += numbugs(ml)
	}

	return num
}

func makeemptygrid(m, n int) [][]string {
	grid := make([][]string, n)
	for y := 0; y < n; y++ {
		grid[y] = make([]string, m)
		for x := 0; x < m; x++ {
			grid[y][x] = "."
		}
	}

	return grid
}

func getreecursivenumadjacentbugs(m [][][]string, x, y, z, baselevel int) int {
	num := 0
	for _, dp := range []geometry.Point{geometry.Point{X: 0, Y: 1}, geometry.Point{X: 0, Y: -1}, geometry.Point{X: 1, Y: 0}, geometry.Point{X: -1, Y: 0}} {
		p := geometry.Point{X: x + dp.X, Y: y + dp.Y}
		cell := "."
		if p.X == -1 {
			if z > 0 {
				cell = m[z-1][2][1]
			}
		} else if p.Y == -1 {
			if z > 0 {
				cell = m[z-1][1][2]
			}
		} else if p.Y == len(m[z]) {
			if z > 0 {
				cell = m[z-1][3][2]
			}
		} else if p.X == len(m[z][p.Y]) {
			if z > 0 {
				cell = m[z-1][2][3]
			}
		} else if p.Y == 2 && p.X == 2 {
			if x == 2 && y == 1 {
				if z < len(m)-1 {
					for _, cell := range m[z+1][0] {
						if cell == "#" {
							num++
						}
					}
				}
			} else if x == 1 && y == 2 {
				if z < len(m)-1 {
					for _, row := range m[z+1] {
						if row[0] == "#" {
							num++
						}
					}
				}
			} else if x == 2 && y == 3 {
				if z < len(m)-1 {
					for _, cell := range m[z+1][len(m[z+1])-1] {
						if cell == "#" {
							num++
						}
					}
				}
			} else {
				if z < len(m)-1 {
					for _, row := range m[z+1] {
						if row[len(m[z+1][0])-1] == "#" {
							num++
						}
					}
				}
			}
		} else {
			cell = m[z][p.Y][p.X]
		}

		if cell == "#" {
			num++
		}
	}

	return num
}

func numbugs(m [][]string) int {
	count := 0
	for _, row := range m {
		for _, cell := range row {
			if cell == "#" {
				count++
			}
		}
	}

	return count
}

func main() {
	runPartA()
	runPartB()
}

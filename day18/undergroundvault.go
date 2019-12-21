package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/scottsquatch/adventofcode2019/utils/fileutils"

	"github.com/scottsquatch/adventofcode2019/utils/geometry"
)

func runPartA(cave [][]string, start geometry.Point, allKeys map[string]geometry.Point) {
	fmt.Println(minStepsToSolve(cave, start, allKeys))
}

func minStepsToSolve(g [][]string, start geometry.Point, allKeys map[string]geometry.Point) int {
	// method inspired by /r/jonathan_paulson, post: https://www.reddit.com/r/adventofcode/comments/ec8090/2019_day_18_solutions/fb9wfnz/?utm_source=share&utm_medium=web2x,
	// solution: https://topaz.github.io/paste/#XQAAAQCFBAAAAAAAAAA0m0pnuFI8c/fBNAn6x25rti77on4e8DYCelyI4Xj/SWYY0/Ntl1UF+LeIl0XT1SetTywqGks/U6aRCB49U12enKDcqoR8FDA7Qi3g59f4xFVpL0j6CGRf9DbelJFU7an7OE1hRCw9ohXPo4pT+PpX094PH4p0YMBYLNq7ai3VQC9C6ukfZw07M/ASruA/9p/Jzxdq3WIeyfEpRidKJybeGVMba9CAzr6nmEbPVCa0Qc4jim7qRRQzEG5j/2umX3l/aFbeE1dSCN3fpyho5xVCtXQYH1DKeouxbVphzC5Cru3KlLZYmPaPrSYqedT/+hpt4uxmdUjkk+JK9GLHplxU4nhNO3xN18Z/qquNgF4260hc7F9uMk1+oNAz3VGPJ5sbiyueMZe9ELPOwEz9WRT3cfsc+W30x0Xb89gkr9hWafJaYnZxQTyR6JdMNsGesR6PIreWAy89AlI2fJ/ztLFsEL+r4/fTtMWSTHSvBz4/H+JwwoZFzPvyyjkaqMpdpUXiBVAPosVE8+ZEk0T8TbmF6EHxeHI+74+VcDnoTlw/NyQiLWXhApd5TK27fjIoki8IW9dycouBhgzfmVXdSgdGFQfLlcTcvdeNq/J5LrdFuOyJKtQkz4yPrWmy80bms/lgQwDqg8BhgGdUiWUyr05EB+fkoN5yq+g3+q90sYbb9ts0x1mxTMnJsi0I6RnE1pDxkIJa3zLO7FuJQyasyGlYU6+OV1fizp7n1RR2+Yv/+s4O5A==
	// Idea is to turn the 2D field into 3D where the 3rd dimension is the keys which have been obtained
	type queueitem struct {
		pos   geometry.Point
		steps int
		keys  map[string]bool
	}

	type visitedkey struct {
		pos  geometry.Point
		keys string
	}
	fk := func(pos geometry.Point, keys map[string]bool) visitedkey {
		var keyset []string
		for k := range keys {
			keyset = append(keyset, k)
		}

		sort.Strings(keyset)
		return visitedkey{pos, strings.Join(keyset, "")}
	}

	visited := make(map[visitedkey]bool)
	keys := make(map[string]bool)
	queue := []queueitem{queueitem{start, 0, make(map[string]bool)}}
	visited[fk(start, keys)] = true
	var next queueitem
	for len(queue) > 0 {
		next, queue = queue[0], queue[1:]

		if hasAllKeys(next.keys, allKeys) {
			break
		}

		for _, p := range getNeighbors(g, next.pos, next.keys) {
			newkeys := make(map[string]bool)
			cell := g[p.Y][p.X]
			copyKeys(newkeys, next.keys)
			if isKey(cell) {
				newkeys[cell] = true
			}
			vk := fk(p, newkeys)
			if !visited[vk] {
				visited[vk] = true
				queue = append(queue, queueitem{p, next.steps + 1, newkeys})
			}

		}
	}

	return next.steps
}

func copyKeys(dst, src map[string]bool) {
	for k, v := range src {
		dst[k] = v
	}
}

func hasAllKeys(keys map[string]bool, allKeys map[string]geometry.Point) bool {
	missing := make(map[string]bool)
	for k := range allKeys {
		if !keys[k] {
			missing[k] = true
		}
	}

	return len(missing) == 0
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
	robot  = "@"
	keys   = "abcdefghijklmnopqrstuvwxyz"
)

func getNeighbors(g [][]string, p geometry.Point, keys map[string]bool) []geometry.Point {
	directions := []geometry.Point{geometry.Point{X: p.X, Y: p.Y - 1}, geometry.Point{X: p.X, Y: p.Y + 1}, geometry.Point{X: p.X - 1, Y: p.Y}, geometry.Point{X: p.X + 1, Y: p.Y}}
	neighbors := []geometry.Point{}
	for _, q := range directions {
		if isValid(g, q) {
			cell := g[q.Y][q.X]
			if cell == ground || isKey(cell) || keys[strings.ToLower(cell)] || cell == robot {
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

func runPartB(cave [][]string, start geometry.Point, allkeys map[string]geometry.Point) {
	// create new map according to problem spec (start is replaced by walls with four bots in diagonal pattern)
	newcave := [][]string{}
	for _, row := range cave {
		newcaverow := []string{}
		for _, cell := range row {
			newcaverow = append(newcaverow, cell)
		}
		newcave = append(newcave, newcaverow)
	}

	newwalls := []geometry.Point{geometry.Point{X: 0, Y: -1}, geometry.Point{X: -1, Y: 0}, geometry.Point{X: 0, Y: 0}, geometry.Point{X: 1, Y: 0}, geometry.Point{X: 0, Y: 1}}
	newrobots := []geometry.Point{geometry.Point{X: -1, Y: -1}, geometry.Point{X: 1, Y: -1}, geometry.Point{X: -1, Y: 1}, geometry.Point{X: 1, Y: 1}}

	for _, dp := range newwalls {
		p := geometry.Point{X: dp.X + start.X, Y: dp.Y + start.Y}
		newcave[p.Y][p.X] = wall
	}

	robotstarts := make([]geometry.Point, len(newrobots))
	for i, dp := range newrobots {
		p := geometry.Point{X: dp.X + start.X, Y: dp.Y + start.Y}
		robotstarts[i] = p
		newcave[p.Y][p.X] = robot
	}

	type queueitem struct {
		pos   []geometry.Point
		steps int
		keys  map[string]bool
	}

	type visitedkey struct {
		pos  geometry.Point
		keys string
	}
	fk := func(pos geometry.Point, keys map[string]bool) visitedkey {
		var keyset []string
		for k := range keys {
			keyset = append(keyset, k)
		}

		sort.Strings(keyset)
		return visitedkey{pos, strings.Join(keyset, "")}
	}

	visited := make([]map[visitedkey]bool, len(robotstarts))
	for i, p := range robotstarts {
		vk := fk(p, make(map[string]bool))
		visited[i] = make(map[visitedkey]bool)
		visited[i][vk] = true
	}
	queue := []queueitem{queueitem{robotstarts, 0, make(map[string]bool)}}
	var next queueitem
	for len(queue) > 0 {
		next, queue = queue[0], queue[1:]
		// fmt.Printf("Visiting %v %v %v %v %v\n", next, newcave[next.pos[0].Y][next.pos[0].X], newcave[next.pos[1].Y][next.pos[1].X], newcave[next.pos[2].Y][next.pos[2].X], newcave[next.pos[3].Y][next.pos[3].X])

		if hasAllKeys(next.keys, allkeys) {
			break
		}

		for i, r := range next.pos {
			pos2 := make([]geometry.Point, len(next.pos))
			copy(pos2, next.pos)
			// fmt.Printf("Robot %d %v pos2 %v\n", i, r, pos2)
			for _, p := range getNeighbors(newcave, r, next.keys) {
				newkeys := make(map[string]bool)
				cell := newcave[p.Y][p.X]
				copyKeys(newkeys, next.keys)
				if isKey(cell) {
					newkeys[cell] = true
				}
				vk := fk(p, newkeys)
				if !visited[i][vk] {
					visited[i][vk] = true
					newPos2 := make([]geometry.Point, len(pos2))
					copy(newPos2, pos2)
					newPos2[i] = p
					qi := queueitem{newPos2, next.steps + 1, newkeys}
					queue = append(queue, qi)
					// fmt.Printf("add %v to queue\n", qi)
				}

			}
		}
	}

	fmt.Printf("Minimum steps to get all keys %d\n", next.steps)
}

func getCave(filepath string) ([][]string, geometry.Point, map[string]geometry.Point) {
	g := [][]string{}
	var start geometry.Point
	allKeys := make(map[string]geometry.Point)
	for y, row := range strings.Split(filepath, "\n") {
		gRow := []string{}
		for x, cell := range row {
			s := string(cell)
			if s == robot {
				start = geometry.Point{X: x, Y: y}
			} else if isKey(s) {
				allKeys[s] = geometry.Point{X: x, Y: y}
			}
			gRow = append(gRow, s)
		}
		g = append(g, gRow)
	}

	return g, start, allKeys
}

func main() {
	cave, start, keys := getCave(fileutils.ReadFile(os.Args[1]))
	runPartA(cave, start, keys)
	runPartB(cave, start, keys)
}

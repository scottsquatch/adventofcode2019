package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/scottsquatch/adventofcode2019/utils/fileutils"
	"github.com/scottsquatch/adventofcode2019/utils/geometry"
)

func runPartA(filePath string) (geometry.Point, []geometry.Point) {
	spaceMap := fileutils.ReadFile(filePath)

	asteroids := []geometry.Point{}

	for y, r := range strings.Split(spaceMap, "\n") {
		for x, c := range strings.Split(r, "") {
			if c == "#" {
				asteroids = append(asteroids, geometry.Point{X: x, Y: y})
			}
		}
	}

	counts := make(map[geometry.Point]int)
	for i, p := range asteroids {
		points := append([]geometry.Point(nil), asteroids[:i]...)
		points = append(points, asteroids[i+1:]...)
		angles := make(map[float64]bool)

		for _, q := range points {
			angle := math.Atan2(float64(q.Y-p.Y), float64(q.X-p.X))
			if !angles[angle] {
				angles[angle] = true
			}
		}

		counts[p] = len(angles)
	}

	max := struct {
		point             geometry.Point
		asteroidsDetected int
	}{
		geometry.Point{},
		0,
	}
	for p, c := range counts {
		if c > max.asteroidsDetected {
			max.asteroidsDetected = c
			max.point = p
		}
	}

	fmt.Printf("Best location: %v. Asteroids detected: %d\n", max.point, max.asteroidsDetected)

	return max.point, asteroids
}

func runPartB(station geometry.Point, asteroids []geometry.Point) {
	angleMap := make(map[float64][]geometry.Point)
	for _, p := range asteroids {
		if p.X != station.X || p.Y != station.Y {
			theta := math.Atan2(float64(p.Y-station.Y), float64(p.X-station.X))
			angleMap[theta] = append(angleMap[theta], p)
			// keep points sorted
			sort.Slice(angleMap[theta], func(i, j int) bool {
				ri := station.EuclideanDistance(angleMap[theta][i])
				rj := station.EuclideanDistance(angleMap[theta][j])
				return ri < rj
			})
		}
	}

	var vaporized geometry.Point
	numVaporized := 0
	for len(angleMap) > 0 {
		angles := make([][]float64, 4)
		// Separate into quadrants
		for angle := range angleMap {
			if angle <= 0 && angle >= -math.Pi/2 {
				angles[0] = append(angles[0], angle)
			} else if angle > 0 && angle <= math.Pi/2 {
				angles[1] = append(angles[1], angle)
			} else if angle > math.Pi/2 && angle <= math.Pi {
				angles[2] = append(angles[2], angle)
			} else {
				angles[3] = append(angles[3], angle)
			}
		}

		for _, quadrant := range angles {
			sort.Float64s(quadrant)
			for _, angle := range quadrant {
				asts := angleMap[angle]
				if len(asts) > 0 {
					vaporized = asts[0]
					numVaporized++
					angleMap[angle] = asts[1:]
					fmt.Printf("Vaporized: %v, num %d\n", vaporized, numVaporized)
				} else {
					delete(angleMap, angle)
				}

				if numVaporized == 200 {
					fmt.Printf("200th asteroid to be vaporized %v, part B answer is %d\n", vaporized, 100*vaporized.X+vaporized.Y)
					return
				}
			}
		}
	}
}

func main() {
	p, asteroids := runPartA(os.Args[1])
	runPartB(p, asteroids)
}

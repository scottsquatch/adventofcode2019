package main

import (
	"fmt"
	"strings"
)

type Point3D struct {
	X int
	Y int
	Z int
}

type Velocity struct {
	X int
	Y int
	Z int
}

type Planet struct {
	point    Point3D
	velocity Velocity
}

func (p *Planet) move() {
	p.point.X += p.velocity.X
	p.point.Y += p.velocity.Y
	p.point.Z += p.velocity.Z
}

func (p *Planet) applyGravity(others []Planet) {
	p.velocity.X += calculateGravity(*p, others, func(p Planet) int { return p.point.X })
	p.velocity.Y += calculateGravity(*p, others, func(p Planet) int { return p.point.Y })
	p.velocity.Z += calculateGravity(*p, others, func(p Planet) int { return p.point.Z })
}

func abs(n int) int {
	if n < 0 {
		return -1 * n
	}

	return n
}

func (p *Planet) pot() int {
	return abs(p.point.X) + abs(p.point.Y) + abs(p.point.Z)
}

func (p *Planet) kin() int {
	return abs(p.velocity.X) + abs(p.velocity.Y) + abs(p.velocity.Z)
}

func (p *Planet) totalEnergy() int {
	return p.pot() * p.kin()
}

func calculateGravity(this Planet, others []Planet, val func(Planet) int) int {
	thisVal := val(this)
	g := 0
	for _, p := range others {
		pVal := val(p)
		switch {
		case thisVal < pVal:
			g++
		case thisVal > pVal:
			g--
		}
	}

	return g
}

func runPartA() {
	planets := []Planet{}
	for _, line := range strings.Split(input, "\n") {
		p := Point3D{}
		fmt.Sscanf(line, "<x=%d, y=%d, z=%d>", &p.X, &p.Y, &p.Z)
		pl := Planet{p, Velocity{0, 0, 0}}
		planets = append(planets, pl)
	}

	const steps = 1000
	for i := 0; i < steps; i++ {
		for i := range planets {
			others := append([]Planet(nil), planets[:i]...)
			others = append(others, planets[i+1:]...)

			planets[i].applyGravity(others)
		}

		for i := range planets {
			others := append([]Planet(nil), planets[:i]...)
			others = append(others, planets[i+1:]...)

			planets[i].move()
		}
	}

	tot := 0
	for _, p := range planets {
		tot += p.totalEnergy()
	}

	fmt.Printf("Total Energy in system after %d steps is %d\n", steps, tot)
}

func any(a []Point3D, f func(Point3D) bool) bool {
	for _, e := range a {
		if f(e) {
			return true
		}
	}

	return false
}

func gcd(a int64, b int64) int64 {
	r := a % b
	if r == 0 {
		return b
	}

	return gcd(b, r)
}

func lcm(a int64, b int64) int64 {
	return (a * b) / gcd(a, b)
}

func runPartB() {
	planets := []Planet{}
	for _, line := range strings.Split(input, "\n") {
		p := Point3D{}
		fmt.Sscanf(line, "<x=%d, y=%d, z=%d>", &p.X, &p.Y, &p.Z)
		pl := Planet{p, Velocity{0, 0, 0}}
		planets = append(planets, pl)
	}

	steps := 0
	orig := append([]Planet(nil), planets...)
	period := Point3D{}
	for period.X == 0 || period.Y == 0 || period.Z == 0 {
		for i := range planets {
			others := append([]Planet(nil), planets[:i]...)
			others = append(others, planets[i+1:]...)

			planets[i].applyGravity(others)
		}

		for i := range planets {
			others := append([]Planet(nil), planets[:i]...)
			others = append(others, planets[i+1:]...)

			planets[i].move()
		}

		steps++

		xPeriod := true
		yPeriod := true
		zPeriod := true
		for i, p := range planets {
			if p.point.X != orig[i].point.X || p.velocity.X != orig[i].velocity.X {
				xPeriod = false
			}
			if p.point.Y != orig[i].point.Y || p.velocity.Y != orig[i].velocity.Y {
				yPeriod = false
			}
			if p.point.Z != orig[i].point.Z || p.velocity.Z != orig[i].velocity.Z {
				zPeriod = false
			}
		}

		if xPeriod && period.X == 0 {
			period.X = steps
		}
		if yPeriod && period.Y == 0 {
			period.Y = steps
		}
		if zPeriod && period.Z == 0 {
			period.Z = steps
		}
	}

	repeatSteps := lcm(int64(period.X), int64(period.Y))
	repeatSteps = lcm(repeatSteps, int64(period.Z))
	fmt.Printf("System repeats itself after %d steps\n", repeatSteps)
}

func main() {
	runPartA()
	runPartB()
}

var input = `<x=14, y=4, z=5>
<x=12, y=10, z=8>
<x=1, y=7, z=-10>
<x=16, y=-5, z=3>`

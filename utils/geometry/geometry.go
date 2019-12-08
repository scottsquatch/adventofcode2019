package geometry

import "math"

// Point represents a point in the Cartesian plane
type Point struct {
	X int
	Y int
}

// Line represents a line.
type Line struct {
	start Point
	end   Point
}

// Path represents a sequence of line segments
type Path struct {
	points []Point
}

// NewLine create a new line from the given points
func NewLine(start Point, end Point) *Line {
	return &Line{start, end}
}

// NewPath create a new Path
func NewPath() *Path {
	p := Path{}
	return &p
}

// AddPoint add point to path
func (path *Path) AddPoint(point Point) {
	path.points = append(path.points, point)
}

// GetLines get lines of a path
func (path *Path) GetLines() []Line {
	var lines []Line

	p := path.points[0]
	for _, q := range path.points[1:] {
		lines = append(lines, Line{p, q})
		p = q
	}

	return lines
}

// ManhattanDistance calculate the manhatten distance to the other point
func (p Point) ManhattanDistance(q Point) int {
	result := math.Abs(float64(p.X-q.X)) + math.Abs(float64(p.Y-q.Y))
	return int(result)
}

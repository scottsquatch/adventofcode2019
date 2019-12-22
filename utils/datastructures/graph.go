package datastructures

import (
	"fmt"

	"github.com/scottsquatch/adventofcode2019/utils/geometry"
)

// PointGraph which uses a point as a vertex key
type PointGraph struct {
	v      int
	keymap map[geometry.Point]int
	idxmap map[int]geometry.Point
	adj    [][]bool
}

// NewPointGraph creates a new point graph
func NewPointGraph(vertices []geometry.Point) *PointGraph {
	v := len(vertices)
	km := make(map[geometry.Point]int, v)
	im := make(map[int]geometry.Point, v)
	for i, v := range vertices {
		km[v] = i
		im[i] = v
	}

	adj := make([][]bool, v)
	for i := 0; i < v; i++ {
		adj[i] = make([]bool, v)
	}
	return &PointGraph{v, km, im, adj}
}

// AddEdge adds an edge from u to v to the graph
func (g *PointGraph) AddEdge(u, v geometry.Point) {
	if !g.isValidVertex(u) {
		panic(fmt.Sprintf("Vertex %v is not valid\n", u))
	} else if !g.isValidVertex(v) {
		panic(fmt.Sprintf("Vertex %v is not valid\n", v))
	}
	g.adj[g.keymap[u]][g.keymap[v]] = true
	g.adj[g.keymap[v]][g.keymap[u]] = true
}

func (g *PointGraph) isValidVertex(u geometry.Point) bool {
	_, found := g.keymap[u]
	return found
}

// Adj returns vertices adjacent to u
func (g *PointGraph) Adj(u geometry.Point) []geometry.Point {
	if !g.isValidVertex(u) {
		panic(fmt.Sprintf("Vertex %v is not valid\n", u))
	}
	var neighbors []geometry.Point

	for i, edge := range g.adj[g.keymap[u]] {
		if edge {
			neighbors = append(neighbors, g.idxmap[i])
		}
	}

	return neighbors
}

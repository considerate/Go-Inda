package main

import (
	graph "./Graph-Go"
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Grapher interface {
	NumVertices() int
	NumEdges() int
	Degree(int) int
	DoNeighbors(int, func(int, interface{}))
	HasEdge(int, int) bool
	Label(int, int) interface{}
	Add(int, int)
	AddLabel(int, int, interface{})
	AddBi(int, int)
	AddBiLabel(int, int, interface{})
	Remove(int, int)
	RemoveBi(int, int)
}

var NewFuncs = map[string]func(int) Grapher{
	"Hash":   func(n int) Grapher { return graph.NewHash(n) },
	"Matrix": func(n int) Grapher { return graph.NewMatrix(n) },
}

// Constructs test graphs using the factory method f.
func setup(f func(int) Grapher, n int) (graph Grapher) {
	graph = f(n)
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	for count := 0; count < n; {
		a := random.Intn(n)
		b := random.Intn(n)
		if !graph.HasEdge(a, b) {
			graph.Add(a, b)
			count++
		}
	}
	return
}

func getDFSInfo(g Grapher) (maxSize, numComponents int64) {
	var currentSize int64
	maxSize = math.MinInt64
	numComponents = 0
	state := make([]bool, g.NumVertices())
	for v, visited := range state {
		if !visited {
			numComponents++
			currentSize = 0
			graph.DFS(g, v, state, func(w int) {
				currentSize++
				if currentSize > maxSize {
					maxSize = currentSize
				}
			})
		}
	}
	return
}

func doDFSWork(g Grapher) {
	dummy := 0
	state := make([]bool, g.NumVertices())
	for v, visited := range state {
		if !visited {
			graph.DFS(g, v, state, func(w int) {
				dummy++
			})
		}
	}
}

func main() {
	sizes := []int{1, 10, 50, 100, 180, 250, 300, 1000, 5000}
	for _, size := range sizes {
		for graphName, f := range NewFuncs {
			g := setup(f, size)
			before := time.Now()
			maxSize, numComponents := getDFSInfo(g)
			for i := 0; i < 100; i++ {
				doDFSWork(g)
			}
			after := time.Now()
			timeUsed := after.Sub(before)
			fmt.Println(graphName, "Size", size, "Components", numComponents, "Max", maxSize, "Time", timeUsed)
		}
	}
}

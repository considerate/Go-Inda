package graph

const initialMapSize = 4

type Hash struct {
	// The map edges[v] contains the mapping {w:x} if there is an edge
	// from v to w; x is the label assigned to this edge.
	// The maps may be nil and are allocated only when needed.
	edges []map[int]interface{}

	numEdges int // total number of directed edges in the graph
}

// NewList constructs a new graph with n vertices and no edges.
func NewHash(n int) *Hash {
	return &Hash{edges: make([]map[int]interface{}, n)}
}

// NumVertices returns the number of vertices in this graph.
// Time complexity: O(1).
func (g *Hash) NumVertices() int {
	return len(g.edges)
}

// NumEdges returns the number of (directed) edges in this graph.
// Time complexity: O(1).
func (g *Hash) NumEdges() int {
	return g.numEdges
}

// Degree returns the degree of vertex v. Time complexity: O(1).
// The degreee returned is the outdegree.
func (g *Hash) Degree(v int) int {
	return len(g.edges[v])
}

// DoNeighbors calls action for each neighbor w of v,
// with x equal to the label of the edge from v to w.
// Time complexity: O(m), where m is the number of neighbors.
func (g *Hash) DoNeighbors(v int, action func(w int, x interface{})) {
	neighbors := g.edges[v]
	for w, label := range neighbors {
		action(w, label)
	}
}

// HasEdge returns true if there is an edge from v to w.
// Time complexity: O(1).
func (g *Hash) HasEdge(v, w int) bool {
	for vert, _ := range g.edges[v] {
		if vert == w {
			return true
		}
	}
	return false
}

// Returns the label for the edge from v to w, NoLabel if the edge has no label,
// or nil if no such edge exists.
// Time complexity: O(1).
func (g *Hash) Label(v, w int) interface{} {
	m := g.edges[v]
	if label, ok := m[w]; ok {
		return label
	}
	return nil
}

// Add inserts a directed edge.
// It removes any previous label if this edge already exists.
// Time complexity: O(1).
func (g *Hash) Add(v, w int) {
	if _, ok := g.edges[v][w]; !ok {
		g.numEdges += 1
	}
	m := g.edges[v]
	if m == nil {
		m = make(map[int]interface{}, initialMapSize)
		g.edges[v] = m
	}
	g.edges[v][w] = NoLabel
}

// AddBi inserts edges between v and w.
// It removes any previous labels if these edges already exists.
// Time complexity: O(1).
func (g *Hash) AddBi(v, w int) {
	g.Add(v, w)
	g.Add(w, v)
}

// AddLabel inserts a directed edge with label x.
// It overwrites any previous label if this edge already exists.
// Time complexity: O(1).
func (g *Hash) AddLabel(v, w int, x interface{}) {
	m := g.edges[v]
	if m == nil {
		m = make(map[int]interface{}, initialMapSize)
		g.edges[v] = m
	}
	if _, ok := m[w]; !ok {
		g.numEdges += 1
	}
	m[w] = x
}

// AddBiLabel inserts edges with label x between v and w.
// It overwrites any previous labels if these edges already exists.
// Time complexity: O(1).
func (g *Hash) AddBiLabel(v, w int, x interface{}) {
	g.AddLabel(v, w, x)
	g.AddLabel(w, v, x)
}

// Remove removes an edge. Time complexity: O(1).
func (g *Hash) Remove(v, w int) {
	if _, ok := g.edges[v][w]; ok {
		g.numEdges -= 1
	}
	delete(g.edges[v], w)
}

// RemoveBi removes all edges between v and w. Time complexity: O(1).
func (g *Hash) RemoveBi(v, w int) {
	g.Remove(v, w)
	g.Remove(w, v)
}

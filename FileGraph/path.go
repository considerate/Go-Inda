package main

import (
	graph "../Graph/Graph-Go"
	"bufio"
	"fmt"
	//"io"
	//stack "./stack"
	"container/heap"
	"flag"
	"math"
	"os"
	"strconv"
	"strings"
)

type Item struct {
	value    int
	priority int
	index    int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) find(value int) *Item {
	for _, item := range *pq {
		if item.value == value {
			return item
		}
	}
	return nil
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *Item, priority int) {
	heap.Remove(pq, item.index)
	item.priority = priority
	heap.Push(pq, item)
}

var filename string
var from int
var to int

func init() {
	flag.StringVar(&filename, "filename", "graph.txt", "Specify a file name")
	flag.Parse()
}

func strToInt(str string) (int, error) {
	tmp := strings.Trim(str, " \n")
	num, err := strconv.Atoi(tmp)
	return num, err
}

func nextNumber(reader *bufio.Reader, delim byte) (int, error) {
	str, err := reader.ReadString(delim)
	if err != nil {
		return 0, err
	}

	num, err := strToInt(str)
	if err != nil {
		return 0, err
	}

	return num, nil

}

func reverse(slice []int) {
	l := len(slice)
	for i := 0; i < l/2; i++ {
		slice[i], slice[l-i-1] = slice[l-i-1], slice[i]
	}
}

func printPath(g *graph.Hash, from, to int, stack []int) {
	var path []int
	totalcost := 0
	found := false
	index := to
	path = append(path, index)
	for stack[index] != -1 {
		prev := stack[index]
		cost := g.Label(prev, index).(int)
		totalcost += cost
		index = prev
		path = append(path, index)

		if index == from {
			found = true
			break
		}
	}
	if found {
		reverse(path)
		fmt.Printf("Path from %d to %d: %v cost(%d)\n", from, to, path, totalcost)
	} else {
		fmt.Printf("Dijkstra... No path found between %d and %d\n", from, to)
	}
}

func main() {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Failed to open file", filename)
		os.Exit(1)
	}
	reader := bufio.NewReader(file)

	line, err := reader.ReadString('\n')
	if err != nil {
		panic("Failed to read first line")
	}
	size, err := strToInt(line)
	if err != nil {
		panic("Failed to convert to int")
	}
	g := graph.NewHash(size)

	reader.ReadString('\n') //Skip one line
	var v, w, cost int
	for {
		v, err = nextNumber(reader, ' ')
		if err != nil {
			break
		}

		w, err = nextNumber(reader, ' ')
		if err != nil {
			break
		}

		cost, err = nextNumber(reader, '\n')
		if err != nil {
			break
		}

		g.AddLabel(v, w, cost)
	}

	from, _ = strToInt(flag.Arg(0))
	to, _ = strToInt(flag.Arg(1))
	visited := make([]bool, g.NumVertices())
	stack := make([]int, g.NumEdges())
	for i, _ := range stack {
		stack[i] = -1
	}
	graph.BFS(g, from, visited, func(parent, w int) {
		stack[w] = parent
	})
	printPath(g, from, to, stack)

	distance := make([]int, g.NumVertices())
	pq := &PriorityQueue{}
	heap.Init(pq)
	for i, _ := range distance {
		if i == from {
			distance[i] = 0
		} else {
			distance[i] = math.MaxInt32
		}
		item := &Item{
			value:    i,
			priority: distance[i],
		}
		heap.Push(pq, item)
	}

	stack = make([]int, g.NumVertices())
	for i, _ := range stack {
		stack[i] = -1
	}

	for pq.Len() > 0 {
		current := heap.Pop(pq).(*Item)
		v := current.value
		g.DoNeighbors(v, func(w int, x interface{}) {
			cost := x.(int)
			dist := distance[v] + cost
			if dist < distance[w] {
				distance[w] = dist
				stack[w] = v
				item := pq.find(w)
				pq.update(item, dist)
			}
		})
	}
	printPath(g, from, to, stack)
}

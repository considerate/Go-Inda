package main

import (
	graph "../Graph/Graph-Go"
	"bufio"
	"fmt"
	//"io"
	//stack "./stack"
	"flag"
	"os"
	"strconv"
	"strings"
)

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

	found := false
	visited := make([]bool, g.NumVertices())
	stack := make([]int, g.NumEdges())
	previous := -1
	graph.BFS(g, from, visited, func(w int) {
		if !found {
			stack[w] = previous
			previous = w
		}
		if w == to {
			found = true
		}
	})
	if found {
		index := to
		path := make([]int, g.NumEdges())
		count := 0
		for index != -1 {
			path[count] = index
			index = stack[index]
			count++
		}
		path = path[0:count]
		reverse(path)
		fmt.Printf("Path from %d to %d: %v\n", from, to, path)
	} else {
		fmt.Printf("No path found between %d and %d\n", from, to)
	}
}

package main

import (
	"fmt"
	"sync"
)

const NUM_SECTIONS = 5

// Add adds the numbers in a and sends the result on res.
func Add(a []int, res chan<- int) {
	sum := 0
	for _, value := range a {
		sum += value
	}
	res <- sum
}

func main() {
	a := []int{1, 2, 3, 4, 5, 6, 7, 82, 2, 3, 5, 1}

	n := len(a)
	ch := make(chan int)
	wg := new(sync.WaitGroup)

	wg.Add(NUM_SECTIONS)
	var head, tail int
	for i := 0; i < NUM_SECTIONS; i++ {
		head = i * n / NUM_SECTIONS
		tail = (i + 1) * n / NUM_SECTIONS
		go Add(a[head:tail], ch)
	}

	count := 0
	total := 0
	go func() {
		for val := range ch {
			count += 1
			fmt.Printf("%d", val)
			total += val
			if count != NUM_SECTIONS {
				fmt.Print(" + ")
			}
			wg.Done()
		}
	}()
	wg.Wait()
	close(ch)
	fmt.Printf(" = %d", total)
}

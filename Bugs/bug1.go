package main

import (
	"fmt"
)

// I want this program to print "Hello world!", but it doesn't work.
// *Ahem* I mean: It works :P
func main() {
	ch := make(chan string)
	go func() {
		ch <- "Hello world!"
	}()
	fmt.Println(<-ch)
}

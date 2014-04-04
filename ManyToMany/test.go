package main

import "fmt"
import "runtime"

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	ch := make(chan int)
	go func() {
		ch <- 4
	}()
	a := <-ch
	fmt.Println(a)
	close(ch)
}

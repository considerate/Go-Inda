package main

import (
	"fmt"
	"math"
)

func Sqrt(x float64) float64 {
	z := 1.5
	var last, next float64
	for {
		last = z
		z = z - (z*z-x)/(2*z)
		next = math.Nextafter(z, last)
		//Delta is one bit
		if next == last {
			break
		}
	}
	return z
}

func main() {
	fmt.Println(Sqrt(2))
	fmt.Println(math.Sqrt(2))
}

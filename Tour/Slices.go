package main

import "code.google.com/p/go-tour/pic"

func Pic(dx, dy int) [][]uint8 {
	rows := make([][]uint8, dy)
	for y, _ := range rows {
		rows[y] = make([]uint8, dx)
	}

	for y, row := range rows {
		for x, _ := range row {
			row[x] = uint8(x*y - x + 42 + x*x)
		}
	}
	return rows
}

func main() {
	pic.Show(Pic)
}

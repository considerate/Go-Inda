// Stefan Nilsson 2013-02-27

// This program creates pictures of Julia sets (en.wikipedia.org/wiki/Julia_set).
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
	"runtime"
	"strconv"
	"sync"
)

type ComplexFunc func(complex128) complex128

var Funcs []ComplexFunc = []ComplexFunc{
	func(z complex128) complex128 { return z*z - 0.61803398875 },
	func(z complex128) complex128 { return z*z + complex(0, 1) },
	func(z complex128) complex128 { return z*z + complex(-0.835, -0.2321) },
	func(z complex128) complex128 { return z*z + complex(0.45, 0.1428) },
	func(z complex128) complex128 { return z*z*z + 0.400 },
	func(z complex128) complex128 { return cmplx.Exp(z*z*z) - 0.621 },
	func(z complex128) complex128 { return (z*z+z)/cmplx.Log(z) + complex(0.268, 0.060) },
	func(z complex128) complex128 { return cmplx.Sqrt(cmplx.Sinh(z*z)) + complex(0.065, 0.122) },
}

type pixel struct {
	x     int
	y     int
	value int
}

var count int

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	count = 0
	errors := make(chan error)
	images := make(chan image.Image)
	var wg sync.WaitGroup
	fmt.Println("Entered main")
	go CreatePng(images, errors, &wg)
	for _, f := range Funcs {
		wg.Add(1)
		img := JuliaRow(f, 16384)
		go func() {
			images <- img
		}()
	}
	wg.Wait()
	close(images)
	close(errors)
	fmt.Println("Done")
}

// CreatePng creates a PNG picture file with a Julia image of size n x n.
func CreatePng(images <-chan image.Image, errors chan error, wg *sync.WaitGroup) {
	for image := range images {
		count++
		filename := "picture-" + strconv.Itoa(count) + ".png"
		fmt.Printf("Writing %s\n", filename)
		file, err := os.Create(filename)
		if err != nil {
			errors <- err
			return
		}
		defer file.Close()

		err = png.Encode(file, image)
		if err != nil {
			errors <- err
		}

		wg.Done()
		fmt.Printf("%s written\n", filename)
	}
}

func JuliaRow(f ComplexFunc, n int) image.Image {
	bounds := image.Rect(-n/2, -n/2, n/2, n/2)
	img := image.NewRGBA(bounds)
	s := float64(n / 4)
	size := n * n
	positions := make(chan struct{ i, j int })
	var wg sync.WaitGroup
	wg.Add(size)

	for i := bounds.Min.X; i < bounds.Max.X; i++ {
		go func(i int) {
			for j := bounds.Min.Y; j < bounds.Max.Y; j++ {
				fillPixel(f, s, i, j, img, &wg)
			}
		}(i)
	}
	close(positions)
	wg.Wait()
	return img
}

// Julia returns an image of size n x n of the Julia set for f.
func JuliaFive(f ComplexFunc, n int) image.Image {
	bounds := image.Rect(-n/2, -n/2, n/2, n/2)
	img := image.NewRGBA(bounds)
	s := float64(n / 4)
	size := n * n
	positions := make(chan struct{ i, j int })
	var wg sync.WaitGroup
	wg.Add(size)

	go fillPixels(f, s, positions, img, &wg)
	go fillPixels(f, s, positions, img, &wg)
	go fillPixels(f, s, positions, img, &wg)
	go fillPixels(f, s, positions, img, &wg)
	go fillPixels(f, s, positions, img, &wg)

	for i := bounds.Min.X; i < bounds.Max.X; i++ {
		for j := bounds.Min.Y; j < bounds.Max.Y; j++ {
			positions <- struct{ i, j int }{i, j}
		}
	}
	close(positions)
	wg.Wait()
	return img
}

func fillPixel(f ComplexFunc, s float64, i int, j int, img *image.RGBA, wg *sync.WaitGroup) {
	value := JuliaPixel(f, i, j, s)
	r := uint8((value * 8) % 256)
	g := uint8(0)
	b := uint8(value % 32 * 8)
	img.Set(i, j, color.RGBA{r, g, b, 255})
	wg.Done()
}

func fillPixels(f ComplexFunc, s float64, positions <-chan struct{ i, j int }, img *image.RGBA, wg *sync.WaitGroup) {
	for pos := range positions {
		i := pos.i
		j := pos.j
		value := JuliaPixel(f, i, j, s)
		r := uint8((value * 8) % 256)
		g := uint8(0)
		b := uint8(value % 32 * 8)
		img.Set(i, j, color.RGBA{r, g, b, 255})
		wg.Done()
	}
}

func JuliaPixel(f ComplexFunc, i int, j int, s float64) int {
	value := Iterate(f, complex(float64(i)/s, float64(j)/s), 256)
	return value
}

// Iterate sets z_0 = z, and repeatedly computes z_n = f(z_{n-1}), n ≥ 1,
// until |z_n| > 2  or n = max and returns this n.
func Iterate(f ComplexFunc, z complex128, max int) (n int) {
	for ; n < max; n++ {
		if real(z)*real(z)+imag(z)*imag(z) > 4 {
			break
		}
		z = f(z)
	}
	return
}

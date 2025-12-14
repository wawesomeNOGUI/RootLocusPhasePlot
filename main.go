package main

// #cgo  LDFLAGS: -lgdi32
/*
#include "main_c.c"
*/
import "C"

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"strconv"
	"time"
)

var poles []complex128
var zeros []complex128
var numCPU int = runtime.NumCPU()

func getInputLoopTransferFunction(poles, zeros []complex128) {
	// get input poles and zeroes of the loop transfer function
	if len(os.Args) <= 1 {
		fmt.Println("ERROR: input args in format e.g.: p 0.5 1+2i (1-2i) z 1i 2 -5")
		os.Exit(1)
	}

	// set variable for if subsequent numbers are poles or zeros
	p := true

	for i, v := range os.Args[1:] {
		if i == 0 {
			if v != "p" && v != "z" {
				fmt.Println("ERROR: input args in format e.g.: p 0.5 1+2i (1-2i) z 1i 2 -5")
				os.Exit(1)
			}
		}

		if v == "p" {
			p = true
			continue
		} else if v == "z" {
			p = false
			continue
		}

		cmplx, err := strconv.ParseComplex(v, 128)
		if err != nil {
			fmt.Println(err)
			fmt.Println("ERROR: input args in format e.g.: p 0.5 1+2i (1-2i) z 1i 2 -5")
			os.Exit(1)
		}

		// else must have gotten a correct number format so store in respective slices
		if p {
			poles = append(poles, cmplx)
		} else {
			zeros = append(zeros, cmplx)
		}
	}
}

func drawPhasePlot(poles, zeros []complex128) {
	for j := 0; j < 200; j++ {
		for i := float64(0); i < 1; i += 0.001 {
			p := 2 * math.Pi * i
			C.DrawPixel(C.int(500/2+float64(j)*math.Cos(p)), C.int(500/2+float64(j)*-math.Sin(p)), C.int(colorFromPhase(p)))

			// fmt.Println(C.int(500/2 + 10*math.Cos(p)))
		}
	}
}

// p: input phase in radians from 0 to 2pi
// returns an integer from 0 to 0xFFFFFF (RGB value)
func colorFromPhase(p float64) int {
	// color wheel associated with phase
	// copying the style from Dr. Adams Controls II notes
	// 0 red, pi/2 light green, pi cyan, 3pi/2 dark blue

	// p = math.Mod(p, 2*math.Pi)
	var r, g, b int

	// 7 intervals
	// https://www.rapidtables.com/convert/color/hsv-to-rgb.html
	switch {
	case p <= 2*math.Pi/6:
		g = int(p / (2 * math.Pi / 6) * 255)
		r = 255
	case p <= 2*2*math.Pi/6:
		r = 255 - int((p-(2*math.Pi/6))/(2*math.Pi/6)*255)
		g = 255
	case p <= 3*2*math.Pi/6:
		g = 255
		b = int((p - (2 * 2 * math.Pi / 6)) / (2 * math.Pi / 6) * 255)
	case p <= 4*2*math.Pi/6:
		g = 255 - int((p-(3*2*math.Pi/6))/(2*math.Pi/6)*255)
		b = 255
	case p <= 5*2*math.Pi/6:
		r = int((p - (4 * 2 * math.Pi / 6)) / (2 * math.Pi / 6) * 255)
		b = 255
	case p <= 6*2*math.Pi/6:
		r = 255
		b = 255 - int((p-(5*2*math.Pi/6))/(2*math.Pi/6)*255)
	}

	return (r << 16) | (g << 8) | b
}

func main() {
	go func() {
		for {
			time.Sleep(time.Millisecond)
			if C.readyToDraw {
				break
			}
		}

		fmt.Println("Number of CPU Logical Processors: ", numCPU)

		fmt.Printf("%X\n", C.int(colorFromPhase(0.5*math.Pi/3)))
		fmt.Printf("%X\n", colorFromPhase(2*math.Pi/3))
		fmt.Printf("%X\n", colorFromPhase(3*math.Pi/3))
		fmt.Printf("%X\n", colorFromPhase(4*math.Pi/3))
		fmt.Printf("%X\n", colorFromPhase(5*math.Pi/3))
		fmt.Printf("%X\n", colorFromPhase(6*math.Pi/3))

		// user should input poles and zeros as command arguments when running program
		getInputLoopTransferFunction(poles, zeros)

		fmt.Println("Poles: ", poles)
		fmt.Println("Zeros: ", zeros)

		// for {
		// C.DrawPixel(5, 5, 0xFF00FF)
		drawPhasePlot(poles, zeros)
		// for i := 0; i < 500; i++ {
		// 	for j := 0; j < 500; j++ {
		// 		C.DrawPixel(C.int(i), C.int(j), C.int(colorFromPhase(float64(i)/500*2*math.Pi)))
		// 	}
		// }

		C.BitBltToWindowDC()

		// fmt.Println("test")

		time.Sleep(time.Millisecond * 250)
		// }

	}()

	C.init()
}

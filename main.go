package main

// #cgo  LDFLAGS: -lgdi32
/*
#include "main_c.c"
*/
import "C"

import (
	"fmt"
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
	for i := C.int(0); i < 500; i++ {
		C.DrawPixel(i, i, i*5)
	}
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

		// user should input poles and zeros as command arguments when running program
		getInputLoopTransferFunction(poles, zeros)

		fmt.Println("Poles: ", poles)
		fmt.Println("Zeros: ", zeros)

		for {
			// C.DrawPixel(5, 5, 0xFF00FF)
			drawPhasePlot(poles, zeros)
			C.BitBltToWindowDC()

			fmt.Println("test")

			time.Sleep(time.Millisecond * 250)
		}

	}()

	C.init()
}

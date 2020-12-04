package main

import "math"

// ...
//var C = [8]float64{-0.6, -0.535, -0.33, -0.185, -0.095, -0.041, -0.0142, -0.0037}
var cs = [8]float32{0.8574929, 0.881742, 0.94962865, 0.9833146, 0.9955178, 0.9991606, 0.9998992, 0.99999315}
var ca = [8]float32{-0.5144958, -0.47173202, -0.31337747, -0.18191321, -0.09457418, -0.040965583, -0.014198569, -0.0036999746}

var imdctWinData = [4][36]float32{}

func init() {
	// Block type 0
	for i := 0; i < 36; i++ {
		imdctWinData[0][i] = float32(math.Sin(math.Pi / 36. * (float64(i) + .5)))
	}
	// Block type 1
	for i := 0; i < 18; i++ {
		imdctWinData[1][i] = float32(math.Sin(math.Pi / 36. * (float64(i) + .5)))
	}
	for i := 18; i < 24; i++ {
		imdctWinData[1][i] = 1.
	}
	for i := 24; i < 30; i++ {
		imdctWinData[1][i] = float32(math.Sin(math.Pi / 12. * (float64(i) - 18. + .5)))
	}
	for i := 30; i < 36; i++ {
		imdctWinData[1][i] = 0.
	}
	// Block type 2
	for i := 0; i < 12; i++ {
		imdctWinData[2][i] = float32(math.Sin(math.Pi / 12. * (float64(i) + .5)))
	}
	for i := 12; i < 36; i++ {
		imdctWinData[2][i] = 0.
	}
	// Block type 3
	for i := 0; i < 6; i++ {
		imdctWinData[3][i] = 0.
	}
	for i := 6; i < 12; i++ {
		imdctWinData[3][i] = float32(math.Sin(math.Pi / 12. * (float64(i) - 6. + .5)))
	}
	for i := 12; i < 18; i++ {
		imdctWinData[3][i] = 1.
	}
	for i := 18; i < 36; i++ {
		imdctWinData[3][i] = float32(math.Sin(math.Pi / 36. * (float64(i) + .5)))
	}
}

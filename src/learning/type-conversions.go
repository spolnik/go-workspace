package main

import (
	"fmt"
	"math"
)

func main() {
	x := 3
	y := 4

	var f float64 = math.Sqrt(float64(x*x + y*y))
	z := uint(f)
	fmt.Println(x, y, z)
}

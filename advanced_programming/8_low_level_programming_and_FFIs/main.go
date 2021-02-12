package main

/*
#include <math.h>
*/
import "C"
import "fmt"

func SquareRoot(num float64) float64{
	return float64(C.sqrt(C.double(num)))
}

func main() {
	fmt.Println(SquareRoot(16))
}
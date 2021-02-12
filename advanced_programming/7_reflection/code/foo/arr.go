package main

import (
	"fmt"
	"unsafe"
)

func main() {
	a := [3]int{1, 2, 3}

	fmt.Println(unsafe.Pointer(&a))
	fmt.Println(unsafe.Pointer(&a[0]))
}

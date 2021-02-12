package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	var f *os.File
	fmt.Println(f == nil)

	var w io.Writer = f
	fmt.Println(w == nil)
}

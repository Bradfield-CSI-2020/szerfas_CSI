package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func main() {

	var w io.Writer

	b := new(bytes.Buffer) // *bytes.Buffer
	f := os.Stdout         // *os.File

	w = b
	w.Write([]byte("hello"))
	fmt.Println(b.String())

	w = f
	w.Write([]byte("hello"))
}

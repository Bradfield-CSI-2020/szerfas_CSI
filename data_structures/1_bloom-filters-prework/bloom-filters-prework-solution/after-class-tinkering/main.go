package main

import (
	"fmt"
	"unsafe"
)

func main() {
	item := "hello, world!"

	voidPointer := unsafe.Pointer(&item)
	stringPointer := (*string)(voidPointer)
	byteSlicePointer := (*[]byte)(voidPointer)
	byteSlice := *byteSlicePointer

	// Obviously, just the string "hello, world!"
	fmt.Println(item)

	// These three pointers all point to the same location
	fmt.Printf("%p\n", voidPointer)
	fmt.Printf("%p\n", stringPointer)
	fmt.Printf("%p\n", byteSlicePointer)

	// Bytes for "hello, world!" (in ASCII)
	fmt.Println(byteSlice)

	// A string has a ptr and a len, but not a cap (unlike a []byte)
	// Thus trying to do `cap` will give us nonsense
	fmt.Println(len(byteSlice))
	fmt.Println(cap(byteSlice))

	// This points to a DIFFERENT location than the three pointers above!
	// It points to the location of the [ptr | len | cap] data used by the byteSlice variable,
	// and NOT the location of the underlying "hello, world" bytes!
	anotherByteSlicePointer := unsafe.Pointer(&byteSlice)
	fmt.Printf("%p\n", anotherByteSlicePointer)
}

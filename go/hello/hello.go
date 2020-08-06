package main

import (
	"fmt"
	"example/stephenzerfas/hello/morestrings"
	"github.com/google/go-cmp/cmp"
)

func main() {
	fmt.Println("Hello, world.")
	fmt.Println(cmp.Diff("Hello World", "Hello Go"))
	fmt.Println(morestrings.ReverseRunes("!oG, olleH"))
}

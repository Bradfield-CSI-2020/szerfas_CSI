package main

import (
	"golang.org/x/tour/tree"
	"fmt"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	// look left
	if t.Left != nil {
		// if node, move to that node and repeat
		Walk(t.Left, ch)
	}
	// add value to the count and return to previous node
	ch <- t.Value

	// look right
	if t.Right != nil {
		// if node, move to that node and repeat
		Walk(t.Right, ch)
	}

}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	go Walk(t1, ch1)
	ch2 := make(chan int)
	go Walk(t2, ch2)
	acc := true
	for i := 0; i < 10; i++ {
		acc = acc && <-ch1 == <-ch2
	}
	return acc
}

func main() {
	fmt.Println(Same(tree.New(1), tree.New(1)))
	fmt.Println(Same(tree.New(1), tree.New(2)))
}

package main

import "fmt"

func mapIntToInt(f func(int) int, arr []int) []int {
	n := len(arr)
	result := make([]int, n)
	for i := 0; i < n; i++ {
		result[i] = f(arr[i])
	}
	return result
}

/*
	Every item in the slice could have a different type,
	some of them could be some unknown / complicated structs
	So f might need to deal with that
*/
func myMap(f func(interface{}) interface{}, arr []interface{}) []interface{} {
	n := len(arr)
	result := make([]interface{}, n)
	for i := 0; i < n; i++ {
		result[i] = f(arr[i])
	}
	return result
}

func Square2(a interface{}) interface{} {
	// Does this automatically coercion to interface{}
	return a.(int) * a.(int)
}

func main() {
	arr := []int{1, 2, 3}
	/*
		result := myMap(Square, arr)
		fmt.Println(result)
	*/

	n := len(arr)
	// Can we assign a []int into an []interface{} variable?
	var x = make([]interface{}, n)
	for i := 0; i < n; i++ {
		x[i] = arr[i]
	}

	x[0] = "hello"

	// fmt.Println(x)
	fmt.Println(myMap(Square2, x))
}

package main

import (
	"fmt"
	"reflect"
)

// WARNING: DO NOT DO THIS IN PRACTICE
func ReflectMap(f interface{}, array interface{}) interface{} {
	fValue := reflect.ValueOf(f)
	fType := fValue.Type()
	// TODO: Ensure that fType is a function with 1 parameter / return value

	aValue := reflect.ValueOf(array)
	// TODO: Ensure that array is either a slice or an array

	rType := fType.Out(0)

	n := aValue.Len()
	result := reflect.MakeSlice(reflect.SliceOf(rType), n, n)
	for i := 0; i < n; i++ {
		params := []reflect.Value{aValue.Index(i)}
		returns := fValue.Call(params)
		result.Index(i).Set(returns[0])
	}

	return result.Interface()
}

func Square(x int) int {
	return x * x
}

func main() {
	array := []int{1, 2, 3}
	result := ReflectMap(Square, array)
	fmt.Println(result)
}

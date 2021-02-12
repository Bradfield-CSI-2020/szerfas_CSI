package main

import (
	"fmt"
	"reflect"
)

func main() {
	//x := 3
	//kind := reflect.ValueOf(x).Kind()
	x := 2
	d := reflect.ValueOf(&x).Elem()
	t := d.Type()
	fmt.Println(t)
	//px := d.Addr().Interface().(*int)
	px := d.Addr().Interface().(*t)
	*px = 3
	fmt.Println(x)
}
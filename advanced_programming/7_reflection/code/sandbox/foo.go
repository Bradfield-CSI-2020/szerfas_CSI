package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	var x interface{} = "hello"
	result, _ := json.Marshal(x)
	fmt.Println(result)
}

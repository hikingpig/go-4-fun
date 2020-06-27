package main

import (
	"fmt"
	"reflect"
	"strings"
)

func example5() {
	got := strings.Split("a:b:c", ":")
	want := []string{"a", "b", "c"}

	fmt.Println("==== got == want?", reflect.DeepEqual(got, want))

	var a, b []string = nil, []string{}
	fmt.Println("nil slice == empty slice? ", reflect.DeepEqual(a, b))

	var c, d map[string]int = nil, make(map[string]int)
	fmt.Println("nil map == empty map? ", reflect.DeepEqual(c, d))
}

// func main() {
// 	example5()
// }

package main

import (
	"fmt"
	"reflect"
)

func example3() {
	v := reflect.ValueOf(3)
	fmt.Println(v)
	fmt.Printf("%v\n", v)
	fmt.Println(v.String())
	// the String method only return the type!

	// get the reflect.Type from the value
	t := v.Type()
	fmt.Println(t.String())

	// get the interface from the value
	x := v.Interface()
	// the underlying concrete type is hidden in empty interface
	// we must know the underlying type and do type assertion
	// to make use the methods of the concrete type
	i := x.(int)
	fmt.Printf("%d\n", i)
}

func main() {
	example3()
}

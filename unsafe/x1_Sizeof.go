package main

import (
	"fmt"
	"unsafe"
)

/*
- size for string, array, map, func, channel, interface are fixed
regardless of whatever data stored under them
*/

func example1() {
	fmt.Println("=== float64 size:", unsafe.Sizeof(float64(0)))
	fmt.Println("=== bool size: ", unsafe.Sizeof(true))
	fmt.Println("=== int8 size: ", unsafe.Sizeof(int8(1)))
	fmt.Println("=== int32 size: ", unsafe.Sizeof(int32(2)))
	fmt.Println("=== unit32 size: ", unsafe.Sizeof(uint32(1)))
	fmt.Println("=== int size: ", unsafe.Sizeof(int(1)))
	fmt.Println("=== uintptr size: ", unsafe.Sizeof(uintptr(1)))
	x := int(8)
	fmt.Println("=== *T size: ", unsafe.Sizeof(&x))
	y := "hellosafdasdasdgsdagdsagdagsdgg"
	fmt.Println("=== string size: ", unsafe.Sizeof(y))
	z := []int{1, 2, 3, 4}
	fmt.Println("=== size of []T: ", unsafe.Sizeof(z))
	a := map[string]string{"a": "b"}
	fmt.Println("=== size of map: ", unsafe.Sizeof(a))
	b := func() {}
	fmt.Println("=== size of func: ", unsafe.Sizeof(b))
	c := make(chan int)
	fmt.Println("=== size of channel: ", unsafe.Sizeof(c))
	var d interface{}
	fmt.Println("=== size of interface: ", unsafe.Sizeof(d))
}

// func main() {
// 	example1()
// }

package main

import (
	"fmt"
	"unsafe"
)

/*
- not sure the meaning of Alignof. what its number is used for?
- Offsetof return how far the field from the beginning of struct. it counts for holes
*/

func example2() {
	var x struct {
		a bool
		b int16
		c []int
	}
	fmt.Printf("=== sizeOf(x) = %v; Alignof(x)= %v\n", unsafe.Sizeof(x), unsafe.Alignof(x))
	fmt.Printf("=== sizeOf(x.a) = %v; Alignof(x.a)= %v; Offsetof(x.a)=%v\n", unsafe.Sizeof(x.a), unsafe.Alignof(x.a), unsafe.Offsetof(x.a))
	fmt.Printf("=== sizeOf(x.b) = %v; Alignof(x.b)= %v; Offsetof(x.b)=%v\n", unsafe.Sizeof(x.b), unsafe.Alignof(x.b), unsafe.Offsetof(x.b))
	fmt.Printf("=== sizeOf(x.c) = %v; Alignof(x.c)= %v; Offsetof(x.c)=%v\n", unsafe.Sizeof(x.c), unsafe.Alignof(x.c), unsafe.Offsetof(x.c))
}

// func main() {
// 	example2()
// }

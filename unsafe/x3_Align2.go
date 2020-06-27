package main

import (
	"fmt"
	"unsafe"
)

/*
- the first struct arrangement has more holes than the 2nd and 3rd
- due to alignment
*/
func example3() {
	var x struct {
		a bool
		b float64
		c int16
	}
	fmt.Println("===== size of x", unsafe.Sizeof(x))
	fmt.Println("===== Alignof x", unsafe.Alignof(x))
	var y struct {
		a float64
		b int16
		c bool
	}
	fmt.Println("===== size of y", unsafe.Sizeof(y))
	fmt.Println("===== AlignOf y", unsafe.Alignof(y))
	var z struct {
		a bool
		b int16
		c float64
	}
	fmt.Println("===== size of z", unsafe.Sizeof(z))
	fmt.Println("===== Alignof z", unsafe.Alignof(z))
}

// func main() {
// 	example3()
// }

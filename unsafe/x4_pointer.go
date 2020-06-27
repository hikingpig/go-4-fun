package main

import (
	"fmt"
	"unsafe"
)

/*
- a pointer is converted to unsafe pointer with unsafe.Pointer(ptr)
- an unsafe.Pointer is converted back to ordinary pointer with *uint64(unsafePtr)
- we can access the converted pointer's variable with *convertedPtr
- an unsafe ptr can hold any type and can be converted to any type of ordinary pointers
*/

func Float64bits(f float64) uint64 {
	return *(*uint64)(unsafe.Pointer(&f))
}

func example4() {
	fmt.Printf("==== uint64 of float64: %#16x\n", Float64bits(1.0))
	var x struct {
		a bool
		b int16
		c []int
	}
	/*
		- pointer, as address are integer. they can be converted to integer value using uintptr. use %d to see their actual integer value.
		- the address of ptr of x must be converted to integer to be added with offset of x.b
			address of b = address of x + offset of x.b
		- we we can perform unsafe.Pointer to convert the integer back to an unsafe pointer
		- we then convert that unsafe pointer to int16 pointer and we use that to change value of x.b
	*/

	pb := (*int16)(unsafe.Pointer(uintptr(unsafe.Pointer(&x)) + unsafe.Offsetof(x.b)))
	*pb = 42
	fmt.Println("==== x.b", x.b)

	/*
		- the code above can run but SUBTLY INCORRECT
		- because garbage collector change variable address all the time. the uintptr may point to the old address of variable, which not able to update that variable anymore
		- this is WRONG too:
		pT := uintptr(unsafe.Pointer(new(T)))
			- after new is called, no variable hold that pointer, it is deleted by garbage collector, so the address is outdated.

	*/
}

func main() {
	example4()
}

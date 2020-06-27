package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

/*
2. equal is with 2 reflect.Value of type string slice
	- x.IsValid and y.IsValid both are true
	- x.Type() and y.Type() are both string slice
	- x.CanAddr and y.CanAddr are false
	- x.Kind() is slice
*/
/*
5. equal is called inside the loop of case slice recursively
	- x, y: reflect.Value of type string, seen is an empty map
	- x.IsValid(), y.IsValid() are true
	- x.Type() and y.Type() are both string
	- x.CanAddr() and y.CanAddr are both true
		- get address of the variable and turn it into an unsafe.Pointer
		- but xptr != yptr
		- we create a comparison struct with xptr, yptr and the type of x. ofcourse it is same as type of y.
		- we set seen for that struct == true
		* but isn't it that the garbage collector can make the address change so it is not that we store anymore?

*/
func equal(x, y reflect.Value, seen map[comparison]bool) bool {
	if !x.IsValid() || !y.IsValid() {
		return x.IsValid() == y.IsValid()
	}
	if x.Type() != y.Type() {
		return false
	}
	if x.CanAddr() && y.CanAddr() {
		xptr := unsafe.Pointer(x.UnsafeAddr())
		yptr := unsafe.Pointer(y.UnsafeAddr())
		if xptr == yptr {
			return true
		}
		c := comparison{xptr, yptr, x.Type()}
		if seen[c] {
			return true
		}
		seen[c] = true
	}
	switch x.Kind() {
	case reflect.Bool:
		return x.Bool() == y.Bool()
	/*
		6. Get to case String
		 - "foo" == "foo", get back to the loop
	*/
	case reflect.String:
		return x.String() == y.String()
	case reflect.Chan, reflect.UnsafePointer, reflect.Func:
		return x.Pointer() == y.Pointer()
	case reflect.Array, reflect.Slice:
		/*
			3. the slice case is called:
				- x.Len() == y.Len() true
				- proceed to loop over each element of x. we don't use range for reflect.Value, it is not slice!
		*/
		if x.Len() != y.Len() {
			return false
		}
		for i := 0; i < x.Len(); i++ {
			/*
				4. for i = 0, x.Index(0), y.Index(0) are "foo", seen is an empty map. call equal recursively
			*/
			/*
				7. the first equal "foo" == "foo" return true, seen has one value
					- similar to "bar" == "bar", seen now has 2 value
			*/
			if !equal(x.Index(i), y.Index(i), seen) {
				return false
			}
		}
		/*
			8. finish the outermost equal, return true
		*/
		return true
	}
	panic("unreachable")
}

/*
1. Equal is called with 2 slice ["foo"] and ["foo"]
	- seen is an empty map
	- reflect.ValueOf(x) is the slice with type reflect.Value
	- equal is call with 2 reflect.Value of type string slice and the empty map
*/
func Equal(x, y interface{}) bool {
	seen := make(map[comparison]bool)
	return equal(reflect.ValueOf(x), reflect.ValueOf(y), seen)
}

type comparison struct {
	x, y unsafe.Pointer
	t    reflect.Type
}

func example6() {
	fmt.Println("=== [\"foo\", \"bar\"] == [\"foo\", \"bar\"]", Equal([]string{"foo", "bar"}, []string{"foo", "bar"}))
	// fmt.Println("=== [\"foo\", \"bar\"] == [\"foo\"]", Equal([]string{"foo", "bar"}, []string{"foo"}))
}

func main() {
	example6()
}

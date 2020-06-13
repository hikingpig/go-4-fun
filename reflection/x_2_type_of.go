package main

import (
	"fmt"
	"io"
	"os"
	"reflect"
)

func example2() {
	t := reflect.TypeOf(3)
	/* assigning concrete value to interface type will
	- performs implicit type conversion
	- creates an interface value consists 2 components:
		- dynamic type (int)
		- dynamic value (3)
	* the concrete value is call the operand
	*/
	fmt.Println(t.String())
	fmt.Println(t)

	var w io.Writer = os.Stdout
	fmt.Println(reflect.TypeOf(w))
	// the returned type is concrete: os.File not os.Writer
}

func main() {
	example2()
}

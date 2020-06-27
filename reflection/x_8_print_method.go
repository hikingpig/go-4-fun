package main

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

/*
	- reflect.ValueOf will return the reflect.Value that has all type info for that struct
	- v.Type() return the type of that struct. like it is defined as time.Duration, strings.Replacer
	- v.NumMethod returns the number of methods
	- v.Method(i).Type() return the type of func for that method, i.e, its input, output
	- and its Strings method give us the description
	- the method name can only be accessed by v.Type().Method(i).Name, not v.Method(i).Name
*/

func Print(x interface{}) {
	v := reflect.ValueOf(x)
	t := v.Type()
	fmt.Printf("type %s\n", t)
	for i := 0; i < v.NumMethod(); i++ {
		methType := v.Method(i).Type()
		fmt.Printf("func (%s) %s%s\n", t, t.Method(i).Name, strings.TrimPrefix(methType.String(), "func"))
	}
}

/*
- strings.Replacer is a Type, we must create a variable from it to use as "expression"
- time.Hour is a const. it can be use as "expression" as input for Print
*/

func main() {
	Print(time.Hour)
	fmt.Println("\n==========================")
	Print(new(strings.Replacer))
}

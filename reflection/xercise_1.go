package main

import (
	"fmt"
	"reflect"
)

func Display2(name string, x interface{}) {
	fmt.Println("\n\n========================================== Display is called")
	fmt.Printf("Display %s (%T):\n", name, x)
	display2(name, reflect.ValueOf(x))
}

// extend Display so it can display maps whose keys are struct or array
func display2(path string, v reflect.Value) string {
	fmt.Printf("\n######### display2 is called with path %s, v %+v\n", path, v)
	fmt.Println("######## v.Kind is", v.Kind())
	res := ""
	switch v.Kind() {
	case reflect.Invalid:
		res += fmt.Sprintf("%s = invalid,", path)
		fmt.Println(res)
		return res
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			res += display2(fmt.Sprintf("%s[%d]", path, i), v.Index(i))
		}
		return res
	case reflect.Struct:
		fmt.Println("########## v.NumField is", v.NumField())
		for i := 0; i < v.NumField(); i++ {
			fieldPath := fmt.Sprintf("%s.%s", path, v.Type().Field(i).Name)
			fmt.Printf("######## i = %d, fieldpath %s\n", i, fieldPath)
			res += display2(fieldPath, v.Field(i))
		}
		return res
	case reflect.Map:
		for _, key := range v.MapKeys() {
			res += display2(fmt.Sprintf("%s[%s]", path, display2("value", key)), v.MapIndex(key))
		}
		return res
	case reflect.Ptr:
		if v.IsNil() {
			res += fmt.Sprintf("%s = nil,", path)
			fmt.Println(res)
		} else {
			res += display2(fmt.Sprintf("(*%s)", path), v.Elem())
		}
		return res
	case reflect.Interface:
		if v.IsNil() {
			res += fmt.Sprintf("%s = nil,", path)
			fmt.Println(res)
		} else {
			fmt.Printf("%s.type = %s,", path, v.Elem().Type())
			res += display2(path+".value", v.Elem())
		}
		return res
	default: // basic types, channels, funcs
		res := fmt.Sprintf("%s = %s,", path, formatAtom(v))
		fmt.Println(res)
		return res
	}
}

func test_exercise1() {
	type weirdKey struct {
		X int
		Y int
	}
	type weirdMap map[weirdKey]string
	myweird := weirdMap{weirdKey{1, 2}: "it's weird"}
	Display2("myweird", myweird)

	eccentric := [2]int{1, 2}
	myweird2 := map[[2]int]string{eccentric: "it's still weird"}
	Display2("myweird v2", myweird2)
}

func main() {
	test_exercise1()
}

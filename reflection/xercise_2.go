package main

import (
	"fmt"
	"os"
	"reflect"
)

/* added log with ######## to follow the flow */
const (
	maxRecursiveStep = 5
)

func Display3(name string, x interface{}) {
	fmt.Println("\n\n========================================== Display3 is called")
	fmt.Printf("Display %s (%T):\n", name, x)
	display3(0, name, reflect.ValueOf(x))
}

func display3(recursiveStep int, path string, v reflect.Value) {
	fmt.Printf("\n######### display3 is called with path %s, v %+v\n", path, v)
	fmt.Println("######## v.Kind is", v.Kind())
	if recursiveStep > maxRecursiveStep {
		fmt.Printf("%s max recursive step reached\n", path)
		return
	}
	recursiveStep += 1
	switch v.Kind() {
	case reflect.Invalid:
		fmt.Printf("%s = invalid\n", path)
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			display3(recursiveStep, fmt.Sprintf("%s[%d]", path, i), v.Index(i))
		}
	case reflect.Struct:
		fmt.Println("########## v.NumField is", v.NumField())
		for i := 0; i < v.NumField(); i++ {
			fieldPath := fmt.Sprintf("%s.%s", path, v.Type().Field(i).Name)
			fmt.Printf("######## i = %d, fieldpath %s\n", i, fieldPath)
			display3(recursiveStep, fieldPath, v.Field(i))
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			display3(recursiveStep, fmt.Sprintf("%s[%s]", path, formatAtom(key)), v.MapIndex(key))
		}
	case reflect.Ptr:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			display3(recursiveStep, fmt.Sprintf("(*%s)", path), v.Elem())
		}
	case reflect.Interface:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			fmt.Printf("%s.type = %s\n", path, v.Elem().Type())
			display3(recursiveStep, path+".value", v.Elem())
		}
	default: // basic types, channels, funcs
		fmt.Printf("%s = %s\n", path, formatAtom(v))
	}
}

func test_exercise2() {
	strangelove := Movie{
		Title:    "Dr. Strangelove",
		Subtitle: "How I Learned to Stop Worrying and Love the Bomb",
		Year:     1964,
		Color:    false,
		Actor: map[string]string{
			"Dr. Strangelove":            "Peter Sellers",
			"Grp. Capt. Lionel Mandrake": "Peter Sellers",
		},
		Oscars: []string{
			"Best Actor (Nomin.",
			"Best Adapted Screenplay (Nomin.)",
		},
	}
	fmt.Println("########### START")
	Display3("strangelove", strangelove)

	Display3("Os.Stderr", os.Stderr)

	var c Cycle
	c = Cycle{42, &c}
	Display3("c", c)
}

// func main() {
// 	fmt.Println("================== 222222222222222222222222222222222")
// 	test_exercise2()
// }

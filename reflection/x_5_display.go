package main

import (
	"fmt"
	"os"
	"reflect"
)

/* added log with ######## to follow the flow */

func Display(name string, x interface{}) {
	fmt.Println("\n\n========================================== Display is called")
	fmt.Printf("Display %s (%T):\n", name, x)
	display(name, reflect.ValueOf(x))
}

func display(path string, v reflect.Value) {
	fmt.Printf("\n######### display is called with path %s, v %+v\n", path, v)
	fmt.Println("######## v.Kind is", v.Kind())
	switch v.Kind() {
	case reflect.Invalid:
		fmt.Printf("%s = invalid\n", path)
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			display(fmt.Sprintf("%s[%d]", path, i), v.Index(i))
		}
	case reflect.Struct:
		fmt.Println("########## v.NumField is", v.NumField())
		for i := 0; i < v.NumField(); i++ {
			fieldPath := fmt.Sprintf("%s.%s", path, v.Type().Field(i).Name)
			fmt.Printf("######## i = %d, fieldpath %s\n", i, fieldPath)
			display(fieldPath, v.Field(i))
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			display(fmt.Sprintf("%s[%s]", path, formatAtom(key)), v.MapIndex(key))
		}
	case reflect.Ptr:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			display(fmt.Sprintf("(*%s)", path), v.Elem())
		}
	case reflect.Interface:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			fmt.Printf("%s.type = %s\n", path, v.Elem().Type())
			display(path+".value", v.Elem())
		}
	default: // basic types, channels, funcs
		fmt.Printf("%s = %s\n", path, formatAtom(v))
	}
}

type Movie struct {
	Title, Subtitle string
	Year            int
	Color           bool
	Actor           map[string]string
	Oscars          []string
	Sequel          *string
}

type Cycle struct {
	Value int
	Tail  *Cycle
}

func example5() {
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
	Display("strangelove", strangelove)

	Display("Os.Stderr", os.Stderr)

	var i interface{} = 3

	Display("i", i)

	Display("&i", &i)

	var c Cycle
	c = Cycle{42, &c}
	// this will run forever!!!
	// Display("c", c)
}

// func main() {
// 	example5()
// }

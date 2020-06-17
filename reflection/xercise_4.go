package main

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

/* we don't have encoding for Bool yet, so remove color field from Movie
else, it will results in an error and we will not see anything
use buf.String() to display the result
*/

func encode3(buf *bytes.Buffer, v reflect.Value, indentDepth int) error {
	fmt.Println("\n========= encode3 is called, v.Kind is", v.Kind())
	fmt.Println("=========== indentdepth", indentDepth)
	switch v.Kind() {
	case reflect.Invalid:
		buf.WriteString("nil")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Fprintf(buf, "%d", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fmt.Fprintf(buf, "%d", v.Uint())
	case reflect.String:
		fmt.Fprintf(buf, "%q", v.String())
	case reflect.Ptr:
		return encode3(buf, v.Elem(), indentDepth+1)
	case reflect.Array, reflect.Slice:
		buf.WriteByte('(')
		for i := 0; i < v.Len(); i++ {
			fmt.Println("################### looping over arrary/slice, i =", i)
			if i > 0 {
				buf.WriteString("\n" + strings.Repeat("\t", indentDepth+1))
			}
			if err := encode3(buf, v.Index(i), indentDepth); err != nil {
				return err
			}
			fmt.Println("==================== buf now is", buf.String())
		}
		buf.WriteByte(')')
	case reflect.Struct:
		fmt.Println("======== struct numfield", v.NumField())
		buf.WriteByte('(')
		for i := 0; i < v.NumField(); i++ {
			fmt.Println("################## looping over numfield, i =", i)
			if i > 0 {
				buf.WriteString("\n " + strings.Repeat(" ", indentDepth+1))
			}
			fmt.Fprintf(buf, "(%s\t", v.Type().Field(i).Name)
			if err := encode3(buf, v.Field(i), indentDepth+1); err != nil {
				return err
			}
			buf.WriteByte(')')
			fmt.Println("========== buf now is", buf.String())
		}
		buf.WriteByte(')')

	case reflect.Map:
		buf.WriteByte('(')
		for i, key := range v.MapKeys() {
			fmt.Println("##################### looping over map keys, i =", i, "key =", key)
			buf.WriteByte('(')
			if i > 0 {
				buf.WriteString("\n" + strings.Repeat("\t", indentDepth+1))
			}
			if err := encode3(buf, key, indentDepth); err != nil {
				return err
			}
			buf.WriteByte(' ')
			if err := encode3(buf, v.MapIndex(key), indentDepth); err != nil {
				return err
			}
			buf.WriteByte(')')
			fmt.Println("============ buf now is", buf.String())
		}
		buf.WriteByte(')')
	default: // float, complex, bool, chan, func, interface
		return fmt.Errorf("unsupported type: %s", v.Type())
	}
	buf.WriteByte(')')
	return nil
}

func Marshal3(v interface{}) ([]byte, error) {
	fmt.Println("======= Marshal3 is called")
	var buf bytes.Buffer
	if err := encode3(&buf, reflect.ValueOf(v), 0); err != nil {
		fmt.Println("=========== encode3 error", err)
		return nil, err
	}
	return buf.Bytes(), nil
}

func test_xercise4() {
	type Movie struct {
		Title, Subtitle string
		Year            int
		Actor           map[string]string
		Oscars          []string
		Sequel          *string
	}

	strangelove := Movie{
		Title:    "Dr. Strangelove",
		Subtitle: "How I Learned to Stop Worrying and Love the Bomb",
		Year:     1964,
		Actor: map[string]string{
			"Dr. Strangelove":            "Peter Sellers",
			"Grp. Capt. Lionel Mandrake": "Peter Sellers",
		},
		Oscars: []string{
			"Best Actor (Nomin.)",
			"Best Adapted Screenplay (Nomin.)",
		},
	}
	fmt.Println("########### START")
	data, err := Marshal3(strangelove)
	if err != nil {
		fmt.Println("======== err", err)
	}
	fmt.Println("=========== FINAL RESULT: \n", string(data))
}

// func main() {
// 	test_xercise4()
// }

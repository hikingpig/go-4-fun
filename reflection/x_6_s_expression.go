package main

import (
	"bytes"
	"fmt"
	"reflect"
)

/* we don't have encoding for Bool yet, so remove color field from Movie
else, it will results in an error and we will not see anything
use buf.String() to display the result
*/

func encode(buf *bytes.Buffer, v reflect.Value) error {
	fmt.Println("\n========= encode is called, v.Kind is", v.Kind())
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
		return encode(buf, v.Elem())
	case reflect.Array, reflect.Slice:
		buf.WriteByte('(')
		for i := 0; i < v.Len(); i++ {
			fmt.Println("################### looping over arrary/slice, i =", i)
			if i > 0 {
				buf.WriteByte(' ')
			}
			if err := encode(buf, v.Index(i)); err != nil {
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
				buf.WriteByte((' '))
			}
			fmt.Fprintf(buf, "(%s ", v.Type().Field(i).Name)
			if err := encode(buf, v.Field(i)); err != nil {
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
				buf.WriteByte(' ')
			}
			if err := encode(buf, key); err != nil {
				return err
			}
			buf.WriteByte(' ')
			if err := encode(buf, v.MapIndex(key)); err != nil {
				return err
			}
			buf.WriteByte(')')
			fmt.Println("============ buf now is", buf.String())
		}
		buf.WriteByte(')')
	default: // float, complex, bool, chan, func, interface
		return fmt.Errorf("unsupported type: %s", v.Type())
	}
	return nil
}

func Marshal(v interface{}) ([]byte, error) {
	fmt.Println("======= Marshal is called")
	var buf bytes.Buffer
	if err := encode(&buf, reflect.ValueOf(v)); err != nil {
		fmt.Println("=========== encode error", err)
		return nil, err
	}
	return buf.Bytes(), nil
}

func example6() {
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
	Marshal(strangelove)
}

func main() {
	example6()
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

/* we don't have encoding for Bool yet, so remove color field from Movie
else, it will results in an error and we will not see anything
use buf.String() to display the result

json marshal invalid type as null, not nil

*/

func encode4(buf *bytes.Buffer, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Invalid:
		buf.WriteString("null")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Fprintf(buf, "%d", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fmt.Fprintf(buf, "%d", v.Uint())
	case reflect.String:
		fmt.Fprintf(buf, "%q", v.String())
	case reflect.Ptr:
		return encode4(buf, v.Elem())
	case reflect.Array, reflect.Slice:
		buf.WriteByte('[')
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				buf.WriteByte(',')
			}
			if err := encode4(buf, v.Index(i)); err != nil {
				return err
			}
		}
		buf.WriteByte(']')
	case reflect.Struct:
		buf.WriteByte('{')
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				buf.WriteByte((','))
			}
			fmt.Fprintf(buf, "%q: ", v.Type().Field(i).Name)
			if err := encode4(buf, v.Field(i)); err != nil {
				return err
			}
		}
		buf.WriteByte('}')

	case reflect.Map:
		buf.WriteByte('{')
		for i, key := range v.MapKeys() {
			if i > 0 {
				buf.WriteByte(',')
			}
			if err := encode4(buf, key); err != nil {
				return err
			}
			buf.WriteByte(':')
			if err := encode4(buf, v.MapIndex(key)); err != nil {
				return err
			}
		}
		buf.WriteByte('}')
	default: // float, complex, bool, chan, func, interface
		return fmt.Errorf("unsupported type: %s", v.Type())
	}
	return nil
}

func Marshal4(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := encode4(&buf, reflect.ValueOf(v)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func test_xercise_5() {
	type Movie struct {
		Title    string            `json:"Title"`
		Subtitle string            `json:"Subtitle"`
		Year     int               `json:"Year"`
		Actor    map[string]string `json:"Actor"`
		Oscars   []string          `json:"Oscars"`
		Sequel   *string           `json:"Sequel"`
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
	data, err := Marshal4(strangelove)
	if err != nil {
		fmt.Println("======= error", err)
	}
	fmt.Println("====== Marsal result: ", string(data))
	mymovie := Movie{}
	err = json.Unmarshal(data, &mymovie)
	if err != nil {
		fmt.Println("======= error", err)
	}
	fmt.Println("======= unmarshal result: ", mymovie)
}

// func main() {
// 	test_xercise_5()
// }

package main

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

func search(resp http.ResponseWriter, req *http.Request) {
	var data struct {
		Labels     []string `http:"l"`
		MaxResults int      `http:"max"`
		Exact      bool     `http:"x"`
	}
	/*
		1. req has the field URL, with the params like: /search?x=true&l=golang&l=programming
	*/
	data.MaxResults = 10
	if err := Unpack(req, &data); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(resp, "Search: %+v\n", data)
}

func Unpack(req *http.Request, ptr interface{}) error {
	/*
		2. req also has Form field, which is empty befor ParseForm called
		after parse form, the Form map is filled with keys and values from the url
		x=true&l=golang&l=programming
		- the values for "l" will be an array [golang, programming]
	*/
	if err := req.ParseForm(); err != nil {
		return err
	}
	fields := make(map[string]reflect.Value)
	/*
		3. we take the variable (addressable) of the address for the output struct so we can change it with reflect
	*/

	/*
		v.Type() will return the type of struct of output struct, including its fields and tags
		v.Type().Field(i) return the type of field in the struct, like its Name, Tag.
		v.Field(i) return the current value of field i

		4. we loop over each field of the output struct, and fill the map with the field name as key and the field value as value. the values are that from the original output struct, not modified by the url yet.

		we need a map to match it with the map of the Form
		so the struct is only 1 level depth, we never go more
	*/
	v := reflect.ValueOf(ptr).Elem()
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i)
		tag := fieldInfo.Tag
		name := tag.Get("http")
		if name == "" {
			name = strings.ToLower(fieldInfo.Name)
		}
		/*
			so when we change fields[name], it will change v.Field and v will also be changed!
		*/
		fields[name] = v.Field(i)
	}

	/*
		5. now we will loop through the req Form that contains the url params
			- the name is key from the form, we check if it is in the fields map, if not, skip
			- all the values are in form of []string
			- Kind contains Type, for bool, Type = Kind. For Kind = slice, Type = []string
			- reflect.New create a zero pointer. For slice, we call f.Type().Elem() to get the String type. Then New().Elem to give us an addressable reflect.Value to that string element
			- we gonna call populate to give that elem the value
	*/

	for name, values := range req.Form {
		f := fields[name]
		if !f.IsValid() {
			continue
		}
		for _, value := range values {
			if f.Kind() == reflect.Slice {
				elem := reflect.New(f.Type().Elem()).Elem()
				if err := populate(elem, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
				/* for slice: f = append(f, elem)
				- but for reflect.Value we use f.Set(reflect.Append(f, elem))
				- they are reflect.Values, not slice or string
				*/
				f.Set(reflect.Append(f, elem))
			} else {
				if err := populate(f, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
			}
		}
	}
	return nil
}

/*
6. Populate will fill the reflect.Value.
	- But for Int and Bool, we will need to parse it with strconv first
	- the change will be reflected in output struct
*/
func populate(v reflect.Value, value string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(value)
	case reflect.Int:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		v.SetBool(b)
	default:
		return fmt.Errorf("unsupported kind %s", v.Type())
	}
	return nil
}

func main() {
	http.HandleFunc("/search", search)
	log.Fatal(http.ListenAndServe(":12345", nil))
}

/* commands to run:
go run x_8_field_tag.go

- should put '' around the url so the params will be read fully

curl 'http://localhost:12345/search'

curl 'http://localhost:12345/search?l=golang&l=programming'

curl 'http://localhost:12345/search?l=golang&l=programming&max=100'

curl 'http://localhost:12345/search?x=true&l=golang&l=programming'

curl 'http://localhost:12345/search?q=hello&x=123'

curl 'http://localhost:12345/search?q=hello&max=lots'
*/

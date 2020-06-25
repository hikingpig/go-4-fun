package main

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"text/scanner"
)

/* scanner is text/scanner not go/scanner */
type lexer struct {
	scan  scanner.Scanner
	token rune
}

func (lex *lexer) next() {
	lex.token = lex.scan.Scan()
}

func (lex *lexer) text() string {
	return lex.scan.TokenText()
}

func (lex *lexer) consume(want rune) {
	if lex.token != want {
		panic(fmt.Sprintf("got %q, want %q", lex.text(), want))
	}
	lex.next()
}

func read(lex *lexer, v reflect.Value) {
	switch lex.token {
	case scanner.Ident:
		if lex.text() == "nil" {
			v.Set(reflect.Zero(v.Type()))
			lex.next()
			return
		}
	case scanner.String:
		s, _ := strconv.Unquote(lex.text())
		v.SetString(s)
		lex.next()
		return
	case scanner.Int:
		i, _ := strconv.Atoi(lex.text())
		v.SetInt(int64(i))
		lex.next()
		return
	case '(':
		/*
			2. the first charater is '(', the second character is still '('
			- we assume it is struct or map or array and call readList

			8. lex.next, now lex at 2nd '(' of the map, which enclose a key-value pair, it will call readList for the Map
		*/
		lex.next()
		readList(lex, v)
		lex.next()
		return
	}
	panic(fmt.Sprintf("unexpected tokens %q", lex.text()))
}

func readList(lex *lexer, v reflect.Value) {
	switch v.Kind() {
	case reflect.Array:
		for i := 0; !endList(lex); i++ {
			read(lex, v.Index(i))
		}
	case reflect.Slice:
		for !endList(lex) {
			item := reflect.New(v.Type().Elem()).Elem()
			read(lex, item)
			v.Set(reflect.Append(v, item))
		}
	case reflect.Struct:
		/*
			3. it is struct, the 2nd character ( now enclosinng a field
				- we will iterate through all the field
		*/
		for !endList(lex) {
			lex.consume('(')
			if lex.token != scanner.Ident {
				panic(fmt.Sprintf("got token %q, want field name", lex.text()))
			}
			name := lex.text()
			lex.next()
			read(lex, v.FieldByName(name))
			lex.consume(')')
			/*
				4. the first field is finished, lex is now at ( of 2nd field. similarly, it is a string
				5. the 2nd field is finished, lex is now at ( of 3rd field, it is an Int
				6. the 3rd field is finished, lex is now at ( of 4th field.
				7. lex consume the first '(' of the field name Actor, get the name, and move to the '(' enclosing
					a Map
					- it will call read with lex at 1st '(' of a Map, which is case '('
				10. similarly for slice.
				11. "nil" is read by Ident
			*/
		}
	/*
		9. the type Map is read, lex at the 2nd '(' enclosing a key-value pair
		- we will loop over the key-value pair one by one, read them into key-value and set into the map
		- endlist finish the map, return to struct readList
	*/
	case reflect.Map:
		/*

			- Some syntax:
					- v.Type is map[string]string. we give v an empty map of that type first
					- v.Type.Key() return "string", which is the type of key
					- reflect.New return a pointer to that value. .Elem makes an addressable Value
					- similarly v.Type().Elem() returns the type of the value, also string
					- we use them to pass into read, as read receive reflect.Value
		*/
		v.Set(reflect.MakeMap(v.Type()))
		fmt.Println("===== v.type.key", v.Type().Key())
		for !endList(lex) {
			lex.consume('(')
			key := reflect.New(v.Type().Key()).Elem()
			read(lex, key)
			value := reflect.New(v.Type().Elem()).Elem()
			read(lex, value)
			v.SetMapIndex(key, value)
			lex.consume(')')
		}
	default:
		panic(fmt.Sprintf("cannot decode list into %v", v.Type()))
	}
}

func endList(lex *lexer) bool {
	switch lex.token {
	case scanner.EOF:
		panic("end of file")
	case ')':
		return true
	}
	return false
}

func Unmarshal(data []byte, out interface{}) (err error) {
	lex := &lexer{scan: scanner.Scanner{Mode: scanner.GoTokens}}
	// use bytes.NewReader to create an io.Reader for scan.Init
	lex.scan.Init(bytes.NewReader(data))
	/*
		1. the first character is '('
		 - we will move to read with out as movie struct
	*/
	lex.next()
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("error at %s: %v", lex.scan.Position, x)
		}
	}()
	read(lex, reflect.ValueOf(out).Elem())
	return nil
}

func example7() {
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
	data, err := Marshal(strangelove)
	if err != nil {
	}

	mymovie := Movie{}
	fmt.Println("\n======================== START UNMARSHALLING")
	if err := Unmarshal(data, &mymovie); err != nil {
	}
	fmt.Println("======== Unmarshalling Result:", mymovie)
}

func main() {
	example7()
}

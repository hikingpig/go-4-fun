package main

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"text/scanner"
)

/* scanner lis text/scanner not go/scanner */
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
	/* 5. the next scan got "Dr. Strangelove" whole as a string
	- so we unquote it and set the value to the field
	- the recursive stop here, return back to readlist
	- we also call next to move to the next character
	*/
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
			3. the first character read is ( which encloses a struct
				the second character is still ( which enclose a field
				we will call readList to read that struct
		*/
		lex.next()
		readList(lex, v)
		/* 7. readList finished. we have read the whole struct.
		- lex now is ')' the last character closing the struct.
		- lex.next will move to EOF
		- we return to Unmarshal now
		*/
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
		/*4. after the first character, we decide what to do next base on type of out, this case a struct
		- so the 2nd character must be ( too, which enclose a field's name and value
		- we call consume to make sure we got the correct string.
		- consume will also make the next scan. We expect it is a word or "Ident"
		- the word is the field name, so we call read again to change the field value
		- v.FieldByFieldName return an addressable to make change on out's value
		*/
		for !endList(lex) {
			lex.consume('(')
			if lex.token != scanner.Ident {
				panic(fmt.Sprintf("got token %q, want field name", lex.text()))
			}
			name := lex.text()
			lex.next()
			read(lex, v.FieldByName(name))
			/*
				6. the recursion return, the field Title is changed.
					-the struct is finished, we expect the next character to be ')', closing the field
					- it is the only field in the struct, so the next character is ')', closing the struct
					- we will call endlist with ')', which return true. the loop stop
					- get back to read
			*/
			lex.consume(')')
		}
	case reflect.Map:
		v.Set(reflect.MakeMap(v.Type()))
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
	lex.scan.Init(bytes.NewReader(data))
	/* 1. read the first character of data, it is (, a non-identifier*/
	lex.next()
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("error at %s: %v", lex.scan.Position, x)
		}
	}()
	/* 2. get access to the variable of the address of the out (Movie struct)
	- the reflect.Value from Elem() can change the value of out
	*/
	read(lex, reflect.ValueOf(out).Elem())
	return nil
}

func example7() {
	type Movie struct {
		Title string
	}

	strangelove := Movie{
		Title: "Dr. Strangelove",
	}
	data, _ := Marshal(strangelove)
	fmt.Println("\n============== START UNMARSHALLING")
	mymovie := Movie{}
	if err := Unmarshal(data, &mymovie); err != nil {
	}
	fmt.Println("====== Unmarshalled result: ", mymovie)
}

func main() {
	example7()
}

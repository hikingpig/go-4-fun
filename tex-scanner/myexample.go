package main

// naming a file with _test will make visual code put Run Test in every function in the file
import (
	"fmt"
	"strings"
	"text/scanner"
	"unicode"
)

/* command: go test -v -run Example/
the / is to run Example() only, not the whole suite
*/
func Example1() {
	const src = `
// comment will be ignored
if a > 10 {
	someParsable = text
}`
	var s scanner.Scanner
	s.Init(strings.NewReader(src))
	/* the Position like example:3:1
	- FileName:Line:Column
	- the line will be counted even if empty or a comment but
	- the scanner only pointed at the first interested position
	*/
	s.Filename = "example"

	/* here's how we move around the scanner in a loop:
		-----------------------------------------------------
		for i := 0; i < 10; i++
		-----------------------------------------------------
		for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan()
		-----------------------------------------------------
	- the first s.Scan run only once and will not be called again
	- the 2nd s.Scan like i++ will move to scanner forward in the loop
	*/

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		fmt.Printf("%s: %s\n", s.Position, s.TokenText())
	}
	/*
		- each token.Scan move to the beginning of a word
		- each word is separated by a space or new line character
		- the s.TokenText() will return that word
	*/
	fmt.Println("=========== END Example1")
}

func Example_isIdentRune() {
	const src = "%var1 var2%"

	var s scanner.Scanner

	s.Init(strings.NewReader(src))
	s.Filename = "default"

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		fmt.Printf("%s: %s\n", s.Position, s.TokenText())
	}

	fmt.Println()

	/* Init will allow us to start scanning the src from the beginning again */

	s.Init(strings.NewReader(src))
	s.Filename = "percent"
	/* use string() to display the rune and byte character
	How the method Scan works:
	- the character " " space is not listed as identifier and will be skipped whenever encountered
	- scan will stop whenever encounter a non-identifier character, the word is between 2 non-identifier
	- i is the position of the character in the new word, not the line of src data
	- The IsIdentRune will define new identifier: the % will not be an identifier if it is not at position i = 0 in the new word.
	- but this % will not be ignored and will be counted again as scan building new word
	- the 3rd word is "%".
	- the final token is EOF character. As a non-identifier, it will stop the scan and then as an EOF token, it will stop the loop. EOF appears twice.
	*/
	s.IsIdentRune = func(ch rune, i int) bool {
		fmt.Println("\nch =", string(ch), ", i =", i)
		if ch == '%' && i == 0 || unicode.IsLetter(ch) || unicode.IsDigit(ch) && i > 0 {
			fmt.Println("IDENT FOUND")
			return true
		}
		fmt.Println("NON-IDENT FOUND")
		return false
	}

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		fmt.Printf("(TokenText) %s: %s\n", s.Position, s.TokenText())
	}
}

func Example_CommentMode() {
	const src = `
		// First comments
		
 some normal words ...

/*
Second
Comment
*/
`

	var s scanner.Scanner
	s.Init(strings.NewReader(src))
	s.Filename = "comments"
	/* default mode has many setup in it, including skip comments
	here, we just turn off the bit for skip comments
	the other bits are still the same
	*/
	s.Mode ^= scanner.SkipComments // don't skip comments

	/* the entire comment will be treated as one word during Scan */
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		txt := s.TokenText()
		if strings.HasPrefix(txt, "//") || strings.HasPrefix(txt, "/*") {
			fmt.Printf("%s: %s\n", s.Position, txt)
		}
	}
}

func Example_whitespace() {
	// tab-separated values
	const src = `aa	ab	ac	ad
ba	bb	bc	bd
ca	cb	cc	cd
da	db	dc	dd`

	var (
		col, row int
		s        scanner.Scanner
		tsv      [4][4]string // large enough for example above
	)
	s.Init(strings.NewReader(src))
	fmt.Printf("whitespce before %d or %b\n", s.Whitespace, s.Whitespace)
	s.Whitespace ^= 1<<'\t' | 1<<'\n' // don't skip tabs and new lines
	fmt.Printf("1<<t %d or %b \n", 1<<'\t', 1<<'\t')
	fmt.Printf("1<<n %d or %b \n", 1<<'\n', 1<<'\n')
	fmt.Printf("whitespce after %d or %b\n", s.Whitespace, s.Whitespace)
	/*
		- we turn off the bit for whitespace with tab and new line, so they are not ignored anymore
		- but they are also not consider the same as other identifiers like letter and digit.
		- they are similar to the % in we define in isIdentRune
		- a new word will always be created for each tab or new line
	*/
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		fmt.Println("==== token is", s.TokenText())
		switch tok {
		case '\n':
			row++
			col = 0
		case '\t':
			col++
		default:
			tsv[row][col] = s.TokenText()
		}
	}
	fmt.Println("========== the array is")
	fmt.Print(tsv)
	fmt.Println()
	// Output:
	// [[aa ab ac ad] [ba bb bc bd] [ca cb cc cd] [da db dc dd]]
}

func main() {
	// Example1()
	// Example_isIdentRune()
	// Example_CommentMode()
	Example_whitespace()
}

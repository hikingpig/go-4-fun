package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
	"text/scanner"
)

type StringReader struct {
	data []string
	step int
}

/*
The Read method will be called by Next method of the scanner
- the p lenth is 1024, it can hold 128 characters
- p also holds some character from previous read at the end. not sure how these characters added
- Next method of scanner seems to call read for next element of the slice string if it has 0 or ONLY 1 character to read.
- Next method SKIP the LAST character of the word from the current read and also include the LAST character of previous read.
- the step starts from 0 to len(data), so we read len(data) + 1 times, the last read return EOF. the last Next return the last character from nth Read.
- the Whitespace character is not omitted.
*/

func (r *StringReader) Read(p []byte) (n int, err error) {
	fmt.Println("\n========== read is called, step", r.step, "len data", len(r.data))
	if r.step < len(r.data) {
		s := r.data[r.step]
		fmt.Println("======== s = ", s)
		n = copy(p, s)
		r.step++
	} else {
		fmt.Println("######### ERROR EOF occurs")
		err = io.EOF
	}
	fmt.Println("====== n = ", n)
	return
}

func readRuneSegments(t *testing.T, segments []string) {
	fmt.Println("======== segments", segments)
	got := ""
	want := strings.Join(segments, "")
	s := new(scanner.Scanner).Init(&StringReader{data: segments})
	for {
		ch := s.Next()
		fmt.Printf("======== ch = %s\n", string(ch))
		if ch == scanner.EOF {
			break
		}
		got += string(ch)
	}
	if got != want {
		t.Errorf("segments=%v got=%s want=%s", segments, got, want)
	}
}

var segmentList = [][]string{
	// {"Hello", ", ", "World", "!"},
	{"Hello123aaabbbccc", ", 456dddeeefff", "World789ggghhhiii", "!kkklllmmm", "X"},
}

func TestNext(t *testing.T) {
	for _, s := range segmentList {
		fmt.Println("======== s =", s)
		readRuneSegments(t, s)
	}
}

type token struct {
	tok  rune
	text string
}

/*
use strings.Repeat to get the precise repeating of string
*/

var f100 = strings.Repeat("f", 100)

var tokenList = []token{
	{scanner.Comment, "// line comments"},
	{scanner.Ident, "foobar"},
	{scanner.Ident, "abc123"},
	{scanner.Int, "1234567890"},
	{scanner.Char, `'\ufA16'`},
}

/*
fmt.Fprintf is to write into buffer or any Writer.
The formatting rule is the same with fmt.Printf
*/

func makeSource(pattern string) *bytes.Buffer {
	var buf bytes.Buffer
	for _, k := range tokenList {
		fmt.Fprintf(&buf, pattern, k.text)
	}
	fmt.Println("======= buff now: ", buf.String())
	return &buf
}

/*
- got, and want are integer that describe the type of token we have. rune is used to for untyped integer type for these values
- scanner.TokenString will return the string representation of these like Ident, ... It retrieves values from a map
- s.Scan() will change the line according to position of the word in the source data
- s.TokenText() returns the word
- in this test we check if:
	- the words are the same
	- the position of words are correct
*/

func checkTok(t *testing.T, s *scanner.Scanner, line int, got, want rune, text string) {
	fmt.Println("\n====== checkTok: got", got, "want: ", want, "line", line)
	if got != want {
		t.Fatalf("tok = %s, want %s for %q", scanner.TokenString(got), scanner.TokenString(want), text)
	}
	fmt.Println("====== S.Line: ", s.Line)
	if s.Line != line {
		t.Errorf("line = %d, want %d for %q", s.Line, line, text)
	}
	stext := s.TokenText()
	fmt.Println("====== stext", stext)
	if stext != text {
		t.Errorf("text = %q, want %q", stext, text)
	} else {
		/* check idempotency of TokenText() call
		- idempotency means the func returns the same result over and over
		*/
		stext = s.TokenText()
		if stext != text {
			t.Errorf("text = %q, want %q (idempotency check)", stext, text)
		}
	}
}

func countNewlines(s string) int {
	n := 0
	for _, ch := range s {
		if ch == '\n' {
			n++
		}
	}
	return n
}

func testScan(t *testing.T, mode uint) {
	fmt.Printf("============== testScan mode %d or %b\n", mode, mode)
	s := new(scanner.Scanner).Init(makeSource(" \t%s\n"))
	s.Mode = mode
	/* s.Scan will skip the comments and return the first non-comment word*/
	tok := s.Scan()
	line := 1
	/* we will loop over the token list and use them to check the values in the scanner */
	for _, k := range tokenList {
		/* if the mode is not skip comments, we will just proceed with any token
		- if the token is not a comment, we don't care about the mode
		- if the token is a comment and mode has skip comment, we skip it but still count the lines in it
		*/
		if mode&scanner.SkipComments == 0 || k.tok != scanner.Comment {
			checkTok(t, s, line, tok, k.tok, k.text)
			tok = s.Scan()
		}
		line += countNewlines(k.text) + 1 // each token is on a new line
	}
	checkTok(t, s, line, tok, scanner.EOF, "")
}

func TestScan(t *testing.T) {
	testScan(t, scanner.GoTokens)

	/* AND NOT:
	^scanner.Skip comment will return a sequence of all 1 except the skip comment bit
	&^scanner.SkipComment will keep all the bit the same except the bit of skip comment is turned off
	*/
	// testScan(t, scanner.GoTokens&^scanner.SkipComments)
}

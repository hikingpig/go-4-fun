package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
	"text/scanner"
	"unicode/utf8"
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

var f10 = strings.Repeat("f", 10)

var tokenList = []token{
	{scanner.Comment, "// line comments"},
	{scanner.Ident, "foobar"},
	{scanner.Ident, "abc123"},
	{scanner.Int, "1234567890"},
	{scanner.Char, `'\ufA16'`},
	{scanner.RawString, "`" + f10 + "`"},
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
	fmt.Println("====== checkTok: got", got, "want: ", want, "line", line)
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

/*
	- s.Scan() inside checkTok will run first. it will detect the error and call s.Error we define implicitly
	- that function we define just to check if the msg is right. maybe the default s.Error does nothing
	- even we have an error, Scan still returns the exact word and the exact type scanner.Float and checkTok still works well
	- the checkToErr is called when we are sure there is an error. if not the ErrorCount will not increase and it will blow up
*/

func checkTokErr(t *testing.T, s *scanner.Scanner, line int, want rune, text string) {
	fmt.Println("\n========= checkToErr is called")
	prevCount := s.ErrorCount
	fmt.Println("======== prevCount:", prevCount)
	fmt.Println("======== calling Scan before checkTok")
	checkTok(t, s, line, s.Scan(), want, text)
	if s.ErrorCount != prevCount+1 {
		t.Fatalf("want error for %q", text)
	}
}

func TestInvalidExponent(t *testing.T) {
	const src = "1.5e 1.5E 1e+ 1e- 1.5z"
	s := new(scanner.Scanner).Init(strings.NewReader(src))
	s.Error = func(s *scanner.Scanner, msg string) {
		fmt.Println("####### s.Error is called msg:", msg)
		const want = "exponent has no digits"
		if msg != want {
			t.Errorf("%s: got error %q; want %q", s.TokenText(), msg, want)
		}
	}
	checkTokErr(t, s, 1, scanner.Float, "1.5e")
	checkTokErr(t, s, 1, scanner.Float, "1.5E")
	checkTokErr(t, s, 1, scanner.Float, "1e+")
	checkTokErr(t, s, 1, scanner.Float, "1e-")
	fmt.Printf("\n======= CHECK TOKEN ONLY \n\n")
	/*
		- for 1.5z, Scan will separate it into 2 words
			- the Float with 1.5
			- the Ident with z
		- the last scan will be EOF type, which results in an empty string for s.TokenText()
	*/
	checkTok(t, s, 1, s.Scan(), scanner.Float, "1.5")
	checkTok(t, s, 1, s.Scan(), scanner.Ident, "z")
	checkTok(t, s, 1, s.Scan(), scanner.EOF, "")
	if s.ErrorCount != 4 {
		t.Errorf("%d errors, want 4", s.ErrorCount)
	}
}

func TestPosition(t *testing.T) {
	src := makeSource("\t\t\t\t%s\n")
	s := new(scanner.Scanner).Init(src)
	s.Mode = scanner.GoTokens &^ scanner.SkipComments
	s.Scan()
	/*
		+ Offset is how far away the word is from the begining of the source data. For the 1st word, Offset is 4 because we have 4 tabs before the 1st character of the word. that word's column is 5
		- S.Offset will returns the Offset, maybe the white spaces before the 1st character
		- For every new word, to get Offset, we must count all the characters before it, including white spaces characters and the text
	*/
	pos := scanner.Position{"", 4, 1, 5}
	for _, k := range tokenList {
		fmt.Println("=== Offset:", s.Offset)
		if s.Offset != pos.Offset {
			t.Errorf("offset = %d, want %d for %q", s.Offset, pos.Offset, k.text)
		}
		if s.Line != pos.Line {
			t.Errorf("line = %d, want %d for %q", s.Line, pos.Line, k.text)
		}
		if s.Column != pos.Column {
			t.Errorf("column = %d, want %d for %q", s.Column, pos.Column, k.text)
		}
		pos.Offset += 4 + len(k.text) + 1 // 4 tabs + token bytes + newline
		pos.Line += countNewlines(k.text) + 1
		s.Scan()
	}
	// make sure there were no token-internal errors reported by scanner
	if s.ErrorCount != 0 {
		t.Errorf("%d errors", s.ErrorCount)
	}
}

/*
- in zero mode, we scan the string "character by character". we return the rune integer of that exact character instead of the token class like Ident
- setting whitespace = 0, we will not skip any white space
*/
func TestScanZeroMode(t *testing.T) {
	src := makeSource("%s\n")
	str := src.String()
	s := new(scanner.Scanner).Init(src)
	s.Mode = 0       // don't recognize any token classes
	s.Whitespace = 0 // don't skip any whitespace
	tok := s.Scan()
	fmt.Println("======= first scan text:", s.TokenText(), "tok: ", tok)
	for i, ch := range str {
		fmt.Println("====== i =", i, "ch = ", string(ch))
		if tok != ch {
			t.Fatalf("%d. tok = %s, want %s", i, scanner.TokenString(tok), scanner.TokenString(ch))
		}
		tok = s.Scan()
	}
	if tok != scanner.EOF {
		t.Fatalf("tok = %s, want EOF", scanner.TokenString(tok))
	}
	if s.ErrorCount != 0 {
		t.Errorf("%d errors", s.ErrorCount)
	}
}

func testScanSelectedMode(t *testing.T, mode uint, class rune) {
	src := makeSource("%s\n")
	s := new(scanner.Scanner).Init(src)
	s.Mode = mode
	tok := s.Scan()
	for tok != scanner.EOF {
		fmt.Println("===== tok:", tok, "text: ", s.TokenText())
		if tok < 0 && tok != class {
			t.Fatalf("tok = %s, want %s", scanner.TokenString(tok), scanner.TokenString(class))
		}
		tok = s.Scan()
	}
	if s.ErrorCount != 0 {
		t.Errorf("%d errors", s.ErrorCount)
	}
}

func TestScanSelectedMask(t *testing.T) {
	/*
		- for mod 0, Scan still skip all the white spaces while return exactly character by character. the class will not be used at all. it can be any value
	*/

	testScanSelectedMode(t, 0, 0)

	/*
		- For scan Idents, the Scan still recognize non-identifier character by character, comments are not skipped
		- /, ``, numbers, \, ' are non-ident. and it returns their exact rune int
		- for idents, it recognized them as words between non-idents and white spaces. and it returns -2, code for Idents
		- so we can get the data of type we want based on negative code tok
	*/
	testScanSelectedMode(t, scanner.ScanIdents, scanner.Ident)

	// Don't test ScanInts and ScanNumbers since some parts of
	// the floats in the source look like (invalid) octal ints
	// and ScanNumbers may return either Int or Float.

	/*Char mode will recognize char representation such as '\u123'
	the other characters are read one by one and return their rune too
	*/
	testScanSelectedMode(t, scanner.ScanChars, scanner.Char)
	/* String Mode will read everything one by one, ignore the white spaces
	- quite similar to mode 0
	*/
	testScanSelectedMode(t, scanner.ScanStrings, scanner.String)
	/* the comments are still read character by character */
	testScanSelectedMode(t, scanner.SkipComments, 0)
	/* ScanComments will read the whole comment as a word
	- others are read character by character, white spaces are skipped
	*/
	testScanSelectedMode(t, scanner.ScanComments, scanner.Comment)
}

/*
- IsIdentRune is the function we can define the ident ourselves
-	Group of idents between 2 non-idents or white spaces is a word
*/
func TestScanCustomIdent(t *testing.T) {
	const src = "faab12345 a12b123 a12 3b"
	s := new(scanner.Scanner).Init(strings.NewReader(src))
	/*
		Idents are defined as:
			- character a or b at position 0 in the word
			- or a character 0, 1, 2, 3 at position 0 < i < 4
	*/
	s.IsIdentRune = func(ch rune, i int) bool {
		res := i == 0 && (ch == 'a' || ch == 'b') || 0 < i && i < 4 && '0' <= ch && ch <= '3'
		fmt.Println("===== is identRune for ch =", string(ch), "i =", i, res)
		return res
	}
	fmt.Println("********************")
	checkTok(t, s, 1, s.Scan(), 'f', "f")
	fmt.Println("********************")
	/* a is Ident but the next a is not, a is the word */
	checkTok(t, s, 1, s.Scan(), scanner.Ident, "a")
	/* 2nd a is not Ident in previous scan but an Ident in this scan, b is not Ident in this scan, a is the word */
	fmt.Println("********************")
	checkTok(t, s, 1, s.Scan(), scanner.Ident, "a")
	fmt.Println("********************")
	checkTok(t, s, 1, s.Scan(), scanner.Ident, "b123")
	fmt.Println("********************")
	/*
		- the Ident rule takes precedence to the Int rule
		- here, the 45 is not Ident will be scan as Int
	*/
	checkTok(t, s, 1, s.Scan(), scanner.Int, "45")
	fmt.Println("********************")
	checkTok(t, s, 1, s.Scan(), scanner.Ident, "a12")
	fmt.Println("********************")
	checkTok(t, s, 1, s.Scan(), scanner.Ident, "b123")
	fmt.Println("********************")
	checkTok(t, s, 1, s.Scan(), scanner.Ident, "a12")
	fmt.Println("********************")
	checkTok(t, s, 1, s.Scan(), scanner.Int, "3")
	fmt.Println("********************")
	checkTok(t, s, 1, s.Scan(), scanner.Ident, "b")
	fmt.Println("********************")
	checkTok(t, s, 1, s.Scan(), scanner.EOF, "")
	fmt.Println("********************")
}

func TestScanNext(t *testing.T) {
	/*
		The Unicode character U+FEFF is the byte order mark, or BOM, and is used to tell the difference between big- and little-endian UTF-16 encoding.
	*/
	const BOM = '\uFEFF'
	BOMs := string(BOM)
	src := BOMs + "if a == bcd /* com" + BOMs + "ment */ {\n\ta += c\n}" + BOMs + "// line comment ending in eof"
	fmt.Println("====== src is", src)

	s := new(scanner.Scanner).Init(strings.NewReader(src))
	fmt.Println("********************")
	checkTok(t, s, 1, s.Scan(), scanner.Ident, "if") // the first BOM is ignored
	fmt.Println("********************")
	checkTok(t, s, 1, s.Scan(), scanner.Ident, "a")
	fmt.Println("********************")
	checkTok(t, s, 1, s.Scan(), '=', "=")
	fmt.Println("********************")
	/* we will call Next(), instead of Scan()
	- Next return the next unicode character but
	- It reset the Line to 0 and the TokenText to ""
	- Is Next related to BOM?
	*/
	checkTok(t, s, 0, s.Next(), '=', "")
	fmt.Println("********************")
	/* Next does not skip white space character */
	checkTok(t, s, 0, s.Next(), ' ', "")
	fmt.Println("********************")
	checkTok(t, s, 0, s.Next(), 'b', "")
	fmt.Println("********************")
	checkTok(t, s, 1, s.Scan(), scanner.Ident, "cd")
	fmt.Println("********************")
	/* Scan skip comments and BOMs inside comment */
	checkTok(t, s, 1, s.Scan(), '{', "{")
	fmt.Println("********************")
	checkTok(t, s, 2, s.Scan(), scanner.Ident, "a")
	fmt.Println("********************")
	checkTok(t, s, 2, s.Scan(), '+', "+")
	fmt.Println("********************")
	checkTok(t, s, 0, s.Next(), '=', "")
	fmt.Println("********************")
	checkTok(t, s, 2, s.Scan(), scanner.Ident, "c")
	fmt.Println("********************")
	checkTok(t, s, 3, s.Scan(), '}', "}")
	fmt.Println("********************")
	/* Scan does not skip the BOM, not the first, not in comment */
	checkTok(t, s, 3, s.Scan(), BOM, BOMs)
	fmt.Println("********************")
	/* -1 means EOF */
	checkTok(t, s, 3, s.Scan(), -1, "")
	fmt.Println("********************")
	if s.ErrorCount != 0 {
		t.Errorf("%d errors", s.ErrorCount)
	}
}

func TestScanWhitespace(t *testing.T) {
	var buf bytes.Buffer
	var ws uint64
	// start at 1, NUL character is not allowed
	/*
		character 0 is null
		character 1 is start heading
		character 32 is space
	*/
	for ch := byte(1); ch < ' '; ch++ {
		buf.WriteByte(ch)
		/* By OR we turn on bit at the position of that character
		Nul character is not turned on
		*/
		ws |= 1 << ch
	}
	/*ws is the list of white space characters in bitwise form
	all the character in the list will be skipped
	- check by: bit shift the character for 1 and use &
	*/

	fmt.Println("====== buf after the loop: ", buf.String())
	const orig = 'x'
	buf.WriteByte(orig)
	fmt.Println("====== buf after add orig", buf.String())
	/* writebyte('x') will not be visible in fmt.Print but writestring is ok
	However, the scanner can return it in text form
	*/
	s := new(scanner.Scanner).Init(&buf)
	s.Mode = 0
	s.Whitespace = ws
	tok := s.Scan()
	fmt.Println("====== tok is", tok)
	fmt.Println("====== text is", s.TokenText())
	if tok != orig {
		t.Errorf("tok = %s, want %s", scanner.TokenString(tok), scanner.TokenString(orig))
	}
}

func testError(t *testing.T, src, pos, msg string, tok rune) {
	s := new(scanner.Scanner).Init(strings.NewReader(src))
	errorCalled := false
	s.Error = func(s *scanner.Scanner, m string) {
		fmt.Println("====== s.Error is called, msg", msg)
		if !errorCalled {
			// only look at first error
			/* default pos filename is input */
			fmt.Println("===== s.Pos", s.Pos().String())
			if p := s.Pos().String(); p != pos {
				t.Errorf("pos = %q, want %q for %q", p, pos, src)
			}
			if m != msg {
				t.Errorf("msg = %q, want %q for %q", m, msg, src)
			}
			errorCalled = true
		}
	}
	/*
		- for non-ident, included NUL, it returns the rune of that character, 0 for NUL
		- Scan will call Error implicitly, but still able to return a value, still has the text*/
	tk := s.Scan()
	fmt.Println("===== token string", scanner.TokenString(tk))
	if tk != tok {
		t.Errorf("tok = %s, want %s for %q", scanner.TokenString(tk), scanner.TokenString(tok), src)
	}
	if !errorCalled {
		t.Errorf("error handler not called for %q", src)
	}
	fmt.Println("======= error count = ", s.ErrorCount)
	if s.ErrorCount == 0 {
		t.Errorf("count = %d, want > 0 for %q", s.ErrorCount, src)
	}
}

func TestError(t *testing.T) {
	/* fmt.Println of "\x00" will not show anything
	it is NUL anyway
	*/
	testError(t, "\x00", "<input>:1:1", "invalid character NUL", 0)
	testError(t, "\x80", "<input>:1:1", "invalid UTF-8 encoding", utf8.RuneError)
	testError(t, "\xff", "<input>:1:1", "invalid UTF-8 encoding", utf8.RuneError)
	testError(t, "a\x00", "<input>:1:2", "invalid character NUL", scanner.Ident)
}

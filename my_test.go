package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
	"text/scanner"
)

// A StringReader delivers its data one string segment at a time via Read.
type StringReader struct {
	data []string
	step int
}

func (r *StringReader) Read(p []byte) (n int, err error) {
	if r.step < len(r.data) {
		s := r.data[r.step]
		n = copy(p, s)
		r.step++
	} else {
		err = io.EOF
	}
	return
}

func readRuneSegments(t *testing.T, segments []string) {
	got := ""
	want := strings.Join(segments, "")
	s := new(scanner.Scanner).Init(&StringReader{data: segments})
	for {
		ch := s.Next()
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
	// {},
	// {""},
	// {"日", "本語"},
	// {"\u65e5", "\u672c", "\u8a9e"},
	// {"\U000065e5", " ", "\U0000672c", "\U00008a9e"},
	// {"\xe6", "\x97\xa5\xe6", "\x9c\xac\xe8\xaa\x9e"},
	{"Hello", ", ", "World", "!"},
	{"Hello", ", ", "", "World", "!"},
}

func TestNext(t *testing.T) {
	for _, s := range segmentList {
		readRuneSegments(t, s)
	}
}

type token struct {
	tok  rune
	text string
}

var f100 = "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"

var tokenList = []token{
	// {scanner.Comment, "// line comments"},
	// {scanner.Comment, "//"},
	// {scanner.Comment, "////"},
	// {scanner.Comment, "// comment"},
	// {scanner.Comment, "// /* comment */"},
	// {scanner.Comment, "// // comment //"},
	// {scanner.Comment, "//" + f100},

	// {scanner.Comment, "// general comments"},
	// {scanner.Comment, "/**/"},
	// {scanner.Comment, "/***/"},
	// {scanner.Comment, "/* comment */"},
	// {scanner.Comment, "/* // comment */"},
	// {scanner.Comment, "/* /* comment */"},
	// {scanner.Comment, "/*\n comment\n*/"},
	// {scanner.Comment, "/*" + f100 + "*/"},

	// {scanner.Comment, "// identifiers"},
	// {scanner.Ident, "a"},
	// {scanner.Ident, "a0"},
	{scanner.Ident, "foobar"},
	// {scanner.Ident, "abc123"},
	// {scanner.Ident, "LGTM"},
	// {scanner.Ident, "_"},
	// {scanner.Ident, "_abc123"},
	// {scanner.Ident, "abc123_"},
	// {scanner.Ident, "_abc_123_"},
	// {scanner.Ident, "_äöü"},
	// {scanner.Ident, "_本"},
	// {scanner.Ident, "äöü"},
	// {scanner.Ident, "本"},
	// {scanner.Ident, "a۰۱۸"},
	// {scanner.Ident, "foo६४"},
	// {scanner.Ident, "bar９８７６"},
	// {scanner.Ident, f100},

	// {scanner.Comment, "// decimal ints"},
	// {scanner.Int, "0"},
	// {scanner.Int, "1"},
	// {scanner.Int, "9"},
	// {scanner.Int, "42"},
	// {scanner.Int, "1234567890"},

	// {scanner.Comment, "// octal ints"},
	// {scanner.Int, "00"},
	// {scanner.Int, "01"},
	// {scanner.Int, "07"},
	// {scanner.Int, "042"},
	// {scanner.Int, "01234567"},

	// {scanner.Comment, "// hexadecimal ints"},
	// {scanner.Int, "0x0"},
	// {scanner.Int, "0x1"},
	// {scanner.Int, "0xf"},
	// {scanner.Int, "0x42"},
	// {scanner.Int, "0x123456789abcDEF"},
	// {scanner.Int, "0x" + f100},
	// {scanner.Int, "0X0"},
	// {scanner.Int, "0X1"},
	// {scanner.Int, "0XF"},
	// {scanner.Int, "0X42"},
	// {scanner.Int, "0X123456789abcDEF"},
	// {scanner.Int, "0X" + f100},

	// {scanner.Comment, "// floats"},
	// {scanner.Float, "0."},
	// {scanner.Float, "1."},
	// {scanner.Float, "42."},
	// {scanner.Float, "01234567890."},
	// {scanner.Float, ".0"},
	// {scanner.Float, ".1"},
	// {scanner.Float, ".42"},
	// {scanner.Float, ".0123456789"},
	// {scanner.Float, "0.0"},
	// {scanner.Float, "1.0"},
	// {scanner.Float, "42.0"},
	// {scanner.Float, "01234567890.0"},
	// {scanner.Float, "0e0"},
	// {scanner.Float, "1e0"},
	// {scanner.Float, "42e0"},
	// {scanner.Float, "01234567890e0"},
	// {scanner.Float, "0E0"},
	// {scanner.Float, "1E0"},
	// {scanner.Float, "42E0"},
	// {scanner.Float, "01234567890E0"},
	// {scanner.Float, "0e+10"},
	// {scanner.Float, "1e-10"},
	// {scanner.Float, "42e+10"},
	// {scanner.Float, "01234567890e-10"},
	// {scanner.Float, "0E+10"},
	// {scanner.Float, "1E-10"},
	// {scanner.Float, "42E+10"},
	// {scanner.Float, "01234567890E-10"},

	// {scanner.Comment, "// chars"},
	// {scanner.Char, `' '`},
	// {scanner.Char, `'a'`},
	// {scanner.Char, `'本'`},
	// {scanner.Char, `'\a'`},
	// {scanner.Char, `'\b'`},
	// {scanner.Char, `'\f'`},
	// {scanner.Char, `'\n'`},
	// {scanner.Char, `'\r'`},
	// {scanner.Char, `'\t'`},
	// {scanner.Char, `'\v'`},
	// {scanner.Char, `'\''`},
	// {scanner.Char, `'\000'`},
	// {scanner.Char, `'\777'`},
	// {scanner.Char, `'\x00'`},
	// {scanner.Char, `'\xff'`},
	// {scanner.Char, `'\u0000'`},
	// {scanner.Char, `'\ufA16'`},
	// {scanner.Char, `'\U00000000'`},
	// {scanner.Char, `'\U0000ffAB'`},

	// {scanner.Comment, "// strings"},
	// {scanner.String, `" "`},
	// {scanner.String, `"a"`},
	// {scanner.String, `"本"`},
	// {scanner.String, `"\a"`},
	// {scanner.String, `"\b"`},
	// {scanner.String, `"\f"`},
	// {scanner.String, `"\n"`},
	// {scanner.String, `"\r"`},
	// {scanner.String, `"\t"`},
	// {scanner.String, `"\v"`},
	// {scanner.String, `"\""`},
	// {scanner.String, `"\000"`},
	// {scanner.String, `"\777"`},
	// {scanner.String, `"\x00"`},
	// {scanner.String, `"\xff"`},
	// {scanner.String, `"\u0000"`},
	// {scanner.String, `"\ufA16"`},
	// {scanner.String, `"\U00000000"`},
	// {scanner.String, `"\U0000ffAB"`},
	// {scanner.String, `"` + f100 + `"`},

	// {scanner.Comment, "// raw strings"},
	// {scanner.RawString, "``"},
	// {scanner.RawString, "`\\`"},
	// {scanner.RawString, "`" + "\n\n/* foobar */\n\n" + "`"},
	// {scanner.RawString, "`" + f100 + "`"},

	// {scanner.Comment, "// individual characters"},
	// // NUL character is not allowed
	// {'\x01', "\x01"},
	// {' ' - 1, string(' ' - 1)},
	// {'+', "+"},
	// {'/', "/"},
	// {'.', "."},
	// {'~', "~"},
	// {'(', "("},
}

func makeSource(pattern string) *bytes.Buffer {
	var buf bytes.Buffer
	for _, k := range tokenList {
		fmt.Fprintf(&buf, pattern, k.text)
	}
	return &buf
}

func checkTok(t *testing.T, s *scanner.Scanner, line int, got, want rune, text string) {
	if got != want {
		t.Fatalf("tok = %s, want %s for %q", scanner.TokenString(got), scanner.TokenString(want), text)
	}
	if s.Line != line {
		t.Errorf("line = %d, want %d for %q", s.Line, line, text)
	}
	stext := s.TokenText()
	if stext != text {
		t.Errorf("text = %q, want %q", stext, text)
	} else {
		// check idempotency of TokenText() call
		stext = s.TokenText()
		if stext != text {
			t.Errorf("text = %q, want %q (idempotency check)", stext, text)
		}
	}
}

func checkTokErr(t *testing.T, s *scanner.Scanner, line int, want rune, text string) {
	prevCount := s.ErrorCount
	checkTok(t, s, line, s.Scan(), want, text)
	if s.ErrorCount != prevCount+1 {
		t.Fatalf("want error for %q", text)
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
	s := new(scanner.Scanner).Init(makeSource(" \t%s\n"))
	s.Mode = mode
	tok := s.Scan()
	line := 1
	for _, k := range tokenList {
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
	testScan(t, scanner.GoTokens&^scanner.SkipComments)
}

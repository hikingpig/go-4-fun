package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

var (
	n = flag.Bool("n", false, "omit trailing newline")
	s = flag.String("s", " ", "separator")
)
var out io.Writer = os.Stdout

/*
- n is a pointer, we access it value with *n and get its negated value with !*n
- flag.Args return all the args in cmd (those that not associate with a flag)
*/

func main() {
	flag.Parse()
	if err := echo(!*n, *s, flag.Args()); err != nil {
		fmt.Fprintf(os.Stderr, "echo: %v\n", err)
		os.Exit(1)
	}
}

func echo(newline bool, sep string, args []string) error {
	fmt.Fprint(out, strings.Join(args, sep))
	if newline {
		fmt.Fprintln(out)
	}
	return nil
}

/*
- we don't test the main func, only the echo func
- we never call flag at all
*/
func TestEcho(t *testing.T) {
	var tests = []struct {
		newline bool
		sep     string
		args    []string
		want    string
	}{
		{true, "", []string{}, "\n"},
		{false, "", []string{}, ""},
		{true, "\t", []string{"one", "two", "three"}, "one\ttwo\tthree\n"},
		{true, ",", []string{"a", "b", "c"}, "a,b,c\n"},
		{false, ":", []string{"1", "2", "3"}, "1:2:3"},
		{true, ",", []string{"a", "b", "c"}, "a b c\n"},
	}
	for _, test := range tests {
		descr := fmt.Sprintf("echo(%v, %q, %q)", test.newline, test.sep, test.args)
		/*
			- the original out is os.Stdout to print the result on screen
			- we replace it with bytes.Buffer to store the printed result
			- and later use that printed result to confirm test case
		*/
		out = new(bytes.Buffer)
		if err := echo(test.newline, test.sep, test.args); err != nil {
			t.Errorf("%s failed: %v", descr, err)
			continue
		}
		got := out.(*bytes.Buffer).String()
		if got != test.want {
			t.Errorf("%s = %q, want %q", descr, got, test.want)
		}
	}
}

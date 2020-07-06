package eval

import (
	"fmt"
	"strconv"
	"strings"
	"text/scanner"
)

type lexer struct {
	scan  scanner.Scanner
	token rune
}

func (lex *lexer) next()        { lex.token = lex.scan.Scan() }
func (lex *lexer) text() string { return lex.scan.TokenText() }

type lexPanic string

func (lex *lexer) describe() string {
	switch lex.token {
	case scanner.EOF:
		return "end of file"
	case scanner.Ident:
		return fmt.Sprintf("identifier %s", lex.text())
	case scanner.Int, scanner.Float:
		return fmt.Sprintf("number %s", lex.text())
	}
	return fmt.Sprintf("%q", rune(lex.token)) // any other rune
}

func precedence(op rune) int {
	switch op {
	case '*', '/':
		return 2
	case '+', '-':
		return 1
	}
	return 0
}

/*
2. call Parse with sqrt(A / pi)
	- lex after next is sqrt
*/
func Parse(input string) (_ Expr, err error) {
	defer func() {
		switch x := recover().(type) {
		case nil:
			// no panic
		case lexPanic:
			err = fmt.Errorf("%s", x)
		default:
			panic(x)
		}
	}()
	lex := new(lexer)
	lex.scan.Init(strings.NewReader(input))
	lex.scan.Mode = scanner.ScanIdents | scanner.ScanInts | scanner.ScanFloats
	lex.next()
	e := parseExpr(lex)
	if lex.token != scanner.EOF {
		return nil, fmt.Errorf("unexpected %s", lex.describe())
	}
	return e, nil
}

func parseExpr(lex *lexer) Expr { return parseBinary(lex, 1) }

/*
4. ParseBinary is called with A.
	- ParseUnary will call ParsePrimary and return string A
		- calling ParsePrimary again, A is Ident, id = A
		- lex.next() is "/"
		- lex.token != '(', return Var(A), which is string(A)
	- lex is now at /, precedence for / is 2. loop 2 times
		- only the 1st loop works, the 2nd doesnt do anything
		- op is /, lex is now pi
		- parseBinary is called again with pi, prec = 3. so the loop will be skipped
	- lhs now is binary{op, A, pi}
	- return it to 3
*/
func parseBinary(lex *lexer, prec1 int) Expr {
	lhs := parseUnary(lex)
	for prec := precedence(lex.token); prec >= prec1; prec-- {
		for precedence(lex.token) == prec {
			op := lex.token
			lex.next()
			rhs := parseBinary(lex, prec+1)
			lhs = binary{op, lhs, rhs}
		}
	}
	return lhs
}

func parseUnary(lex *lexer) Expr {
	if lex.token == '+' || lex.token == '-' {
		op := lex.token
		lex.next()
		return unary{op, parseUnary(lex)}
	}
	return parsePrimary(lex)
}

/*
3. ParsePrimary is called with lex at sqrt in sqrt(A /pi)
	- id = sqrt
	- lex.next is (
	- lex.next is A
	- inside for loop, parseExpr is called with A
	- go to parseBinary
	- parseBinary return binary{/, A, pi}
	- return the call{sqrt, {/, A, pi}}
*/
func parsePrimary(lex *lexer) Expr {
	switch lex.token {
	case scanner.Ident:
		id := lex.text()
		lex.next()
		if lex.token != '(' {
			return Var(id)
		}
		lex.next()
		var args []Expr
		if lex.token != ')' {
			for {
				args = append(args, parseExpr(lex))
				if lex.token != ',' {
					break
				}
				lex.next()
			}
			if lex.token != ')' {
				msg := fmt.Sprintf("got %s, want ')'", lex.describe())
				panic(lexPanic(msg))
			}
		}
		lex.next()
		return call{id, args}

	case scanner.Int, scanner.Float:
		f, err := strconv.ParseFloat(lex.text(), 64)
		if err != nil {
			panic(lexPanic(err.Error()))
		}
		lex.next()
		return literal(f)

	case '(':
		lex.next()
		e := parseExpr(lex)
		if lex.token != ')' {
			msg := fmt.Sprintf("got %s, want ')'", lex.describe())
			panic(lexPanic(msg))
		}
		lex.next()
		return e
	}
	msg := fmt.Sprintf("unexpected %s", lex.describe())
	panic(lexPanic(msg))
}

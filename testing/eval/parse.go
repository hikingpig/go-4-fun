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

/*
13. return 0, for ","
*/

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
	2. Inside Parse, input: "sqrt(1, 2)"
		- lexer is init with "sqrt(1, 2)", Mode: Idents(word), Int, Floats
		- lex.netx(), now text is "sqrt", the first word recognized by Idents
		- call parseExpr with lex

	20. continue from 2. receive call{sqrt, [1, 2]} from 19. finish Parse
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

/*
3. parseExpr is called with lex, "sqrt(1, 2), now text is "sqrt"

7. parseExpr is called again with lex now at 1 in sqrt(1, 2)

15. receive 1, lex at ",", go back to 6

19. continue from 3. receive call{sqrt, [1,2]} from 18. return to 2
*/
func parseExpr(lex *lexer) Expr { return parseBinary(lex, 1) }

/*
4. parseBinary is called with lex ("sqrt(1, 2)", at "sqrt")
	- call parseUnary

8. parseBinary is called again with lex now at 1 in sqrt(1, 2)
12. receive lhs=1 from 11. tex is now at , in sqrt(1, 2). call precedence of lex.token (","), receive 0
	- skip the loop
14. return lhs = 1, lex at ,. go back to 7

18. continue from 4. received call{id: sqrt, args: [1, 2]} from 17.
	- precedence == 0, skip for loop, go back to 3

*/
func parseBinary(lex *lexer, prec1 int) Expr {
	lhs := parseUnary(lex)
	for prec := precedence(lex.token); prec >= prec1; prec-- {
		for precedence(lex.token) == prec {
			fmt.Println("===== predence: ", precedence(lex.token), ",text:", lex.text(), ",token:", string(lex.token))
			op := lex.token
			lex.next()
			rhs := parseBinary(lex, prec+1)
			lhs = binary{op, lhs, rhs}
		}
	}
	return lhs
}

/*
5. parseUnary is called with lex "sqrt(1, 2)", at "sqrt"
	- lex.token is -2, for Idents "sqrt", call parsePrimary

9. parseUnary is calle again with lex at "1" in sqrt(1, 2)
11. receive float64(1) from 10, return to 8.

17. continue from 5: received "call{id: "sqrt", args: [1, 2]} from 16.
	- got back to 4

*/
func parseUnary(lex *lexer) Expr {
	if lex.token == '+' || lex.token == '-' {
		op := lex.token
		lex.next()
		return unary{op, parseUnary(lex)}
	}
	return parsePrimary(lex)
}

/*
6. parsePrimary called with lex "sqrt(1, 2)", at "sqrt"
	- lex.token = -2, Ident
		- id = "sqrt"
		- lex.next() is called
		- now at '(', token = '('
		- lex.next is called, now at 1 in sqrt(1, 2)
		- starting for loop
			- first loop with 1, call parseExpr

10. parsePrimary is called again with lex at "1", in sqrt(1, 2)
	- now Int, Float
		- lex.next is called, lex is now at ","
		- return float64(1), go back to 9

16. continue from 6. received literal 1 from parseExpr
	- lex.token == "," so we dont break
	- lex.next is called, now, lex at 2
	- similarly we got 2 literal from parseExpr
	- args now is [1, 2]
	- lex is now at ')' != ",", for loop break
	- call lex. next again. now lex at EOF
	- "call" is returned with id: sqrt, args: 1, 2
	- go back to 5
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

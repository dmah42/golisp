package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type typ string

const (
	TYPE_INT typ = "int"
	TYPE_FLOAT typ = "float"
	TYPE_SYMBOL typ = "symbol"
	TYPE_LIST typ = "list"
)

type object struct {
	t typ
	i int64
	f float64
	s string
	l []object
}

func (o object) String() string {
	switch o.t {
	case TYPE_INT:
		return fmt.Sprintf("{%s: %d}", o.t, o.i)
	case TYPE_FLOAT:
		return fmt.Sprintf("{%s: %f}", o.t, o.f)
	case TYPE_SYMBOL:
		return fmt.Sprintf("{%s: %s}", o.t, o.s)
	case TYPE_LIST:
		return fmt.Sprintf("%s", o.l)
	}
	return fmt.Sprintf("<unknown type>: %+v", o)
}

func removeEmpty(tokens []string) []string {
	b := tokens[:0]
	for _, t := range tokens {
		if len(t) !=  0 {
			b = append(b, t)
		}
	}
	return b
}

func tokenize(program string) []string {
	return removeEmpty(strings.Split(strings.Replace(strings.Replace(program, "(", " ( ", -1), ")", " ) ", -1), " "))
}

func atom(token string) (object, error) {
	if token == "" {
		return object{}, errors.New("unexpected empty token")
	}
	valInt, err := strconv.ParseInt(token, 10, 64)
	if err == nil {
		return object{t: TYPE_INT, i: valInt}, nil
	}
	valFloat, err := strconv.ParseFloat(token, 64)
	if err == nil {
		return object{t: TYPE_FLOAT, f: valFloat}, nil
	}
	return object{t: TYPE_SYMBOL, s: token}, nil
}

func lex(tokens []string) ([]string, object, error) {
	fmt.Printf("lex called with %#v\n", tokens)
	if len(tokens) == 0 {
		return nil, object{}, errors.New("unexpected EOF")
	}

	var token string
	token, tokens = tokens[0], tokens[1:]
	fmt.Printf("-- %s .. %#v\n", token, tokens)
	if token == "(" {
		l := object{t: TYPE_LIST, l: []object{}}
		for len(tokens) != 0 && tokens[0] != ")" {
			var ls object
			var err error
			tokens, ls, err = lex(tokens)
			if err != nil {
				// TODO: maybe return the tokens to this point?
				return nil, object{}, err
			}
			fmt.Printf("---- %#v\n", tokens)
			l.l = append(l.l, ls)
		}
		if len(tokens) == 0 {
			return nil, object{}, errors.New("unexpected EOF")
		}
		// Pop off the ")"
		_, tokens = tokens[0], tokens[1:]
		return tokens, l, nil
	}
	if token == ")" {
		// TODO: add the line/column to the error.
		return nil, object{}, errors.New("unexpected ')'")
	}
	a, err := atom(token)
	return tokens, a, err
}

// Parse a program into an AST.
func buildAST(program string) (object, error) {
	tokens := tokenize(program)
	fmt.Printf("tokens: %#v\n", tokens)

	tokens, ast, err := lex(tokens)
	if len(tokens) != 0 {
		return object{}, errors.New("unexpected leftover tokens")
	}
	return ast, err
}

func main() {
	program := os.Args[1]
	ast, err := buildAST(program)
	if err != nil {
		fmt.Printf("ERROR: %s while parsing %q\n", err, program)
		os.Exit(1)
	}

	fmt.Printf("ast: %+v\n", ast)
}

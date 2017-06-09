package golisp

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func Repl() error {
	for {
		fmt.Print("golisp> ")
		scanner := bufio.NewScanner(os.Stdin)
		comment := 0
		for scanner.Scan() {
			var res *object
			var err error
			in := scanner.Text()
			if in == "" || in[0] == ';' {
				goto prompt
			}
			if strings.HasPrefix(in, "#|") {
				comment++
				continue
			}
			if comment > 0 {
				if strings.HasPrefix(in, "|#") {
					comment--
				}

				if comment == 0 {
					goto prompt
				}
				continue
			}

			log.Printf("executing %q\n", in)
			res, err = Exec(in)
			if err != nil {
				fmt.Printf("ERROR: %s\n", err)
				goto prompt
			}
			if res != nil {
				fmt.Printf("%s\n", res.String())
			}
		prompt:
			fmt.Print("golisp> ")
		}
		return scanner.Err()
	}
}

func Exec(program string) (*object, error) {
	ast, err := buildAST(program)
	if err != nil {
		return nil, fmt.Errorf("%s while parsing %q\n", err, program)
	}

	log.Printf("ast: %+v\n", ast)
	return eval(globalEnv, ast)
}

func removeEmpty(tokens []string) []string {
	b := tokens[:0]
	for _, t := range tokens {
		if len(t) != 0 {
			b = append(b, t)
		}
	}
	return b
}

func tokenize(program string) []string {
	return removeEmpty(strings.Split(strings.Replace(strings.Replace(program, "(", " ( ", -1), ")", " ) ", -1), " "))
}

func atom(token string) (*object, error) {
	if token == "" {
		return nil, errors.New("unexpected empty token")
	}
	valInt, err := strconv.ParseInt(token, 10, 64)
	if err == nil {
		return newObject(valInt), nil
	}
	valFloat, err := strconv.ParseFloat(token, 64)
	if err == nil {
		return newObject(valFloat), nil
	}
	return newObject(token), nil
}

func lex(tokens []string) ([]string, *object, error) {
	log.Printf("lex called with %#v\n", tokens)
	if len(tokens) == 0 {
		return nil, nil, errors.New("unexpected EOF")
	}

	var token string
	token, tokens = tokens[0], tokens[1:]
	log.Printf("-- %s .. %#v\n", token, tokens)
	if token == "(" {
		l := newObject([]*object{})
		for len(tokens) != 0 && tokens[0] != ")" {
			var ls *object
			var err error
			tokens, ls, err = lex(tokens)
			if err != nil {
				// TODO: maybe return the tokens to this point?
				return nil, nil, err
			}
			log.Printf("---- %#v\n", tokens)
			l.l = append(l.l, ls)
		}
		if len(tokens) == 0 {
			return nil, nil, errors.New("unexpected EOF")
		}
		// Pop off the ")"
		_, tokens = tokens[0], tokens[1:]
		return tokens, l, nil
	}
	if token == ")" {
		// TODO: add the line/column to the error.
		return nil, nil, errors.New("unexpected ')'")
	}
	a, err := atom(token)
	return tokens, a, err
}

// Parse a program into an AST.
func buildAST(program string) (*object, error) {
	tokens := tokenize(program)
	log.Printf("tokens: %#v\n", tokens)

	tokens, ast, err := lex(tokens)
	if len(tokens) != 0 {
		return nil, errors.New("unexpected leftover tokens")
	}
	return ast, err
}

func eval(e env, o ...*object) (*object, error) {
	log.Printf("eval called with %+v\n", o)
	x := o[0]
	log.Printf("checking operand %+v\n", x)
	switch {
	case x.t == TYPE_SYMBOL:
		log.Printf("SYMBOL %q\n", x.s)
		v, err := e.get(x.s)
		if err != nil {
			return nil, err
		}
		return v, nil
	case x.t != TYPE_LIST:
		log.Printf("CONSTANT %v\n", x)
		return x, nil
	case x.l[0].t == TYPE_BUILTIN:
		log.Printf("BUILTIN %q\n", x.l[0].s)
		switch x.l[0].s {
		case "quote":
			return x.l[1], nil
		case "if":
			test, conseq, alt := x.l[1], x.l[2], x.l[3]
			// TODO: type check on everything
			res, err := eval(e, test)
			if err != nil {
				return nil, err
			}
			log.Printf("test result: %#v", res)
			exp := conseq
			if !res.isTruthy() {
				exp = alt
			}
			return eval(e, exp)
		case "define":
			v, exp := x.l[1], x.l[2]
			// TODO: type check on v
			ev, err := eval(e, exp)
			if err != nil {
				return nil, err
			}
			e.define(v.s, ev)
			return nil, err
		case "set!":
			v, exp := x.l[1], x.l[2]
			ev, err := eval(e, exp)
			if err != nil {
				return nil, err
			}
			e.set(v.s, ev)
			return nil, err
		case "lambda":
			params, body := x.l[1], x.l[2]
			l, err := newLambda(params, body, &e)
			if err != nil {
				return nil, err
			}
			return newObject(l), nil
		default:
			return nil, fmt.Errorf("unknown builtin: %q", x.s)
		}
	default:
		log.Printf("LIST %+v\n", x.l)
		if len(x.l) == 0 {
			log.Printf("returning empty\n")
			return nil, nil
		}
		proc, err := eval(e, x.l[0])
		log.Printf("-- got proc %+v\n", proc)
		if err != nil {
			return nil, err
		}
		// Evaluate the arguments.
		opargs := x.l[1:]
		args := make([]*object, len(opargs))
		for i := range opargs {
			args[i], err = eval(e, opargs[i])
			if err != nil {
				return nil, err
			}
		}
		switch proc.t {
		case TYPE_FN:
			return proc.fn(args...)
		case TYPE_LAMBDA:
			return proc.lambda.call(args...)
		default:
			return nil, errors.New("expected lambda or fn")
		}
	}
}


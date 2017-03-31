package main

import (
	"errors"
	"fmt"
	"log"
	"math"
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
	TYPE_FN typ = "fn"
)

type object struct {
	t typ
	i int64
	f float64
	s string
	l []object
	fn func(v ...object) (object, error)
}

func (o object) ToString() string {
	switch o.t {
	case TYPE_INT:
		return fmt.Sprintf("%d", o.i)
	case TYPE_FLOAT:
		return fmt.Sprintf("%f", o.f)
	case TYPE_SYMBOL:
		return fmt.Sprintf("%s", o.s)
	case TYPE_LIST:
		return fmt.Sprintf("%s", o.l)
	case TYPE_FN:
		return fmt.Sprintf("%s", o.fn)
	default:
		return "nil"
	}
}

type env map[string]object

var globalEnv env = env{
	"+": object{t: TYPE_FN, fn: func(o ...object) (object, error) {
		if len(o) != 2 {
			return object{}, errors.New("expected two arguments to +")
		}
		var a float64
		var ta typ
		switch o[0].t {
		case TYPE_FLOAT:
			a = o[0].f
			ta = TYPE_FLOAT
		case TYPE_INT:
			a = float64(o[0].i)
			ta = TYPE_INT
		default:
			return object{}, fmt.Errorf("cannot add object of type %s", o[0].t)
		}

		var b float64
		var tb typ
		switch o[1].t {
		case TYPE_FLOAT:
			b = o[1].f
			tb = TYPE_FLOAT
		case TYPE_INT:
			b = float64(o[1].i)
			tb = TYPE_INT
		default:
			return object{}, fmt.Errorf("cannot add object of type %s", o[1].t)
		}

		res := a + b
		if ta == TYPE_INT && tb == TYPE_INT {
			return object{t: TYPE_INT, i: int64(res)}, nil

		}
		return object{t: TYPE_FLOAT, f: res}, nil
	}},
	"sin": object{t: TYPE_FN, fn: func(o ...object) (object, error) {
		if len(o) != 1 {
			return object{}, errors.New("expected one argument to sin")
		}
		var f float64
		if o[0].t == TYPE_FLOAT {
			f = o[0].f
		} else if o[0].t == TYPE_INT {
			f = float64(o[0].i)
		} else {
			return object{}, errors.New("expected float or int argument to sin")
		}
		return object{t: TYPE_FLOAT, f: math.Sin(f)}, nil
	}},
	"list": object{t: TYPE_FN, fn: func(o ...object) (object, error) {
		if len(o) != 1 {
			return object{}, errors.New("expected one argument to cos")
		}
		return object{t: TYPE_LIST, l: []object{o[0]}}, nil
	}},
	"list?": object{t: TYPE_FN, fn: func(o ...object) (object, error) {
		if len(o) != 1 {
			return object{}, errors.New("expected one argument to list?")
		}
		if o[0].t == TYPE_LIST {
			return object{t: TYPE_INT, i: 1}, nil
		} else {
			return object{t: TYPE_INT, i: 0}, nil
		}
	}},
}

func eval(e env, o ...object) (object, error) {
	log.Printf("eval called with %+v\n", o)
	var err error
	x := o[0]
	log.Printf("checking operand %+v\n", x)
	switch {
	case x.t == TYPE_SYMBOL:
		fmt.Println("SYMBOL")
		switch o[0].s {
		case "if": 
			fmt.Println("-- IF")
			test, conseq, alt := o[1], o[2], o[3]
			// TODO: type check on everything
			res, err := eval(e, test)
			if err != nil {
				return object{}, err
			}
			exp := conseq
			if res.i == 0 {
				exp = alt
			}
			return eval(e, exp)
		case "define":
			fmt.Println("-- DEFINE")
			v, exp := o[1], o[2]
			// TODO: type check on v
			e[v.s], err = eval(e, exp)
			if err != nil {
				return object{}, err
			}
		default:
			fmt.Println("-- DEFAULT")
			log.Printf(".. returning symbol %q from env\n", x.s)
			return e[x.s], nil
		}
	case x.t != TYPE_LIST:
		log.Println("NOT LIST")
		log.Printf("-- returning %+v directly\n", x)
		return x, nil
	default:
		log.Println("LIST")
		log.Printf("+++ %+v\n", x.l)
		if len(x.l) == 0 {
			log.Printf("returning empty\n")
			return object{}, nil
		}
		log.Printf("-- getting proc for %+v\n", x.l[0])
		proc, err := eval(e, x.l[0])
		log.Printf("-- got proc %+v\n", proc)
		if err != nil {
			return object{}, err
		}
		if proc.t != TYPE_FN {
			return object{}, errors.New("expected function")
		}

		opargs := x.l[1:]
		args := make([]object, len(opargs))
		for i := range opargs {
			log.Printf("---- evaluating arg %d from %+v\n", i, opargs[i])
			args[i], err = eval(e, opargs[i])
			if err != nil {
				return object{}, err
			}
		}
		log.Printf("calling function %+v with args %+v", proc.fn, args)
		return proc.fn(args...)
	}

	return object{}, fmt.Errorf("unhandled case: %+v", o)
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
	log.Printf("lex called with %#v\n", tokens)
	if len(tokens) == 0 {
		return nil, object{}, errors.New("unexpected EOF")
	}

	var token string
	token, tokens = tokens[0], tokens[1:]
	log.Printf("-- %s .. %#v\n", token, tokens)
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
			log.Printf("---- %#v\n", tokens)
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
	log.Printf("tokens: %#v\n", tokens)

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
		log.Fatalf("%s while parsing %q\n", err, program)
	}

	log.Printf("ast: %+v\n", ast)

	res, err := eval(globalEnv, ast)
	log.Printf("%+v %+v\n", res, err)
	if err != nil {
		log.Fatalf("%s while evaling %+v\n", err, ast)
	}
	fmt.Printf("%s\n", res.ToString())
}

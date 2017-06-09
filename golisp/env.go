package golisp

import (
	"errors"
	"fmt"
	"log"
	"math"
	"reflect"
)

type env struct {
	outer *env
	m     map[string]*object
}

// TODO: test
// get returns the value of the key from the innermost scope.
func (e *env) get(key string) (*object, error) {
	ee, err := e.find(key)
	if err != nil {
		return nil, err
	}
	return ee.m[key], nil
}

// find returns the innermost scope that contains the key
func (e *env) find(key string) (*env, error) {
	_, ok := e.m[key]
	if ok {
		return e, nil
	}
	if e.outer != nil {
		return e.outer.find(key)
	}
	return nil, fmt.Errorf("%q not found", key)
}

// define creates a new key in the current scope.
func (e *env) define(key string, value *object) {
	e.m[key] = value
}

// set overrides the value of an existing key wherever it is in scope.
func (e *env) set(key string, value *object) error {
	ee, err := e.find(key)
	if err != nil {
		return err
	}
	ee.define(key, value)
	return nil
}

var globalEnv env = env{
	outer: nil,
	m: map[string]*object{
		// operators
		"+": newObject(func(o ...*object) (*object, error) {
			if len(o) == 1 {
				return nil, errors.New("expected at least two arguments to +")
			}

			var res float64
			var isint bool
			for _, v := range o {
				isint = isint && (v.t == TYPE_INT)
				f, err := v.toFloat()
				if err != nil {
					return nil, err
				}
				res += f
			}

			if isint {
				return newObject(int64(res)), nil
			}
			return newObject(res), nil
		}),
		"-": newObject(func(o ...*object) (*object, error) {
			if len(o) != 2 {
				return nil, errors.New("expected two arguments to -")
			}

			a, err := o[0].toFloat()
			if err != nil {
				return nil, err
			}

			b, err := o[1].toFloat()
			if err != nil {
				return nil, err
			}

			res := a - b
			if o[0].t == TYPE_INT && o[1].t == TYPE_INT {
				return newObject(int64(res)), nil

			}
			return newObject(res), nil
		}),
		"*": newObject(func(o ...*object) (*object, error) {
			if len(o) != 2 {
				return nil, errors.New("expected two arguments to *")
			}

			a, err := o[0].toFloat()
			if err != nil {
				return nil, err
			}

			b, err := o[1].toFloat()
			if err != nil {
				return nil, err
			}

			res := a * b
			if o[0].t == TYPE_INT && o[1].t == TYPE_INT {
				return newObject(int64(res)), nil

			}
			return newObject(res), nil
		}),
		"/": newObject(func(o ...*object) (*object, error) {
			if len(o) != 2 {
				return nil, errors.New("expected two arguments to /")
			}

			a, err := o[0].toFloat()
			if err != nil {
				return nil, err
			}

			b, err := o[1].toFloat()
			if err != nil {
				return nil, err
			}

			res := a / b
			if o[0].t == TYPE_INT && o[1].t == TYPE_INT {
				return newObject(int64(res)), nil

			}
			return newObject(res), nil
		}),
		">": newObject(func(o ...*object) (*object, error) {
			if len(o) != 2 {
				return nil, errors.New("expected two arguments to >")
			}

			a, err := o[0].toFloat()
			if err != nil {
				return nil, err
			}

			b, err := o[1].toFloat()
			if err != nil {
				return nil, err
			}

			return newObject(a > b), nil
		}),
		"<": newObject(func(o ...*object) (*object, error) {
			if len(o) != 2 {
				return nil, errors.New("expected two arguments to <")
			}

			a, err := o[0].toFloat()
			if err != nil {
				return nil, err
			}

			b, err := o[1].toFloat()
			if err != nil {
				return nil, err
			}

			return newObject(a < b), nil
		}),
		">=": newObject(func(o ...*object) (*object, error) {
			if len(o) != 2 {
				return nil, errors.New("expected two arguments to >=")
			}

			a, err := o[0].toFloat()
			if err != nil {
				return nil, err
			}

			b, err := o[1].toFloat()
			if err != nil {
				return nil, err
			}

			return newObject(a >= b), nil
		}),
		"<=": newObject(func(o ...*object) (*object, error) {
			if len(o) != 2 {
				return nil, errors.New("expected two arguments to <=")
			}

			a, err := o[0].toFloat()
			if err != nil {
				return nil, err
			}

			b, err := o[1].toFloat()
			if err != nil {
				return nil, err
			}

			return newObject(a <= b), nil
		}),
		"=": newObject(func(o ...*object) (*object, error) {
			if len(o) != 2 {
				return nil, errors.New("expected two arguments to =")
			}

			a, err := o[0].toFloat()
			if err != nil {
				return nil, err
			}

			b, err := o[1].toFloat()
			if err != nil {
				return nil, err
			}

			return newObject(a == b), nil
		}),

		// math
		"abs": newObject(func(o ...*object) (*object, error) {
			if len(o) != 1 {
				return nil, errors.New("expected one argument to abs")
			}
			if o[0].t == TYPE_FLOAT {
				return newObject(math.Abs(o[0].f)), nil
			} else if o[0].t == TYPE_INT {
				return newObject(math.Abs(float64(o[0].i))), nil
			}
			return nil, errors.New("expected float or int argument to abs")
		}),
		"pow": newObject(func(o ...*object) (*object, error) {
			f0, err := o[0].toFloat()
			if err != nil {
				return nil, err
			}
			f1, err := o[1].toFloat()
			if err != nil {
				return nil, err
			}

			return newObject(math.Pow(f0, f1)), nil
		}),
		"sqrt": newObject(func(o ...*object) (*object, error) {
			f, err := o[0].toFloat()
			if err != nil {
				return nil, err
			}
			return newObject(math.Sqrt(f)), nil
		}),
		"round": newObject(func(o ...*object) (*object, error) {
			f, err := o[0].toFloat()
			if err != nil {
				return nil, err
			}
			return newObject(math.Trunc(f)), nil
		}),
		"sin": newObject(func(o ...*object) (*object, error) {
			if len(o) != 1 {
				return nil, errors.New("expected one argument to sin")
			}
			if o[0].t == TYPE_FLOAT {
				return newObject(math.Sin(o[0].f)), nil
			} else if o[0].t == TYPE_INT {
				return newObject(math.Sin(float64(o[0].i))), nil
			}
			return nil, errors.New("expected float or int argument to sin")
		}),
		"cos": newObject(func(o ...*object) (*object, error) {
			if len(o) != 1 {
				return nil, errors.New("expected one argument to cos")
			}
			if o[0].t == TYPE_FLOAT {
				return newObject(math.Cos(o[0].f)), nil
			} else if o[0].t == TYPE_INT {
				return newObject(math.Cos(float64(o[0].i))), nil
			}
			return nil, errors.New("expected float or int argument to cos")
		}),
		"pi": newObject(math.Pi),

		"begin": newObject(func(o ...*object) (*object, error) {
			return newObject(o), nil
		}),

		// list manipulation
		"car": newObject(func(o ...*object) (*object, error) {
			if len(o) != 1 {
				return nil, errors.New("expected one argument to car")
			}
			x := o[0]
			if x.t != TYPE_LIST {
				return nil, errors.New("expected list as argument to car")
			}
			return x.l[0], nil
		}),
		"cdr": newObject(func(o ...*object) (*object, error) {
			if len(o) != 1 {
				return nil, errors.New("expected one argument to cdr")
			}
			x := o[0]
			if x.t != TYPE_LIST {
				return nil, errors.New("expected list as argument to cdr")
			}
			return newObject(x.l[1:]), nil
		}),
		"cons": newObject(func(o ...*object) (*object, error) {
			if len(o) != 2 {
				return nil, errors.New("expected two arguments to cons")
			}
			if o[1].t != TYPE_LIST {
				return nil, errors.New("expected list as second argument to cons")
			}
			l := append([]*object{o[0]}, o[1].l...)
			return newObject(l), nil
		}),
		"eq?": newObject(func(o ...*object) (*object, error) {
			return newObject(&o[0] == &o[1]), nil
		}),
		"equal?": newObject(func(o ...*object) (*object, error) {
			return newObject(reflect.DeepEqual(o[0], o[1])), nil
		}),
		"length": newObject(func(o ...*object) (*object, error) {
			if len(o) != 1 {
				return nil, errors.New("expected one argument to len")
			}
			if o[0].t != TYPE_LIST {
				return nil, errors.New("expected list as argument to len")
			}
			return newObject(len(o[0].l)), nil
		}),
		"list": newObject(func(o ...*object) (*object, error) {
			return newObject(o), nil
		}),
		"list?": newObject(func(o ...*object) (*object, error) {
			if len(o) != 1 {
				return nil, errors.New("expected one argument to list?")
			}
			return newObject(o[0].t == TYPE_LIST), nil
		}),
		"map": newObject(func(o ...*object) (*object, error) {
			fn := o[0]
			if fn.t != TYPE_FN && fn.t != TYPE_LAMBDA {
				return nil, errors.New("expected callable for first argument to map")
			}

			args := o[1]
			if args.t != TYPE_LIST {
				return nil, errors.New("expected list for second argument to map")
			}

			res := []*object{}
			for _, arg := range args.l {
				var r *object
				var err error
				log.Printf("mapping with arg %+v", arg)
				if fn.t == TYPE_FN {
					r, err = fn.fn(arg)
				} else if fn.t == TYPE_LAMBDA {
					r, err = fn.lambda.call(arg)
				}
				if err != nil {
					return nil, err
				}

				res = append(res, r)
			}
			return newObject(res), nil
		}),
		"null?": newObject(func(o ...*object) (*object, error) {
			return newObject(reflect.DeepEqual(o[0], nil)), nil
		}),
		"number?": newObject(func(o ...*object) (*object, error) {
			return newObject(o[0].t == TYPE_INT || o[0].t == TYPE_FLOAT), nil
		}),
		"procedure?": newObject(func(o ...*object) (*object, error) {
			if len(o) != 1 {
				return nil, errors.New("expected one argument to procedure?")
			}
			return newObject(o[0].t == TYPE_FN || o[0].t == TYPE_LAMBDA), nil
		}),
		"symbol?": newObject(func(o ...*object) (*object, error) {
			if len(o) != 1 {
				return nil, errors.New("expected one argument to symbol?")
			}
			return newObject(o[0].t == TYPE_SYMBOL), nil
		}),
	},
}

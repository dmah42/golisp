package main

import (
	"errors"
	"fmt"
	"math"
)

type env struct {
	outer *env
	m     map[string]object
}

// TODO: test
// get returns the value of the key from the innermost scope.
func (e *env) get(key string) (object, error) {
	ee, err := e.find(key)
	if err != nil {
		return object{}, err
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
func (e *env) define(key string, value object) {
	e.m[key] = value
}

// set overrides the value of an existing key wherever it is in scope.
func (e *env) set(key string, value object) error {
	ee, err := e.find(key)
	if err != nil {
		return err
	}
	ee.define(key, value)
	return nil
}

var globalEnv env = env{
	outer: nil,
	m: map[string]object{
		// operators
		"+": newObject(func(o ...object) (object, error) {
			if len(o) != 2 {
				return object{}, errors.New("expected two arguments to +")
			}

			a, err := o[0].toFloat()
			if err != nil {
				return object{}, err
			}

			b, err := o[1].toFloat()
			if err != nil {
				return object{}, err
			}

			res := a + b
			if o[0].t == TYPE_INT && o[1].t == TYPE_INT {
				return newObject(int64(res)), nil

			}
			return newObject(res), nil
		}),
		"-": newObject(func(o ...object) (object, error) {
			if len(o) != 2 {
				return object{}, errors.New("expected two arguments to -")
			}

			a, err := o[0].toFloat()
			if err != nil {
				return object{}, err
			}

			b, err := o[1].toFloat()
			if err != nil {
				return object{}, err
			}

			res := a - b
			if o[0].t == TYPE_INT && o[1].t == TYPE_INT {
				return newObject(int64(res)), nil

			}
			return newObject(res), nil
		}),
		"*": newObject(func(o ...object) (object, error) {
			if len(o) != 2 {
				return object{}, errors.New("expected two arguments to *")
			}

			a, err := o[0].toFloat()
			if err != nil {
				return object{}, err
			}

			b, err := o[1].toFloat()
			if err != nil {
				return object{}, err
			}

			res := a * b
			if o[0].t == TYPE_INT && o[1].t == TYPE_INT {
				return newObject(int64(res)), nil

			}
			return newObject(res), nil
		}),
		"/": newObject(func(o ...object) (object, error) {
			if len(o) != 2 {
				return object{}, errors.New("expected two arguments to /")
			}

			a, err := o[0].toFloat()
			if err != nil {
				return object{}, err
			}

			b, err := o[1].toFloat()
			if err != nil {
				return object{}, err
			}

			res := a / b
			if o[0].t == TYPE_INT && o[1].t == TYPE_INT {
				return newObject(int64(res)), nil

			}
			return newObject(res), nil
		}),
		">": newObject(func(o ...object) (object, error) {
			if len(o) != 2 {
				return object{}, errors.New("expected two arguments to >")
			}

			a, err := o[0].toFloat()
			if err != nil {
				return object{}, err
			}

			b, err := o[1].toFloat()
			if err != nil {
				return object{}, err
			}

			return newObject(a > b), nil
		}),
		"<": newObject(func(o ...object) (object, error) {
			if len(o) != 2 {
				return object{}, errors.New("expected two arguments to <")
			}

			a, err := o[0].toFloat()
			if err != nil {
				return object{}, err
			}

			b, err := o[1].toFloat()
			if err != nil {
				return object{}, err
			}

			return newObject(a < b), nil
		}),
		">=": newObject(func(o ...object) (object, error) {
			if len(o) != 2 {
				return object{}, errors.New("expected two arguments to >=")
			}

			a, err := o[0].toFloat()
			if err != nil {
				return object{}, err
			}

			b, err := o[1].toFloat()
			if err != nil {
				return object{}, err
			}

			return newObject(a >= b), nil
		}),
		"<=": newObject(func(o ...object) (object, error) {
			if len(o) != 2 {
				return object{}, errors.New("expected two arguments to <=")
			}

			a, err := o[0].toFloat()
			if err != nil {
				return object{}, err
			}

			b, err := o[1].toFloat()
			if err != nil {
				return object{}, err
			}

			return newObject(a <= b), nil
		}),
		"=": newObject(func(o ...object) (object, error) {
			if len(o) != 2 {
				return object{}, errors.New("expected two arguments to =")
			}

			a, err := o[0].toFloat()
			if err != nil {
				return object{}, err
			}

			b, err := o[1].toFloat()
			if err != nil {
				return object{}, err
			}

			return newObject(a == b), nil
		}),
		
		// math
		"abs": newObject(func(o ...object) (object, error) {
			if len(o) != 1 {
				return object{}, errors.New("expected one argument to abs")
			}
			if o[0].t == TYPE_FLOAT {
				return newObject(math.Abs(o[0].f)), nil
			} else if o[0].t == TYPE_INT {
				return newObject(math.Abs(float64(o[0].i))), nil
			}
			return object{}, errors.New("expected float or int argument to abs")
		}),
		"sin": newObject(func(o ...object) (object, error) {
			if len(o) != 1 {
				return object{}, errors.New("expected one argument to sin")
			}
			if o[0].t == TYPE_FLOAT {
				return newObject(math.Sin(o[0].f)), nil
			} else if o[0].t == TYPE_INT {
				return newObject(math.Sin(float64(o[0].i))), nil
			}
			return object{}, errors.New("expected float or int argument to sin")
		}),
		"cos": newObject(func(o ...object) (object, error) {
			if len(o) != 1 {
				return object{}, errors.New("expected one argument to cos")
			}
			if o[0].t == TYPE_FLOAT {
				return newObject(math.Cos(o[0].f)), nil
			} else if o[0].t == TYPE_INT {
				return newObject(math.Cos(float64(o[0].i))), nil
			}
			return object{}, errors.New("expected float or int argument to cos")
		}),
		"pi": newObject(math.Pi),

		// list manipulation
		"car":  newObject(func(o ...object) (object, error) {
			if len(o) != 1 {
				return object{}, errors.New("expected one argument to car")
			}
			x := o[0]
			if x.t != TYPE_LIST {
				return object{}, errors.New("expected list as argument to car")
			}
			return x.l[0], nil
		}),
		"cdr":  newObject(func(o ...object) (object, error) {
			if len(o) != 1 {
				return object{}, errors.New("expected one argument to cdr")
			}
			x := o[0]
			if x.t != TYPE_LIST {
				return object{}, errors.New("expected list as argument to cdr")
			}
			return newObject(x.l[1:]), nil
		}),
		"cons":  newObject(func(o ...object) (object, error) {
			if len(o) != 2 {
				return object{}, errors.New("expected two arguments to cons")
			}
			x := o[0]
			y := o[1]
			if y.t != TYPE_LIST {
				return object{}, errors.New("expected list as second argument to cons")
			}
			return newObject(append([]object{x}, y.l...)), nil
		}),
		"len": newObject(func(o ...object) (object, error) {
			if len(o) != 1 {
				return object{}, errors.New("expected one argument to len")
			}
			if o[0].t != TYPE_LIST {
				return object{}, errors.New("expected list as argument to len")
			}
			return newObject(len(o[0].l)), nil
		}),
		"list": newObject(func(o ...object) (object, error) {
			if len(o) != 1 {
				return object{}, errors.New("expected one argument to list")
			}
			return newObject([]object{o[0]}), nil
		}),
		"list?": newObject(func(o ...object) (object, error) {
			if len(o) != 1 {
				return object{}, errors.New("expected one argument to list?")
			}
			return newObject(o[0].t == TYPE_LIST), nil
		}),
		"procedure?": newObject(func(o ...object) (object, error) {
			if len(o) != 1 {
				return object{}, errors.New("expected one argument to procedure?")
			}
			return newObject(o[0].t == TYPE_FN || o[0].t == TYPE_LAMBDA), nil
		}),
		"symbol?": newObject(func(o ...object) (object, error) {
			if len(o) != 1 {
				return object{}, errors.New("expected one argument to symbol?")
			}
			return newObject(o[0].t == TYPE_SYMBOL), nil
		}),
	}}

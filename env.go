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

			var res int
			if a > b {
				res = 1
			}
			return newObject(res), nil
		}),
		">=": newObject(func(o ...object) (object, error) {
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

			var res int
			if a >= b {
				res = 1
			}
			return newObject(res), nil
		}),
		"<=": newObject(func(o ...object) (object, error) {
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

			var res int
			if a <= b {
				res = 1
			}
			return newObject(res), nil
		}),
		"sin": newObject(func(o ...object) (object, error) {
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
			return newObject(math.Sin(f)), nil
		}),
		"cos": newObject(func(o ...object) (object, error) {
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
			return newObject(math.Cos(f)), nil
		}),
		"pi": newObject(math.Pi),
		"list": newObject(func(o ...object) (object, error) {
			if len(o) != 1 {
				return object{}, errors.New("expected one argument to cos")
			}
			return newObject([]object{o[0]}), nil
		}),
		"list?": newObject(func(o ...object) (object, error) {
			if len(o) != 1 {
				return object{}, errors.New("expected one argument to list?")
			}
			if o[0].t == TYPE_LIST {
				return newObject(1), nil
			}
			return newObject(0), nil
		}),
	}}

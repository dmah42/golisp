package main

import (
	"fmt"
)

type typ string

const (
	TYPE_INT typ = "int"
	TYPE_FLOAT typ = "float"
	TYPE_SYMBOL typ = "symbol"
	TYPE_LIST typ = "list"
	TYPE_FN typ = "fn"
	TYPE_BUILTIN typ = "builtin"
)

var builtins = []string {
	"define",
	"if",
}

type object struct {
	t typ
	i int64
	f float64
	s string
	l []object
	fn func(...object) (object, error)
}

func isBuiltin(s string) bool {
	for _, b := range builtins {
		if b == s {
			return true
		}
	}
	return false
}

func newObject(v interface{}) object {
	switch v.(type) {
	case float64:
		return object{t: TYPE_FLOAT, f: v.(float64)}
	case float32:
		return object{t: TYPE_FLOAT, f: float64(v.(float32))}
	case int64:
		return object{t: TYPE_INT, i: v.(int64)}
	case int32:
		return object{t: TYPE_INT, i: int64(v.(int32))}
	case int:
		return object{t: TYPE_INT, i: int64(v.(int))}
	case string:
		if isBuiltin(v.(string)) {
			return object{t: TYPE_BUILTIN, s: v.(string)}
		}
		return object{t: TYPE_SYMBOL, s: v.(string)}
	case []object:
		return object{t: TYPE_LIST, l: v.([]object)}
	case func(...object) (object, error):
		return object{t: TYPE_FN, fn: v.(func(...object) (object, error))}
	default:
		return object{}
	}
}

func (o object) toFloat() (float64, error) {
	switch o.t {
	case TYPE_FLOAT:
		return o.f, nil
	case TYPE_INT:
		return float64(o.i), nil
	default:
		return 0.0, fmt.Errorf("cannot convert %q to float", o.t)
	}
}

func (o object) toString() string {
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


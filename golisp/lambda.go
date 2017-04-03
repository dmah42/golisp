package golisp

import (
	"errors"
	"fmt"
)

type lambda struct {
	params *object
	body   *object
	outer  *env
}

func newLambda(params, body *object, env *env) (*lambda, error) {
	if params == nil || body == nil {
		return nil, errors.New("nil params or body")
	}
	if params.t != TYPE_LIST {
		return nil, errors.New("invalid params. expected list.")
	}
	for _, p := range params.l {
		if p.t != TYPE_SYMBOL {
			return nil, fmt.Errorf("unexpected non-symbolic param: %s", p)
		}
	}
	return &lambda{params, body, env}, nil
}

func (l *lambda) call(args ...*object) (*object, error) {
	if len(args) != len(l.params.l) {
		return nil, fmt.Errorf("mismatch number of args %d to params %d.", len(args), len(l.params.l))
	}

	e := env{
		outer: l.outer,
		m:     map[string]*object{},
	}

	for i, _ := range l.params.l {
		e.define(l.params.l[i].s, args[i])
	}

	return eval(e, l.body)
}

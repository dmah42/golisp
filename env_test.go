package main

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {
	e := &env{
		outer: &env{
			m: map[string]*object{
				"foo": newObject("foo"),
				"bar": newObject("bar"),
			},
		},
		m: map[string]*object{
			"bar": newObject("baz"),
		},
	}

	cases := []struct {
		key     string
		want    *object
		wantErr error
	}{
		{
			key:  "bar",
			want: newObject("baz"),
		},
		{
			key:  "foo",
			want: newObject("foo"),
		},
		{
			key:     "baz",
			wantErr: fmt.Errorf("%q not found", "baz"),
		},
	}

	for _, tt := range cases {
		got, err := e.get(tt.key)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("got %+v, want %+v", got, tt.want)
		}
		if !reflect.DeepEqual(err, tt.wantErr) {
			t.Errorf("got err %q, want err %q", err, tt.wantErr)
		}
	}
}

func TestFind(t *testing.T) {
	outer := &env{
		m: map[string]*object{
			"foo": newObject("foo"),
			"bar": newObject("bar"),
		},
	}
	e := &env{
		outer: outer,
		m: map[string]*object{
			"bar": newObject("baz"),
		},
	}

	cases := []struct {
		key     string
		want    *env
		wantErr error
	}{
		{
			key:  "bar",
			want: e,
		},
		{
			key:  "foo",
			want: outer,
		},
		{
			key:     "baz",
			wantErr: fmt.Errorf("%q not found", "baz"),
		},
	}

	for _, tt := range cases {
		got, err := e.find(tt.key)
		if got != tt.want {
			t.Errorf("got %+v, want %+v", got, tt.want)
		}
		if !reflect.DeepEqual(err, tt.wantErr) {
			t.Errorf("got err %q, want err %q", err, tt.wantErr)
		}
	}
}

func TestGlobalEnv(t *testing.T) {
	cases := []struct {
		key     string
		args    []*object
		want    *object
		wantErr error
	}{
		{
			key:  "+",
			args: []*object{newObject(4), newObject(2)},
			want: newObject(6),
		},
		{
			key:     "+",
			args:    []*object{newObject(4)},
			wantErr: errors.New("expected two arguments to +"),
		},
		{
			key:  "-",
			args: []*object{newObject(4), newObject(2)},
			want: newObject(2),
		},
		{
			key:     "-",
			args:    []*object{newObject(4)},
			wantErr: errors.New("expected two arguments to -"),
		},
		{
			key:  "*",
			args: []*object{newObject(4), newObject(2)},
			want: newObject(8),
		},
		{
			key:     "*",
			args:    []*object{newObject(4)},
			wantErr: errors.New("expected two arguments to *"),
		},
		{
			key:  "/",
			args: []*object{newObject(4), newObject(2)},
			want: newObject(2),
		},
		{
			key:     "/",
			args:    []*object{newObject(4)},
			wantErr: errors.New("expected two arguments to /"),
		},
		{
			key:  ">",
			args: []*object{newObject(4), newObject(2)},
			want: newObject(true),
		},
		{
			key:  ">",
			args: []*object{newObject(2), newObject(4)},
			want: newObject(false),
		},
		{
			key:     ">",
			args:    []*object{newObject(4)},
			wantErr: errors.New("expected two arguments to >"),
		},
		{
			key:  ">=",
			args: []*object{newObject(4), newObject(2)},
			want: newObject(true),
		},
		{
			key:  ">=",
			args: []*object{newObject(4), newObject(4)},
			want: newObject(true),
		},
		{
			key:  ">=",
			args: []*object{newObject(2), newObject(4)},
			want: newObject(false),
		},
		{
			key:     ">=",
			args:    []*object{newObject(4)},
			wantErr: errors.New("expected two arguments to >="),
		},
		{
			key:  "<",
			args: []*object{newObject(4), newObject(2)},
			want: newObject(false),
		},
		{
			key:  "<",
			args: []*object{newObject(2), newObject(4)},
			want: newObject(true),
		},
		{
			key:     "<",
			args:    []*object{newObject(4)},
			wantErr: errors.New("expected two arguments to <"),
		},
		{
			key:  "<=",
			args: []*object{newObject(4), newObject(2)},
			want: newObject(false),
		},
		{
			key:  "<=",
			args: []*object{newObject(4), newObject(4)},
			want: newObject(true),
		},
		{
			key:  "<=",
			args: []*object{newObject(2), newObject(4)},
			want: newObject(true),
		},
		{
			key:     "<=",
			args:    []*object{newObject(4)},
			wantErr: errors.New("expected two arguments to <="),
		},
		{
			key:  "=",
			args: []*object{newObject(4), newObject(2)},
			want: newObject(false),
		},
		{
			key:  "=",
			args: []*object{newObject(4), newObject(4)},
			want: newObject(true),
		},
		{
			key:     "=",
			args:    []*object{newObject(4)},
			wantErr: errors.New("expected two arguments to ="),
		},
		{
			key:  "abs",
			args: []*object{newObject(-42)},
			want: newObject(42.0),
		},
		{
			key:  "abs",
			args: []*object{newObject(-42.0)},
			want: newObject(42.0),
		},
		{
			key:     "abs",
			args:    []*object{newObject("42")},
			wantErr: errors.New("expected float or int argument to abs"),
		},
		{
			key:     "abs",
			args:    []*object{newObject(4), newObject(2)},
			wantErr: errors.New("expected one argument to abs"),
		},
		{
			key:  "sin",
			args: []*object{newObject(42)},
			want: newObject(math.Sin(42.0)),
		},
		{
			key:  "sin",
			args: []*object{newObject(42.0)},
			want: newObject(math.Sin(42.0)),
		},
		{
			key:     "sin",
			args:    []*object{newObject("42")},
			wantErr: errors.New("expected float or int argument to sin"),
		},
		{
			key:     "sin",
			args:    []*object{newObject(4), newObject(2)},
			wantErr: errors.New("expected one argument to sin"),
		},
		{
			key:  "cos",
			args: []*object{newObject(42)},
			want: newObject(math.Cos(42.0)),
		},
		{
			key:  "cos",
			args: []*object{newObject(42.0)},
			want: newObject(math.Cos(42.0)),
		},
		{
			key:     "cos",
			args:    []*object{newObject("42")},
			wantErr: errors.New("expected float or int argument to cos"),
		},
		{
			key:     "cos",
			args:    []*object{newObject(4), newObject(2)},
			wantErr: errors.New("expected one argument to cos"),
		},
		{
			key:  "pi",
			want: newObject(math.Pi),
		},
		{
			key:  "car",
			args: []*object{newObject([]*object{newObject("foo"), newObject("bar")})},
			want: newObject("foo"),
		},
		{
			key:     "car",
			args:    []*object{newObject(4), newObject(2)},
			wantErr: errors.New("expected one argument to car"),
		},
		{
			key:     "car",
			args:    []*object{newObject(4)},
			wantErr: errors.New("expected list as argument to car"),
		},
		{
			key:  "cdr",
			args: []*object{newObject([]*object{newObject("foo"), newObject("bar")})},
			want: newObject([]*object{newObject("bar")}),
		},
		{
			key:     "cdr",
			args:    []*object{newObject(4), newObject(2)},
			wantErr: errors.New("expected one argument to cdr"),
		},
		{
			key:     "cdr",
			args:    []*object{newObject(4)},
			wantErr: errors.New("expected list as argument to cdr"),
		},
		{
			key: "cons",
			args: []*object{
				newObject("baz"),
				newObject([]*object{newObject("foo"), newObject("bar")}),
			},
			want: newObject([]*object{newObject("baz"), newObject("foo"), newObject("bar")}),
		},
		{
			key: "cons",
			args: []*object{
				newObject(42),
				newObject([]*object{newObject("foo"), newObject("bar")}),
			},
			want: newObject([]*object{newObject(42), newObject("foo"), newObject("bar")}),
		},
		{
			key:     "cons",
			args:    []*object{newObject(4)},
			wantErr: errors.New("expected two arguments to cons"),
		},
		{
			key: "length",
			args: []*object{
				newObject([]*object{newObject("foo"), newObject("bar")}),
			},
			want: newObject(2),
		},
		{
			key:     "length",
			args:    []*object{},
			wantErr: errors.New("expected one argument to len"),
		},
		{
			key:     "length",
			args:    []*object{newObject("baz")},
			wantErr: errors.New("expected list as argument to len"),
		},
		{
			key:  "list",
			args: []*object{newObject(42)},
			want: newObject([]*object{newObject(42)}),
		},
		{
			key:  "list",
			args: []*object{newObject(42), newObject("foo"), newObject(64)},
			want: newObject([]*object{newObject(42), newObject("foo"), newObject(64)}),
		},
		{
			key:     "list?",
			args:    []*object{},
			wantErr: errors.New("expected one argument to list?"),
		},
		{
			key:  "list?",
			args: []*object{newObject(42)},
			want: newObject(false),
		},
		{
			key:  "list?",
			args: []*object{newObject([]*object{newObject(42), newObject("foo")})},
			want: newObject(true),
		},
		{
			key: "map",
			args: []*object{
				newObject(func(o ...*object) (*object, error) {
					return newObject(o[0].i * 2), nil
				}),
				newObject([]*object{
					newObject(0), newObject(1), newObject(2),
				}),
			},
			want: newObject([]*object{
				newObject(0), newObject(2), newObject(4),
			}),
		},
		{
			key:     "procedure?",
			args:    []*object{},
			wantErr: errors.New("expected one argument to procedure?"),
		},
		{
			key:  "procedure?",
			args: []*object{newObject(42)},
			want: newObject(false),
		},
		{
			key:  "procedure?",
			args: []*object{newObject(func(o ...*object) (*object, error) { return nil, nil })},
			want: newObject(true),
		},
		{
			key:  "procedure?",
			args: []*object{newObject(&lambda{newObject(42), newObject(64), nil})},
			want: newObject(true),
		},
		{
			key:     "symbol?",
			args:    []*object{},
			wantErr: errors.New("expected one argument to symbol?"),
		},
		{
			key:  "symbol?",
			args: []*object{newObject("foo")},
			want: newObject(true),
		},
		{
			key:  "symbol?",
			args: []*object{newObject(42)},
			want: newObject(false),
		},
	}

	for _, tt := range cases {
		o, ok := globalEnv.m[tt.key]
		if !ok {
			t.Fatalf("key %q not found", tt.key)
		}

		var got *object
		var err error

		switch o.t {
		case TYPE_FN:
			got, err = o.fn(tt.args...)
		case TYPE_LAMBDA:
			got, err = o.lambda.call(tt.args...)
		default:
			got = o
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%s: got %+v, want %+v", tt.key, got, tt.want)
		}
		if !reflect.DeepEqual(err, tt.wantErr) {
			t.Errorf("%s: got err %q, want err %q", tt.key, err, tt.wantErr)
		}
	}
}

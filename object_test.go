package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewObject(t *testing.T) {
	cases := []struct {
		v interface{}
		want object
	}{
		{
			v: false,
			want: object{t: TYPE_INT, i: 0},
		},
		{
			v: true,
			want: object{t: TYPE_INT, i: 1},
		},
		{
			v: 42,
			want: object{t: TYPE_INT, i: 42},
		},
		{
			v: int32(42),
			want: object{t: TYPE_INT, i: 42},
		},
		{
			v: int64(42),
			want: object{t: TYPE_INT, i: 42},
		},
		{
			v: 42.0,
			want: object{t: TYPE_FLOAT, f: 42.0},
		},
		{
			v: float32(42.0),
			want: object{t: TYPE_FLOAT, f: 42.0},
		},
		{
			v: float64(42.0),
			want: object{t: TYPE_FLOAT, f: 42.0},
		},
		{
			v: "foo",
			want: object{t: TYPE_SYMBOL, s: "foo"},
		},
		{
			v: "define",
			want: object{t: TYPE_BUILTIN, s: "define"},
		},
		{
			v: []object{newObject(42), newObject("foo")},
			want: object{t: TYPE_LIST, l: []object{newObject(42), newObject("foo")}},
		},
	}

	for _, tt := range cases {
		if got := newObject(tt.v); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("got %+v, want %+v", got, tt.want)
		}
	}

}

func TestToFloat(t *testing.T) {
	cases := []struct {
		o object
		want float64
		wantErr error
	}{
		{
			o: newObject(nil),
			wantErr: fmt.Errorf("cannot convert %q to float", ""),
		},
		{
			o: newObject(42.0),
			want: 42.0,
		},
		{
			o: newObject(42),
			want: 42.0,
		},
		{
			o: newObject("42"),
			wantErr: fmt.Errorf("cannot convert %q to float", "symbol"),
		},
	}

	for _, tt := range cases {
		got, err := tt.o.toFloat()
		if got != tt.want {
			t.Errorf("got %f, want %f", got, tt.want)
		}
		if !reflect.DeepEqual(err, tt.wantErr) {
			t.Errorf("got err %q, want err %q", err, tt.wantErr)
		}
	}
}

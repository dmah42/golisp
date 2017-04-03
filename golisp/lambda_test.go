package golisp

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestNewLambda(t *testing.T) {
	cases := []struct {
		params, body *object
		env          *env
		want         *lambda
		wantErr      error
	}{
		{
			wantErr: errors.New("nil params or body"),
		},
		{
			params:  newObject("foo"),
			body: newObject(42),
			wantErr: errors.New("invalid params. expected list."),
		},
		{
			params:  newObject([]*object{newObject(42)}),
			body: newObject(42),
			wantErr: fmt.Errorf("unexpected non-symbolic param: %s", "42"),
		},
		{
			params: newObject([]*object{newObject("foo")}),
			body:   newObject(42),
			env:    &globalEnv,
			want: &lambda{
				newObject([]*object{newObject("foo")}),
				newObject(42),
				&globalEnv,
			},
		},
	}

	for _, tt := range cases {
		got, err := newLambda(tt.params, tt.body, tt.env)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("got %+v, want %+v", got, tt.want)
		}
		if !reflect.DeepEqual(err, tt.wantErr) {
			t.Errorf("got err %q, want err %q", err, tt.wantErr)
		}
	}
}

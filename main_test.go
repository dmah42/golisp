package main 

import (
	"errors"
	"reflect"
	"testing"
)

func TestRemoveEmpty(t *testing.T) {
	cases := []struct{
		tokens, want []string
	}{
		{},
		{[]string{}, []string{}},
		{[]string{"", "foo", "", "bar", ""}, []string{"foo", "bar"}},
		{[]string{"foo", "bar"}, []string{"foo", "bar"}},
	}

	for _, tt := range cases {
		got := removeEmpty(tt.tokens)
		if !reflect.DeepEqual(got, tt.want)  {
			t.Errorf("got %#v, want %#v", got, tt.want)
		}
	}
}

func TestTokenize(t *testing.T) {
	cases := []struct{
		program string
		want []string
	}{
		{
			want: []string{},
		},
		{
			program: "",
			want: []string{},
		},
		{
			program: "(begin (* pi (* r r)))",
			want: []string{"(", "begin", "(", "*", "pi", "(", "*", "r", "r", ")", ")", ")"},
		},
	}

	for _, tt := range cases {
		got := tokenize(tt.program)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("got %#v, want %#v", got, tt.want)
		}
	}
}

func TestAtom(t *testing.T) {
	cases := []struct {
		token string
		want object
		wantErr error
	}{
		{token: "", want: object{}, wantErr: errors.New("unexpected empty token")},
		{token: "42", want: object{t: TYPE_INT, i: 42}},
		{token: "42.3", want: object{t: TYPE_FLOAT, f: 42.3}},
		{token: "answer", want: object{t: TYPE_SYMBOL, s: "answer"}},
	}

	for _, tt := range cases {
		got, err := atom(tt.token)
		if !reflect.DeepEqual(err, tt.wantErr) {
			t.Errorf("got err %q, want err %q", err, tt.wantErr)
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("got %#v, want %#v", got, tt.want)
		}
	}
}

func TestLex(t *testing.T) {
	cases := []struct {
		name string
		tokens []string
		want object
		wantErr error
	}{
		{
			name: "eof",
			tokens: []string{},
			wantErr: errors.New("unexpected EOF"),
		},
		{
			name: "int",
			tokens: []string{"42"},
			want: object{t: TYPE_INT, i: 42},
		},
		{
			name: "unexpected ')'",
			tokens: []string{")"},
			wantErr: errors.New("unexpected ')'"),
		},
		{
			name: "unexpected EOF 2",
			tokens:  []string{"(", "begin", "(", "define", "r", "10", ")", "(", "*", "pi", "(", "*", "r", "r", ")", ")"},
			wantErr: errors.New("unexpected EOF"),
		},
		{
			name: "full",
			tokens: []string{"(","begin","(","define","r","10",")","r",")"},
			want: object{
				t: TYPE_LIST, l: []object{
					{t: TYPE_SYMBOL, s: "begin"},
					{t: TYPE_LIST, l: []object{
						{t: TYPE_SYMBOL, s: "define"},
						{t: TYPE_SYMBOL, s: "r"},
						{t: TYPE_INT, i: 10},
					}},
					{t: TYPE_SYMBOL, s: "r"},
				},
			},
		},
	}

	for _, tt := range cases {
		tokens, got, err := lex(tt.tokens)
		if len(tokens) != 0 {
			t.Fatalf("%s: unexpected extra tokens", tt.name)
		}

		if !reflect.DeepEqual(err, tt.wantErr) {
			t.Errorf("%s: got err %q, want err %q", tt.name, err, tt.wantErr)
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%s: got %#v, want %#v", tt.name, got, tt.want)
		}
	}
}

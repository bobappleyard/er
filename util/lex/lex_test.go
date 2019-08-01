package lex

import (
	"testing"
)

func TestParse(t *testing.T) {
	l, err := New(`a+`, `(ab)+`)
	if err != nil {
		t.Fatal(err)
	}
	bs := []byte("ababaa")
	toks := []Token{
		{2, 0, 4},
		{1, 4, 6},
	}
	var tok Token
	for _, exp := range toks {
		tok, err = l.Match(bs, tok)
		if err != nil {
			t.Fatal(err)
		}
		if tok != exp {
			t.Errorf("expecting %v, got %v", exp, tok)
		}
	}
	if tok.End < len(bs) {
		t.Error("didn't match whole input")
	}
}

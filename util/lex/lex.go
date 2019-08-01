package lex

import (
	"errors"
	"regexp"
	"regexp/syntax"
)

var ErrFailedToMatch = errors.New("failed to match")

type Token struct {
	ID         int
	Start, End int
}

type Lexer struct {
	re *regexp.Regexp
}

func New(toks ...string) (*Lexer, error) {
	alt := syntax.Regexp{
		Op: syntax.OpAlternate,
	}
	for i, tok := range toks {
		tokre, err := syntax.Parse(tok, syntax.Perl)
		if err != nil {
			return nil, err
		}
		decapture(tokre)
		alt.Sub = append(alt.Sub, &syntax.Regexp{
			Op:  syntax.OpCapture,
			Cap: i + 1,
			Sub: []*syntax.Regexp{tokre},
		})
	}
	res, err := regexp.Compile(alt.String())
	if err != nil {
		return nil, err
	}
	res.Longest()
	return &Lexer{res}, nil
}

func (l *Lexer) Match(bs []byte, after Token) (Token, error) {
	match := l.re.FindSubmatchIndex(bs[after.End:])
	if match == nil || match[0] > 0 {
		return Token{}, ErrFailedToMatch
	}
	for i := 2; i < len(match); i += 2 {
		if match[i] >= 0 {
			return Token{
				ID:    i / 2,
				Start: match[i] + after.End,
				End:   match[i+1] + after.End,
			}, nil
		}
	}
	panic("unreachable")
}

func (t Token) Text(bs []byte) string {
	return string(bs[t.Start:t.End])
}

func decapture(re *syntax.Regexp) {
	if re.Op == syntax.OpCapture {
		re.Op = syntax.OpConcat
	}
	for _, sub := range re.Sub {
		decapture(sub)
	}
}

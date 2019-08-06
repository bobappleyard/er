package path

import (
	"github.com/bobappleyard/er/util/lex"
	"github.com/pkg/errors"
)

var ErrFailedToMatch = errors.New("failed to parse path")

func Parse(bs []byte) (Path, error) {
	p := &parser{bs: bs}
	path := p.parse(0)
	return path, p.err
}

const (
	doneTok = iota
	idTok
	openTok
	closeTok
	invTok
	joinTok
	andTok
	orTok
)

var lexicon, _ = lex.New(
	`\*|[a-zA-Z_][a-zA-Z0-9_]*`,
	`\(`,
	`\)`,
	`~`,
	`/`,
	`&`,
	`\|`,
)

type parser struct {
	p    *parser
	bs   []byte
	last lex.Token
	err  error
}

func (p *parser) setErr(err error) {
	if p.err == nil {
		p.err = err
	}
}

func (p *parser) parse(prec int) Path {
	if p.err != nil {
		return nil
	}
	t := p.next()
	pref := prefixDefs[t.ID]
	if pref == nil {
		p.setErr(ErrFailedToMatch)
		return nil
	}
	left := pref(p)
	for {
		if p.err != nil || p.last.End >= len(p.bs) {
			break
		}
		t = p.peek()
		inf := infixDefs[t.ID]
		if inf.impl == nil || inf.precedence <= prec {
			break
		}
		p.last = t
		left = inf.impl(p, left)
	}
	return left
}

func (p *parser) peek() lex.Token {
	t, err := lexicon.Match(p.bs, p.last)
	p.setErr(err)
	return t
}

func (p *parser) next() lex.Token {
	p.last = p.peek()
	return p.last
}

func (p *parser) text() string {
	return p.last.Text(p.bs)
}

type prefix func(*parser) Path

type infix struct {
	precedence int
	impl       func(*parser, Path) Path
}

var prefixDefs map[int]prefix
var infixDefs map[int]infix

func init() {
	prefixDefs = map[int]prefix{
		joinTok: func(p *parser) Path {
			if p.next().ID != idTok {
				p.setErr(errors.Wrap(ErrFailedToMatch, "expecting ID"))
			}
			return Term{"/" + p.text()}
		},
		idTok: func(p *parser) Path {
			return Term{p.text()}
		},
		invTok: func(p *parser) Path {
			res := p.parse(100)
			if res == nil {
				return res
			}
			return Inverse{res}
		},
		openTok: func(p *parser) Path {
			path := p.parse(0)
			if p.next().ID != closeTok {
				p.setErr(errors.Wrap(ErrFailedToMatch, "expecting ')'"))
			}
			return path
		},
	}
	infixDefs = map[int]infix{
		joinTok: infix{80, func(p *parser, left Path) Path {
			right := p.parse(79)
			return Join{left, right}
		}},
		andTok: infix{70, func(p *parser, left Path) Path {
			right := p.parse(70)
			return Intersection{left, right}
		}},
		orTok: infix{70, func(p *parser, left Path) Path {
			right := p.parse(70)
			return Union{left, right}
		}},
	}
}

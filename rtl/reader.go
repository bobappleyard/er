package rtl

import (
	"github.com/bobappleyard/er"
	"github.com/pkg/errors"
	"strconv"
	"unicode"
	"unicode/utf8"
)

type Reader struct {
	parent  *Reader
	src     []byte
	pos     int
	attrsOK bool
	name    string
	err     error
}

func NewReader(src []byte) *Reader {
	return &Reader{src: src}
}

func (p *Reader) Next() bool {
	if !p.skipSpace() {
		return false
	}
	if p.parent != nil {
		if p.readChar() == '}' {
			p.parent.pos = p.pos
			p.parent.attrsOK = false
			return false
		}
		p.unreadChar()
	}
	if !p.parseName() {
		return false
	}
	return true
}

func (p *Reader) Name() string {
	return p.name
}

func (p *Reader) Err() error {
	return p.err
}

func (p *Reader) SetErr(err error) {
	if p.err != nil {
		return
	}
	p.err = errors.WithStack(err)
	if p.parent != nil {
		p.parent.SetErr(err)
	}
}

func (p *Reader) Record() *Reader {
	if !p.skipSpace() {
		p.SetErr(er.ErrBadSyntax)
		return nil
	}
	if p.readChar() != '{' {
		p.SetErr(er.ErrBadSyntax)
		return nil
	}
	return &Reader{
		parent:  p,
		src:     p.src,
		pos:     p.pos,
		attrsOK: true,
	}
}

func (p *Reader) StringAttr() string {
	res, err := strconv.Unquote(p.parseAttr())
	if err != nil {
		p.SetErr(err)
	}
	return res
}

func (p *Reader) IntAttr() int {
	attr := p.parseAttr()
	if p.err != nil {
		return 0
	}
	res, err := strconv.Atoi(attr[1 : len(attr)-1])
	if err != nil {
		p.SetErr(err)
	}
	return res
}

func (p *Reader) BoolAttr() bool {
	attr := p.parseAttr()
	if p.err != nil {
		return false
	}
	res, err := strconv.ParseBool(attr[1 : len(attr)-1])
	if err != nil {
		p.SetErr(err)
	}
	return res
}

func (p *Reader) ExpectEOF() {
	if p.pos < len(p.src) {
		p.SetErr(er.ErrBadSyntax)
	}
}

func (p *Reader) running() bool {
	return p != nil && p.err == nil && p.pos < len(p.src)
}

func (p *Reader) readChar() rune {
	if !p.running() {
		return 0
	}
	r, n := utf8.DecodeRune(p.src[p.pos:])
	p.pos += n
	return r
}

func (p *Reader) unreadChar() {
	_, n := utf8.DecodeLastRune(p.src[:p.pos])
	p.pos -= n
}

func (p *Reader) skipSpace() bool {
	for {
		if !p.running() {
			return false
		}
		if !unicode.IsSpace(p.readChar()) {
			p.unreadChar()
			return true
		}
	}
}

func (p *Reader) parseName() bool {
	nameStart := p.pos
	for {
		if !p.running() {
			return false
		}
		r := p.readChar()
		if r == '_' || unicode.IsLetter(r) {
			continue
		}
		if p.pos != nameStart && unicode.IsDigit(r) {
			continue
		}
		break
	}
	p.unreadChar()
	p.name = string(p.src[nameStart:p.pos])
	return true
}

func (p *Reader) parseAttr() string {
	if !p.attrsOK {
		p.SetErr(er.ErrBadSyntax)
		return ""
	}
	if !p.skipSpace() {
		p.SetErr(er.ErrBadSyntax)
		return ""
	}
	if p.readChar() != ':' {
		p.SetErr(er.ErrBadSyntax)
		return ""
	}
	if !p.skipSpace() {
		p.SetErr(er.ErrBadSyntax)
		return ""
	}
	attrStart := p.pos
	if p.readChar() != '"' {
		p.SetErr(er.ErrBadSyntax)
		return ""
	}
	for r := p.readChar(); p.running() && r != '"'; r = p.readChar() {
		if r == '\\' {
			p.readChar()
		}
	}
	return string(p.src[attrStart:p.pos])
}

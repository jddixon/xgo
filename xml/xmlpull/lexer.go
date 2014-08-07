package xmlpull

// xgo/xml/xmlpull/lexer.go

import (
	"fmt"
	"io"
)

// These functions wrap the xgo/lex functions to provide line and
// column numbers in error messages.

func (p *Parser) AcceptStr(what string) (found bool, err error) {
	lx := p.GetLexer()
	found, lxErr := lx.AcceptStr(what)
	if lxErr != nil {
		err = p.NewXmlPullError(lxErr.Error())
	}
	return
}

func (p *Parser) ExpectCh(ch rune) (err error) {
	var lxErr error
	lxErr = p.GetLexer().ExpectCh(ch)
	if lxErr != nil {
		err = p.NewXmlPullError(lxErr.Error())
	}
	return
}

// Expect at least one XML space.  If one such space is seen, skip
// any others silently.  This implementation never returns io.EOF.
//
func (p *Parser) ExpectS() (err error) {
	var lxErr error
	ch, lxErr := p.NextCh()
	if lxErr == nil || lxErr == io.EOF {
		if !p.IsS(ch) {
			msg := fmt.Sprintf("expected XML space buf found '%c'", ch)
			err = p.NewXmlPullError(msg)
		} else {
			p.SkipS()
		}
	} else {
		err = p.NewXmlPullError(lxErr.Error())
	}
	return
}

func (p *Parser) ExpectStr(what string) (err error) {
	lx    := p.GetLexer()
	lxErr := lx.ExpectStr(what)
	if lxErr != nil {
		err = p.NewXmlPullError(lxErr.Error())
	}
	return
}

func (p *Parser) Encoding() string {
	return p.GetLexer().Encoding()
}

func (p *Parser) GetOffset() int {
	return p.GetLexer().GetOffset()
}

func (p *Parser) IsS(ch rune) bool {
	return ch == '\t' || ch == '\n' || ch == '\r' || ch == ' '
}

// Get the next character from the input.
func (p *Parser) NextCh() (ch rune, err error) {
	ch, err = p.GetLexer().NextCh()
	if err != nil && err != io.EOF {
		err = p.NewXmlPullError(err.Error())
	}
	return
}

func (p *Parser) PeekCh() (r rune, err error) {
	var lxErr error
	r, err = p.GetLexer().PeekCh()
	if err != nil && err != io.EOF {
		err = p.NewXmlPullError(lxErr.Error())
	}
	return
}

// Push a character back on input.  The user must ensure that this
// operation makes sense.
//
// XXX Line and column number changes are NOT HANDLED.
func (p *Parser) PushBack(ch rune) {
	p.GetLexer().PushBack(ch)
}

// If an error is seen other than EOF, break.  Skip any of the four XML
// space characters (tab, linefeed, carriage return, space), but if a
// character other than  else is encountered, silently push it back
// and break.  If an EOF has been seen, break.
//
func (p *Parser) SkipS() {
	ch, err := p.NextCh()
	for {
		if err != nil && err != io.EOF {
			break
		}
		if !p.IsS(ch) {
			p.GetLexer().PushBack(ch)
			break
		}
		if err == io.EOF {
			break
		}
		ch, err = p.NextCh()
	}
}

package xmlpull

// xgo/xml/xmlpull/lexer.go

import ()

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

func (p *Parser) ExpectS() (err error) {
	var lxErr error
	lxErr = p.GetLexer().ExpectS()
	if lxErr != nil {
		err = p.NewXmlPullError(lxErr.Error())
	}
	return
}

func (p *Parser) ExpectStr(what string) (err error) {
	lx := p.GetLexer()
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

// Get the next character from the input.
func (p *Parser) NextCh() (ch rune, err error) {
	var lxErr error
	ch, lxErr = p.GetLexer().NextCh()
	if lxErr != nil {
		err = p.NewXmlPullError(lxErr.Error())
	}
	return
}

func (p *Parser) PeekCh() (r rune, err error) {
	var lxErr error
	r, lxErr = p.PeekCh()
	if lxErr != nil {
		err = p.NewXmlPullError(lxErr.Error())
	}
	return
}

func (p *Parser) SkipS() {
	p.GetLexer().SkipS()
}

package xmlpull

// xgo/xml/xmlpull/lexer.go

import ()

// These functions wrap the xgo/lex functions to provide line and
// column numbers in error messages.

func (p *Parser) ExpectStr(what string) (err error) {
	lx := p.GetLexer()
	lxErr := lx.ExpectStr(what)
	if lxErr != nil {
		err = p.NewXmlPullError(lxErr.Error())
	}
	return
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

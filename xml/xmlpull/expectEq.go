package xmlpull

// xgo/xml/xmlpull/expectEq.go

import (
	e "errors"
	"fmt"
)

// [25] Eq ::= S? '=' S?
//
// Return an error if the production is not matched.  This routine makes
// no change to the Parser struct.
//
func (p *Parser) expectEq() (err error) {
	p.SkipS() // skip zero or more spaces
	ch, err := p.NextCh()
	if err == nil && ch != '=' {
		msg := fmt.Sprintf(
			"parseXmlDecl.expectEq: expected = , found %c\n", ch)
		err = e.New(msg)
	}
	if err == nil {
		p.SkipS() // zero or more
	}
	return
}

package xmlpull

// xgo/xml/xmlpull/xmlpull_error.go

import (
	"fmt"
	// gl "github.com/jddixon/xgo/lex"
)

type XmlPullError struct {
	text string
}

func (lxErr *XmlPullError) Error() string {
	return lxErr.text
}

//
func (p *Parser) NewXmlPullError(msg string) (xppErr error) {
	lx := &p.LexInput
	fullMsg := fmt.Sprintf("line %d col %d: %s", lx.LineNo(), lx.ColNo(), msg)
	xppErr = &XmlPullError{fullMsg}
	return
}

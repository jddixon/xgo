package xmlpull

// xgo/xml/xmlPull/xmlpullError.go

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
	fullMsg := fmt.Sprintf("line %d col %d: %s", p.LineNo(), p.ColNo(), msg)
	xppErr = &XmlPullError{fullMsg}
	return
}

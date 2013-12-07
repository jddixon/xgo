package xmlpull

import (
	e "errors"
	"fmt"
)

// The last is of course a Java name, and will be changed or will be replaced
// by a more specific error such as UnsupportedFeature.
var (
	EmptyFeatureString       = e.New("XmlPullParser empty feature string")
	IllegalArgument          = e.New("XmlPullParser illegal argument")
	IndexOutOfBounds         = e.New("XmlPullParser index out of bounds")
	IOError                  = e.New("XmlPullParser io error")
	MissingDeclClosingQuote  = e.New("missing closing quote in xml decl")
	MustBeYesOrNo            = e.New("standalone choice must be yes or no")
	NilReader                = e.New("nil reader argument")
	NotOnStartTag            = e.New("XmlPullParser not on start tag")
	NotElementContentRelated = e.New("XmlPullParser not element-content-related")
	OnlyVersion1_0           = e.New("Only XML version 1.0 supported")
	PosOutOfRange            = e.New("XmlPullParser pos out of range")
	UnsupportedFeature       = e.New("XmlPullParser unsupported feature")
	XmlDeclMustBeAtStart	= e.New("XmlDecl must be at start of file")
	XmlInDeclMustBeLowerCase = e.New("xml in XmlDecl must be lower case")
)

func (p *Parser) NotClosedErr(what string) error {
	msg := fmt.Sprintf("%s started line %d column %d not closed", 
		what, p.startLine, p.startCol)
	return e.New(msg)
}

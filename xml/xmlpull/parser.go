package xmlpull

import (
	gl "github.com/jddixon/xgo/lex"
	gu "github.com/jddixon/xgo/util"
	"io"
)

// Any fields added here also should be added to reset()
//
type Parser struct {
	si gu.StrIntern

	// parser state
	curEvent PullToken // defined in const.go

	// element stack
	elmDepth int

	// attribute stack

	// namespaces
	nsCount  int
	nsPrefix []string
	nsUri    []string

	lineNo int // line number
	colNo  int // column number

	// buffer management
	gl.LexInput
}

func NewParser() (p *Parser, err error) {

	si := gu.NewStrIntern()

	if err == nil {
		p = &Parser{
			si: si,
		}
	}
	return
}

func (xpp *Parser) reset() {
	xpp.elmDepth = 0

	// XXX STUB XXX

}

// ==================================================================
// IMPLMENTATION OF XmlPullParserI
// ==================================================================

func (xpp *Parser) SetFeature(name string, state bool) (err error) {
	// XXX STUB XXX
	return
}
func (xpp *Parser) GetFeature(name string) (whether bool, err error) {
	// XXX STUB XXX
	return
}
func (xpp *Parser) SetProperty(name string, value interface{}) (err error) {
	// XXX STUB XXX
	return
}
func (xpp *Parser) GetProperty(name string) (prop interface{}) {
	// XXX STUB XXX
	return
}
func (xpp *Parser) SetInput(in *io.Reader) (err error) {
	if in == nil {
		err = NilReader
	} else {
		xpp.reset()
		xpp.LexInput.Reset()
	}
	return
}
func (xpp *Parser) DefineEntityReplacementText(entityName, replacementText string) (err error) {
	// XXX STUB XXX
	return
}
func (xpp *Parser) GetNamespaceCount(depth int) (ret int, err error) {
	// XXX STUB XXX
	return
}
func (xpp *Parser) GetNamespacePrefix(pos int) (ns string, err error) {
	// XXX STUB XXX
	return
}
func (xpp *Parser) GetNamespaceUri(pos int) (ns string, err error) {
	// XXX STUB XXX
	return
}
func (xpp *Parser) GetNamespaceForPrefix(prefix string) (ns string, err error) {
	// XXX STUB XXX
	return
}

//
//
func (xpp *Parser) GetDepth() int {
	return xpp.elmDepth
}
func (xpp *Parser) GetPositionDescription() (desc string) {
	// XXX STUB XXX
	return
}
func (xpp *Parser) GetLineNumber() int {
	return xpp.lineNo
}
func (xpp *Parser) GetColumnNumber() int {
	return xpp.colNo
}
func (xpp *Parser) IsWhitespace() (whether bool, err error) {
	// XXX STUB XXX
	return
}
func (xpp *Parser) GetText() (t string) {
	// XXX STUB XXX
	return
}

// Should probably return a string.
//
func (xpp *Parser) GetTextCharacters(holderForStartAndLength []int) (
	chars []byte) {
	// XXX STUB XXX
	return
}
func (xpp *Parser) GetNamespace() (ns string, err error) {
	// XXX STUB XXX
	return
}
func (xpp *Parser) GetName() (name string) {
	// XXX STUB XXX
	return
}
func (xpp *Parser) GetPrefix() (prefix string) {
	// XXX STUB XXX
	return
}
func (xpp *Parser) IsEmptyElementTag() (whether bool, err error) {
	// XXX STUB XXX
	return
}
func (xpp *Parser) GetAttributeCount() (count int) {
	// XXX STUB XXX
	return
}
func (xpp *Parser) GetAttributeNamespace(index int) (
	ns string, err error) {

	// XXX STUB XXX
	return
}
func (xpp *Parser) GetAttributeName(index int) (name string) {
	// XXX STUB XXX
	return
}
func (xpp *Parser) GetAttributePrefix(index int) (p string, err error) {
	// XXX STUB XXX
	return
}
func (xpp *Parser) GetAttributeValue(index int) (val string, err error) {
	// XXX STUB XXX
	return
}
func (xpp *Parser) GetAttributeValueNS(namespace, name string) (
	ns string, err error) {

	// XXX STUB XXX
	return
}

// Return the type of the current parser event.  The spec requires the
// error return.
//
func (xpp *Parser) GetEventType() (PullToken, error) {
	return xpp.curEvent, nil
}

func (xpp *Parser) Next() (ret int, err error) {
	// XXX STUB XXX
	return
}
func (xpp *Parser) NextToken() (ret int, err error) {
	// XXX STUB XXX
	return
}
func (xpp *Parser) Require(type_ int, namespace, name string) {
	// XXX STUB XXX
	return
}
func (xpp *Parser) ReadText() (t string) {
	// XXX STUB XXX
	return
}

package xmlpull

import (
	gl "github.com/jddixon/xgo/lex"
	gu "github.com/jddixon/xgo/util"
	"io"
)

// Any fields added here also should be added to reset()
//
type Parser struct {
	xmlDeclVersion, xmlDeclEncoding string
	xmlDeclStandalone               bool
	docTypeDecl                     string

	tokenizing         bool
	roundtripSupported bool

	startLine, startCol int // where a syntactic element begins

	// accumulated characters of various types -- yes, kludgey
	text         []rune
	cDataChars   string
	commentChars string
	piChars      string
	piTarget     string

	// parser state
	curEvent PullEvent // aka eventType; PullEvent defined in const.go

	afterLT bool // have encountered a left angle bracket (<)

	// element stack
	elmDepth int

	// attribute stack

	// namespaces
	nsCount  int
	nsPrefix []string
	nsUri    []string

	si gu.StrIntern

	// buffer management
	lineNo int // line number		// redundant
	colNo  int // column number		// redundant

	gl.LexInput
}

func (p *Parser) start() {
	p.startLine = p.LineNo()
	p.startCol = p.ColNo()
}

// Return an XmlPullParser with the default encoding
func NewNewParser(reader io.Reader) (*Parser, error) {
	return NewParser(reader, "")
}

func NewParser(reader io.Reader, encoding string) (p *Parser, err error) {

	var lx *gl.LexInput
	if reader == nil {
		err = NilReader
	} else {
		lx, err = gl.NewLexInput(reader, encoding)
	}

	if err == nil {
		si := gu.NewStrIntern()
		p = &Parser{
			si:       si,
			LexInput: *lx,
		}
		p.reset()
	}
	return
}

// Return a pointer to the parser's lexer.
func (xpp *Parser) GetLexer() *gl.LexInput {
	return &xpp.LexInput
}

// All fields in the Parser struct should be reinitialized here.
func (xpp *Parser) reset() {
	xpp.afterLT = false
	xpp.colNo = 0
	xpp.elmDepth = 0
	xpp.lineNo = 1
	xpp.nsCount = 0

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
func (xpp *Parser) GetEventType() (PullEvent, error) {
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

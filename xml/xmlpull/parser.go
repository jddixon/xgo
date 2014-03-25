package xmlpull

import (
	gl "github.com/jddixon/xgo/lex"
	"io"
)

const (
	READ_CHUNK_SIZE = 4096 // or multiple
)

// Any fields added here also should be added to reset()
//
type Parser struct {
	xmlDeclVersion, xmlDeclEncoding string
	xmlDeclStandalone               bool
	docTypeDecl                     string

	// parser behavior
	normalizeIgnorableWS bool
	processNamespaces    bool
	roundtripSupported   bool

	startLine, startCol int // where a syntactic element begins

	haveRootTag bool

	// accumulated characters of various types -- yes, kludgey
	cDataChars   []rune
	commentChars []rune
	piChars      []rune
	piTarget     []rune

	// transient variables for NextToken()
	tokenizing    bool
	text          []rune
	entityRefName []rune

	// global parser state
	seenStartTag  bool
	seenEndTag    bool
	pastEndTag    bool
	seenAmpersand bool
	afterLT       bool // have encountered a left angle bracket (<)
	seenDocdecl   bool

	curEvent       PullEvent // aka eventType; PullEvent defined in const.go
	isEmptyElement bool

	// element stack
	elmDepth      int
	elRawName     [][]rune
	elRawNameLine []int

	elName           [][]rune
	elPrefix         []string // []rune
	elUri            []string // [][]rune
	elValue          [][]rune
	elNamespaceCount []int

	// attribute stack
	attributeCount    int
	attributeName     []string //[]rune
	attributeNameHash []uint32 // stores FastHash output
	attributePrefix   []string // []rune
	attributeUri      []string // [][]rune
	attributeValue    []string // []rune

	// namespace stack
	namespaceCount      int
	namespacePrefix     []string //[]rune
	namespacePrefixHash []uint32 // more FastHash output
	namespaceUri        []string // []rune

	namespaceEnd int

	// entity replacement stack ---------------------------
	entityEnd            int
	entityName           [][]rune
	entityNameBuf        [][]rune
	entityReplacement    [][]rune
	entityReplacementBuf [][]rune
	entityNameHash       []uint32

	// buffer management ----------------------------------
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
		p = &Parser{
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

// Reset all global variables; also used to set initial values at
// program start.
//
// All fields in the Parser struct should be reinitialized here.

func (xpp *Parser) reset() {
	xpp.afterLT = false
	xpp.colNo = 0
	xpp.elmDepth = 0
	xpp.lineNo = 1
	xpp.namespaceCount = 0

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

func (xpp *Parser) Require(type_ int, namespace, name string) {
	// XXX STUB XXX
	return
}
func (xpp *Parser) ReadText() (t string) {
	// XXX STUB XXX
	return
}

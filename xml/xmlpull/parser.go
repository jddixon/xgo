package xmlpull

import (
	"fmt"
	gl "github.com/jddixon/xgo/lex"
	"io"
)

var _ = fmt.Println

const (
	READ_CHUNK_SIZE = 4096 // or multiple
)

// Any fields added here also should be added to reset()
//
type Parser struct {
	// VETTED FIELDS ////////////////////////////////////////////////
	curEvent PullEvent   //  defined in const.go
	state    ParserState // states also defined in const.go

	// defined in the prolog
	xmlDeclVersion, xmlDeclEncoding string
	xmlDeclContent                  string
	xmlDeclStandalone               bool
	docTypeDecl                     string

	// FIELDS MADE UNNECESSARY BY Parser.state //////////////////////

	// doNext State -------------------------------------------------
	isEmptyElement bool
	pastEndTag     bool
	rootElmSeen    bool
	seenStartTag   bool

	// OTHER FIELDS /////////////////////////////////////////////////

	// parser behavior
	normalizeIgnorableWS bool
	processNamespaces    bool
	roundtripSupported   bool

	startLine, startCol int // where a syntactic element begins

	// accumulated characters of various types -- yes, kludgey
	cDataChars   []rune
	commentChars []rune
	piChars      []rune
	piTarget     []rune

	// transient variables for NextToken()
	tokenizing    bool
	text          []rune
	entityRefName string // []rune

	// global parser state
	lineNo int // line number		// redundant
	colNo  int // column number		// redundant

	location      string
	reachedEnd    bool // used only in parseEpilog?
	seenEndTag    bool
	seenAmpersand bool
	afterLT       bool // have encountered a left angle bracket (<)
	seenDocdecl   bool

	// element stack
	elmDepth      uint
	elRawName     [][]rune
	elRawNameLine []int

	elName           []string // []rune
	elPrefix         []string // []rune
	elUri            []string // [][]rune
	elValue          [][]rune
	elNamespaceCount []uint

	// attribute stack
	attributeCount    uint
	attributeName     []string //[]rune
	attributeNameHash []uint32 // stores FastHash output
	attributePrefix   []string // []rune
	attributeUri      []string // [][]rune
	attributeValue    []string // []rune

	// namespace stack
	namespaceCount      uint
	namespacePrefix     []string //[]rune
	namespacePrefixHash []uint32 // more FastHash output
	namespaceUri        []string // []rune

	namespaceEnd uint

	// entity replacement stack ---------------------------
	entityEnd            int
	entityName           []string // [][]rune
	entityNameBuf        [][]rune
	entityReplacement    []string // ][]rune
	entityReplacementBuf [][]rune
	entityNameHash       []uint32

	// XXX hmmm
	inputEncoding string
	gl.LexInput
}

// Called at the beginning of a syntactic construct.
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
		// --------------------------------------------------------------
		// BOM (Byte Order Mark) - not in the syntax graph, examines the
		// very first byte in a document.  This code should precede the
		// call to parseProlog(); look at the character and then panic,
		// push it back, or discard it.
		// --------------------------------------------------------------
		// This block analyzes the very first byte in a document.
		var ch rune
		ch, err = p.NextCh()
		// This is the first character of input, and so might be the
		// unicode byte order mark (BOM)
		if ch == '\uFFFE' {
			panic("data in wrong byte order!")
		} else if ch != '\uFEFF' {
			p.PushBack(ch)
		}
		p.state = PRE_START_DOC
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
	// VETTED FIELDS ------------------------------------------------
	xpp.curEvent = START_DOCUMENT

	// OTHER FIELDS -------------------------------------------------
	xpp.afterLT = false
	xpp.colNo = 0
	xpp.elmDepth = 0
	xpp.lineNo = 1
	xpp.namespaceCount = 0

	// XXX STUB XXX

}
